package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"
	"net/url"
	"strings"
	"yip/src/api/auth/verifier"
	"yip/src/api/services/dto"
	"yip/src/config"
	"yip/src/contracts"
	"yip/src/providers"
	"yip/src/repositories"
	"yip/src/slyerrors"
)

type SIWEService struct {
	config           *config.Config
	verifier         *verifier.Verifier
	userDB           repositories.Database
	ethProvider      *providers.EthProvider
	slyWalletManager *contracts.WalletManager
}

func NewSIWEService(
	config *config.Config,
	verifier *verifier.Verifier,
	useDB repositories.Database,
	ethProvider *providers.EthProvider,
	slyWalletManager *contracts.WalletManager,
) SIWEService {
	return SIWEService{
		config:           config,
		verifier:         verifier,
		userDB:           useDB,
		ethProvider:      ethProvider,
		slyWalletManager: slyWalletManager,
	}
}

func (s SIWEService) Challenge(data *dto.ChallengeRequestDTO) (*dto.ChallengeResponse, error) {
	domainURL, err := url.Parse(data.Domain)
	if err != nil {
		return nil, err
	}
	m, err := siwe.InitMessage(domainURL.Host, data.Address, data.Domain, siwe.GenerateNonce(), map[string]interface{}{"chainId": data.ChainId})
	if err != nil {
		return nil, err
	}

	return &dto.ChallengeResponse{
		Challenge: m.String(),
		Address:   data.Address,
		Domain:    data.Domain,
		ChainId:   data.ChainId,
	}, err
}

func (s SIWEService) Nonce() dto.NonceResponse {
	return dto.NonceResponse{
		Nonce: siwe.GenerateNonce(),
	}
}

func (s SIWEService) Verify(data *dto.SubmitRequestDTO) (*dto.VerifyResponse, error) {
	m, err := siwe.ParseMessage(data.Message)
	if err != nil {
		return nil, err
	}

	recoveredKey, err := m.Verify(data.Signature, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	u := m.GetURI()
	m.GetDomain()
	return &dto.VerifyResponse{
		Domain:           m.GetDomain(),
		URI:              u.String(),
		OriginalAddress:  m.GetAddress().String(),
		RecoveredAddress: crypto.PubkeyToAddress(*recoveredKey).String(),
	}, nil
}

func (s SIWEService) Submit(data *dto.SubmitRequestDTO) (*dto.VerifyResponse, error) {
	vr, err := s.Verify(data)
	if err != nil {
		return nil, err
	}
	return vr, nil
}

func (s SIWEService) CreateToken(audiences []string, accountId string, ecdsaAddress string, slyWalletAddress string, role string) (*verifier.Token, error) {
	return s.verifier.CreateToken(audiences, accountId, ecdsaAddress, slyWalletAddress, role)
}

func (s SIWEService) GetOrCreateAccount(context context.Context, address string) (repositories.ECDSAKey, error) {
	return s.userDB.GetOrCreateECDSAKey(context, address)
}

type AuthenticationResult struct {
	IsAuthenticated bool
	StatusCode      int
	ErrorCode       string
	ErrorMessage    string
	ErrorDetails    string
}

func (s SIWEService) Authenticate(eoa common.Address, slyWalletAddress common.Address) AuthenticationResult {
	slyWallet, err := s.slyWalletManager.GetSLYWalletContractAtAddress(slyWalletAddress)

	if err != nil {
		if strings.Contains(err.Error(), "no contract code at given address") {
			return AuthenticationResult{
				IsAuthenticated: false,
				StatusCode:      400,
				ErrorCode:       slyerrors.ErrCodeNoContractAtGivenAddress,
				ErrorMessage:    err.Error(),
				ErrorDetails:    fmt.Sprintf("no contract at %s", slyWalletAddress.Hex()),
			}
		}
		return AuthenticationResult{
			IsAuthenticated: false,
			StatusCode:      500,
			ErrorCode:       slyerrors.ErrCodeUnknown,
			ErrorMessage:    err.Error(),
			ErrorDetails:    "",
		}
	}

	role, err := slyWallet.GetKeyRole(&bind.CallOpts{}, eoa)
	if err != nil {
		return AuthenticationResult{
			IsAuthenticated: false,
			StatusCode:      500,
			ErrorCode:       slyerrors.ErrCodeGetSLYAuthentication,
			ErrorMessage:    err.Error(),
			ErrorDetails:    "",
		}
	}

	if role == 0 {
		return AuthenticationResult{
			IsAuthenticated: false,
			StatusCode:      403,
			ErrorCode:       slyerrors.ErrCodeNotAControllerKey,
			ErrorMessage:    "not a controller key",
			ErrorDetails:    fmt.Sprintf("%s is not a controller key of SLY %s", eoa.Hex(), slyWalletAddress.Hex()),
		}
	}

	return AuthenticationResult{
		IsAuthenticated: true,
		StatusCode:      200,
	}
}
