package repo

import (
	"context"
	"fmt"
	"time"
	"yip/src/common"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"yip/.gen/slyip/slyip/model"
	"yip/.gen/slyip/slyip/table"
)

// SlyWalletRepository handles all SlyWallet related database operations
type SlyWalletRepository struct {
	db *Database
}

// NewSlyWalletRepository creates a new SlyWallet repository
func NewSlyWalletRepository(db *Database) *SlyWalletRepository {
	return &SlyWalletRepository{
		db: db,
	}
}

// Create creates a new SlyWallet
func (r *SlyWalletRepository) Create(ctx context.Context, slyWallet *SlyWalletModel) (*SlyWalletModel, error) {
	now := time.Now()
	slyWallet.CreatedAt = now
	slyWallet.UpdatedAt = now

	stmt := table.SlyWallet.INSERT(
		table.SlyWallet.Address,
		table.SlyWallet.Chainid,
		table.SlyWallet.AccountID,
		table.SlyWallet.TransactionHash,
		table.SlyWallet.TransactionStatus,
		table.SlyWallet.InvitationCode,
		table.SlyWallet.CreatedAt,
		table.SlyWallet.UpdatedAt,
	).VALUES(
		slyWallet.Address,
		slyWallet.Chainid,
		slyWallet.AccountID,
		slyWallet.TransactionHash,
		slyWallet.TransactionStatus,
		slyWallet.InvitationCode,
		slyWallet.CreatedAt,
		slyWallet.UpdatedAt,
	).RETURNING(
		table.SlyWallet.AllColumns,
	)

	var dbSlyWallet model.SlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create SlyWallet: %w", err)
	}

	return mapSlyWalletToModel(dbSlyWallet), nil
}

// GetByAddress retrieves a SlyWallet by address
func (r *SlyWalletRepository) GetByAddress(ctx context.Context, address string) (*SlyWalletModel, error) {
	stmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.Address.EQ(postgres.String(address)),
	)

	var dbSlyWallet model.SlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallet)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get SlyWallet by address: %w", err)
	}

	return mapSlyWalletToModel(dbSlyWallet), nil
}

// GetByAccountID retrieves all SlyWallets for an account
func (r *SlyWalletRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]SlyWalletModel, error) {
	stmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.AccountID.EQ(postgres.UUID(accountID)),
	)

	var dbSlyWallets []model.SlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallets)
	if err != nil {
		return nil, fmt.Errorf("failed to get SlyWallets by account ID: %w", err)
	}

	slyWallets := make([]SlyWalletModel, len(dbSlyWallets))
	for i, dbSlyWallet := range dbSlyWallets {
		slyWallets[i] = *mapSlyWalletToModel(dbSlyWallet)
	}

	return slyWallets, nil
}

// GetByInvitationCode retrieves all SlyWallets with a specific invitation code
func (r *SlyWalletRepository) GetByInvitationCode(ctx context.Context, code string) ([]SlyWalletModel, error) {
	stmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.InvitationCode.EQ(postgres.String(code)),
	)

	var dbSlyWallets []model.SlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallets)
	if err != nil {
		return nil, fmt.Errorf("failed to get SlyWallets by invitation code: %w", err)
	}

	slyWallets := make([]SlyWalletModel, len(dbSlyWallets))
	for i, dbSlyWallet := range dbSlyWallets {
		slyWallets[i] = *mapSlyWalletToModel(dbSlyWallet)
	}

	return slyWallets, nil
}

// GetWithEcdsaConnections retrieves a SlyWallet with its ECDSA connections
func (r *SlyWalletRepository) GetWithEcdsaConnections(ctx context.Context, address string) (*SlyWalletModel, error) {
	// First get the SlyWallet
	slyWallet, err := r.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	// Then get the EcdsaSlyWallet connections
	connectionStmt := postgres.SELECT(
		table.EcdsaSlyWallet.AllColumns,
	).FROM(
		table.EcdsaSlyWallet,
	).WHERE(
		table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(address)),
	)

	var dbConnections []model.EcdsaSlyWallet
	err = connectionStmt.QueryContext(ctx, r.db.GetDB(), &dbConnections)
	if err != nil {
		return nil, fmt.Errorf("failed to get SlyWallet ECDSA connections: %w", err)
	}

	// Map connections
	connections := make([]EcdsaSlyWalletModel, len(dbConnections))
	for i, dbConn := range dbConnections {
		connections[i] = *mapEcdsaSlyWalletToModel(dbConn)
	}

	slyWallet.EcdsaConnections = connections

	return slyWallet, nil
}

