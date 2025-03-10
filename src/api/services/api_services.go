package services

import (
	"yip/src/api/auth/pin"
	"yip/src/app"
	"yip/src/repositories/repo"
)

type Services struct {
	PinService            pin.Service
	UserService           UserService
	TokenService          TokenService
	SIWEService           SIWEService
	AccountService        AccountService
	SLYWalletService      *SLYWalletService
	Repos                 *repo.Repositories
	InvitationCodeService InvitationCodeService
}

func GenerateApiServices(app *app.App) Services {
	repos := repo.NewRepositories(app.DB)

	return Services{
		PinService:            pin.NewService(app.Config, app.Verifier, app.UserDB, &app.EmailProvider, repos),
		UserService:           NewUserService(app.Config, app.Verifier, app.UserDB, repos),
		TokenService:          NewTokenService(app.Config, app.Verifier, app.UserDB, app.EthProvider),
		SIWEService:           NewSIWEService(app.Config, app.Verifier, app.UserDB, app.EthProvider, app.SLYWalletManager),
		AccountService:        NewAccountService(repos),
		InvitationCodeService: NewInvitationCodeService(repos),
		SLYWalletService:      NewSLYWalletService(app.Config, app.SLYWalletManager, repos),
		Repos:                 repos,
	}
}
