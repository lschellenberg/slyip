package contracts

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WalletManager handles creating and managing SLY smart wallets
type WalletManager struct {
	client        *ethclient.Client
	factoryAddr   common.Address
	factory       *SLYWalletFactory
	defaultSigner *bind.TransactOpts
}

// NewWalletManager creates a new SLYWallet manager
func NewWalletManager(client *ethclient.Client, factoryAddress common.Address, defaultSigner *bind.TransactOpts) (*WalletManager, error) {
	factory, err := NewSLYWalletFactory(factoryAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to SLYWalletFactory at %s: %w", factoryAddress.Hex(), err)
	}

	return &WalletManager{
		client:        client,
		factoryAddr:   factoryAddress,
		factory:       factory,
		defaultSigner: defaultSigner,
	}, nil
}

// CreateWallet creates a new SLYWallet with the provided owner
func (m *WalletManager) CreateWallet(ctx context.Context, owner common.Address) (common.Address, *types.Transaction, error) {
	if owner == (common.Address{}) {
		return common.Address{}, nil, errors.New("owner address cannot be zero")
	}

	tx, err := m.factory.CreateSLYWallet(m.defaultSigner, owner)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Wait for transaction to be mined and get receipt
	receipt, err := bind.WaitMined(ctx, m.client, tx)
	if err != nil {
		return common.Address{}, tx, fmt.Errorf("failed to wait for wallet creation: %w", err)
	}

	// Parse events to get the created wallet address
	var walletAddress common.Address
	eventSig := []byte("SLYWalletCreated(address,address,address)")
	eventSigHash := crypto.Keccak256Hash(eventSig)

	for _, log := range receipt.Logs {
		// Check if this log belongs to the factory contract
		if log.Address == m.factoryAddr {
			// The SLYWalletCreated event has the signature:
			// event SLYWalletCreated(address indexed diamond, address indexed owner, address indexed creator)
			// We need to check the topics
			if len(log.Topics) == 4 && log.Topics[0] == eventSigHash {
				// The first topic is the event signature, the second is the wallet address
				walletAddress = common.HexToAddress(log.Topics[1].Hex())
				return walletAddress, tx, nil
			}
		}
	}

	return common.Address{}, tx, errors.New("wallet created but address not found in logs")
}

// SpawnWallet creates a new SLYWallet with the provided owner
func (m *WalletManager) SpawnWallet(owner common.Address) (*types.Transaction, error) {
	if owner == (common.Address{}) {
		return nil, errors.New("owner address cannot be zero")
	}

	tx, err := m.factory.CreateSLYWallet(m.defaultSigner, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return tx, nil
}

func (m *WalletManager) GetSLYWalletAddressFromReceipt(ctx context.Context, receipt *types.Receipt) (common.Address, error) {

	// Parse events to get the created wallet address
	var walletAddress common.Address
	eventSig := []byte("SLYWalletCreated(address,address,address)")
	eventSigHash := crypto.Keccak256Hash(eventSig)

	for _, log := range receipt.Logs {
		// Check if this log belongs to the factory contract
		if log.Address == m.factoryAddr {
			// The SLYWalletCreated event has the signature:
			// event SLYWalletCreated(address indexed diamond, address indexed owner, address indexed creator)
			// We need to check the topics
			if len(log.Topics) == 4 && log.Topics[0] == eventSigHash {
				// The first topic is the event signature, the second is the wallet address
				walletAddress = common.HexToAddress(log.Topics[1].Hex())
				return walletAddress, nil
			}
		}
	}

	return common.Address{}, errors.New("wallet created but address not found in logs")
}

// CreateWalletWithSalt creates a new SLYWallet with a specific salt for deterministic address
func (m *WalletManager) CreateWalletWithSalt(ctx context.Context, owner common.Address, salt [32]byte) (common.Address, *types.Transaction, error) {
	if owner == (common.Address{}) {
		return common.Address{}, nil, errors.New("owner address cannot be zero")
	}

	// First check if we can predict the address
	predictedAddr, err := m.factory.PredictWalletAddress(&bind.CallOpts{}, salt)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to predict wallet address: %w", err)
	}

	// Create the wallet with salt
	tx, err := m.factory.CreateSLYWalletWithSalt(m.defaultSigner, owner, salt)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to create wallet with salt: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, m.client, tx)
	if err != nil {
		return predictedAddr, tx, fmt.Errorf("failed to wait for wallet creation with salt: %w", err)
	}

	// Parse events to verify created wallet address
	var walletAddress common.Address
	eventSig := []byte("SLYWalletCreated(address,address,address)")
	eventSigHash := crypto.Keccak256Hash(eventSig)

	for _, log := range receipt.Logs {
		// Check if this log belongs to the factory contract
		if log.Address == m.factoryAddr {
			// The SLYWalletCreated event has the signature:
			// event SLYWalletCreated(address indexed diamond, address indexed owner, address indexed creator)
			// We need to check the topics
			if len(log.Topics) == 4 && log.Topics[0] == eventSigHash {
				// The first topic is the event signature, the second is the wallet address
				walletAddress = common.HexToAddress(log.Topics[1].Hex())

				// Verify it matches our prediction
				if walletAddress != predictedAddr {
					return walletAddress, tx, fmt.Errorf("created wallet address %s does not match predicted address %s",
						walletAddress.Hex(), predictedAddr.Hex())
				}

				return walletAddress, tx, nil
			}
		}
	}

	// If we don't find the event, return the predicted address
	return predictedAddr, tx, nil
}

