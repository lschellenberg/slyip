package pin

import (
	"fmt"
	"testing"
	"time"
	"yip/src/cryptox"
)

const (
	test_email = "some@email.com"
)

func TestPinWrongSignature(t *testing.T) {
	pool := NewPool(time.Minute)

	key, _ := cryptox.GenerateNewKey()
	pin := pool.request("", test_email, key.Address.String())
	signature, err := cryptox.Sign(pin.Pin, key, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		t.Error(err)
		return
	}
	calcPin, err := pool.redeem(pin.Pin, signature.Signature)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(calcPin)
}
