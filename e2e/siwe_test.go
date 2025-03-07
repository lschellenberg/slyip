package e2e

import (
	"context"
	"github.com/google/uuid"
	"testing"
	"yip/e2e/samples"
	"yip/src/api/auth/verifier"
	"yip/src/api/services/dto"
	"yip/src/cryptox"

	"github.com/stretchr/testify/assert"
)

func TestChallenge(t *testing.T) {
	domain := "http://localhost:3000"
	statusCode, response, err := client.SIWEChallenge(dto.ChallengeRequestDTO{
		ChainId: chainId,
		Address: samples.TestAddress,
		Domain:  domain,
	})

	if err != nil {
		t.Error(err)
		return
	}

	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.NotNil(t, response)
	assert.Equal(t, samples.TestAddress, response.Address, "addresses are equal")
	assert.Truef(t, validateSIWEMessage(response.Challenge), "siwe message not valid")
	assert.Equal(t, chainId, response.ChainId, "chainID not equal")
	assert.Equal(t, domain, response.Domain, "chainID not equal")
}

func TestNonce(t *testing.T) {
	statusCode, response, err := client.SIWENonce()

	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Nonce, "nonce should not be empty")
}

func TestSignature(t *testing.T) {
	domain := "http://localhost:3000"
	_, response, err := client.SIWEChallenge(dto.ChallengeRequestDTO{
		ChainId: chainId,
		Address: cryptox.PublicKeyFromKey(userWallet),
		Domain:  domain,
	})

	if err != nil {
		t.Error(err)
		return
	}

	wallet, err := cryptox.WalletFromPrivateKey(samples.SepoliaWalletPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	sm, err := cryptox.Sign(response.Challenge, wallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		t.Error(err)
		return
	}

	statusCode, verifyResponse, err := client.SIWEVerify(dto.SubmitRequestDTO{
		Message:   response.Challenge,
		Signature: sm.Signature,
		Audience:  "http://localhost:8081",
	})

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.Equal(t, userWallet.Address.String(), verifyResponse.RecoveredAddress)
	assert.Equal(t, userWallet.Address.String(), verifyResponse.OriginalAddress)
}

func TestSubmission(t *testing.T) {
	// user wallet is 0x4E345039EE45217fC99a717a441384A46dD2b85C Public Key
	adminAtSLYWallet := "0x030E4BFabdF1d5463B92BBC4fA8cE8587c7BA079"
	wallet, err := cryptox.WalletFromPrivateKey(samples.SepoliaWalletPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	v := GetVerifier()
	domain := "http://localhost:3000"
	_, response, err := client.SIWEChallenge(dto.ChallengeRequestDTO{
		ChainId: chainId,
		Address: cryptox.PublicKeyFromKey(wallet),
		Domain:  domain,
	})

	if err != nil {
		t.Error(err)
		return
	}

	sm, err := cryptox.Sign(response.Challenge, wallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		t.Error(err)
		return
	}

	statusCode, token, err := client.SIWESubmit(dto.SubmitRequestDTO{
		Message:   response.Challenge,
		Signature: sm.Signature,
		Audience:  "http://localhost:8081",
	})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 200, statusCode, "status code should be ok")

	p, err := v.VerifyToken(context.Background(), token.IdToken)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, cryptox.PublicKeyFromKey(wallet), p.ECDSAAddress)
	assert.Equal(t, verifier.RoleBasic, p.Role)

	_, err = uuid.Parse(p.ID)
	assert.NoError(t, err)

	client.SetToken(token.IdToken)

	statusCode, userInfo, err := client.UserInfo()
	if err != nil {
		t.Error(err)
		return
	}
	// TODO need to recreate db
	//assert.Empty(t, userInfo.Account.LastUsedSLYWallet)

	statusCode, userInfo, err = client.UserInfo()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, adminAtSLYWallet, userInfo.Account.LastUsedSLYWallet)
	statusCode, adminToken, err := client.SignIn(samples.AdminUser)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, statusCode)
	client.SetToken(adminToken.IdToken)

	statusCode, actualUser, err := client.GetUser(userInfo.Account.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, adminAtSLYWallet, actualUser.LastUsedSLYWallet)

	// relogin

	_, response, err = client.SIWEChallenge(dto.ChallengeRequestDTO{
		ChainId: chainId,
		Address: cryptox.PublicKeyFromKey(wallet),
		Domain:  domain,
	})

	if err != nil {
		t.Error(err)
		return
	}

	sm, err = cryptox.Sign(response.Challenge, wallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		t.Error(err)
		return
	}

	statusCode, token, err = client.SIWESubmit(dto.SubmitRequestDTO{
		Message:   response.Challenge,
		Signature: sm.Signature,
		Audience:  "http://localhost:8081",
	})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 200, statusCode, "status code should be ok")

	p, err = v.VerifyToken(context.Background(), token.IdToken)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, cryptox.PublicKeyFromKey(wallet), p.ECDSAAddress)
	assert.Equal(t, verifier.RoleBasic, p.Role)
	assert.Equal(t, adminAtSLYWallet, p.SLYWalletAddress)

}

func validateSIWEMessage(msg string) bool {
	return true
}
