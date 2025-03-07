package services

import (
	"yip/src/api/auth/pin"
	"yip/src/app"
	"yip/src/repositories/repo"
)

type Services struct {
	PinService       pin.Service
	UserService      UserService
	TokenService     TokenService
	SIWEService      SIWEService
	SLYWalletService *SLYWalletService
	Repos            *repo.Repositories
}

func GenerateApiServices(app *app.App) Services {
	repos := repo.NewRepositories(app.DB)

	slyWalletService := NewSLYWalletService(app.Config, app.EthProvider, repos)

	return Services{
		PinService:       pin.NewService(app.Config, app.Verifier, app.UserDB, &app.EmailProvider),
		UserService:      NewUserService(app.Config, app.Verifier, app.UserDB),
		TokenService:     NewTokenService(app.Config, app.Verifier, app.UserDB, app.EthProvider),
		SIWEService:      NewSIWEService(app.Config, app.Verifier, app.UserDB, app.EthProvider),
		SLYWalletService: &slyWalletService,
		Repos:            repos,
	}
}
