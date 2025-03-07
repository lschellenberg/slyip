package cryptox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecover(t *testing.T) {

	// Pass 2
	address := "0x69232eE7b5Bc9928f449B674DDA6A63Cd1d8140f"
	messageSignatureHash := "0x74b19777a3966ff19fac1ab22c357ab5e42c3867f5898f841ebeda358dc4309f1f6e243f010397389eff38905276d7f7b2de12460f321371197e64e12bc110051b"
	originalMessage := "a46ff1fb19e4c1be9eb4e0bee38fdad279385d5602884adaf6f376bc3530058428a130bed1a676cdad85b5667b45224d4f710784e0ceb054080de7e06cd2ec1401"

	recoveredAddress, err := Recover(originalMessage, messageSignatureHash, SignMethodGoDefault, false)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, recoveredAddress.String(), address)

	// Pass 1
	address2 := "0x829814B6E4dfeC4b703F2c6fDba28F1724094D11"
	messageSignatureHash2 := "0x53edb561b0c1719e46e1e6bbbd3d82ff798762a66d0282a9adf47a114e32cbc600c248c247ee1f0fb3a6136a05f0b776db4ac82180442d3a80f3d67dde8290811c"
	originalMessage2 := "hello"

	recoveredAddress2, err := Recover(originalMessage2, messageSignatureHash2, SignMethodEthereumPrefix, false)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, recoveredAddress2.String(), address2)

	address3 := "0xF40a876aF2c47FCBF21b937ad59b8a5dc9332CbA"
	messageSignatureHash3 := "0x515e22f1ef482be012651e0880e4360152b8532e1165d48b4cb82993c0dd4ef2206146d1bb35c4d5a3c79439486f497098d65045cc577321570d7dfcb451bf9f1b"
	originalMessage3 := "jqoMShNuDJsNP28Douw2O6JlBmPsv0wyRNSUxaL3"

	recoveredAddress3, err := Recover(originalMessage3, messageSignatureHash3, SignMethodGoDefault, false)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, recoveredAddress3.String(), address3)
}
