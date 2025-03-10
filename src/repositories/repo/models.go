package repo

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"yip/src/api/auth/verifier"
)

var DBItemNotFound = fmt.Errorf("not found")

// AccountModel represents an account with JSON annotations and relations
type AccountModel struct {
	ID                uuid.UUID `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Phone             string    `json:"phone"`
	Email             string    `json:"email"`
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

func (a AccountModel) IsAdmin() bool {
	return a.Role == verifier.RoleAdmin
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

func (ic *InvitationCodeModel) IsValid() bool {
	return len(ic.TransactionHash) == 0
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

// String returns a formatted string representation of AccountModel
func (a AccountModel) String() string {
	var ecdsasStr, slyWalletsStr string

	// Format Ecdsas slice
	ecdsasStr = "  Ecdsas: [\n"
	if len(a.Ecdsas) == 0 {
		ecdsasStr += "    (empty)\n"
	} else {
		for i, ecdsa := range a.Ecdsas {
			ecdsasStr += fmt.Sprintf("    %d: %s\n", i, ecdsa.Address)
		}
	}
	ecdsasStr += "  ]\n"

	// Format SlyWallets slice
	slyWalletsStr = "  SlyWallets: [\n"
	if len(a.SlyWallets) == 0 {
		slyWalletsStr += "    (empty)\n"
	} else {
		for i, wallet := range a.SlyWallets {
			slyWalletsStr += fmt.Sprintf("    %d: %s (%s)\n", i, wallet.Address, wallet.Chainid)
		}
	}
	slyWalletsStr += "  ]\n"

	return fmt.Sprintf("Account{\n"+
		"  ID: %s\n"+
		"  Name: %s %s\n"+
		"  Email: %s\n"+
		"  Phone: %s\n"+
		"  Role: %s\n"+
		"  EmailVerified: %t\n"+
		"  PhoneVerified: %t\n"+
		"  InvitationCode: %s\n"+
		"%s"+
		"%s"+
		"  CreatedAt: %s\n"+
		"  UpdatedAt: %s\n"+
		"}",
		a.ID, a.FirstName, a.LastName, a.Email, a.Phone, a.Role,
		a.IsEmailVerified, a.IsPhoneVerified, a.InvitationCode,
		slyWalletsStr, ecdsasStr,
		a.CreatedAt.Format(time.RFC3339), a.UpdatedAt.Format(time.RFC3339))
}

// String returns a formatted string representation of EcdsaModel
func (e EcdsaModel) String() string {
	// Format SlyWallets slice
	slyWalletsStr := "  SlyWallets: [\n"
	if len(e.SlyWallets) == 0 {
		slyWalletsStr += "    (empty)\n"
	} else {
		for i, conn := range e.SlyWallets {
			slyWalletsStr += fmt.Sprintf("    %d: %s (permissions: %d)\n",
				i, conn.OnChainAccountAddress, conn.OnChainPermissions)
		}
	}
	slyWalletsStr += "  ]\n"

	return fmt.Sprintf("Ecdsa{\n"+
		"  Address: %s\n"+
		"  AccountID: %s\n"+
		"%s"+
		"  CreatedAt: %s\n"+
		"  UpdatedAt: %s\n"+
		"}",
		e.Address, e.AccountID, slyWalletsStr,
		e.CreatedAt.Format(time.RFC3339), e.UpdatedAt.Format(time.RFC3339))
}

// String returns a formatted string representation of EcdsaSlyWalletModel
func (e EcdsaSlyWalletModel) String() string {
	return fmt.Sprintf("EcdsaSlyWallet{\n"+
		"  EcdsaAddress: %s\n"+
		"  OnChainAccountAddress: %s\n"+
		"  OnChainPermissions: %d\n"+
		"}",
		e.EcdsaAddress, e.OnChainAccountAddress, e.OnChainPermissions)
}

// String returns a formatted string representation of SlyWalletModel
func (s SlyWalletModel) String() string {
	// Format EcdsaConnections slice
	connectionsStr := "  EcdsaConnections: [\n"
	if len(s.EcdsaConnections) == 0 {
		connectionsStr += "    (empty)\n"
	} else {
		for i, conn := range s.EcdsaConnections {
			connectionsStr += fmt.Sprintf("    %d: %s (permissions: %d)\n",
				i, conn.EcdsaAddress, conn.OnChainPermissions)
		}
	}
	connectionsStr += "  ]\n"

	return fmt.Sprintf("SlyWallet{\n"+
		"  Address: %s\n"+
		"  ChainID: %s\n"+
		"  AccountID: %s\n"+
		"  TransactionStatus: %s\n"+
		"  TransactionHash: %s\n"+
		"  InvitationCode: %s\n"+
		"%s"+
		"  CreatedAt: %s\n"+
		"  UpdatedAt: %s\n"+
		"}",
		s.Address, s.Chainid, s.AccountID, s.TransactionStatus,
		s.TransactionHash, s.InvitationCode, connectionsStr,
		s.CreatedAt.Format(time.RFC3339), s.UpdatedAt.Format(time.RFC3339))
}

// String returns a formatted string representation of InvitationCodeModel
func (i InvitationCodeModel) String() string {
	// Format Accounts slice
	accountsStr := "  Accounts: [\n"
	if len(i.Accounts) == 0 {
		accountsStr += "    (empty)\n"
	} else {
		for j, account := range i.Accounts {
			accountsStr += fmt.Sprintf("    %d: %s (%s %s)\n",
				j, account.ID, account.FirstName, account.LastName)
		}
	}
	accountsStr += "  ]\n"

	// Format SlyWallets slice
	walletsStr := "  SlyWallets: [\n"
	if len(i.SlyWallets) == 0 {
		walletsStr += "    (empty)\n"
	} else {
		for j, wallet := range i.SlyWallets {
			walletsStr += fmt.Sprintf("    %d: %s (%s)\n",
				j, wallet.Address, wallet.Chainid)
		}
	}
	walletsStr += "  ]\n"

	return fmt.Sprintf("InvitationCode{\n"+
		"  Code: %s\n"+
		"  TransactionHash: %s\n"+
		"  ExpiresAt: %s\n"+
		"%s"+
		"%s"+
		"}",
		i.Code, i.TransactionHash, i.ExpiresAt.Format(time.RFC3339),
		accountsStr, walletsStr)
}

// String returns a formatted string representation of PaginatedResponse
func (p PaginatedResponse[T]) String() string {
	// Format Data slice
	dataStr := "  Data: [\n"
	if len(p.Data) == 0 {
		dataStr += "    (empty)\n"
	} else {
		for i, item := range p.Data {
			dataStr += fmt.Sprintf("    %d: %v\n", i, item)
		}
	}
	dataStr += "  ]\n"

	// Format Filter map
	filterStr := "  Filter: {\n"
	if len(p.Filter) == 0 {
		filterStr += "    (empty)\n"
	} else {
		for key, value := range p.Filter {
			filterStr += fmt.Sprintf("    %s: %s\n", key, value)
		}
	}
	filterStr += "  }\n"

	return fmt.Sprintf("PaginatedResponse{\n"+
		"%s"+
		"  PageSize: %d\n"+
		"  Offset: %d\n"+
		"  SortBy: %s\n"+
		"  SortOrder: %s\n"+
		"  Total: %d\n"+
		"  Count: %d\n"+
		"%s"+
		"}",
		dataStr, p.PageSize, p.Offset, p.SortBy, p.SortOrder,
		p.Total, p.Count, filterStr)
}
