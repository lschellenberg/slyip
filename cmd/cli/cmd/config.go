package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"yip/pkg"
	"yip/src/api/auth/verifier"
	"yip/src/cryptox"
)

const ConfigLocation = "config.json"

const Local = "http://localhost:8081"
const Prod = "https://auth.singularry.xyz"

type CLIContext struct {
	apiClient      *pkg.ApiClient
	wallet         *cryptox.Wallet
	testUserToken  *verifier.Token
	adminUserToken *verifier.Token
}

type Config struct {
	Email         string `json:"email"`
	YIPUrl        string `json:"yipUrl"`
	Password      string `json:"password"`
	SingularryAPI string `json:"singularryApi"`
	TestWallet    string `json:"testWallet"`
	ChainId       int64  `json:"chainId"`
}

func PrepareApiClient() (*CLIContext, error) {
	c, err := readConfigFile(ConfigLocation)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(c)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", Prod, "api/v1")
	apiClient := pkg.NewApiClient(url)
	apiClient.SetToken(token.IdToken)

	wallet := &cryptox.Wallet{}
	if err = wallet.FromPrivateKey(c.TestWallet); err != nil {
		return nil, err
	}

	testToken, err := GetTokenFromEOA(c, wallet)
	if err != nil {
		return nil, err
	}

	cc := CLIContext{
		apiClient:      &apiClient,
		wallet:         wallet,
		testUserToken:  testToken,
		adminUserToken: token,
	}
	return &cc, nil
}

func readConfigFile(path string) (Config, error) {
	var c Config

	content, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(content, &c)
	return c, err
}

func GetToken(c Config) (*verifier.Token, error) {
	httpClient := pkg.NewHttpClient(c.YIPUrl)
	token := &verifier.Token{}
	_, err := httpClient.Post(
		SignInRequest{
			Email:     c.Email,
			Password:  c.Password,
			Audiences: []string{c.SingularryAPI}},
		token,
		"",
		"api/v1/admin/accounts/token")
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetTokenFromEOA(c Config, wallet *cryptox.Wallet) (*verifier.Token, error) {
	httpClient := pkg.NewHttpClient(c.YIPUrl)
	domain := "http://localhost:3000"

	response := &ChallengeResponse{}
	_, err := httpClient.Post(ChallengeRequestDTO{
		ChainId: fmt.Sprintf("%d", c.ChainId),
		Address: wallet.Public.Hex(),
		Domain:  domain,
	}, response, "", "api/v1/auth/siwe/challenge")
	if err != nil {
		return nil, err
	}

	keyStore := cryptox.NewKeyFromECDSA(wallet.Private)

	sm, err := cryptox.Sign(response.Challenge, keyStore, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)

	if err != nil {
		return nil, err
	}

	token := &verifier.Token{}
	_, err = httpClient.Post(SubmitRequestDTO{
		Signature: sm.Signature,
		Message:   sm.Message,
		Audience:  c.SingularryAPI,
	}, token, "", "api/v1/auth/siwe/submit")
	return token, nil
}

type SignInRequest struct {
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Audiences []string `json:"audiences"`
}

type ChallengeRequestDTO struct {
	ChainId string `json:"chainId"`
	Address string `json:"address"`
	Domain  string `json:"domain"`
}

type SubmitRequestDTO struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	Audience  string `json:"audience"`
}

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
	Address   string `json:"address"`
	Domain    string `json:"domain"`
	ChainId   string `json:"chainId"`
}
