package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"yip/e2e/samples"
	"yip/src/api/auth/pin"
	"yip/src/cryptox"
)

func TestPinRequest(t *testing.T) {
	_, response, err := client.SignIn(samples.AdminUser)
	if err != nil {
		t.Error(err)
		return
	}
	client.SetToken(response.IdToken)

	email := "adamshulman@gmail.com" //gofakeit.Email()
	wallet, err := cryptox.GenerateNewKey()
	if err != nil {
		t.Error(err)
		return
	}

	status, pinResponse, err := client.RequestPin(pin.PinRequestDTO{
		Email:       email,
		ECDSAPubKey: wallet.Address.String(),
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pinResponse)
	assert.Equal(t, 200, status, "status code should be ok")
	assert.Equal(t, email, pinResponse.Email)
	assert.Equal(t, wallet.Address.String(), pinResponse.ECDSAPubKey)
	assert.Equal(t, 6, len(pinResponse.Pin))

	status1, response1, err := client.GetAccount(pinResponse.AccountId)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status1, "status code should be ok")
	assert.Equal(t, 0, len(response1.Keys), "status code should be ok")
	assert.Equal(t, false, response1.IsEmailVerified)

	signature, err := cryptox.Sign(pinResponse.Pin, wallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)

	status2, response2, err := client.RedeemPin(pin.PinRedeemDTO{
		Pin:          pinResponse.Pin,
		PinSignature: signature.Signature,
		Audiences:    samples.TestAudiences,
	})
	assert.Equal(t, 200, status2, "status code should be ok")
	assert.NotEmpty(t, response2.IdToken)

	// check if ecdsa key was added
	status3, response3, err := client.GetAccount(pinResponse.AccountId)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status3, "status code should be ok")
	assert.Equal(t, 1, len(response3.Keys), "status code should be ok")
	assert.Equal(t, true, response3.IsEmailVerified)
}
