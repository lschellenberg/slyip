package cryptox

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"strings"
)

// Recover
// Given a message msg and a signature sig based on the method it will return
//
//		signature string
//		pubKey    string
//		err 	  error
//
//		Methods:
//			SignMethodEthereumPrefix: adds "\x19Ethereum Signed Message\n" +len(msg)  to the message
//			SignMethodGoDefault: no prefix
//		Types:
//	  	SignTypeWeb3JS: sig[65] = 27|28
//	  	SignTypeGo: 	sig[65] = 0|1
func Recover(msg string, sigHex string, method SignMethod, web3JSStrict bool) (*common.Address, error) {
	preparedSignature := sigHex
	if strings.Contains(sigHex, "0x") {
		preparedSignature = sigHex[2:]
	}

	sig := common.Hex2Bytes(preparedSignature)
	if len(sig) != 65 {
		return nil, fmt.Errorf("wrong signature length")
	}
	if sig[64] != 27 && sig[64] != 28 {
		// must be 0 or 1 then which is GO standard
		if web3JSStrict {
			return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
		}
	} else {
		sig[64] -= 27
	}

	switch method {
	case SignMethodGoDefault:
		return RecoverGoDefault(msg, sig)
	case SignMethodEthereumPrefix:
		return RecoverPrefix(msg, sig)
	default:
		return nil, fmt.Errorf("could not recover: signing method not found: %d", method)
	}
}

func RecoverGoDefault(msg string, sig []byte) (*common.Address, error) {
	msgBytes := []byte(msg) // its hex but must be treated as a string
	// recover signature
	hash := crypto.Keccak256(msgBytes)
	pubKeyBytes, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	// transform pub key
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*pubKey)

	return &address, nil
}

func RecoverPrefix(msg string, sig []byte) (*common.Address, error) {
	prefix := "\x19Ethereum Signed Message:\n"
	preSigned := prefix + strconv.Itoa(len(msg)) + msg
	//preSigned := prefix + msg
	hash := crypto.Keccak256Hash([]byte(preSigned))

	// recover signature
	pubKeyBytes, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		return nil, err
	}

	// transform pub key
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*pubKey)

	return &address, nil
}
