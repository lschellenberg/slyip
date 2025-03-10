package services

import (
	"context"
	"log"
	"time"
	"yip/src/common"
	"yip/src/cryptox"
	"yip/src/repositories/repo"
)

const CodeLength = 16

type InvitationCodeService struct {
	repos *repo.Repositories
}

func NewInvitationCodeService(repos *repo.Repositories) InvitationCodeService {
	return InvitationCodeService{
		repos: repos,
	}
}

func (s InvitationCodeService) GetAllCodes(ctx context.Context) (*repo.PaginatedResponse[repo.InvitationCodeModel], error) {
	codes, err := s.repos.InvitationCodeRepo.ListValidCodes(ctx, &common.PaginationQuery{
		QueryFilters: common.QueryFilters{},
		PageSize:     1000,
		Offset:       0,
	})

	if err != nil {
		return nil, err
	}

	if len(codes.Data) <= 20 {
		log.Println("need to add invitation codes...")
		for i := 0; i <= 20; i++ {

			_, err := s.repos.InvitationCodeRepo.Create(ctx, &repo.InvitationCodeModel{
				Code:      cryptox.GenerateNumberCode(CodeLength),
				ExpiresAt: time.Now().Add(100000 * time.Hour),
			})
			if err != nil {
				return nil, err
			}
		}
		return s.GetAllCodes(ctx)
	}

	return codes, nil
}

func (s InvitationCodeService) ValidateCode(ctx context.Context, code string) (bool, error) {
	iCode, err := s.repos.InvitationCodeRepo.GetByCode(ctx, code)
	if err != nil {
		return false, err
	}

	return iCode.IsValid(), nil
}
