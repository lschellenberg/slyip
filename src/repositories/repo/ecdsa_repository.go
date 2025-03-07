package repo

import (
	"context"
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

// EcdsaRepository handles all ECDSA related database operations
type EcdsaRepository struct {
	db *Database
}

// NewEcdsaRepository creates a new ECDSA repository
func NewEcdsaRepository(db *Database) *EcdsaRepository {
	return &EcdsaRepository{
		db: db,
	}
}

// Create creates a new ECDSA key
func (r *EcdsaRepository) Create(ctx context.Context, ecdsa *EcdsaModel) (*EcdsaModel, error) {
	now := time.Now()
	ecdsa.CreatedAt = now
	ecdsa.UpdatedAt = now

	stmt := table.Ecdsa.INSERT(
		table.Ecdsa.Address,
		table.Ecdsa.AccountID,
		table.Ecdsa.CreatedAt,
		table.Ecdsa.UpdatedAt,
	).VALUES(
		ecdsa.Address,
		ecdsa.AccountID,
		ecdsa.CreatedAt,
		ecdsa.UpdatedAt,
	).RETURNING(
		table.Ecdsa.AllColumns,
	)

	var dbEcdsa model.Ecdsa
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsa)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECDSA key: %w", err)
	}

	return mapEcdsaToModel(dbEcdsa), nil
}

// GetByAddress retrieves an ECDSA key by address
func (r *EcdsaRepository) GetByAddress(ctx context.Context, address string) (*EcdsaModel, error) {
	stmt := postgres.SELECT(
		table.Ecdsa.AllColumns,
	).FROM(
		table.Ecdsa,
	).WHERE(
		table.Ecdsa.Address.EQ(postgres.String(address)),
	)

	var dbEcdsa model.Ecdsa
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsa)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get ECDSA key by address: %w", err)
	}

	return mapEcdsaToModel(dbEcdsa), nil
}

// GetByAccountID retrieves all ECDSA keys for an account
func (r *EcdsaRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]EcdsaModel, error) {
	stmt := postgres.SELECT(
		table.Ecdsa.AllColumns,
	).FROM(
		table.Ecdsa,
	).WHERE(
		table.Ecdsa.AccountID.EQ(postgres.UUID(accountID)),
	)

	var dbEcdsas []model.Ecdsa
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsas)
	if err != nil {
		return nil, fmt.Errorf("failed to get ECDSA keys by account ID: %w", err)
	}

	ecdsas := make([]EcdsaModel, len(dbEcdsas))
	for i, dbEcdsa := range dbEcdsas {
		ecdsas[i] = *mapEcdsaToModel(dbEcdsa)
	}

	return ecdsas, nil
}

// GetWithSlyWallets retrieves an ECDSA key with its connected SlyWallets
func (r *EcdsaRepository) GetWithSlyWallets(ctx context.Context, address string) (*EcdsaModel, error) {
	// First get the ECDSA key
	ecdsa, err := r.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	// Then get the EcdsaSlyWallet connections
	connectionStmt := postgres.SELECT(
		table.EcdsaSlyWallet.AllColumns,
	).FROM(
		table.EcdsaSlyWallet,
	).WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(address)),
	)

	var dbConnections []model.EcdsaSlyWallet
	err = connectionStmt.QueryContext(ctx, r.db.GetDB(), &dbConnections)
	if err != nil {
		return nil, fmt.Errorf("failed to get ECDSA SlyWallet connections: %w", err)
	}

	// Map connections
	connections := make([]EcdsaSlyWalletModel, len(dbConnections))
	for i, dbConn := range dbConnections {
		connections[i] = *mapEcdsaSlyWalletToModel(dbConn)
	}

	ecdsa.SlyWallets = connections

	return ecdsa, nil
}

// List retrieves a paginated list of ECDSA keys
func (r *EcdsaRepository) List(ctx context.Context, query *common.PaginationQuery) (*PaginatedResponse[EcdsaModel], error) {
	if query == nil {
		query = NewPaginationQuery(AllowedEcdsaFilters)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Count total
	countStmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("total"),
	).FROM(
		table.Ecdsa,
	)

	var totalCount struct {
		Total uint64 `sql:"total"`
	}
	err := countStmt.QueryContext(ctx, r.db.GetDB(), &totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count ECDSA keys: %w", err)
	}

	// Get data with pagination
	selectStmt := postgres.SELECT(
		table.Ecdsa.AllColumns,
	).FROM(
		table.Ecdsa,
	)

	// Apply sorting
	switch query.SortBy {
	case "address":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.Address.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.Address.DESC())
		}
	case "updatedAt":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.UpdatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.UpdatedAt.DESC())
		}
	default: // createdAt is default
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.CreatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.CreatedAt.DESC())
		}
	}

	// Apply pagination
	selectStmt = selectStmt.
		OFFSET(int64(query.Offset)).
		LIMIT(int64(query.PageSize))

	var dbEcdsas []model.Ecdsa
	err = selectStmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsas)
	if err != nil {
		return nil, fmt.Errorf("failed to list ECDSA keys: %w", err)
	}

	ecdsas := make([]EcdsaModel, len(dbEcdsas))
	for i, dbEcdsa := range dbEcdsas {
		ecdsas[i] = *mapEcdsaToModel(dbEcdsa)
	}

	return &PaginatedResponse[EcdsaModel]{
		Data:      ecdsas,
		PageSize:  query.PageSize,
		Offset:    query.Offset,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Total:     totalCount.Total,
		Count:     uint64(len(ecdsas)),
		Filter:    make(map[string]string),
	}, nil
}

