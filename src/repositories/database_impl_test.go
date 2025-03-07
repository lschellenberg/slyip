package repositories

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseImpl_GetECDSAAndContractByAddress(t *testing.T) {
	d := NewDatabase(nil)
	d.GetSLYWallets(context.Background(), uuid.New())
}
