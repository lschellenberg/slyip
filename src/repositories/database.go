package repositories

import (
	"context"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type Database interface {
	RegisterAccount(context context.Context, email string, passwordHashed string, role string) (UserAccount, error)
	RegisterAccountWithEmail(context context.Context, email string) (UserAccount, error)
	GetAccounts(context context.Context, limit int, offset int) (ListUsersAccountResponse, error)
	GetAccountByEmail(context context.Context, email string) (UserAccount, error)
	GetAccountById(context context.Context, id uuid.UUID) (UserAccount, error)
	GetLastAccount(context context.Context) (UserAccount, error)
	SetRole(context context.Context, id uuid.UUID, role string) (UserAccount, error)
	GetOrCreateECDSAKey(context context.Context, address string) (ECDSAKey, error)
	SetEmail(ctx context.Context, id uuid.UUID, email string) (UserAccount, error)
	AddDevice(ctx context.Context, userId uuid.UUID, pubKey string) (ECDSAKey, error)
	GetDevices(ctx context.Context, userId uuid.UUID) ([]ECDSAKey, error)
	GetSLYWallets(context context.Context, accountId uuid.UUID) ([]SLYWalletWithControllerKeys, error)
	SetLastUsedSLYWallet(ctx context.Context, id uuid.UUID, address common2.Address) (UserAccount, error)
}
