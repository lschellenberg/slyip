package base

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/api/middleware"
	"yip/src/api/services"
	"yip/src/api/services/dto"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

type Controller struct {
	yipAdminMiddleware *verifier.TokenVerifierMiddleware
	slyService         *services.SLYWalletService
	accountMiddleware  middleware.EntityMiddleware[*dto.SLYBase]
	allServices        *services.Services
}

func NewController(
	tokenMiddleware *verifier.TokenVerifierMiddleware,
	allServices *services.Services,
) Controller {

	c := Controller{
		yipAdminMiddleware: tokenMiddleware,
		slyService:         allServices.SLYWalletService,
		allServices:        allServices,
	}

	return c
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {

			r.Use(c.yipAdminMiddleware.PrincipalCtx)
			r.Post("/spawn", c.spawnSLYWallet)
			r.Get("/receipt/{hash}", c.GetSLYWalletReceipt)

		})
	}
}

func (c Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	principal, err := verifier.GetPrincipal(r.Context())
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}
	httpx.RespondWithJSON(w, httpx.OK(principal.ID))
	//httpx.RespondWithJSON(w, httpx.OK(c.meMiddleware.EntityFromCtx(r)))
}

func (c Controller) spawnSLYWallet(w http.ResponseWriter, r *http.Request) {
	principal, err := verifier.GetPrincipal(r.Context())
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}

	// check if ecdsa is valid eth address
	if !slyerrors.IsValidEthAddress(principal.ECDSAAddress) {
		httpx.RespondWithJSON(w, httpx.BadRequest("ecdsa address in token is not real ecdsa"))
		return
	}

	//check if ecdsa address is already a controller key of a sly wallet
	wallet, err := c.allServices.Repos.EcdsaSlyWalletRepo.GetByEcdsaAddress(r.Context(), principal.ECDSAAddress)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError("ecdsa address in token is not real ecdsa"))
		return
	}
	if len(wallet) != 0 {
		httpx.RespondWithJSON(w, httpx.BadRequest("ecdsa address is already controller key of a wallet"))
		return
	}

	// parse body
	body := &dto.CreateSLYWalletRequest{}
	if err := body.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}

	// check validity of invitation code
	ok, err := c.allServices.InvitationCodeService.ValidateCode(r.Context(), body.InvitationCode)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError(err.Error()))
		return
	}
	if !ok {
		httpx.RespondWithJSON(w, httpx.InternalError("invitation code is not valid"))
		return
	}

	// check if SLYWallet already exists
	account, err := c.allServices.AccountService.GetCompleteAccount(r.Context(), principal.ID)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	if len(account.SlyWallets) > 0 {
		httpx.RespondWithJSON(w, httpx.BadRequest("account already has a SLYWallet attached"))
		return
	}

	// spawn wallet and return ticket
	ticket, err := c.slyService.SpawnSLYWallet(r.Context(), common.HexToAddress(principal.ECDSAAddress), body.InvitationCode)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	_, err = c.allServices.Repos.AccountRepo.SetInvitationCode(r.Context(), account.ID, body.InvitationCode)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(ticket))
}

func (c Controller) GetSLYWalletReceipt(w http.ResponseWriter, r *http.Request) {
	if transactionHash := chi.URLParam(r, "hash"); transactionHash != "" {
		status, err := c.slyService.GetReceipt(r.Context(), transactionHash)

		if err != nil {
			fmt.Println(err.Error())
			httpx.RespondWithJSON(w, httpx.MapServiceError(err))
			return
		}

		httpx.RespondWithJSON(w, httpx.OK(status))
	} else {
		httpx.RespondWithJSON(w, httpx.BadRequest("no transactionHash given"))
	}
}
