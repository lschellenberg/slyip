package base

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/api/middleware"
	"yip/src/api/services"
	"yip/src/api/services/dto"
	"yip/src/httpx"
)

type Controller struct {
	yipAdminMiddleware *verifier.TokenVerifierMiddleware
	slyService         *services.SLYWalletService
	meMiddleware       middleware.MeMiddleware[*dto.SLYBase]
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

			r.Route("/me", func(r chi.Router) {
				r.Use(c.yipAdminMiddleware.PrincipalCtx)
				r.Use(c.meMiddleware.EntityContext)

				r.Get("/", c.GetMe)
			})

		})
	}
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

func (c Controller) GetSLYWalletById(w http.ResponseWriter, r *http.Request) {
	m := c.slyService.ByIdMiddleware.EntityFromCtx(r)
	httpx.RespondWithJSON(w, httpx.OK(m))
}

func (c Controller) spawnSLYWallet(w http.ResponseWriter, r *http.Request) {
	principal, err := verifier.GetPrincipal(r.Context())
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}
	// TODO
	// parse body
	// check validity of invitation code
	// check if SLYWallet already exists
	// spawn SLYWallet
	// return transaction hash

	fmt.Println(principal)
}

func (c Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	httpx.RespondWithJSON(w, httpx.OK(c.meMiddleware.EntityFromCtx(r)))
}
