package repo

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"yip/.gen/slyip/slyip/model"
	"yip/.gen/slyip/slyip/table"
)

// EcdsaSlyWalletRepository handles all EcdsaSlyWallet related database operations
type EcdsaSlyWalletRepository struct {
	db *Database
}

// NewEcdsaSlyWalletRepository creates a new EcdsaSlyWallet repository
func NewEcdsaSlyWalletRepository(db *Database) *EcdsaSlyWalletRepository {
	return &EcdsaSlyWalletRepository{
		db: db,
	}
}

// Create creates a new EcdsaSlyWallet connection
func (r *EcdsaSlyWalletRepository) Create(ctx context.Context, connection *EcdsaSlyWalletModel) (*EcdsaSlyWalletModel, error) {
	stmt := table.EcdsaSlyWallet.INSERT(
		table.EcdsaSlyWallet.EcdsaAddress,
		table.EcdsaSlyWallet.OnChainAccountAddress,
		table.EcdsaSlyWallet.OnChainPermissions,
	).VALUES(
		connection.EcdsaAddress,
		connection.OnChainAccountAddress,
		connection.OnChainPermissions,
	).RETURNING(
		table.EcdsaSlyWallet.AllColumns,
	)

	var dbConnection model.EcdsaSlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to create EcdsaSlyWallet connection: %w", err)
	}

	return mapEcdsaSlyWalletToModel(dbConnection), nil
}

// GetByAddresses retrieves an EcdsaSlyWallet connection by ECDSA address and SlyWallet address
func (r *EcdsaSlyWalletRepository) GetByAddresses(ctx context.Context, ecdsaAddress string, slyWalletAddress string) (*EcdsaSlyWalletModel, error) {
	stmt := postgres.SELECT(
		table.EcdsaSlyWallet.AllColumns,
	).FROM(
		table.EcdsaSlyWallet,
	).WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(ecdsaAddress)).
			AND(table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(slyWalletAddress))),
	)

	var dbConnection model.EcdsaSlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbConnection)
	if err != nil {
		if err == qrm.ErrNoRows {
			return nil, DBItemNotFound
		}
		return nil, fmt.Errorf("failed to get EcdsaSlyWallet connection: %w", err)
	}

	return mapEcdsaSlyWalletToModel(dbConnection), nil
}

// GetByEcdsaAddress retrieves all EcdsaSlyWallet connections for an ECDSA address
func (r *EcdsaSlyWalletRepository) GetByEcdsaAddress(ctx context.Context, ecdsaAddress string) ([]EcdsaSlyWalletModel, error) {
	stmt := postgres.SELECT(
		table.EcdsaSlyWallet.AllColumns,
	).FROM(
		table.EcdsaSlyWallet,
	).WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(ecdsaAddress)),
	)

	var dbConnections []model.EcdsaSlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbConnections)
	if err != nil {
		return nil, fmt.Errorf("failed to get EcdsaSlyWallet connections by ECDSA address: %w", err)
	}

	connections := make([]EcdsaSlyWalletModel, len(dbConnections))
	for i, dbConnection := range dbConnections {
		connections[i] = *mapEcdsaSlyWalletToModel(dbConnection)
	}

	return connections, nil
}

// GetBySlyWalletAddress retrieves all EcdsaSlyWallet connections for a SlyWallet address
func (r *EcdsaSlyWalletRepository) GetBySlyWalletAddress(ctx context.Context, slyWalletAddress string) ([]EcdsaSlyWalletModel, error) {
	stmt := postgres.SELECT(
		table.EcdsaSlyWallet.AllColumns,
	).FROM(
		table.EcdsaSlyWallet,
	).WHERE(
		table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(slyWalletAddress)),
	)

	var dbConnections []model.EcdsaSlyWallet
	err := stmt.QueryContext(ctx, r.db.GetDB(), &dbConnections)
	if err != nil {
		return nil, fmt.Errorf("failed to get EcdsaSlyWallet connections by SlyWallet address: %w", err)
	}

	connections := make([]EcdsaSlyWalletModel, len(dbConnections))
	for i, dbConnection := range dbConnections {
		connections[i] = *mapEcdsaSlyWalletToModel(dbConnection)
	}

	return connections, nil
}

