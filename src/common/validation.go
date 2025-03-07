package common

import (
	"gopkg.in/square/go-jose.v2/json"
	"net/http"
	"yip/src/slyerrors"
)

type Validatible interface {
	Validate() error
}

func ReadAndValidate(a Validatible, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return a.Validate()
}
