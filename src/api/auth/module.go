package auth

import (
	"github.com/go-chi/chi/v5"
	"yip/src/api/auth/pin"
	"yip/src/api/auth/session"
	"yip/src/api/auth/siwe"
	"yip/src/api/auth/token"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/config"
)

type Module struct {
	TokenController   token.Controller
	SIWEController    siwe.Controller
	PinController     pin.Controller
	SessionController session.Controller
}

func NewAuthModule(
	config *config.Config,
	services *services.Services,
	middleware *verifier.TokenVerifierMiddleware,
) Module {
	return Module{
		TokenController:   token.NewController(&services.TokenService, &services.UserService, middleware),
		SIWEController:    siwe.NewController(config, &services.SIWEService, &services.UserService),
		PinController:     pin.NewController(&services.PinService),
		SessionController: session.NewController(config, &services.SIWEService, &services.UserService),
	}
}

func (a Module) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/token", a.TokenController.Routes())
		r.Route("/siwe", a.SIWEController.Routes())
		r.Route("/pin", a.PinController.Routes())
		r.Route("/session", a.SessionController.Routes())
	}
}
