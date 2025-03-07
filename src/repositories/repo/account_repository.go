package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"yip/src/common"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"yip/.gen/slyip/slyip/model"
	"yip/.gen/slyip/slyip/table"
)

// AccountRepository handles all account related database operations
type AccountRepository struct {
	db *Database
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *Database) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// Create creates a new account
func (r *AccountRepository) Create(ctx context.Context, account *AccountModel) (*AccountModel, error) {
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	stmt := table.Account.INSERT(
		table.Account.ID,
		table.Account.FirstName,
		table.Account.LastName,
		table.Account.Phone,
		table.Account.Email,
		table.Account.IsEmailVerified,
		table.Account.IsPhoneVerified,
		table.Account.PasswordHashed,
		table.Account.InvitationCode,
		table.Account.Role,
		table.Account.LastUsedSlyWallet,
		table.Account.CreatedAt,
		table.Account.UpdatedAt,
	).VALUES(
		account.ID,
		account.FirstName,
		account.LastName,
		account.Phone,
		account.Email,
		account.IsEmailVerified,
		account.IsPhoneVerified,
		account.PasswordHashed,
		account.InvitationCode,
		account.Role,
		account.LastUsedSlyWallet,
		account.CreatedAt,
		account.UpdatedAt,
	).RETURNING(
		table.Account.AllColumns,
	)

	var dbAccount model.Account
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return mapAccountToModel(dbAccount), nil
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*AccountModel, error) {
	stmt := postgres.SELECT(
		table.Account.AllColumns,
	).FROM(
		table.Account,
	).WHERE(
		table.Account.ID.EQ(postgres.UUID(id)),
	)

	var dbAccount model.Account
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbAccount)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}

	return mapAccountToModel(dbAccount), nil
}

// GetByEmail retrieves an account by email
func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*AccountModel, error) {
	stmt := postgres.SELECT(
		table.Account.AllColumns,
	).FROM(
		table.Account,
	).WHERE(
		table.Account.Email.EQ(postgres.String(email)),
	)

	var dbAccount model.Account
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbAccount)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get account by email: %w", err)
	}

	return mapAccountToModel(dbAccount), nil
}

// GetByPhone retrieves an account by phone
func (r *AccountRepository) GetByPhone(ctx context.Context, phone string) (*AccountModel, error) {
	stmt := postgres.SELECT(
		table.Account.AllColumns,
	).FROM(
		table.Account,
	).WHERE(
		table.Account.Phone.EQ(postgres.String(phone)),
	)

	var dbAccount model.Account
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbAccount)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get account by phone: %w", err)
	}

	return mapAccountToModel(dbAccount), nil
}

// GetByInvitationCode retrieves accounts by invitation code
func (r *AccountRepository) GetByInvitationCode(ctx context.Context, code string) ([]*AccountModel, error) {
	stmt := postgres.SELECT(
		table.Account.AllColumns,
	).FROM(
		table.Account,
	).WHERE(
		table.Account.InvitationCode.EQ(postgres.String(code)),
	)

	var dbAccounts []model.Account
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts by invitation code: %w", err)
	}

	accounts := make([]*AccountModel, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		accounts[i] = mapAccountToModel(dbAccount)
	}

	return accounts, nil
}

