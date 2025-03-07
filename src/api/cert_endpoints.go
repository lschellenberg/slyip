package api

import (
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"yip/src/httpx"
)

// swagger:route GET /.well-known/openid-configuration Auth well-known-configuration
// returns oidc configuration of this service
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: ProviderJSON
func (a Api) WellKnownConfiguration(w http.ResponseWriter, r *http.Request) {
	httpx.RespondWithJSON(w, httpx.OK(ProviderJSON{
		Issuer:      a.App.Config.JWT.Issuer,
		AuthURL:     "",
		TokenURL:    "",
		JWKSURL:     fmt.Sprintf("%s/.well-known/jwks", a.App.Config.JWT.Issuer),
		UserInfoURL: "",
		Algorithms:  []string{"HS256"},
	}))
}

// swagger:route GET /.well-known/jwks Auth well-known-jswks
// returns oidc public keys of token signer
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: JWKSResponse
func (a Api) WellKnownJWKS(w http.ResponseWriter, r *http.Request) {
	httpx.RespondWithJSON(w, httpx.OK(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		a.App.Verifier.JSONWebKey(),
	}}))
}

// swagger:response JWKSResponse
type JWKSResponse jose.JSONWebKey

// swagger:response ProviderJSON
type ProviderJSON struct {
	Issuer      string   `json:"issuer"`
	AuthURL     string   `json:"authorization_endpoint"`
	TokenURL    string   `json:"token_endpoint"`
	JWKSURL     string   `json:"jwks_uri"`
	UserInfoURL string   `json:"userinfo_endpoint"`
	Algorithms  []string `json:"id_token_signing_alg_values_supported"`
}
