package services

import (
	"context"
	"github.com/google/uuid"
	"yip/src/httpx"
	"yip/src/repositories/repo"
)

type AccountService struct {
	Repo *repo.Repositories
}

func NewAccountService(repo *repo.Repositories) AccountService {
	return AccountService{
		repo,
	}
}

func (accountService AccountService) GetCompleteAccount(ctx context.Context, accountId string) (*repo.AccountModel, error) {
	aid, err := uuid.Parse(accountId)
	if err != nil {
		return nil, httpx.BadRequest("accountId not uuid")
	}
	return accountService.Repo.AccountRepo.GetCompleteAccount(ctx, aid)
}
