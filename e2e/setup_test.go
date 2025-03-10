package e2e

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"yip/pkg"
	"yip/src/api/auth/verifier"
	"yip/src/config"
	"yip/src/cryptox"
)

var localServerURL = "http://localhost:8080/api/v1"
var url = localServerURL
var client = pkg.NewApiClient(url)
var userWallet = GetWallet()
var chainId = "31337"

var clientId = "d798c41c-1afe-465e-9143-75c5f111a1cb"

func GetVerifier() verifier.Verifier {
	return verifier.NewVerifier(config.JWTTokenConfig{
		TokenExpirationInSec:        1000,
		RefreshTokenExpirationInSec: 10000,
		CertificatePrivate:          "./samples/test_certs/app.rsa",
		CertificatePublic:           "./samples/test_certs/app.rsa.pub",
		Issuer:                      "https://yip.net",
	})
}

func GetWallet() *keystore.Key {
	p, _ := cryptox.ReadPasswordFile("./samples/test_wallets/user_password.txt")
	key, _ := cryptox.KeyFromWalletAndPasswordFile("./samples/test_wallets/user.json", p)
	return key
}
