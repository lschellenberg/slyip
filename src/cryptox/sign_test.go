package cryptox

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignWOPrefixMessage(t *testing.T) {
	// prepare
	message := "some message"
	key, err := getTestWallet()
	if err != nil {
		t.Error(err)
		return
	}

	// sign
	signedMessage, err := Sign(message, key, SignMethodGoDefault, SignTypeGo)
	if err != nil {
		t.Error(err)
		return
	}

	// validate
	isValid, err := signedMessage.Validate()
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, isValid, "could not recover")

	// recover
	address, err := Recover(message, signedMessage.Signature, SignMethodGoDefault, false)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, address.String(), signedMessage.Address)

	result, _ := json.MarshalIndent(signedMessage, "", "   ")
	fmt.Println(string(result))
}

func TestSignPrefixMessage(t *testing.T) {
	// prepare
	message := "some message"
	key, err := getTestWallet()
	if err != nil {
		t.Error(err)
		return
	}

	// sign
	signedMessage, err := Sign(message, key, SignMethodEthereumPrefix, SignTypeGo)
	if err != nil {
		t.Error(err)
		return
	}

	// validate
	isValid, err := signedMessage.Validate()
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, isValid, "could not recover")

	// recover
	address, err := Recover(message, signedMessage.Signature, SignMethodEthereumPrefix, false)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, address.String(), signedMessage.Address)

	result, _ := json.MarshalIndent(signedMessage, "", "   ")
	fmt.Println(string(result))
}

func getTestWallet() (*keystore.Key, error) {
	return WalletFromPrivateKey("1c5cf04945e924df10d6956119b0cfb872b691c2a5480c4b2c4133300ac60da1")
}
