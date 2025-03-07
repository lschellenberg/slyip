package dto

import (
	"gopkg.in/square/go-jose.v2/json"
	"net/http"
	"yip/src/repositories"
	"yip/src/slyerrors"
)

// swagger:model VerifiedUser
type VerifiedUser struct {
	ID        string   `json:"id"`
	Role      string   `json:"role"`
	Audiences []string `json:"audiences"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refreshToken"`
}

func (r RefreshTokenRequestDTO) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("refreshToken", r.RefreshToken).
		Error()
}

type ROPFRequest struct {
	Email     string   `json:"email"`    // the desired email
	Password  string   `json:"password"` // the desired password
	Audiences []string `json:"audiences"`
}

func (a *ROPFRequest) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return a.Validate()
}

func (r ROPFRequest) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("email", r.Email).
		ValidateNotEmpty("password", r.Password).
		ValidateAtLeastOneElement("audiences", r.Audiences).
		Error()
}

type VerifyTokenRequestDTO struct {
	Token string `json:"token"`
}

func (r *VerifyTokenRequestDTO) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("token", r.Token).
		Error()
}

type UserInfoResponse struct {
	Account    repositories.UserAccount                   `json:"account"`
	SLYWallets []repositories.SLYWalletWithControllerKeys `json:"SLYWallets"`
}

type SwapRequestDTO struct {
	Address string `json:"address"`
}

func (a *SwapRequestDTO) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return a.Validate()
}

func (r *SwapRequestDTO) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateEthAddress("address", r.Address).
		Error()
}
