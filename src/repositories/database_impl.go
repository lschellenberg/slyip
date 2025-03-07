package repositories

import (
	"context"
	"database/sql"
	"fmt"
	common2 "github.com/ethereum/go-ethereum/common"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"yip/.gen/slyip/slyip/model"
	. "yip/.gen/slyip/slyip/table"
	"yip/src/slyerrors"
)

const LogDb = true

type DatabaseImpl struct {
	db *sql.DB
}

func log(msg string) {
	if !LogDb {
		return
	}
	fmt.Println("DB Log: ", msg)
}

func NewDatabase(db *sql.DB) DatabaseImpl {
	return DatabaseImpl{db: db}
}

func (u DatabaseImpl) RegisterAccount(context context.Context, email string, passwordHashed string, role string) (UserAccount, error) {
	log("registering user")
	user := UserAccount{}
	dest := &model.Account{}
	insertStmt := Account.INSERT(Account.Email, Account.PasswordHashed, Account.Role).
		VALUES(email, passwordHashed, role).
		RETURNING(Account.AllColumns)

	err := insertStmt.Query(u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)
	return user, nil
}

func (u DatabaseImpl) RegisterAccountWithEmail(context context.Context, email string) (UserAccount, error) {
	log("registering user with email")
	user := UserAccount{}
	dest := &model.Account{}
	insertStmt := Account.INSERT(Account.Email).
		VALUES(email).
		RETURNING(Account.AllColumns)

	err := insertStmt.Query(u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)
	return user, nil
}

func (u DatabaseImpl) RegisterEmptyAccount(context context.Context) (UserAccount, error) {
	user := UserAccount{}
	dest := &model.Account{}
	insertStmt := Account.INSERT(Account.Role).VALUES("basic").RETURNING(Account.AllColumns)

	err := insertStmt.Query(u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)
	return user, nil
}

func (u DatabaseImpl) GetAccounts(context context.Context, limit int, offset int) (ListUsersAccountResponse, error) {
	response := ListUsersAccountResponse{}
	var dest []model.Account
	log(fmt.Sprintf("Getting users limit: %d, offser: %d", limit, offset))
	stmt := Account.SELECT(Account.AllColumns).LIMIT(int64(limit)).OFFSET(int64(offset))
	err := stmt.Query(u.db, &dest)
	if err != nil {
		return response, err
	}
	users := make([]UserAccount, len(dest))

	for k, v := range dest {
		err = users[k].parseFromDatabase(&v)
		if err != nil {
			return response, err
		}
	}
	response.Users = users
	response.Count = uint64(len(users))
	response.Offset = offset
	response.PageSize = limit

	return response, nil
}

func (u DatabaseImpl) GetAccountByEmail(context context.Context, email string) (UserAccount, error) {
	user := UserAccount{}
	dest := &model.Account{}
	stmt := Account.
		SELECT(Account.AllColumns).
		WHERE(Account.Email.EQ(String(email)))

	err := stmt.Query(u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)
	return user, err
}

func (u DatabaseImpl) GetAccountById(context context.Context, id uuid.UUID) (UserAccount, error) {
	account := UserAccount{}
	dest := &model.Account{}

	stmt := Account.
		SELECT(Account.AllColumns).
		WHERE(Account.ID.EQ(UUID(id)))
	err := stmt.Query(u.db, dest)
	if err != nil {
		return account, err
	}
	err = account.parseFromDatabase(dest)
	if err != nil {
		return account, err
	}

	ecdsa, err := u.GetDevices(context, id)
	if err != nil {
		log(err.Error())
		return account, err
	}

	account.Keys = ecdsa

	return account, err
}

func (u DatabaseImpl) GetSLYWallets(context context.Context, accountId uuid.UUID) ([]SLYWalletWithControllerKeys, error) {
	result := make([]SLYWalletWithControllerKeys, 0)

	var dest []struct {
		model.Ecdsa

		ECDSASLYWallet []struct {
			model.EcdsaSlyWallet

			SlyWallet model.SlyWallet
		}
	}
	stmt := SELECT(Ecdsa.AllColumns, EcdsaSlyWallet.AllColumns, Ecdsa.AllColumns).
		FROM(
			Ecdsa.
				INNER_JOIN(EcdsaSlyWallet, EcdsaSlyWallet.EcdsaAddress.EQ(SlyWallet.Address)).
				INNER_JOIN(SlyWallet, EcdsaSlyWallet.OnChainAccountAddress.EQ(SlyWallet.Address)),
		).
		WHERE(Ecdsa.AccountID.EQ(UUID(accountId)))

	err := stmt.Query(u.db, &dest)
	if err != nil {
		return result, err
	}

	slyWallets := make(map[string]SLYWalletWithControllerKeys, 0)

	for _, v := range dest {
		key := ControllerKey{
			Key: ECDSAKey{
				Address:   v.Address,
				AccountId: accountId.String(),
			},
			Permission: 0,
		}
		for _, w := range v.ECDSASLYWallet {
			m, ok := slyWallets[w.SlyWallet.Address]
			if !ok {
				m = SLYWalletWithControllerKeys{
					Address:        w.SlyWallet.Address,
					ChainId:        w.SlyWallet.Chainid,
					ControllerKeys: make([]ControllerKey, 0),
					CreatedAt:      w.SlyWallet.CreatedAt.Unix(),
					UpdatedAt:      w.SlyWallet.UpdatedAt.Unix(),
				}
			}
			key.Permission = w.OnChainPermissions
			m.ControllerKeys = append(m.ControllerKeys, key)
			slyWallets[w.SlyWallet.Address] = m
		}
	}
	for _, m := range slyWallets {
		result = append(result, m)
	}
	return result, nil
}

