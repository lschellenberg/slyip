package e2e

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"yip/e2e/samples"
	"yip/src/api/auth/verifier"
	"yip/src/api/services/dto"
	"yip/src/cryptox"
	"yip/src/repositories/repo"
)

func TestCreateSLYWallet(t *testing.T) {
	codes := createInvitationCodes()
	if len(codes.Data) == 0 {
		t.Errorf("info codes are empty")
		return
	}

	wallet, err := cryptox.WalletFromPrivateKey(samples.SepoliaWalletPrivateKey)
	require.NoError(t, err)

	token, err := signInWithAddress(cryptox.PublicKeyFromKey(wallet), wallet)
	require.NoError(t, err)

	client.SetToken(token.IdToken)
	_, user, err := client.UserInfo()
	require.NoError(t, err)
	assert.Empty(t, user.SLYWallets)

	reqBody := &dto.CreateSLYWalletRequest{
		InvitationCode: codes.Data[0].Code,
	}
	code, receipt, err := client.SpawnSLYWallet(reqBody)
	require.NoError(t, err)

	fmt.Println("SpawnSLYWallet:Code", code)
	fmt.Println("SpawnSLYWallet:Receipt", receipt)

	time.Sleep(1 * time.Second)
	code, status, err := client.GetTransactionStatus(receipt.TransactionHash)
	require.NoError(t, err)

	fmt.Println("SpawnSLYWallet:Code", code)
	fmt.Println("SpawnSLYWallet:Statuss", status)

}

//func TestTryAgainToCreateSLYWallet(t *testing.T) {
//	fmt.Println("TestCreateSLYWallet")
//	wallet, err := cryptox.WalletFromPrivateKey(samples.SepoliaWalletPrivateKey)
//	require.NoError(t, err)
//
//	token, err := signInWithAddress(cryptox.PublicKeyFromKey(wallet), wallet)
//	require.NoError(t, err)
//
//	client.SetToken(token.IdToken)
//
//	code, _, err := client.SpawnSLYWallet()
//	assert.Equal(t, code, 400)
//	require.NoError(t, err)
//}

// Helpers

func signInWithAddress(address string, wallet *keystore.Key) (*verifier.Token, error) {
	domain := "http://localhost:3000"
	_, response, err := client.SIWEChallenge(dto.ChallengeRequestDTO{
		ChainId: chainId,
		Address: address,
		Domain:  domain,
	})

	if err != nil {
		return nil, err
	}

	sm, err := cryptox.Sign(response.Challenge, wallet, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
	if err != nil {
		return nil, err
	}

	_, token, err := client.SIWESubmit(dto.SubmitRequestDTO{
		Message:   response.Challenge,
		Signature: sm.Signature,
		Audience:  "http://localhost:8081",
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func createInvitationCodes() *repo.PaginatedResponse[repo.InvitationCodeModel] {
	_, response, err := client.SignIn(samples.AdminUser)
	if err != nil {
		panic(err)
	}
	client.SetToken(response.IdToken)
	_, codes, err := client.GetCodes()
	if err != nil {
		panic(err)
	}
	return codes
}