// ListByAccountID retrieves a paginated list of ECDSA keys for an account
func (r *EcdsaRepository) ListByAccountID(ctx context.Context, accountID uuid.UUID, query *common.PaginationQuery) (*PaginatedResponse[EcdsaModel], error) {
	if query == nil {
		query = NewPaginationQuery(AllowedEcdsaFilters)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Count total for this account
	countStmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("total"),
	).FROM(
		table.Ecdsa,
	).WHERE(
		table.Ecdsa.AccountID.EQ(postgres.UUID(accountID)),
	)

	var totalCount struct {
		Total uint64 `sql:"total"`
	}
	err := countStmt.QueryContext(ctx, r.db.GetDB(), &totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count account's ECDSA keys: %w", err)
	}

	// Get data with pagination
	selectStmt := postgres.SELECT(
		table.Ecdsa.AllColumns,
	).FROM(
		table.Ecdsa,
	).WHERE(
		table.Ecdsa.AccountID.EQ(postgres.UUID(accountID)),
	)

	// Apply sorting
	switch query.SortBy {
	case "address":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.Address.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.Address.DESC())
		}
	case "updatedAt":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.UpdatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.UpdatedAt.DESC())
		}
	default: // createdAt is default
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.CreatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.Ecdsa.CreatedAt.DESC())
		}
	}

	// Apply pagination
	selectStmt = selectStmt.
		OFFSET(int64(query.Offset)).
		LIMIT(int64(query.PageSize))

	var dbEcdsas []model.Ecdsa
	err = selectStmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsas)
	if err != nil {
		return nil, fmt.Errorf("failed to list account's ECDSA keys: %w", err)
	}

	ecdsas := make([]EcdsaModel, len(dbEcdsas))
	for i, dbEcdsa := range dbEcdsas {
		ecdsas[i] = *mapEcdsaToModel(dbEcdsa)
	}

	return &PaginatedResponse[EcdsaModel]{
		Data:      ecdsas,
		PageSize:  query.PageSize,
		Offset:    query.Offset,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Total:     totalCount.Total,
		Count:     uint64(len(ecdsas)),
		Filter:    make(map[string]string),
	}, nil
}

// Update updates an ECDSA key
func (r *EcdsaRepository) Update(ctx context.Context, ecdsa *EcdsaModel) (*EcdsaModel, error) {

	stmt := table.Ecdsa.UPDATE().
		SET(
			table.Ecdsa.AccountID.SET(postgres.UUID(ecdsa.AccountID)),
		).WHERE(
		table.Ecdsa.Address.EQ(postgres.String(ecdsa.Address)),
	).RETURNING(
		table.Ecdsa.AllColumns,
	)

	var dbEcdsa model.Ecdsa
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbEcdsa)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to update ECDSA key: %w", err)
	}

	return mapEcdsaToModel(dbEcdsa), nil
}

// Delete deletes an ECDSA key
func (r *EcdsaRepository) Delete(ctx context.Context, address string) error {
	// First delete any EcdsaSlyWallet connections
	connStmt := table.EcdsaSlyWallet.DELETE().WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(address)),
	)

	_, err := connStmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to delete ECDSA SlyWallet connections: %w", err)
	}

	// Then delete the ECDSA key
	stmt := table.Ecdsa.DELETE().WHERE(
		table.Ecdsa.Address.EQ(postgres.String(address)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to delete ECDSA key: %w", err)
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

// UpsertECDSA creates or updates an ECDSA key
func (r *EcdsaRepository) UpsertECDSA(ctx context.Context, ecdsa *EcdsaModel) (*EcdsaModel, error) {
	// Check if ECDSA already exists
	existing, err := r.GetByAddress(ctx, ecdsa.Address)
	if err != nil {
		if err == DBItemNotFound {
			// If not found, create new
			return r.Create(ctx, ecdsa)
		}
		return nil, fmt.Errorf("failed to check existing ECDSA key: %w", err)
	}

	// If found, update it
	ecdsa.CreatedAt = existing.CreatedAt
	ecdsa.UpdatedAt = time.Now()
	return r.Update(ctx, ecdsa)
}

// Helper function to map Ecdsa model to EcdsaModel
func mapEcdsaToModel(ecdsa model.Ecdsa) *EcdsaModel {
	return &EcdsaModel{
		Address:   ecdsa.Address,
		AccountID: ecdsa.AccountID,
		CreatedAt: ecdsa.CreatedAt,
		UpdatedAt: ecdsa.UpdatedAt,
	}
}
