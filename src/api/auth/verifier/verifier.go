package verifier

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"gopkg.in/square/go-jose.v2"
	"time"
	"yip/src/config"
	"yip/src/slyerrors"
)

const (
	BearerTokenType = "bearer"
)

// swagger:model Token
type Token struct {
	IdToken      string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	Type         string `json:"type"`
}

type Claims struct {
	Scopes []string `json:"scopes"`
	Aud    []string `json:"aud"`
	Role   string   `json:"role"`
	ECDSA  string   `json:"ecdsa"`
	SLY    string   `json:"sly"`
	jwt.StandardClaims
}

type Verifier struct {
	config config.JWTTokenConfig
	certs  Certs
}

func NewVerifier(c config.JWTTokenConfig) Verifier {
	certs, err := NewCerts("", c.CertificatePrivate, "", c.CertificatePublic)
	if err != nil {
		panic(err)
	}
	return Verifier{
		config: c,
		certs:  *certs,
	}
}

// VerifyToken ensures that the token is signed with the SecretKey then returns a Principal based on the token content
func (a Verifier) VerifyToken(ctx context.Context, tokenString string) (*Principal, error) {
	sc := &Claims{}
	_, err := a.parseClaimsToken(tokenString, sc)

	if err != nil {
		if vErr, ok := err.(*jwt.ValidationError); ok {
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return nil, slyerrors.Unauthorized(slyerrors.ErrCodeTokenExpired, vErr.Error())
			case jwt.ValidationErrorMalformed:
				return nil, slyerrors.BadRequest(slyerrors.ErrCodeMalformedToken, "token malformed")
			}
		}

		return nil, slyerrors.Unauthorized(slyerrors.ErrCodeUnknownTokenVerificationError, err.Error())
	}

	return &Principal{
		ID:               sc.Subject,
		SLYWalletAddress: sc.SLY,
		ECDSAAddress:     sc.ECDSA,
		Scopes:           sc.Scopes,
		Role:             sc.Role,
		Audiences:        sc.Aud,
	}, nil
}

// NewClaims generates a oidc claims struct with JWT *StandardClaims included.
func (a Verifier) NewClaims(audience []string, accountId string, ecdsaAddress string, slyWalletAddress string, role string, expirationTimeInSec int64) *Claims {
	expirationTime := time.Now().Add(time.Duration(expirationTimeInSec) * time.Second)

	return &Claims{
		ECDSA:  ecdsaAddress,
		SLY:    slyWalletAddress,
		Scopes: []string{},
		Role:   role,
		Aud:    audience,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    a.config.Issuer,
			NotBefore: 0,
			Subject:   accountId,
		},
	}
}

func (a Verifier) parseClaimsToken(tokenString string, sc *Claims) (*jwt.Token, error) {
	if a.certs.VerifyKey == nil {
		return nil, slyerrors.Unexpected("cannot validate token", "secret key is empty", nil)
	}

	return jwt.ParseWithClaims(tokenString, sc, func(token *jwt.Token) (interface{}, error) {
		return a.certs.VerifyKey, nil
	})
}

// SignClaimsToken creates a signed token from claims
func (a Verifier) SignClaimsToken(c *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	signedToken, err := token.SignedString(a.certs.SignKey)
	if err != nil {
		return signedToken, slyerrors.Unexpected("could not sign token", err.Error(), err)
	}
	return signedToken, nil
}

// CreateToken returns a signed JWT token for the userId.
// The token is anonymous and has no scope, so can be used only for endpoints usable by anonymous users
func (a Verifier) CreateToken(audience []string, accountId string, ecdsaAddress string, slyWalletAddress string, role string) (*Token, error) {
	token, err := a.SignClaimsToken(a.NewClaims(audience, accountId, ecdsaAddress, slyWalletAddress, role, a.config.TokenExpirationInSec))

	if err != nil {
		return nil, slyerrors.Unexpected("could not create token", "SignatureHex creation failed", err)
	}

	refreshToken, err := a.SignClaimsToken(a.NewClaims(audience, accountId, ecdsaAddress, slyWalletAddress, role, a.config.RefreshTokenExpirationInSec))

	if err != nil {
		return nil, slyerrors.Unexpected("could not create refresh token", "SignatureHex creation failed", err)
	}

	return &Token{
		IdToken:      token,
		RefreshToken: refreshToken,
		ExpiresIn:    a.config.TokenExpirationInSec,
		Type:         BearerTokenType,
	}, nil
}

// RefreshToken creates a new access token and refresh token pair for the user in the current refresh token
func (a Verifier) RefreshToken(refreshToken string) (*Token, error) {
	claims := &Claims{}

	_, err := a.parseClaimsToken(refreshToken, claims)

	if err != nil {
		if vErr, ok := err.(*jwt.ValidationError); ok {
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return nil, slyerrors.Unauthorized("expired token", vErr.Error(), vErr)
			case jwt.ValidationErrorMalformed:
				return nil, slyerrors.BadRequest("could not parse the token", vErr.Error(), vErr)
			}
		}

		return nil, slyerrors.Unauthorized("could not validate token", err.Error(), err)
	}

	token, err := a.CreateToken(claims.Aud, claims.Subject, claims.ECDSA, claims.SLY, claims.Role)
	if err != nil {
		return &Token{}, slyerrors.Unexpected("could not update token", "Refresh token creation failed", err)
	}

	return token, nil
}

func (a Verifier) JSONWebKey() jose.JSONWebKey {
	key, err := jwk.New(a.certs.VerifyKey)

	if err != nil {
		panic(err)
	}

	return jose.JSONWebKey{
		Key:                         a.certs.VerifyKey,
		KeyID:                       key.KeyID(),
		Algorithm:                   key.Algorithm(),
		Use:                         key.KeyUsage(),
		Certificates:                nil,
		CertificatesURL:             nil,
		CertificateThumbprintSHA1:   nil,
		CertificateThumbprintSHA256: nil,
	}
}
