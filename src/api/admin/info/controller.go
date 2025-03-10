package info

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/config"
	"yip/src/httpx"
	"yip/src/providers"
)

type Controller struct {
	conf               *config.Config
	service            *services.InvitationCodeService
	ethProvider        *providers.EthProvider
	yipAdminMiddleware *verifier.TokenVerifierMiddleware
}

func NewController(
	conf *config.Config,
	service *services.InvitationCodeService,
	ethProvider *providers.EthProvider,
	tokenMiddleware *verifier.TokenVerifierMiddleware) Controller {
	return Controller{
		conf:               conf,
		service:            service,
		ethProvider:        ethProvider,
		yipAdminMiddleware: tokenMiddleware,
	}
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(c.yipAdminMiddleware.PrincipalCtx)
			r.Use(AdminCtx)
			r.Get("/chain", c.GetChainInfo)
			r.Get("/codes", c.GetCodes)
		})
	}
}

// swagger:route GET /admin/info/chain info
// Returns the blockchain setup
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: ChainInfoDTO
func (c Controller) GetChainInfo(w http.ResponseWriter, r *http.Request) {
	result, err := c.ethProvider.BalanceOf(r.Context(), c.ethProvider.PubKeySignerWallet)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError(err.Error()))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(ChainInfoDTO{
		WalletAddress: c.ethProvider.PubKeySignerWallet.Hex(),
		RPCUrl:        c.conf.EthConfig.Chain.RPCUrl,
		ChainId:       c.conf.EthConfig.Chain.ID,
		WalletName:    c.conf.EthConfig.Wallet.Name,
		WalletValue:   fmt.Sprintf("%s ETH", result.String()),
	}))
}

// swagger:route GET /admin/info/codes info
// All Codes
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: []Pin
func (c Controller) GetCodes(w http.ResponseWriter, r *http.Request) {
	codes, err := c.service.GetAllCodes(r.Context())
	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError(err.Error()))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(codes))
}
