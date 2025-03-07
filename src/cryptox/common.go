package cryptox

type SignMethod int64
type SignType int64

const (

	//https://ethereum.stackexchange.com/questions/64667/signed-message-verification-failing-for-signed-messages-with-multiple-passes
	SignMethodEthereumPrefix SignMethod = 1 // can be verified with https://app.mycrypto.com/verify-message    PASS 1
	SignMethodGoDefault      SignMethod = 2 // PASS 2
	SignTypeWeb3JS           SignType   = 1
	SignTypeGo               SignType   = 2
)

type SignedMessage struct {
	Address   string     `json:"address"`
	Version   string     `json:"version"`
	Message   string     `json:"msg"`
	Signature string     `json:"sig"`
	Method    SignMethod `json:"-"`
	Type      SignType   `json:"-"`
}

func (sm SignedMessage) Validate() (bool, error) {
	address, err := Recover(sm.Message, sm.Signature, sm.Method, false)
	if err != nil {
		return false, err
	}

	return address.String() == sm.Address, nil
}
