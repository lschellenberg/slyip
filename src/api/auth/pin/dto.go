package pin

import (
	"gopkg.in/square/go-jose.v2/json"
	"net/http"
	"time"
	"yip/src/slyerrors"
)

type PinRequestDTO struct {
	Email       string `json:"email"`
	ECDSAPubKey string `json:"ecdsaPubKey"`
}

func (p *PinRequestDTO) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("email", p.Email).
		ValidateNotEmpty("deviceId", p.ECDSAPubKey).
		Error()
}

type PinRedeemDTO struct {
	Pin          string   `json:"pin"`
	PinSignature string   `json:"pinSignature"`
	Audiences    []string `json:"audiences"`
}

func (p *PinRedeemDTO) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("email", p.Pin).
		ValidateNotEmpty("deviceId", p.PinSignature).
		ValidateAtLeastOneElement("audiences", p.Audiences).
		Error()
}

// swagger : model PinRequestResponse
type PinRequestResponse struct {
	AccountId   string    `json:"accountId"`
	Email       string    `json:"email"`
	ECDSAPubKey string    `json:"ecdsaPubKey"`
	Expiration  time.Time `json:"expiration"`
	Pin         string    `json:"pin"`
}
