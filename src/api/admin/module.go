package admin

import (
	"github.com/go-chi/chi/v5"
	"yip/src/api/admin/info"
	"yip/src/api/admin/user"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/config"
	"yip/src/providers"
)

type AdminModule struct {
	UserController user.Controller
	InfoController info.Controller
}

func NewAdminModule(
	config *config.Config,
	services *services.Services,
	middleware *verifier.TokenVerifierMiddleware,
	ethProvider *providers.EthProvider,
) AdminModule {
	return AdminModule{
		UserController: user.NewController(&services.UserService, &services.PinService, middleware),
		InfoController: info.NewController(config, &services.InvitationCodeService, ethProvider, middleware),
	}
}

func (a AdminModule) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/accounts", a.UserController.Routes())
		r.Route("/info", a.InfoController.Routes())
	}
}