// List retrieves a paginated list of SlyWallets
func (r *SlyWalletRepository) List(ctx context.Context, query *common.PaginationQuery) (*PaginatedResponse[SlyWalletModel], error) {
	if query == nil {
		query = NewPaginationQuery(AllowedSlyWalletFilters)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Count total
	countStmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("total"),
	).FROM(
		table.SlyWallet,
	)

	var totalCount struct {
		Total uint64 `sql:"total"`
	}
	err := countStmt.QueryContext(ctx, r.db.GetDB(), &totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count SlyWallets: %w", err)
	}

	// Get data with pagination
	selectStmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	)

	// Apply sorting
	switch query.SortBy {
	case "address":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Address.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Address.DESC())
		}
	case "chainid":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Chainid.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Chainid.DESC())
		}
	case "transactionStatus":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.TransactionStatus.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.TransactionStatus.DESC())
		}
	case "updatedAt":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.UpdatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.UpdatedAt.DESC())
		}
	default: // createdAt is default
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.CreatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.CreatedAt.DESC())
		}
	}

	// Apply pagination
	selectStmt = selectStmt.
		OFFSET(int64(query.Offset)).
		LIMIT(int64(query.PageSize))

	var dbSlyWallets []model.SlyWallet
	err = selectStmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallets)
	if err != nil {
		return nil, fmt.Errorf("failed to list SlyWallets: %w", err)
	}

	slyWallets := make([]SlyWalletModel, len(dbSlyWallets))
	for i, dbSlyWallet := range dbSlyWallets {
		slyWallets[i] = *mapSlyWalletToModel(dbSlyWallet)
	}

	return &PaginatedResponse[SlyWalletModel]{
		Data:      slyWallets,
		PageSize:  query.PageSize,
		Offset:    query.Offset,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Total:     totalCount.Total,
		Count:     uint64(len(slyWallets)),
		Filter:    make(map[string]string),
	}, nil
}

// ListByAccountID retrieves a paginated list of SlyWallets for an account
func (r *SlyWalletRepository) ListByAccountID(ctx context.Context, accountID uuid.UUID, query *common.PaginationQuery) (*PaginatedResponse[SlyWalletModel], error) {
	if query == nil {
		query = NewPaginationQuery(AllowedSlyWalletFilters)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Count total for this account
	countStmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("total"),
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.AccountID.EQ(postgres.UUID(accountID)),
	)

	var totalCount struct {
		Total uint64 `sql:"total"`
	}
	err := countStmt.QueryContext(ctx, r.db.GetDB(), &totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count account's SlyWallets: %w", err)
	}

	// Get data with pagination
	selectStmt := postgres.SELECT(
		table.SlyWallet.AllColumns,
	).FROM(
		table.SlyWallet,
	).WHERE(
		table.SlyWallet.AccountID.EQ(postgres.UUID(accountID)),
	)

	// Apply sorting
	switch query.SortBy {
	case "address":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Address.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Address.DESC())
		}
	case "chainid":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Chainid.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.Chainid.DESC())
		}
	case "transactionStatus":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.TransactionStatus.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.TransactionStatus.DESC())
		}
	case "updatedAt":
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.UpdatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.UpdatedAt.DESC())
		}
	default: // createdAt is default
		if query.SortOrder == common.SortingOrderASC {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.CreatedAt.ASC())
		} else {
			selectStmt = selectStmt.ORDER_BY(table.SlyWallet.CreatedAt.DESC())
		}
	}

	// Apply pagination
	selectStmt = selectStmt.
		OFFSET(int64(query.Offset)).
		LIMIT(int64(query.PageSize))

	var dbSlyWallets []model.SlyWallet
	err = selectStmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallets)
	if err != nil {
		return nil, fmt.Errorf("failed to list account's SlyWallets: %w", err)
	}

	slyWallets := make([]SlyWalletModel, len(dbSlyWallets))
	for i, dbSlyWallet := range dbSlyWallets {
		slyWallets[i] = *mapSlyWalletToModel(dbSlyWallet)
	}

	return &PaginatedResponse[SlyWalletModel]{
		Data:      slyWallets,
		PageSize:  query.PageSize,
		Offset:    query.Offset,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Total:     totalCount.Total,
		Count:     uint64(len(slyWallets)),
		Filter:    make(map[string]string),
	}, nil
}

