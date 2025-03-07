package repo

import (
	"context"
	"database/sql"
	"log"
	common2 "yip/src/common"
)

// Database is the main repository struct for database interactions
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new Database instance
func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// GetDB returns the underlying database connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// WithTransaction executes a function within a transaction and handles commit/rollback
func (d *Database) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback transaction in case of panic
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
			// Re-throw panic
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		// Rollback transaction on error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Error rolling back transaction: %v", rollbackErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// AllowedAccountFilters defines the fields that can be used for sorting/filtering accounts
var AllowedAccountFilters = map[string]struct{}{
	"createdAt":       {},
	"updatedAt":       {},
	"firstName":       {},
	"lastName":        {},
	"email":           {},
	"isEmailVerified": {},
	"isPhoneVerified": {},
	"role":            {},
}

// AllowedEcdsaFilters defines the fields that can be used for sorting/filtering ECDSA keys
var AllowedEcdsaFilters = map[string]struct{}{
	"createdAt": {},
	"updatedAt": {},
	"address":   {},
}

// AllowedSlyWalletFilters defines the fields that can be used for sorting/filtering SlyWallets
var AllowedSlyWalletFilters = map[string]struct{}{
	"createdAt":         {},
	"updatedAt":         {},
	"address":           {},
	"chainid":           {},
	"transactionStatus": {},
}

// AllowedInvitationCodeFilters defines the fields that can be used for sorting/filtering invitation codes
var AllowedInvitationCodeFilters = map[string]struct{}{
	"expiresAt": {},
	"code":      {},
}

// NewPaginationQuery creates a pagination query with appropriate filters
func NewPaginationQuery(filters common2.AllowedFilters) *common2.PaginationQuery {
	return &common2.PaginationQuery{
		QueryFilters: common2.QueryFilters{
			Filters: filters,
		},
		PageSize:  common2.DefaultDefaultPageSize,
		Offset:    0,
		SortBy:    common2.DefaultSortBy,
		SortOrder: common2.DefaultSortOrder,
	}
}
