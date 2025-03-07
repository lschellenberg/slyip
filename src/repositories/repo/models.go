package repo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var DBItemNotFound = fmt.Errorf("not found")

// AccountModel represents an account with JSON annotations and relations
type AccountModel struct {
	ID                uuid.UUID `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Phone             string    `json:"phone"`
	Email             *string   `json:"email,omitempty"`
	IsEmailVerified   bool      `json:"isEmailVerified"`
	IsPhoneVerified   bool      `json:"isPhoneVerified"`
	PasswordHashed    string    `json:"-"` // Hide password in JSON responses
	InvitationCode    string    `json:"invitationCode"`
	Role              string    `json:"role"`
	LastUsedSlyWallet string    `json:"lastUsedSlyWallet"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Relations
	Ecdsas     []EcdsaModel     `json:"ecdsas,omitempty"`
	SlyWallets []SlyWalletModel `json:"slyWallets,omitempty"`
}

// EcdsaModel represents an ECDSA key with JSON annotations and relations
type EcdsaModel struct {
	Address   string    `json:"address"`
	AccountID uuid.UUID `json:"accountId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relations
	Account    *AccountModel         `json:"account,omitempty"`
	SlyWallets []EcdsaSlyWalletModel `json:"slyWallets,omitempty"`
}

// EcdsaSlyWalletModel represents a relationship between ECDSA keys and SlyWallets
type EcdsaSlyWalletModel struct {
	EcdsaAddress          string `json:"ecdsaAddress"`
	OnChainAccountAddress string `json:"onChainAccountAddress"`
	OnChainPermissions    int32  `json:"onChainPermissions"`

	// Relations
	Ecdsa     *EcdsaModel     `json:"ecdsa,omitempty"`
	SlyWallet *SlyWalletModel `json:"slyWallet,omitempty"`
}

// SlyWalletModel represents a Sly wallet with JSON annotations and relations
type SlyWalletModel struct {
	Address           string    `json:"address"`
	Chainid           string    `json:"chainId"`
	AccountID         uuid.UUID `json:"accountId"`
	TransactionHash   string    `json:"transactionHash"`
	TransactionStatus string    `json:"transactionStatus"`
	InvitationCode    string    `json:"invitationCode"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Relations
	Account          *AccountModel         `json:"account,omitempty"`
	EcdsaConnections []EcdsaSlyWalletModel `json:"ecdsaConnections,omitempty"`
}

// InvitationCodeModel represents an invitation code with JSON annotations
type InvitationCodeModel struct {
	Code            string    `json:"code"`
	TransactionHash string    `json:"transactionHash"`
	ExpiresAt       time.Time `json:"expiresAt"`

	// Relations - these aren't directly in the model but could be useful
	Accounts   []AccountModel   `json:"accounts,omitempty"`
	SlyWallets []SlyWalletModel `json:"slyWallets,omitempty"`
}

// PaginatedResponse is a generic paginated response for any model
type PaginatedResponse[T any] struct {
	Data      []T               `json:"data"`
	PageSize  int               `json:"pageSize"`
	Offset    int               `json:"offset"`
	SortBy    string            `json:"sortBy"`
	SortOrder string            `json:"order"`
	Total     uint64            `json:"total"`
	Count     uint64            `json:"count"`
	Filter    map[string]string `json:"filter"`
}
