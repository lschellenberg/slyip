package services

import (
	"context"
	"yip/src/api/auth/verifier"
	"yip/src/config"
	"yip/src/providers"
	"yip/src/repositories"
)

type TokenService struct {
	config      *config.Config
	verifier    *verifier.Verifier
	db          repositories.Database
	ethProvider *providers.EthProvider
}

func NewTokenService(
	config *config.Config,
	verifier *verifier.Verifier,
	db repositories.Database,
	ethProvider *providers.EthProvider,
) TokenService {
	return TokenService{
		config:      config,
		verifier:    verifier,
		db:          db,
		ethProvider: ethProvider,
	}
}

func (s TokenService) RefreshToken(token string) (*verifier.Token, error) {
	return s.verifier.RefreshToken(token)
}

func (s TokenService) VerifyToken(token string) (*verifier.Principal, error) {
	return s.verifier.VerifyToken(context.Background(), token)
}
