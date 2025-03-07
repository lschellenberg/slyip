package repo

import "database/sql"

type Repositories struct {
	AccountRepo        *AccountRepository
	EcdsaRepo          *EcdsaRepository
	SlyWalletRepo      *SlyWalletRepository
	InvitationCodeRepo *InvitationCodeRepository
	EcdsaSlyWalletRepo *EcdsaSlyWalletRepository
}

func NewRepositories(database *sql.DB) *Repositories {
	db := NewDatabase(database)

	accountRepo := NewAccountRepository(db)
	ecdsaRepo := NewEcdsaRepository(db)
	slyWalletRepo := NewSlyWalletRepository(db)
	invitationCodeRepo := NewInvitationCodeRepository(db)
	ecdsaSlyWalletRepo := NewEcdsaSlyWalletRepository(db)
	return &Repositories{
		AccountRepo:        accountRepo,
		EcdsaRepo:          ecdsaRepo,
		SlyWalletRepo:      slyWalletRepo,
		InvitationCodeRepo: invitationCodeRepo,
		EcdsaSlyWalletRepo: ecdsaSlyWalletRepo,
	}
}