// Update updates a SlyWallet
func (r *SlyWalletRepository) Update(ctx context.Context, slyWallet *SlyWalletModel) (*SlyWalletModel, error) {

	stmt := table.SlyWallet.UPDATE().
		SET(
			table.SlyWallet.AccountID.SET(postgres.UUID(slyWallet.AccountID)),
			table.SlyWallet.Chainid.SET(postgres.String(slyWallet.Chainid)),
			table.SlyWallet.TransactionHash.SET(postgres.String(slyWallet.TransactionHash)),
			table.SlyWallet.TransactionStatus.SET(postgres.String(slyWallet.TransactionStatus)),
			table.SlyWallet.InvitationCode.SET(postgres.String(slyWallet.InvitationCode)),
		).WHERE(
		table.SlyWallet.Address.EQ(postgres.String(slyWallet.Address)),
	).RETURNING(
		table.SlyWallet.AllColumns,
	)

	var dbSlyWallet model.SlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbSlyWallet)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to update SlyWallet: %w", err)
	}

	return mapSlyWalletToModel(dbSlyWallet), nil
}

// UpdateTransactionStatus updates a SlyWallet's transaction status
func (r *SlyWalletRepository) UpdateTransactionStatus(ctx context.Context, address string, status string) error {

	stmt := table.SlyWallet.UPDATE().
		SET(
			table.SlyWallet.TransactionStatus.SET(postgres.String(status)),
		).WHERE(
		table.SlyWallet.Address.EQ(postgres.String(address)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to update SlyWallet transaction status: %w", err)
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

// Delete deletes a SlyWallet
func (r *SlyWalletRepository) Delete(ctx context.Context, address string) error {
	// First delete any EcdsaSlyWallet connections
	connStmt := table.EcdsaSlyWallet.DELETE().WHERE(
		table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(address)),
	)

	_, err := connStmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to delete SlyWallet ECDSA connections: %w", err)
	}

	// Then delete the SlyWallet
	stmt := table.SlyWallet.DELETE().WHERE(
		table.SlyWallet.Address.EQ(postgres.String(address)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to delete SlyWallet: %w", err)
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

// ConnectToEcdsa connects a SlyWallet to an ECDSA key
func (r *SlyWalletRepository) ConnectToEcdsa(ctx context.Context, connection *EcdsaSlyWalletModel) error {
	stmt := table.EcdsaSlyWallet.INSERT(
		table.EcdsaSlyWallet.EcdsaAddress,
		table.EcdsaSlyWallet.OnChainAccountAddress,
		table.EcdsaSlyWallet.OnChainPermissions,
	).VALUES(
		connection.EcdsaAddress,
		connection.OnChainAccountAddress,
		connection.OnChainPermissions,
	)

	_, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to connect SlyWallet to ECDSA key: %w", err)
	}

	return nil
}

// DisconnectFromEcdsa disconnects a SlyWallet from an ECDSA key
func (r *SlyWalletRepository) DisconnectFromEcdsa(ctx context.Context, slyWalletAddress string, ecdsaAddress string) error {
	stmt := table.EcdsaSlyWallet.DELETE().WHERE(
		table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(slyWalletAddress)).
			AND(table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(ecdsaAddress))),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to disconnect SlyWallet from ECDSA key: %w", err)
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

// Helper function to map SlyWallet model to SlyWalletModel
func mapSlyWalletToModel(slyWallet model.SlyWallet) *SlyWalletModel {
	return &SlyWalletModel{
		Address:           slyWallet.Address,
		Chainid:           slyWallet.Chainid,
		AccountID:         slyWallet.AccountID,
		TransactionHash:   slyWallet.TransactionHash,
		TransactionStatus: slyWallet.TransactionStatus,
		InvitationCode:    slyWallet.InvitationCode,
		CreatedAt:         slyWallet.CreatedAt,
		UpdatedAt:         slyWallet.UpdatedAt,
	}
}

// Helper function to map EcdsaSlyWallet model to EcdsaSlyWalletModel
func mapEcdsaSlyWalletToModel(ecdsaSlyWallet model.EcdsaSlyWallet) *EcdsaSlyWalletModel {
	return &EcdsaSlyWalletModel{
		EcdsaAddress:          ecdsaSlyWallet.EcdsaAddress,
		OnChainAccountAddress: ecdsaSlyWallet.OnChainAccountAddress,
		OnChainPermissions:    ecdsaSlyWallet.OnChainPermissions,
	}
}