// List retrieves a paginated list of accounts
func (r *AccountRepository) List(ctx context.Context, query *common.PaginationQuery) (*PaginatedResponse[AccountModel], error) {
	if query == nil {
		query = NewPaginationQuery(AllowedAccountFilters)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Count total
	countStmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("total"),
	).FROM(
		table.Account,
	)

	var totalCount struct {
		Total uint64 `sql:"total"`
	}
	err := countStmt.QueryContext(ctx, r.db.GetDB(), &totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count accounts: %w", err)
	}

	// Get data with pagination
	selectStmt := postgres.SELECT(
		table.Account.AllColumns,
	).FROM(
		table.Account,
	)

	// Apply sorting
	switch query.SortBy {
	case "firstName":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Account.FirstName.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Account.FirstName.DESC())
		}
	case "lastName":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Account.LastName.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Account.LastName.DESC())
		}
	case "email":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Account.Email.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Account.Email.DESC())
		}
	case "role":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Account.Role.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Account.Role.DESC())
		}
	case "updatedAt":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Account.UpdatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Account.UpdatedAt.DESC())
		}
	default: // createdAt is default
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Account.CreatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Account.CreatedAt.DESC())
		}
	}

	// Apply pagination
	selectStmt = selectStmt.
		OFFSET(int64(query.Offset)).
		LIMIT(int64(query.PageSize))

	var dbAccounts []model.Account
	err = selectStmt.QueryContext(ctx, r.db.GetDB(), &dbAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	accounts := make([]AccountModel, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		accounts[i] = *mapAccountToModel(dbAccount)
	}

	return &PaginatedResponse[AccountModel]{
		Data:      accounts,
		PageSize:  query.PageSize,
		Offset:    query.Offset,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Total:     totalCount.Total,
		Count:     uint64(len(accounts)),
		Filter:    make(map[string]string),
	}, nil
}

// Update updates an account
func (r *AccountRepository) Update(ctx context.Context, account *AccountModel) (*AccountModel, error) {
	account.UpdatedAt = time.Now()

	var stmt postgres.UpdateStatement
	if account.Email != nil {
		stmt = table.Account.UPDATE().
			SET(
				table.Account.FirstName.SET(postgres.String(account.FirstName)),
				table.Account.LastName.SET(postgres.String(account.LastName)),
				table.Account.Phone.SET(postgres.String(account.Phone)),
				table.Account.Email.SET(postgres.String(*account.Email)),
				table.Account.IsEmailVerified.SET(postgres.Bool(account.IsEmailVerified)),
				table.Account.IsPhoneVerified.SET(postgres.Bool(account.IsPhoneVerified)),
				table.Account.Role.SET(postgres.String(account.Role)),
				table.Account.LastUsedSlyWallet.SET(postgres.String(account.LastUsedSlyWallet)),
			).WHERE(
			table.Account.ID.EQ(postgres.UUID(account.ID)),
		).RETURNING(
			table.Account.AllColumns,
		)
	} else {
		stmt = table.Account.UPDATE().
			SET(
				table.Account.FirstName.SET(postgres.String(account.FirstName)),
				table.Account.LastName.SET(postgres.String(account.LastName)),
				table.Account.Phone.SET(postgres.String(account.Phone)),
				table.Account.IsEmailVerified.SET(postgres.Bool(account.IsEmailVerified)),
				table.Account.IsPhoneVerified.SET(postgres.Bool(account.IsPhoneVerified)),
				table.Account.Role.SET(postgres.String(account.Role)),
				table.Account.LastUsedSlyWallet.SET(postgres.String(account.LastUsedSlyWallet)),
			).WHERE(
			table.Account.ID.EQ(postgres.UUID(account.ID)),
		).RETURNING(
			table.Account.AllColumns,
		)
	}

	var dbAccount model.Account
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbAccount)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return mapAccountToModel(dbAccount), nil
}

// UpdatePassword updates an account's password
func (r *AccountRepository) UpdatePassword(ctx context.Context, accountID uuid.UUID, hashedPassword string) error {
	stmt := table.Account.UPDATE().
		SET(
			table.Account.PasswordHashed.SET(postgres.String(hashedPassword)),
		).WHERE(
		table.Account.ID.EQ(postgres.UUID(accountID)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return DBItemNotFound
	}

	return nil
}

// Delete deletes an account and all associated data (can be used in a transaction)
func (r *AccountRepository) Delete(ctx context.Context, accountID uuid.UUID) error {
	return r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Delete associated ECDSA keys
		ecdsaStmt := table.Ecdsa.DELETE().WHERE(
			table.Ecdsa.AccountID.EQ(postgres.UUID(accountID)),
		)

		_, err := ecdsaStmt.ExecContext(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to delete account's ECDSA keys: %w", err)
		}

		// Delete associated SlyWallets
		slyWalletStmt := table.SlyWallet.DELETE().WHERE(
			table.SlyWallet.AccountID.EQ(postgres.UUID(accountID)),
		)

		_, err = slyWalletStmt.ExecContext(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to delete account's SlyWallets: %w", err)
		}

		// Delete the account
		accountStmt := table.Account.DELETE().WHERE(
			table.Account.ID.EQ(postgres.UUID(accountID)),
		)

		result, err := accountStmt.ExecContext(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to delete account: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return DBItemNotFound
		}

		return nil
	})
}