// GenerateRandomSalt generates a random 32-byte salt for CREATE2 deployment
func (m *WalletManager) GenerateRandomSalt() ([32]byte, error) {
	var salt [32]byte
	_, err := rand.Read(salt[:])
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to generate random salt: %w", err)
	}
	return salt, nil
}

// GetWallet returns an instance of SLYWallet bound to the given address
func (m *WalletManager) GetWallet(walletAddress common.Address) (*SLYWallet, error) {
	if walletAddress == (common.Address{}) {
		return nil, errors.New("wallet address cannot be zero")
	}

	wallet, err := NewSLYWallet(walletAddress, m.client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to SLY at %s: %w", walletAddress.Hex(), err)
	}

	return wallet, nil
}

// GetOwnedWallets returns all wallets owned by the specified address
func (m *WalletManager) GetOwnedWallets(ctx context.Context, owner common.Address) ([]common.Address, error) {
	if owner == (common.Address{}) {
		return nil, errors.New("owner address cannot be zero")
	}

	wallets, err := m.factory.GetWalletsByOwner(&bind.CallOpts{Context: ctx}, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets for owner %s: %w", owner.Hex(), err)
	}

	return wallets, nil
}

// AddKey adds a new key with the specified role to the wallet
func (m *WalletManager) AddKey(ctx context.Context, walletAddress common.Address, key common.Address, role Role) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.AddKey(m.defaultSigner, key, uint8(role))
	if err != nil {
		return nil, fmt.Errorf("failed to add key %s with role %d: %w", key.Hex(), role, err)
	}

	return tx, nil
}

// RemoveKey removes a key from the wallet
func (m *WalletManager) RemoveKey(ctx context.Context, walletAddress common.Address, key common.Address) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.RemoveKey(m.defaultSigner, key)
	if err != nil {
		return nil, fmt.Errorf("failed to remove key %s: %w", key.Hex(), err)
	}

	return tx, nil
}

// ExecuteTransaction executes a transaction from the wallet
func (m *WalletManager) ExecuteTransaction(ctx context.Context, walletAddress common.Address, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	// Set value in transaction options
	txOpts := *m.defaultSigner
	txOpts.Value = value

	tx, err := wallet.Execute(&txOpts, to, value, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	return tx, nil
}

// ExecuteBatch executes a batch of transactions from the wallet
func (m *WalletManager) ExecuteBatch(ctx context.Context, walletAddress common.Address, to []common.Address, values []*big.Int, data [][]byte) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	// Calculate total value for the transaction
	totalValue := big.NewInt(0)
	for _, val := range values {
		totalValue.Add(totalValue, val)
	}

	// Set value in transaction options
	txOpts := *m.defaultSigner
	txOpts.Value = totalValue

	tx, err := wallet.ExecuteBatch(&txOpts, to, values, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch transaction: %w", err)
	}

	return tx, nil
}

// GetNonce gets the current nonce for a signer in the wallet
func (m *WalletManager) GetNonce(ctx context.Context, walletAddress common.Address, signer common.Address) (*big.Int, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	nonce, err := wallet.GetNonce(&bind.CallOpts{Context: ctx}, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce for signer %s: %w", signer.Hex(), err)
	}

	return nonce, nil
}

