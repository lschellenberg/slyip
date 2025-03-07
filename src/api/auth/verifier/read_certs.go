package verifier

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

type Certs struct {
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
}

func NewCerts(privKey string, privKeyPath string, pubKey string, pubKeyPath string) (*Certs, error) {
	var err error
	var signBytes []byte
	var verifyBytes []byte

	// check if private key is in env
	if privKey == "" {
		// check if path to private key is in path
		if privKeyPath == "" {
			return nil, fmt.Errorf("no path for private key of certs given")
		}
		// ok: private key path is given
		signBytes, err = ioutil.ReadFile(privKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// ok: private key is given
		signBytes = []byte(privKey)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	// check if public key is in env
	if pubKey == "" {
		// check if path to public key is in path
		if pubKeyPath == "" {
			return nil, fmt.Errorf("no path for public key of certs given")
		}
		// ok: public key path is given
		verifyBytes, err = ioutil.ReadFile(pubKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// ok: public key is given
		verifyBytes = []byte(pubKey)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return &Certs{
		VerifyKey: verifyKey,
		SignKey:   signKey,
	}, nil
}
