package app

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"log"
	migration "yip/internal/goose"
	"yip/src/api/auth/verifier"
	"yip/src/config"
	"yip/src/contracts"
	"yip/src/providers"
	"yip/src/repositories"
)

type App struct {
	Config           *config.Config
	Verifier         *verifier.Verifier
	DB               *sql.DB
	UserDB           repositories.Database
	EmailProvider    providers.EmailProvider
	EthProvider      *providers.EthProvider
	SLYWalletManager *contracts.WalletManager
}

func InitApp(c *config.Config) (*App, error) {
	// open db
	db, err := OpenSqlDB(&c.DB)
	if err != nil {
		return nil, err
	}

	// create verifier
	v := verifier.NewVerifier(c.JWT)

	ep, err := providers.InitEmailProvider(&c.Email)
	if err != nil {
		return nil, err
	}

	ethProvider, err := providers.InitEthProvider(&c.EthConfig)
	if err != nil {
		return nil, err
	}

	wm, err := contracts.NewWalletManager(common.HexToAddress(c.EthConfig.Sly.FactoryAddress), &ethProvider)
	if err != nil {
		return nil, err
	}
	return &App{
		c,
		&v,
		db,
		repositories.NewDatabase(db),
		ep,
		&ethProvider,
		wm,
	}, nil
}

func OpenSqlDB(info *config.SqlDBInfo) (*sql.DB, error) {
	// open db
	db, err := sql.Open("postgres", info.String())
	if err != nil {
		return nil, err
	}
	// check if connection is on
	if _, err = db.Exec("SELECT 1"); err != nil {
		return nil, err
	}
	// migrate
	log.Println("Migrating database (if necessary)...")
	if err = migration.MigrateDB(db); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	return db, nil
}