func (u DatabaseImpl) GetLastAccount(context context.Context) (UserAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (u DatabaseImpl) SetRole(context context.Context, id uuid.UUID, role string) (UserAccount, error) {
	user := UserAccount{}
	dest := &model.Account{}
	stmt := Account.
		UPDATE(Account.Role).
		SET(Account.Role.SET(String(role))).
		WHERE(Account.ID.EQ(UUID(id))).
		RETURNING(Account.AllColumns)

	err := stmt.Query(u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)

	return user, err
}

func (u DatabaseImpl) SetLastUsedSLYWallet(ctx context.Context, accountId uuid.UUID, address common2.Address) (UserAccount, error) {
	user := UserAccount{}
	dest := &model.Account{}
	stmt := Account.
		UPDATE(Account.Role).
		SET(Account.LastUsedSlyWallet.SET(String(address.Hex()))).
		WHERE(Account.ID.EQ(UUID(accountId))).
		RETURNING(Account.AllColumns)

	err := stmt.QueryContext(ctx, u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)

	return user, err
}

func (u DatabaseImpl) SetEmail(ctx context.Context, id uuid.UUID, email string) (UserAccount, error) {
	log(fmt.Sprintf("Setting email of %s to %s", id, email))
	user := UserAccount{}
	dest := &model.Account{}
	stmt := Account.
		UPDATE(Account.Email, Account.IsEmailVerified).
		SET(
			Account.Email.SET(String(email)),
			Account.IsEmailVerified.SET(Bool(true)),
		).
		WHERE(Account.ID.EQ(UUID(id))).
		RETURNING(Account.AllColumns)

	err := stmt.QueryContext(ctx, u.db, dest)
	if err != nil {
		return user, err
	}

	err = user.parseFromDatabase(dest)

	return user, err
}

func (u DatabaseImpl) GetECDSAByAddress(context context.Context, address string) (ECDSAKey, error) {
	ecdsa := ECDSAKey{}
	dest := &model.Ecdsa{}
	stmt := Ecdsa.SELECT(Ecdsa.AllColumns).
		WHERE(Ecdsa.Address.EQ(String(address)))

	err := stmt.Query(u.db, dest)

	if err != nil {
		return ecdsa, err
	}

	err = ecdsa.parseFromDatabase(dest)
	return ecdsa, err
}

func (u DatabaseImpl) RegisterECDSA(context context.Context, address string) (ECDSAKey, error) {
	ecdsa := ECDSAKey{}
	newUser, err := u.RegisterEmptyAccount(context)
	if err != nil {
		log(fmt.Sprintf("error: %v", err))
		return ecdsa, err
	}
	log(fmt.Sprintf("register account: %v", newUser))
	dest := &model.Ecdsa{}
	insertStmt := Ecdsa.INSERT(Ecdsa.Address, Ecdsa.AccountID).
		VALUES(address, newUser.ID).
		RETURNING(Ecdsa.AllColumns)

	err = insertStmt.Query(u.db, dest)
	if err != nil {
		return ecdsa, err
	}

	err = ecdsa.parseFromDatabase(dest)
	return ecdsa, nil
}

func (u DatabaseImpl) GetOrCreateECDSAKey(context context.Context, address string) (ECDSAKey, error) {
	key, err := u.GetECDSAByAddress(context, address)
	if err == nil {
		// Case 1: received key -> return it
		return key, nil
	} else if slyerrors.IsNoRowsError(err) {
		// Case 2: no key -> create db entry
		log("user is not known yet -> lets create one")
		return u.RegisterECDSA(context, address)
	}
	// Case 3: just a db error
	return key, err
}

func (u DatabaseImpl) AddDevice(context context.Context, accountId uuid.UUID, pubKey string) (ECDSAKey, error) {
	ecdsa := ECDSAKey{
		Address:   pubKey,
		AccountId: accountId.String(),
	}
	dest := &model.Ecdsa{}
	insertStmt := Ecdsa.INSERT(Ecdsa.Address, Ecdsa.AccountID).
		VALUES(pubKey, accountId).
		RETURNING(Ecdsa.AllColumns)

	err := insertStmt.Query(u.db, dest)
	if err != nil {
		return ecdsa, err
	}

	err = ecdsa.parseFromDatabase(dest)
	return ecdsa, nil
}

func (u DatabaseImpl) GetDevices(ctx context.Context, accountId uuid.UUID) ([]ECDSAKey, error) {
	log(fmt.Sprintf("Getting ecdsa for account %s", accountId))
	var dest []model.Ecdsa
	stmt := Ecdsa.
		SELECT(Ecdsa.AllColumns).
		WHERE(Ecdsa.AccountID.EQ(UUID(accountId)))

	err := stmt.Query(u.db, &dest)
	if err != nil {
		if slyerrors.IsNoRowsError(err) {
			return make([]ECDSAKey, 0), nil
		}
		return nil, err
	}

	keys := make([]ECDSAKey, len(dest))
	for k, v := range dest {
		err = keys[k].parseFromDatabase(&v)
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}
