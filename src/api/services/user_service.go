package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"yip/src/api/auth/verifier"
	"yip/src/api/services/dto"
	"yip/src/common"
	"yip/src/config"
	"yip/src/cryptox"
	"yip/src/repositories"
	"yip/src/repositories/repo"
	"yip/src/slyerrors"
)

type UserService struct {
	Config   *config.Config
	verifier *verifier.Verifier
	userDB   repositories.Database
	repos    *repo.Repositories
}

func NewUserService(
	config *config.Config,
	verifier *verifier.Verifier,
	useDB repositories.Database,
	repos *repo.Repositories,
) UserService {
	return UserService{
		Config:   config,
		verifier: verifier,
		userDB:   useDB,
		repos:    repos,
	}
}

func (s UserService) SignInYIPAdmin(context context.Context, data *dto.SignInRequest) (*verifier.Token, error) {
	admin := s.Config.API.Admin
	if data.Email != admin.Username {
		return nil, slyerrors.Forbidden("403", "user is not allowed")
	}

	fmt.Println("body password", data.Password)
	fmt.Println("hashed", admin.PasswordHashed)
	if cryptox.CheckPasswordHash(data.Password, admin.PasswordHashed) {
		return s.verifier.CreateToken(data.Audiences, "0000-0000-0000", "", "", verifier.RoleAdmin)
	}

	return nil, slyerrors.BadRequest("400", "password is incorrect")
}

func (s UserService) SignInUser(context context.Context, data *dto.SignInRequest) (*verifier.Token, error) {
	user, err := s.userDB.GetAccountByEmail(context, data.Email)
	if err != nil {
		return nil, err
	}

	if cryptox.CheckPasswordHash(data.Password, user.PasswordHashed) {
		return s.verifier.CreateToken(data.Audiences, user.ID, "", "", user.Role)
	}

	return nil, slyerrors.BadRequest("400", "password is incorrect")
}

func (s UserService) RegisterUser(ctx context.Context, request *dto.RegisterRequest) (*repo.AccountModel, error) {
	passwordHash, err := cryptox.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	role := verifier.RoleAdmin

	return s.repos.AccountRepo.Create(ctx, &repo.AccountModel{
		Email:           request.Email,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		PasswordHashed:  passwordHash,
		Role:            role,
	})
}

func (s UserService) GetAccountById(context context.Context, id uuid.UUID) (repositories.UserAccount, error) {
	return s.userDB.GetAccountById(context, id)
}

func (s UserService) GetSLYWalletsByAccountId(ctx context.Context, id uuid.UUID) (dto.UserInfoResponse, error) {
	result := dto.UserInfoResponse{}
	var err error

	result.Account, err = s.userDB.GetAccountById(ctx, id)
	if err != nil {
		return result, err
	}

	result.SLYWallets, err = s.userDB.GetSLYWallets(ctx, id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s UserService) GetUsers(context context.Context, query common.PaginationQuery) (repositories.ListUsersAccountResponse, error) {
	return s.userDB.GetAccounts(context, query.PageSize, query.Offset)
}

func (s UserService) SetRole(context context.Context, id uuid.UUID, role string) (repositories.UserAccount, error) {
	//if !repositories.RoleExists(role) {
	//	return repositories.UserAccount{}, slyerrors.BadRequest("400", "role does not exist")
	//}

	return s.userDB.SetRole(context, id, role)
}

func (s UserService) SetEmail(context context.Context, id uuid.UUID, email string) (repositories.UserAccount, error) {
	return s.userDB.SetEmail(context, id, email)
}
