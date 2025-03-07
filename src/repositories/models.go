package repositories

import (
	"yip/.gen/slyip/slyip/model"
	"yip/src/api/auth/verifier"
	"yip/src/common"
)

type UserAccount struct {
	ID                string                        `json:"id"`
	Email             string                        `json:"email"`
	Phone             string                        `json:"phone"`
	Role              string                        `json:"role"`
	FirstName         string                        `json:"firstName"`
	LastName          string                        `json:"lastName"`
	PasswordHashed    string                        `json:"-"`
	CreatedAt         int64                         `json:"createdAt"`
	UpdatedAt         int64                         `json:"updatedAt"`
	Keys              []ECDSAKey                    `json:"ecdsaKeys"`
	LastUsedSLYWallet string                        `json:"lastUsedSLYWallet"`
	OnChainAccounts   []SLYWalletWithControllerKeys `json:"onChainAccounts"`
	IsEmailVerified   bool                          `json:"isEmailVerified"`
	IsPhoneVerified   bool                          `json:"isPhoneVerified"`
}

func (u *UserAccount) IsAdmin() bool {
	return u.Role == verifier.RoleAdmin
}

// swagger:model ListUsersAccountResponse
type ListUsersAccountResponse struct {
	Users []UserAccount `json:"accounts"`
	common.PaginationResponse
}

func (u *UserAccount) parseFromDatabase(a *model.Account) error {
	if a.Email != nil {
		u.Email = *a.Email
	} else {
		u.Email = ""
	}

	u.ID = a.ID.String()
	u.Role = a.Role
	u.PasswordHashed = a.PasswordHashed
	u.CreatedAt = a.CreatedAt.Unix()
	u.UpdatedAt = a.UpdatedAt.Unix()
	u.FirstName = a.FirstName
	u.LastName = a.LastName
	u.Phone = a.Phone
	u.IsEmailVerified = a.IsEmailVerified
	u.IsPhoneVerified = a.IsPhoneVerified
	u.LastUsedSLYWallet = a.LastUsedSlyWallet
	return nil
}

type ECDSAKey struct {
	Address   string `json:"address"`
	AccountId string `json:"accountId"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (u *ECDSAKey) parseFromDatabase(a *model.Ecdsa) error {
	u.Address = a.Address
	u.AccountId = a.AccountID.String()
	u.CreatedAt = a.CreatedAt.Unix()
	u.UpdatedAt = a.UpdatedAt.Unix()
	return nil
}

type SLYWalletWithControllerKeys struct {
	Address        string          `json:"address"`
	ChainId        string          `json:"chainId"`
	ControllerKeys []ControllerKey `json:"controllerKeys"`
	CreatedAt      int64           `json:"createdAt"`
	UpdatedAt      int64           `json:"updatedAt"`
}

func (u *SLYWalletWithControllerKeys) parseFromDatabase(a *model.SlyWallet) error {
	u.Address = a.Address
	u.ChainId = a.Chainid
	u.CreatedAt = a.CreatedAt.Unix()
	u.UpdatedAt = a.UpdatedAt.Unix()
	return nil
}

func (u *UserAccount) parseFullFromDatabase(a *model.Account, ecdsa []model.Ecdsa, onChain []model.SlyWallet) error {
	err := u.parseFromDatabase(a)
	if err != nil {
		return err
	}
	u.Keys = make([]ECDSAKey, len(ecdsa))
	for k, e := range ecdsa {
		err = u.Keys[k].parseFromDatabase(&e)
		if err != nil {
			return err
		}
	}

	u.OnChainAccounts = make([]SLYWalletWithControllerKeys, len(onChain))
	for k, e := range onChain {
		err = u.OnChainAccounts[k].parseFromDatabase(&e)
		if err != nil {
			return err
		}
	}
	return nil
}

type ControllerKey struct {
	Key        ECDSAKey `json:"key"`
	Permission int32    `json:"permission"`
}
