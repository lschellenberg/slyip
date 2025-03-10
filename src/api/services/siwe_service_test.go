package services

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"testing"
	"yip/src/api/services/dto"
	"yip/src/cryptox"
)

var SepoliaWalletPrivateKey = "a5ec116d46f3ec04c7337be7af353a3c20eb567ba764304023eef108fd2f5aa2"
var userWallet = GetWallet()

func TestSIWE(t *testing.T) {
	s := NewSIWEService(nil, nil, nil, nil)
	domain := "http://localhost:3000"
	c, err := s.Challenge(&dto.ChallengeRequestDTO{
		ChainId: "111155551111",
		Address: cryptox.PublicKeyFromKey(userWallet),
		Domain:  domain,
	})

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("SIWE Message:\n", c)

	sm, err := cryptox.Sign(c.Challenge, userWallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		t.Error(err)
		return
	}

	verifyResponse, err := s.Verify(&dto.SubmitRequestDTO{
		Message:   c.Challenge,
		Signature: sm.Signature,
		Audience:  "http://localhost:8081",
	})

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, userWallet.Address.String(), verifyResponse.RecoveredAddress)
	assert.Equal(t, userWallet.Address.String(), verifyResponse.OriginalAddress)

}

func GetWallet() *keystore.Key {
	p, _ := cryptox.ReadPasswordFile("./samples/test_wallets/user_password.txt")
	key, _ := cryptox.KeyFromWalletAndPasswordFile("./samples/test_wallets/user.json", p)
	return key
}
