package slywallet

import (
	"github.com/go-chi/chi/v5"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/api/slywallet/base"
)

type Module struct {
	BaseController base.Controller
}

func NewModule(
	services *services.Services,
	middleware *verifier.TokenVerifierMiddleware,
) Module {
	return Module{
		BaseController: base.NewController(middleware, services),
	}
}

func (a Module) Routes() func(r chi.Router) {

	return func(r chi.Router) {
		r.Route("/", a.BaseController.Routes())
	}
}