// GetWithEcdsas retrieves an account with its ECDSA keys
func (r *AccountRepository) GetWithEcdsas(ctx context.Context, accountID uuid.UUID) (*AccountModel, error) {
	// First get the account
	account, err := r.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Then get the ECDSA keys
	ecdsaStmt := postgres.SELECT(
		table.Ecdsa.AllColumns,
	).FROM(
		table.Ecdsa,
	).WHERE(
		table.Ecdsa.AccountID.EQ(postgres.UUID(accountID)),
	)

	var dbEcdsas []model.Ecdsa
	err = ecdsaStmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsas)
	if err != nil {
		return nil, fmt.Errorf("failed to get account's ECDSA keys: %w", err)
	}

	// Map ECDSA keys
	ecdsas := make([]EcdsaModel, len(dbEcdsas))
	for i, dbEcdsa := range dbEcdsas {
		ecdsas[i] = *mapEcdsaToModel(dbEcdsa)
	}

	account.Ecdsas = ecdsas

	return account, nil
}

// GetWithSlyWallets retrieves an account with its SlyWallets
func (r *AccountRepository) GetWithSlyWallets(ctx context.Context, accountID uuid.UUID) (*AccountModel, error) {
	// First get the account
	account, err := r.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Then get the SlyWallets
	slyWalletStmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.AccountID.EQ(postgres.UUID(accountID)),
	)

	var dbSlyWallets []model.SlyWallet
	err = slyWalletStmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallets)
	if err != nil {
		return nil, fmt.Errorf("failed to get account's SlyWallets: %w", err)
	}

	// Map SlyWallets
	slyWallets := make([]SlyWalletModel, len(dbSlyWallets))
	for i, dbSlyWallet := range dbSlyWallets {
		slyWallets[i] = *mapSlyWalletToModel(dbSlyWallet)
	}

	account.SlyWallets = slyWallets

	return account, nil
}

// GetCompleteAccount retrieves an account with all its related data
func (r *AccountRepository) GetCompleteAccount(ctx context.Context, accountID uuid.UUID) (*AccountModel, error) {
	// First get the account
	account, err := r.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Get ECDSA keys
	ecdsaRepo := NewEcdsaRepository(r.db)
	ecdsas, err := ecdsaRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Get SlyWallets
	slyWalletRepo := NewSlyWalletRepository(r.db)
	slyWallets, err := slyWalletRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Assign relations
	account.Ecdsas = ecdsas
	account.SlyWallets = slyWallets

	return account, nil
}

// Helper function to map Account model to AccountModel
func mapAccountToModel(account model.Account) *AccountModel {
	return &AccountModel{
		ID:                account.ID,
		FirstName:         account.FirstName,
		LastName:          account.LastName,
		Phone:             account.Phone,
		Email:             account.Email,
		IsEmailVerified:   account.IsEmailVerified,
		IsPhoneVerified:   account.IsPhoneVerified,
		PasswordHashed:    account.PasswordHashed,
		InvitationCode:    account.InvitationCode,
		Role:              account.Role,
		LastUsedSlyWallet: account.LastUsedSlyWallet,
		CreatedAt:         account.CreatedAt,
		UpdatedAt:         account.UpdatedAt,
	}
}
