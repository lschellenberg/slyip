package dto

import (
	"encoding/json"
	"net/http"
	"yip/src/slyerrors"
)

// swagger:model SignInRequest
type SignInRequest struct {
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Audiences []string `json:"audiences"`
}

func (a *SignInRequest) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return a.Validate()
}

func (a *SignInRequest) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("email", a.Email).
		ValidateNotEmpty("password", a.Password).
		ValidateAtLeastOneElement("audiences", a.Audiences).
		Error()
}

// swagger:model RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email"`    // the desired email
	Password string `json:"password"` // the desired password
}

// swagger:model SetRoleRequest
type SetRoleRequest struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

// swagger:model SetRoleRequest
type SetEmailRequest struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func (a *SetEmailRequest) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return slyerrors.NewValidation("400").
		ValidateNotEmpty("email", a.UserID).
		ValidateNotEmpty("password", a.Email).
		Error()
}
