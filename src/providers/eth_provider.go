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
	"math"
	"math/big"
	"strconv"
	"yip/src/config"
	"yip/src/cryptox"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

const (
	TransactionStatusPending = "pending"
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)

type EthProvider struct {
	Client             *ethclient.Client
	hubCreationBlock   int64
	hubAddress         common.Address
	config             *config.EthConfig
	signingWallet      *cryptox.Wallet
	PubKeySignerWallet common.Address
	ChainId            *big.Int
}

func InitEthProvider(c *config.EthConfig) (p EthProvider, err error) {
	p.config = c

	id, err := strconv.ParseInt(c.Chain.ID, 10, 64)
	if err != nil {
		return
	}

	p.ChainId = big.NewInt(id)
	p.hubCreationBlock = c.Contracts.HubCreationBlock

	if p.Client, err = ethclient.Dial(c.Chain.RPCUrl); err != nil {
		return
	}

	p.signingWallet = &cryptox.Wallet{}
	if err = p.signingWallet.FromPrivateKey(c.Wallet.Private); err != nil {
		return
	}
	p.PubKeySignerWallet = p.signingWallet.Public

	return p, nil
}

func (p EthProvider) BalanceOf(ctc context.Context, address common.Address) (*big.Float, error) {
	value, err := p.Client.BalanceAt(ctc, address, nil)
	if err != nil {
		return nil, err
	}

	fbalance := new(big.Float)
	fbalance.SetString(value.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue, nil
}

func (p EthProvider) PreparedSigner(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := p.Client.PendingNonceAt(ctx, p.signingWallet.Public)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeCantDetermineNonce, err.Error())
	}
	gasPrice, err := p.Client.SuggestGasPrice(ctx)
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
	return auth, nil
}

func (p EthProvider) DefaultSigner(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := p.Client.PendingNonceAt(ctx, p.signingWallet.Public)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeCantDetermineNonce, err.Error())
	}
	auth, err := bind.NewKeyedTransactorWithChainID(p.signingWallet.Private, p.ChainId)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeCantCreateTransactor, err.Error())
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	return auth, nil
}

func (p EthProvider) HubCreationBlock() int64 {
	return p.hubCreationBlock
}

func (p EthProvider) GetRPCClient() *ethclient.Client {
	return p.Client
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

func (p EthProvider) GetTransactionReceipt(transHash common.Hash) (*types.Receipt, error) {
	r, err := p.Client.TransactionReceipt(context.Background(), transHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, slyerrors.NotFound(slyerrors.ErrCodeAccountTransactionNotFound, fmt.Sprintf("transaction hash: %s not found", transHash.Hex()))
		}
		return nil, httpx.InternalError(err.Error())
	}
	return r, nil
}
