//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type SlyWallet struct {
	Address           string `sql:"primary_key"`
	Chainid           string
	AccountID         uuid.UUID
	TransactionHash   string
	TransactionStatus string
	InvitationCode    string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
