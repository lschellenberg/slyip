package dto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yip/src/slyerrors"
)

const RoleOwner = 1
const RoleAdmin = 2
const RoleAuthenticator = 3

type SLYBase struct {
	ID             string          `json:"id"`
	Address        string          `json:"address"`
	ImageUrl       string          `json:"imageUrl"`
	Description    string          `json:"description"`
	Name           string          `json:"username"`
	State          int64           `json:"state"`
	Category       string          `json:"category"`
	AccountId      string          `json:"accountId"`
	Owners         []ControllerKey `json:"owners"`
	Admins         []ControllerKey `json:"admins"`
	Authenticators []ControllerKey `json:"authenticators"`
}

type SpawnSLYWalletRequest struct {
	InvitationCode string `json:"invitationCode"`
	OwnerPublicKey string `json:"ownerPublicKey"`
}

func (p *SpawnSLYWalletRequest) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("invitationCode", p.InvitationCode).
		ValidateEthAddress("ownerPublicKey", p.OwnerPublicKey).
		Error()
}

type ControllerKey struct {
	Address   string `json:"address"`
	AccountId string `json:"accountId"`
	Role      int    `json:"role"`
}

type SyncState struct {
	LastSyncedBlock int64 `json:"lastSyncedBlock"`
}

func (sly *SLYBase) AddOwner(address string) (*SLYBase, error) {
	sly.Owners = append(sly.Owners, ControllerKey{Address: address, Role: RoleOwner})
	return sly, nil
}

func (sly *SLYBase) AddAdmin(address string) (*SLYBase, error) {
	sly.Admins = append(sly.Admins, ControllerKey{Address: address, Role: RoleAdmin})
	return sly, nil
}

func (sly *SLYBase) AddAuthenticator(address string) (*SLYBase, error) {
	sly.Authenticators = append(sly.Authenticators, ControllerKey{Address: address, Role: RoleAuthenticator})
	return sly, nil
}

func (sly *SLYBase) AddController(key ControllerKey) *SLYBase {
	switch key.Role {
	case RoleOwner:
		sly.Owners = append(sly.Owners, key)
	case RoleAdmin:
		sly.Admins = append(sly.Admins, key)
	case RoleAuthenticator:
		sly.Authenticators = append(sly.Authenticators, key)
	}
	return sly
}

func (sly *SLYBase) RemoveController(address string) (*SLYBase, error) {
	index := findIndex(address, sly.Owners)
	if index >= 0 {
		sly.Owners = remove(sly.Owners, index)
		return sly, nil
	}
	index = findIndex(address, sly.Admins)
	if index >= 0 {
		sly.Admins = remove(sly.Admins, index)
		return sly, nil
	}
	index = findIndex(address, sly.Authenticators)
	if index == -1 {
		return nil, fmt.Errorf("address %s not found [not authenticator]", address)
	}
	sly.Authenticators = remove(sly.Authenticators, index)
	return sly, nil
}

func findIndex(key string, s []ControllerKey) int {
	index := -1
	for k, a := range s {
		if key == a.Address {
			index = k
			break
		}
	}
	return index
}

func remove(s []ControllerKey, i int) []ControllerKey {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
