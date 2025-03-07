package repo

import (
	"context"
	"fmt"
	"time"
	"yip/src/common"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"yip/.gen/slyip/slyip/model"
	"yip/.gen/slyip/slyip/table"
)

// InvitationCodeRepository handles all InvitationCode related database operations
type InvitationCodeRepository struct {
	db *Database
}

// NewInvitationCodeRepository creates a new InvitationCode repository
func NewInvitationCodeRepository(db *Database) *InvitationCodeRepository {
	return &InvitationCodeRepository{
		db: db,
	}
}

// Create creates a new InvitationCode
func (r *InvitationCodeRepository) Create(ctx context.Context, invitationCode *InvitationCodeModel) (*InvitationCodeModel, error) {
	stmt := table.InvitationCode.INSERT(
		table.InvitationCode.Code,
		table.InvitationCode.TransactionHash,
		table.InvitationCode.ExpiresAt,
	).VALUES(
		invitationCode.Code,
		invitationCode.TransactionHash,
		invitationCode.ExpiresAt,
	).RETURNING(
		table.InvitationCode.AllColumns,
	)

	var dbInvitationCode model.InvitationCode
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbInvitationCode)
	if err != nil {
		return nil, fmt.Errorf("failed to create InvitationCode: %w", err)
	}

	return mapInvitationCodeToModel(dbInvitationCode), nil
}

// GetByCode retrieves an InvitationCode by code
func (r *InvitationCodeRepository) GetByCode(ctx context.Context, code string) (*InvitationCodeModel, error) {
	stmt := postgres.SELECT(
		table.InvitationCode.AllColumns,
	).FROM(
		table.InvitationCode,
	).WHERE(
		table.InvitationCode.Code.EQ(postgres.String(code)),
	)

	var dbInvitationCode model.InvitationCode
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbInvitationCode)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get InvitationCode by code: %w", err)
	}

	return mapInvitationCodeToModel(dbInvitationCode), nil
}

// GetByCode retrieves an InvitationCode by code
func (r *InvitationCodeRepository) GetByTransactionHash(ctx context.Context, hash string) (*InvitationCodeModel, error) {
	stmt := postgres.SELECT(
		table.InvitationCode.AllColumns,
	).FROM(
		table.InvitationCode,
	).WHERE(
		table.InvitationCode.TransactionHash.EQ(postgres.String(hash)),
	)

	var dbInvitationCode model.InvitationCode
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbInvitationCode)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get InvitationCode by code: %w", err)
	}

	return mapInvitationCodeToModel(dbInvitationCode), nil
}

// GetValid retrieves all valid (not expired) invitation codes
func (r *InvitationCodeRepository) GetValid(ctx context.Context) ([]InvitationCodeModel, error) {
	now := time.Now()

	stmt := postgres.SELECT(
		table.InvitationCode.AllColumns,
	).FROM(
		table.InvitationCode,
	).WHERE(
		table.InvitationCode.ExpiresAt.GT(postgres.TimestampzT(now)),
	)

	var dbInvitationCodes []model.InvitationCode
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbInvitationCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid InvitationCodes: %w", err)
	}

	invitationCodes := make([]InvitationCodeModel, len(dbInvitationCodes))
	for i, dbInvitationCode := range dbInvitationCodes {
		invitationCodes[i] = *mapInvitationCodeToModel(dbInvitationCode)
	}

	return invitationCodes, nil
}

// IsValid checks if an invitation code is valid (exists and not expired)
func (r *InvitationCodeRepository) IsValid(ctx context.Context, code string) (bool, error) {
	now := time.Now()

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("count"),
	).FROM(
		table.InvitationCode,
	).WHERE(
		table.InvitationCode.Code.EQ(postgres.String(code)).
			AND(table.InvitationCode.ExpiresAt.GT(postgres.TimestampzT(now))),
	)

	var count struct {
		Count int `sql:"count"`
	}
	err := stmt.QueryContext(ctx, r.db.GetDB(), &count)
	if err != nil {
		return false, fmt.Errorf("failed to check if InvitationCode is valid: %w", err)
	}

	return count.Count > 0, nil
}

// List retrieves a paginated list of InvitationCodes
func (r *InvitationCodeRepository) List(ctx context.Context, query *common.PaginationQuery) (*PaginatedResponse[InvitationCodeModel], error) {
	if query == nil {
		query = NewPaginationQuery(AllowedInvitationCodeFilters)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Count total
	countStmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("total"),
	).FROM(
		table.InvitationCode,
	)

	var totalCount struct {
		Total uint64 `sql:"total"`
	}
	err := countStmt.QueryContext(ctx, r.db.GetDB(), &totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count InvitationCodes: %w", err)
	}

	// Get data with pagination
	selectStmt := postgres.SELECT(
		table.InvitationCode.AllColumns,
	).FROM(
		table.InvitationCode,
	)

	// Apply sorting
	switch query.SortBy {
	case "code":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.InvitationCode.Code.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.InvitationCode.Code.DESC())
		}
	default: // expiresAt is default
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.InvitationCode.ExpiresAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.InvitationCode.ExpiresAt.DESC())
		}
	}

	// Apply pagination
	selectStmt = selectStmt.
		OFFSET(int64(query.Offset)).
		LIMIT(int64(query.PageSize))

	var dbInvitationCodes []model.InvitationCode
	err = selectStmt.QueryContext(ctx, r.db.GetDB(), &dbInvitationCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to list InvitationCodes: %w", err)
	}

	invitationCodes := make([]InvitationCodeModel, len(dbInvitationCodes))
	for i, dbInvitationCode := range dbInvitationCodes {
		invitationCodes[i] = *mapInvitationCodeToModel(dbInvitationCode)
	}

	return &PaginatedResponse[InvitationCodeModel]{
		Data:      invitationCodes,
		PageSize:  query.PageSize,
		Offset:    query.Offset,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Total:     totalCount.Total,
		Count:     uint64(len(invitationCodes)),
		Filter:    make(map[string]string),
	}, nil
}

