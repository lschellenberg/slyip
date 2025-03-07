package cryptox

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

// Sign
// Given a message msg and a key based on the method it will return
// signature []byte
// pubKey    string
// err 	  error
// Methods:
//
//	SignMethodEthereumPrefix: adds "\x19Ethereum Signed Message\n" +len(msg)  to the message
//	SignMethodGoDefault:nothing is prefixed
func Sign(msg string, key *keystore.Key, method SignMethod, signType SignType) (*SignedMessage, error) {
	var signature []byte
	var err error

	switch method {
	case SignMethodGoDefault:
		signature, err = SignGoDefault(msg, key)
	case SignMethodEthereumPrefix:
		signature, err = SignEthereumPrefixed(msg, key)
	default:
		return nil, fmt.Errorf("signing method not found")
	}

	if err != nil {
		return nil, err
	}

	if signType == SignTypeWeb3JS {
		signature[64] += 27
	}

	return &SignedMessage{
		Address:   key.Address.String(),
		Version:   "2.0",
		Message:   msg,
		Signature: "0x" + common.Bytes2Hex(signature),
		Method:    method,
	}, nil
}

func SignGoDefault(msg string, key *keystore.Key) ([]byte, error) {
	hash := crypto.Keccak256([]byte(msg))
	return crypto.Sign(hash, key.PrivateKey)
}

func SignEthereumPrefixed(msg string, key *keystore.Key) ([]byte, error) {
	prefix := "\x19Ethereum Signed Message:\n"
	preSigned := prefix + strconv.Itoa(len(msg)) + msg
	hash := crypto.Keccak256Hash([]byte(preSigned))

	return crypto.Sign(hash.Bytes(), key.PrivateKey)
}
