package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"yip/e2e/samples"
	"yip/src/api/auth/session"
	"yip/src/api/services/dto"
	"yip/src/cryptox"
	"yip/src/utils"
)

func TestSession(t *testing.T) {
	wallet, err := cryptox.WalletFromPrivateKey(samples.SepoliaWalletPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	// CREATE Session
	status, scr, err := client.Session(session.CreateSessionMessage(clientId, session.SessionTypeAuth))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status, "status code should be ok")
	fmt.Println(scr.Payload)

	// PING
	p := ping(t, scr.SessionId)
	if p == nil {
		return
	}
	assert.Equal(t, "pending", p.AuthState)
	assert.Nil(t, p.Token)

	// CONNECT to Session

	status, challenge, err := client.Session(session.CreateAccountResponse(scr.SessionId, wallet.Address.String(), samples.SLYWalletAddress, chainId))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status, "status code should be ok")
	assert.Equal(t, "eth_sign", challenge.MessageType)

	r := &dto.ChallengeResponse{}
	err = utils.MapToStruct(challenge.Payload, r)
	if err != nil {
		t.Error(err)
		return
	}

	// PING
	p = ping(t, scr.SessionId)
	if p == nil {
		return
	}
	assert.Equal(t, "pending", p.AuthState)
	assert.Nil(t, p.Token)

	// SEND Signature
	sm, err := cryptox.Sign(r.Challenge, wallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		t.Error(err)
		return
	}
	status, _, err = client.Session(session.CreateSubmitSignature(scr.SessionId, sm.Message, sm.Signature))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status, "status code should be ok")

	// PING
	p = ping(t, scr.SessionId)
	if p == nil {
		return
	}
	assert.Equal(t, "success", p.AuthState)
	assert.NotNil(t, p.Token)

	// CLOSE Session
	status, _, err = client.Session(session.CreateCloseSessionRequest(scr.SessionId))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status, "status code should be ok")

	// pings should then result to slyerrors
	status, wm, err := client.Session(session.CreatePingRequest(scr.SessionId))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, status, "status code should be ok")
	assert.Equal(t, "session_error", wm.MessageType)
	e, err := wm.ParseSessionError()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "600001", e.Code)
	assert.Equal(t, "session closed", e.Message)
}

func ping(t *testing.T, sessionId string) *session.PayloadPingTokenResponse {
	status, wm, err := client.Session(session.CreatePingRequest(sessionId))
	if err != nil {
		t.Error(err)
		return nil
	}
	assert.Equal(t, 200, status, "status code should be ok")
	pingResult, err := wm.ParsePingResponse()
	if err != nil {
		t.Error(err)
		return nil
	}

	return pingResult
}

func printJson(o interface{}) {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println("Error marshalling o", o)
	}

	fmt.Println(string(b))
}
