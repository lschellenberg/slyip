//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Ecdsa = newEcdsaTable("slyip", "ecdsa", "")

type ecdsaTable struct {
	postgres.Table

	//Columns
	Address   postgres.ColumnString
	AccountID postgres.ColumnString
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type EcdsaTable struct {
	ecdsaTable

	EXCLUDED ecdsaTable
}

// AS creates new EcdsaTable with assigned alias
func (a EcdsaTable) AS(alias string) *EcdsaTable {
	return newEcdsaTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new EcdsaTable with assigned schema name
func (a EcdsaTable) FromSchema(schemaName string) *EcdsaTable {
	return newEcdsaTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new EcdsaTable with assigned table prefix
func (a EcdsaTable) WithPrefix(prefix string) *EcdsaTable {
	return newEcdsaTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new EcdsaTable with assigned table suffix
func (a EcdsaTable) WithSuffix(suffix string) *EcdsaTable {
	return newEcdsaTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newEcdsaTable(schemaName, tableName, alias string) *EcdsaTable {
	return &EcdsaTable{
		ecdsaTable: newEcdsaTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newEcdsaTableImpl("", "excluded", ""),
	}
}

func newEcdsaTableImpl(schemaName, tableName, alias string) ecdsaTable {
	var (
		AddressColumn   = postgres.StringColumn("address")
		AccountIDColumn = postgres.StringColumn("account_id")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn = postgres.TimestampzColumn("updated_at")
		allColumns      = postgres.ColumnList{AddressColumn, AccountIDColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{AccountIDColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return ecdsaTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Address:   AddressColumn,
		AccountID: AccountIDColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