// Update updates an InvitationCode
func (r *InvitationCodeRepository) Update(ctx context.Context, invitationCode *InvitationCodeModel) (*InvitationCodeModel, error) {
	stmt := table.InvitationCode.UPDATE().
		SET(
			table.InvitationCode.TransactionHash.SET(postgres.String(invitationCode.TransactionHash)),
			table.InvitationCode.ExpiresAt.SET(postgres.TimestampzT(invitationCode.ExpiresAt)),
		).WHERE(
		table.InvitationCode.Code.EQ(postgres.String(invitationCode.Code)),
	).RETURNING(
		table.InvitationCode.AllColumns,
	)

	var dbInvitationCode model.InvitationCode
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbInvitationCode)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to update InvitationCode: %w", err)
	}

	return mapInvitationCodeToModel(dbInvitationCode), nil
}

// ExtendExpiration extends an invitation code's expiration date
func (r *InvitationCodeRepository) ExtendExpiration(ctx context.Context, code string, newExpiration time.Time) error {
	stmt := table.InvitationCode.UPDATE().
		SET(
			table.InvitationCode.ExpiresAt.SET(postgres.TimestampzT(newExpiration)),
		).WHERE(
		table.InvitationCode.Code.EQ(postgres.String(code)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to extend InvitationCode expiration: %w", err)
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

// Delete deletes an InvitationCode
func (r *InvitationCodeRepository) Delete(ctx context.Context, code string) error {
	stmt := table.InvitationCode.DELETE().WHERE(
		table.InvitationCode.Code.EQ(postgres.String(code)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to delete InvitationCode: %w", err)
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

// DeleteExpired deletes all expired invitation codes
func (r *InvitationCodeRepository) DeleteExpired(ctx context.Context) (int64, error) {
	now := time.Now()

	stmt := table.InvitationCode.DELETE().WHERE(
		table.InvitationCode.ExpiresAt.LT(postgres.TimestampzT(now)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired InvitationCodes: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetWithAccounts retrieves an InvitationCode with all accounts that used it
func (r *InvitationCodeRepository) GetWithAccounts(ctx context.Context, code string) (*InvitationCodeModel, error) {
	// First get the InvitationCode
	invitationCode, err := r.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Then get all accounts that used this invitation code
	accountStmt := postgres.SELECT(
		table.Account.AllColumns,
	).FROM(
		table.Account,
	).WHERE(
		table.Account.InvitationCode.EQ(postgres.String(code)),
	)

	var dbAccounts []model.Account
	err = accountStmt.QueryContext(ctx, r.db.GetDB(), &dbAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts using InvitationCode: %w", err)
	}

	// Map accounts
	accounts := make([]AccountModel, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		accounts[i] = *mapAccountToModel(dbAccount)
	}

	invitationCode.Accounts = accounts

	return invitationCode, nil
}

// GetWithSlyWallets retrieves an InvitationCode with all SlyWallets that used it
func (r *InvitationCodeRepository) GetWithSlyWallets(ctx context.Context, code string) (*InvitationCodeModel, error) {
	// First get the InvitationCode
	invitationCode, err := r.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Then get all SlyWallets that used this invitation code
	slyWalletStmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.InvitationCode.EQ(postgres.String(code)),
	)

	var dbSlyWallets []model.SlyWallet
	err = slyWalletStmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallets)
	if err != nil {
		return nil, fmt.Errorf("failed to get SlyWallets using InvitationCode: %w", err)
	}

	// Map SlyWallets
	slyWallets := make([]SlyWalletModel, len(dbSlyWallets))
	for i, dbSlyWallet := range dbSlyWallets {
		slyWallets[i] = *mapSlyWalletToModel(dbSlyWallet)
	}

	invitationCode.SlyWallets = slyWallets

	return invitationCode, nil
}

// Helper function to map InvitationCode model to InvitationCodeModel
func mapInvitationCodeToModel(invitationCode model.InvitationCode) *InvitationCodeModel {
	return &InvitationCodeModel{
		Code:            invitationCode.Code,
		TransactionHash: invitationCode.TransactionHash,
		ExpiresAt:       invitationCode.ExpiresAt,
	}
}
