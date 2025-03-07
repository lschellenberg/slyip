package providers

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"strconv"
	"strings"
	"yip/src/config"
	"yip/src/contracts"
	"yip/src/cryptox"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

const (
	TransactionTypeSpawnSLYWallet = "SpawnSLYWallet"

	TransactionStatusPending = "pending"
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)

type EthProvider struct {
	client             *ethclient.Client
	hubCreationBlock   int64
	hubAddress         common.Address
	config             *config.EthConfig
	signingWallet      *cryptox.Wallet
	PubKeySignerWallet common.Address
	ChainId            *big.Int
	WalletManager      *contracts.WalletManager
}

func InitEthProvider(c *config.EthConfig) (p EthProvider, err error) {
	p.config = c

	id, err := strconv.ParseInt(c.Chain.ID, 10, 64)
	if err != nil {
		return
	}

	p.ChainId = big.NewInt(id)
	p.hubCreationBlock = c.Contracts.HubCreationBlock

	if p.client, err = ethclient.Dial(c.Chain.RPCUrl); err != nil {
		return
	}

	p.signingWallet = &cryptox.Wallet{}
	if err = p.signingWallet.FromPrivateKey(c.Wallet.Private); err != nil {
		return
	}
	p.PubKeySignerWallet = p.signingWallet.Public

	if p.WalletManager, err = contracts.NewWalletManager(p.client, common.HexToAddress(c.Sly.FactoryAddress), nil); err != nil {
		return
	}

	return p, nil
}

func (p EthProvider) BalanceOf(ctc context.Context, address common.Address) (*big.Float, error) {
	value, err := p.client.BalanceAt(ctc, address, nil)
	if err != nil {
		return nil, err
	}

	fbalance := new(big.Float)
	fbalance.SetString(value.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue, nil
}

func (p EthProvider) CreateSLYWallet(ctx context.Context, owner common.Address) (*TransactionTicket, error) {
	nonce, err := p.client.PendingNonceAt(ctx, p.signingWallet.Public)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeCantDetermineNonce, err.Error())
	}
	gasPrice, err := p.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeCantEstimateGasPrice, err.Error())
	}
	auth, err := bind.NewKeyedTransactorWithChainID(p.signingWallet.Private, p.ChainId)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeCantCreateTransactor, err.Error())
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	log.Info("Creating SLY with: ", owner)
	trans, err := p.WalletManager.SpawnWallet(owner)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeTransact, err.Error())
	}

	return &TransactionTicket{
		TransactionType: TransactionTypeSpawnSLYWallet,
		TransactionHash: trans.Hash().Hex(),
	}, nil
}

func (p EthProvider) TransactionStatus(transHash common.Hash) (*TransactionStatus, error) {
	receipt, err := p.client.TransactionReceipt(context.Background(), transHash)

	if err != nil {
		if err.Error() == "not found" {
			return &TransactionStatus{
				TransactionHash: transHash.Hex(),
				Status:          TransactionStatusPending,
				ContractAddress: "",
			}, nil
		}
		fmt.Println(err.Error())
		return nil, err
	}

	status := TransactionStatusSuccess
	contractAddress := ""
	if receipt.Status == 0 {
		status = TransactionStatusFailed
	} else {
		address, err := p.WalletManager.GetSLYWalletAddressFromReceipt(context.Background(), receipt)
		if err != nil {
			return nil, err
		}
		contractAddress = address.Hex()
	}

	return &TransactionStatus{
		TransactionHash: transHash.Hex(),
		Status:          status,
		ContractAddress: contractAddress,
	}, nil
}

func (p EthProvider) GetSLYWalletContractAtAddress(address common.Address) (*contracts.SLYWallet, error) {
	return p.WalletManager.GetWallet(address)
}

func (p EthProvider) HubCreationBlock() int64 {
	return p.hubCreationBlock
}

func (p EthProvider) GetRPCClient() *ethclient.Client {
	return p.client
}

type TransactionTicket struct {
	TransactionType string `json:"type"`
	TransactionHash string `json:"transactionHash"`
}

type TransactionStatus struct {
	TransactionHash string `json:"transactionHash"`
	Status          string `json:"status"`
	ContractAddress string `json:"contractAddress"`
}

type TransactionState struct {
	TransactionHash string    `json:"transactionHash"`
	Status          string    `json:"status"`
	ContractAddress string    `json:"contractAddress"`
	DatabaseId      uuid.UUID `json:"databaseId"`
}

func (s TransactionState) IsPending() bool {
	return s.Status == TransactionStatusPending
}

func (s TransactionState) Failed() bool {
	return s.Status == TransactionStatusFailed
}

func (s TransactionState) Succeeded() bool {
	return s.Status == TransactionStatusSuccess
}

func (p EthProvider) AuthenticateControllerKeyOfSLYWallet(controllerKey common.Address, slyWalletAddress common.Address) AuthenticationResult {
	slyWallet, err := p.WalletManager.GetWallet(slyWalletAddress)

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

	isAuthenticated, err := slyWallet.KeyExists(&bind.CallOpts{}, controllerKey)
	if err != nil {
		return AuthenticationResult{
			IsAuthenticated: false,
			StatusCode:      500,
			ErrorCode:       slyerrors.ErrCodeGetSLYAuthentication,
			ErrorMessage:    err.Error(),
			ErrorDetails:    "",
		}
	}

	if !isAuthenticated {
		return AuthenticationResult{
			IsAuthenticated: false,
			StatusCode:      403,
			ErrorCode:       slyerrors.ErrCodeNotAControllerKey,
			ErrorMessage:    "not a controller key",
			ErrorDetails:    fmt.Sprintf("%s is not a controller key of SLYWallet %s", controllerKey.Hex(), slyWalletAddress.Hex()),
		}
	}

	return AuthenticationResult{
		IsAuthenticated: true,
		StatusCode:      200,
	}
}

type AuthenticationResult struct {
	IsAuthenticated bool
	StatusCode      int
	ErrorCode       string
	ErrorMessage    string
	ErrorDetails    string
}

func (p EthProvider) GetTransactionReceipt(transHash common.Hash) (*types.Receipt, error) {
	r, err := p.client.TransactionReceipt(context.Background(), transHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, slyerrors.NotFound(slyerrors.ErrCodeAccountTransactionNotFound, fmt.Sprintf("transaction hash: %s not found", transHash.Hex()))
		}
		return nil, httpx.InternalError(err.Error())
	}
	return r, nil
}

func (p EthProvider) GetTransactionStatusByReceipt(transactionHash common.Hash, receipt *types.Receipt, err error) (*TransactionState, error) {
	if err != nil {
		if err.Error() == "not found" {
			return &TransactionState{
				TransactionHash: transactionHash.Hex(),
				Status:          TransactionStatusPending,
				ContractAddress: "",
			}, nil
		}
		return nil, err
	}

	status := TransactionStatusSuccess
	contractAddress := ""
	if receipt.Status == 0 {
		status = TransactionStatusFailed
	} else {
		address, err := p.WalletManager.GetSLYWalletAddressFromReceipt(context.Background(), receipt)
		if err != nil {
			return nil, err
		}
		contractAddress = address.Hex()
	}

	return &TransactionState{
		TransactionHash: receipt.TxHash.Hex(),
		Status:          status,
		ContractAddress: contractAddress,
	}, nil
}
