package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"log"
	"yip/src/api/middleware"
	"yip/src/api/services/dto"
	"yip/src/config"
	"yip/src/contracts"
	"yip/src/httpx"
	"yip/src/providers"
	"yip/src/repositories/repo"
)

type SLYWalletService struct {
	slyWalletManager *contracts.WalletManager
	ByIdMiddleware   middleware.EntityMiddleware[*dto.SLYBase]
	repos            *repo.Repositories
	Config           *config.Config
}

func NewSLYWalletService(
	config *config.Config,
	slyWalletManager *contracts.WalletManager,
	repos *repo.Repositories,
) *SLYWalletService {
	s := SLYWalletService{
		slyWalletManager: slyWalletManager,
		repos:            repos,
		Config:           config,
	}

	return &s
}

func (s SLYWalletService) SpawnSLYWallet(
	ctx context.Context,
	userOwnerKey common.Address,
	invitationCode string,
) (*providers.TransactionTicket, error) {
	ticket, err := s.slyWalletManager.SpawnWalletWithGas(ctx, userOwnerKey, 0, nil)
	if err != nil {
		return nil, err
	}

	_, err = s.repos.InvitationCodeRepo.Update(ctx, &repo.InvitationCodeModel{
		Code:            invitationCode,
		TransactionHash: ticket.TransactionHash,
	})
	if err != nil {
		// no reason to waste a transaction
		log.Println(err.Error())
	}

	return ticket, nil
}

func (s SLYWalletService) GetReceipt(ctx context.Context, hash string) (*providers.TransactionState, error) {
	h := common.HexToHash(hash)
	r, err := s.slyWalletManager.GetTransactionReceipt(h)
	if err != nil {
		return nil, err
	}
	state, err := s.slyWalletManager.GetTransactionStatusByReceipt(h, r, nil)
	if err != nil {
		return nil, err
	}

	if state.ContractAddress != "" {
		// 1. check if wallet already exist
		_, err := s.repos.SlyWalletRepo.GetByAddress(ctx, state.ContractAddress)

		if err != nil {
			if err != repo.DBItemNotFound {
				return nil, fmt.Errorf("failed to get SlyWallet by address: %w", err)
			}

			code, err := s.repos.InvitationCodeRepo.GetByTransactionHash(ctx, state.TransactionHash)
			if err != nil {
				return nil, fmt.Errorf("failed to get SlyWallet by address: %w", err)
			}

			codeFull, err := s.repos.InvitationCodeRepo.GetWithAccounts(ctx, code.Code)
			if err != nil {
				return nil, fmt.Errorf("no accounts attached to the transaction hash (via invitation code): %w", err)
			}
			if codeFull.Accounts == nil || len(codeFull.Accounts) == 0 {
				return nil, fmt.Errorf("no accounts attached to the transaction hash (via invitation code): %w", err)
			}

			accountId := codeFull.Accounts[0].ID

			_, err = s.CreateSLYWalletEntry(ctx, state, state.TransactionHash, state.ContractAddress, accountId, codeFull.Code)
			if err != nil {
				return nil, fmt.Errorf("failed to create SlyWallet entry: %w", err)
			}
		}
	}
	return state, nil
}

func (s SLYWalletService) CreateSLYWalletEntry(
	ctx context.Context,
	transactionStatus *providers.TransactionState,
	transactionHash string,
	contractAddress string,
	accountId uuid.UUID,
	invitationCode string,
) (*repo.SlyWalletModel, error) {
	_, err := s.repos.SlyWalletRepo.GetByAddress(ctx, contractAddress)

	if err != nil {
		if err != repo.DBItemNotFound {
			return nil, fmt.Errorf("failed to get SlyWallet by address: %w", err)
		}

		wallets, err := s.slyWalletManager.GetWalletKeys(common.HexToAddress(contractAddress))
		if err != nil {
			return nil, fmt.Errorf("could not create wallet from contract address")
		}

		slyWallet, err := s.repos.SlyWalletRepo.Create(ctx, &repo.SlyWalletModel{
			Address:           contractAddress,
			Chainid:           s.Config.EthConfig.Chain.ID,
			AccountID:         accountId,
			TransactionHash:   transactionHash,
			TransactionStatus: transactionStatus.Status,
			InvitationCode:    invitationCode,
		})
		if err != nil {
			return nil, fmt.Errorf("could not create a slywallet : %w", err)
		}

		for _, owner := range wallets.Owners {
			_, err = s.repos.EcdsaRepo.UpsertECDSA(ctx, &repo.EcdsaModel{
				Address:   owner.Hex(),
				AccountID: accountId,
			})
			if err != nil {
				return nil, fmt.Errorf("could not add ecdsa key : %w", err)
			}
			_, err = s.repos.EcdsaSlyWalletRepo.AddECDSAToSlyWallet(ctx, owner.Hex(), slyWallet.Address, dto.RoleOwner)
			if err != nil {
				return nil, fmt.Errorf("could not add ecdsa key : %w", err)
			}
		}

		for _, admin := range wallets.Admins {
			_, err = s.repos.EcdsaRepo.UpsertECDSA(ctx, &repo.EcdsaModel{
				Address:   admin.Hex(),
				AccountID: accountId,
			})
			if err != nil {
				return nil, fmt.Errorf("could not add ecdsa key : %w", err)
			}
			_, err = s.repos.EcdsaSlyWalletRepo.AddECDSAToSlyWallet(ctx, admin.Hex(), slyWallet.Address, dto.RoleAdmin)
			if err != nil {
				return nil, fmt.Errorf("could not add ecdsa key : %w", err)
			}
		}

		for _, authenticator := range wallets.Authenticators {
			_, err = s.repos.EcdsaRepo.UpsertECDSA(ctx, &repo.EcdsaModel{
				Address:   authenticator.Hex(),
				AccountID: accountId,
			})
			if err != nil {
				return nil, fmt.Errorf("could not add ecdsa key : %w", err)
			}
			_, err = s.repos.EcdsaSlyWalletRepo.AddECDSAToSlyWallet(ctx, authenticator.Hex(), slyWallet.Address, dto.RoleAdmin)
			if err != nil {
				return nil, fmt.Errorf("could not add ecdsa key : %w", err)
			}
		}
	}

	return s.repos.SlyWalletRepo.GetByAddress(ctx, contractAddress)
}

func (s SLYWalletService) IsControllerKeyOf(controllerKey common.Address, slyWalletAddress common.Address) *httpx.Response {
	result := s.slyWalletManager.AuthenticateControllerKeyOfSLYWallet(controllerKey, slyWalletAddress)
	switch result.StatusCode {
	case 400:
		return httpx.BadRequest(result.ErrorDetails)
	case 500:
		return httpx.InternalError(result.ErrorDetails)
	case 403:
		return httpx.Forbidden(result.ErrorDetails)
	default:
		return nil
	}
}