// ExecuteWithSignature executes a transaction via meta-transaction (with signature)
func (m *WalletManager) ExecuteWithSignature(
	ctx context.Context,
	walletAddress common.Address,
	to common.Address,
	value *big.Int,
	data []byte,
	signer common.Address,
	signerPrivateKey *ecdsa.PrivateKey,
) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	// Get current nonce for signer
	nonce, err := wallet.GetNonce(&bind.CallOpts{Context: ctx}, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Create and sign the meta-transaction
	signature, err := m.createMetaTxSignature(ctx, walletAddress, to, value, data, signer, nonce, signerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	// Execute with signature
	tx, err := wallet.ExecuteWithSignature(m.defaultSigner, to, value, data, signer, nonce, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to execute with signature: %w", err)
	}

	return tx, nil
}

func (m *WalletManager) createMetaTxSignature(
	ctx context.Context,
	walletAddress common.Address,
	to common.Address,
	value *big.Int,
	data []byte,
	signer common.Address,
	nonce *big.Int,
	privateKey *ecdsa.PrivateKey,
) ([]byte, error) {
	// Calculate the domain separator
	// Note: In a real implementation, you should fetch the domain separator from the contract
	// or calculate it using the correct domain parameters
	chainID, err := m.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Calculate the EIP-712 digest for the execute transaction
	digest, err := m.calculateExecuteDigest(walletAddress, to, value, data, signer, nonce, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate message digest: %w", err)
	}

	// Sign the digest
	signature, err := crypto.Sign(digest[:], privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign digest: %w", err)
	}

	// Convert to the format expected by the contract
	// Ethereum signature is [R || S || V] where V is 0 or 1 and we need to add 27 to it
	signature[64] += 27

	return signature, nil
}

// CalculateExecuteDigest calculates the EIP-712 digest for the execute function
func (m *WalletManager) calculateExecuteDigest(
	walletAddress common.Address,
	to common.Address,
	value *big.Int,
	data []byte,
	signer common.Address,
	nonce *big.Int,
	chainID *big.Int,
) ([32]byte, error) {
	// EIP-712 domain separator parameters
	// Note: In production, you should get these values from the contract or config
	domainSeparator := crypto.Keccak256Hash(
		[]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
		crypto.Keccak256Hash([]byte("SLY")).Bytes(),
		crypto.Keccak256Hash([]byte("1")).Bytes(),
		common.LeftPadBytes(chainID.Bytes(), 32),
		walletAddress.Bytes(),
	)

	// Hash the message data
	dataHash := crypto.Keccak256Hash(data)

	// Type hash for Execute(address to,uint256 value,bytes data,address signer,uint256 nonce)
	typeHash := crypto.Keccak256Hash(
		[]byte("Execute(address to,uint256 value,bytes data,address signer,uint256 nonce)"),
	)

	// Create the struct hash
	structHash := crypto.Keccak256Hash(
		typeHash.Bytes(),
		common.LeftPadBytes(to.Bytes(), 32),
		common.LeftPadBytes(value.Bytes(), 32),
		dataHash.Bytes(),
		common.LeftPadBytes(signer.Bytes(), 32),
		common.LeftPadBytes(nonce.Bytes(), 32),
	)

	// Calculate the final digest
	digest := crypto.Keccak256Hash(
		[]byte("\x19\x01"),
		domainSeparator.Bytes(),
		structHash.Bytes(),
	)

	var result [32]byte
	copy(result[:], digest.Bytes())
	return result, nil
}

// WithdrawETH withdraws ETH from the wallet
func (m *WalletManager) WithdrawETH(ctx context.Context, walletAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.WithdrawETH(m.defaultSigner, to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw ETH: %w", err)
	}

	return tx, nil
}

// WithdrawERC20 withdraws ERC20 tokens from the wallet
func (m *WalletManager) WithdrawERC20(ctx context.Context, walletAddress common.Address, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	tx, err := wallet.WithdrawERC20(m.defaultSigner, token, to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw ERC20 tokens: %w", err)
	}

	return tx, nil
}

// CheckPermission checks if a key has a specific permission
func (m *WalletManager) CheckPermission(ctx context.Context, walletAddress common.Address, key common.Address, permission Permission) (bool, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return false, err
	}

	hasPermission, err := wallet.HasPermission(&bind.CallOpts{Context: ctx}, key, uint8(permission))
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return hasPermission, nil
}

// GetKeyRole gets the role of a key in the wallet
func (m *WalletManager) GetKeyRole(ctx context.Context, walletAddress common.Address, key common.Address) (Role, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return RoleNone, err
	}

	roleUint8, err := wallet.GetKeyRole(&bind.CallOpts{Context: ctx}, key)
	if err != nil {
		return RoleNone, fmt.Errorf("failed to get key role: %w", err)
	}

	return Role(roleUint8), nil
}

// GetKeysByRole gets all keys with a specific role
func (m *WalletManager) GetKeysByRole(ctx context.Context, walletAddress common.Address, role Role) ([]common.Address, error) {
	wallet, err := m.GetWallet(walletAddress)
	if err != nil {
		return nil, err
	}

	keys, err := wallet.GetKeysByRole(&bind.CallOpts{Context: ctx}, uint8(role))
	if err != nil {
		return nil, fmt.Errorf("failed to get keys by role: %w", err)
	}

	return keys, nil
}

// WalletKeys holds the different types of keys in the wallet
type WalletKeys struct {
	Owners         []common.Address
	Admins         []common.Address
	Authenticators []common.Address
}

// GetWalletKeys retrieves all keys categorized by role from a SLYWallet contract
func (m *WalletManager) GetWalletKeys(walletAddress common.Address) (*WalletKeys, error) {
	// Create a new instance of the SLYWallet contract binding
	wallet, err := NewSLYWallet(walletAddress, m.client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wallet contract: %v", err)
	}

	// Create default call options
	opts := &bind.CallOpts{Context: context.Background()}

	// Retrieve keys for each role
	owners, err := wallet.GetKeysByRole(opts, uint8(RoleOwner))
	if err != nil {
		return nil, fmt.Errorf("failed to get owner keys: %v", err)
	}

	admins, err := wallet.GetKeysByRole(opts, uint8(RoleAdmin))
	if err != nil {
		return nil, fmt.Errorf("failed to get admin keys: %v", err)
	}

	authenticators, err := wallet.GetKeysByRole(opts, uint8(RoleAuthenticator))
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticator keys: %v", err)
	}

	return &WalletKeys{
		Owners:         owners,
		Admins:         admins,
		Authenticators: authenticators,
	}, nil
}
