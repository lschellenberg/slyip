package dto

import (
	"gopkg.in/square/go-jose.v2/json"
	"net/http"
	"yip/src/slyerrors"
)

type ChallengeRequestDTO struct {
	ChainId string `json:"chainId"`
	Address string `json:"address"`
	Domain  string `json:"domain"`
}

func (a *ChallengeRequestDTO) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return a.Validate()
}

func (a *ChallengeRequestDTO) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("chainId", a.ChainId).
		ValidateNotEmpty("address", a.Address).
		ValidateNotEmpty("domain", a.Domain).
		Error()
}

// swagger:model ChallengeResponse
type ChallengeResponse struct {
	Challenge string `json:"challenge"`
	Address   string `json:"address"`
	Domain    string `json:"domain"`
	ChainId   string `json:"chainId"`
}

// swagger:model NonceResponse
type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type SubmitRequestDTO struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	Audience  string `json:"audience,omitempty"`
}

func (a *SubmitRequestDTO) ReadAndValidate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(a)

	if err != nil {
		return slyerrors.NewValidation("400").Add("json is not readable", slyerrors.ValidationCodeCannotValidate, err.Error()).Error()
	}

	return a.Validate()
}

func (a *SubmitRequestDTO) Validate() error {
	return slyerrors.NewValidation("400").
		ValidateNotEmpty("message", a.Message).
		ValidateNotEmpty("signature", a.Signature).
		ValidateNotEmpty("audience", a.Audience).
		Error()
}

// swagger:model VerifyResponse
type VerifyResponse struct {
	Domain           string `json:"domain"`
	URI              string `json:"uri"`
	OriginalAddress  string `json:"originalAddress"`
	RecoveredAddress string `json:"recoveredAddress"`
}
