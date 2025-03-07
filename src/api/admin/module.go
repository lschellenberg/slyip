package admin

import (
	"github.com/go-chi/chi/v5"
	"yip/src/api/admin/user"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
)

type AdminModule struct {
	UserController user.Controller
}

func NewAdminModule(
	services *services.Services,
	middleware *verifier.TokenVerifierMiddleware,
) AdminModule {
	return AdminModule{
		UserController: user.NewController(&services.UserService, &services.PinService, middleware),
	}
}

func (a AdminModule) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/accounts", a.UserController.Routes())
	}
}