// UpdatePermissions updates an EcdsaSlyWallet connection's permissions
func (r *EcdsaSlyWalletRepository) UpdatePermissions(ctx context.Context, ecdsaAddress string, slyWalletAddress string, permissions int32) error {
	stmt := table.EcdsaSlyWallet.UPDATE().
		SET(
			table.EcdsaSlyWallet.OnChainPermissions.SET(postgres.Int32(permissions)),
		).WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(ecdsaAddress)).
			AND(table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(slyWalletAddress))),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to update EcdsaSlyWallet permissions: %w", err)
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

// Delete deletes an EcdsaSlyWallet connection
func (r *EcdsaSlyWalletRepository) Delete(ctx context.Context, ecdsaAddress string, slyWalletAddress string) error {
	stmt := table.EcdsaSlyWallet.DELETE().WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(ecdsaAddress)).
			AND(table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(slyWalletAddress))),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return fmt.Errorf("failed to delete EcdsaSlyWallet connection: %w", err)
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

// DeleteByEcdsaAddress deletes all EcdsaSlyWallet connections for an ECDSA address
func (r *EcdsaSlyWalletRepository) DeleteByEcdsaAddress(ctx context.Context, ecdsaAddress string) (int64, error) {
	stmt := table.EcdsaSlyWallet.DELETE().WHERE(
		table.EcdsaSlyWallet.EcdsaAddress.EQ(postgres.String(ecdsaAddress)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return 0, fmt.Errorf("failed to delete EcdsaSlyWallet connections by ECDSA address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// DeleteBySlyWalletAddress deletes all EcdsaSlyWallet connections for a SlyWallet address
func (r *EcdsaSlyWalletRepository) DeleteBySlyWalletAddress(ctx context.Context, slyWalletAddress string) (int64, error) {
	stmt := table.EcdsaSlyWallet.DELETE().WHERE(
		table.EcdsaSlyWallet.OnChainAccountAddress.EQ(postgres.String(slyWalletAddress)),
	)

	result, err := stmt.ExecContext(ctx, r.db.GetDB())
	if err != nil {
		return 0, fmt.Errorf("failed to delete EcdsaSlyWallet connections by SlyWallet address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// AddECDSAToSlyWallet connects an ECDSA key to a SlyWallet
func (r *EcdsaSlyWalletRepository) AddECDSAToSlyWallet(ctx context.Context, ecdsaAddress string, slyWalletAddress string, permissions int32) (*EcdsaSlyWalletModel, error) {
	connection := &EcdsaSlyWalletModel{
		EcdsaAddress:          ecdsaAddress,
		OnChainAccountAddress: slyWalletAddress,
		OnChainPermissions:    permissions,
	}

	// First check if connection already exists
	existing, err := r.GetByAddresses(ctx, ecdsaAddress, slyWalletAddress)
	if err != nil && err != DBItemNotFound {
		return nil, fmt.Errorf("failed to check existing connection: %w", err)
	}

	if existing != nil {
		// If exists and permissions are different, update permissions
		if existing.OnChainPermissions != permissions {
			err = r.UpdatePermissions(ctx, ecdsaAddress, slyWalletAddress, permissions)
			if err != nil {
				return nil, fmt.Errorf("failed to update connection permissions: %w", err)
			}
			// Fetch updated connection
			return r.GetByAddresses(ctx, ecdsaAddress, slyWalletAddress)
		}
		// If exists with same permissions, just return it
		return existing, nil
	}

	// Otherwise create new connection
	return r.Create(ctx, connection)
}

// RemoveECDSAFromSlyWallet disconnects an ECDSA key from a SlyWallet
func (r *EcdsaSlyWalletRepository) RemoveECDSAFromSlyWallet(ctx context.Context, ecdsaAddress string, slyWalletAddress string) error {
	// Check if connection exists
	_, err := r.GetByAddresses(ctx, ecdsaAddress, slyWalletAddress)
	if err != nil {
		if err == DBItemNotFound {
			return nil // Already not connected, not an error
		}
		return fmt.Errorf("failed to check existing connection: %w", err)
	}

	// Delete the connection
	err = r.Delete(ctx, ecdsaAddress, slyWalletAddress)
	if err != nil {
		return fmt.Errorf("failed to remove ECDSA from SlyWallet: %w", err)
	}

	return nil
}
