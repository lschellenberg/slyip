package pkg

import (
	"fmt"
	info2 "yip/src/api/admin/info"
	"yip/src/api/auth/pin"
	"yip/src/api/auth/session"
	"yip/src/api/auth/verifier"
	"yip/src/api/services/dto"
	"yip/src/providers"
	"yip/src/repositories"
	"yip/src/repositories/repo"
)

type EmptyBody struct{}

type ApiClient struct {
	httpClient HttpClient
	token      string
}

func NewApiClient(baseUrl string) ApiClient {
	return ApiClient{
		httpClient: NewHttpClient(baseUrl),
	}
}

func (c *ApiClient) SetToken(token string) {
	c.token = token
}

func (c *ApiClient) HealthCheck() (statusCode int, message string, err error) {
	statusCode, err = c.httpClient.Get(message, c.token, "healtcheck")
	return
}

/*
SIWE Endpoints
*/
func (c *ApiClient) SIWEChallenge(body dto.ChallengeRequestDTO) (statusCode int, response *dto.ChallengeResponse, err error) {
	response = &dto.ChallengeResponse{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/siwe/challenge")
	return
}

func (c *ApiClient) SIWENonce() (statusCode int, response *dto.NonceResponse, err error) {
	response = &dto.NonceResponse{}
	statusCode, err = c.httpClient.Get(response, c.token, "auth/siwe/nonce")
	return
}

func (c *ApiClient) SIWEVerify(body dto.SubmitRequestDTO) (statusCode int, response *dto.VerifyResponse, err error) {
	response = &dto.VerifyResponse{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/siwe/verify")
	return
}

func (c *ApiClient) SIWESubmit(body dto.SubmitRequestDTO) (statusCode int, response *verifier.Token, err error) {
	response = &verifier.Token{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/siwe/submit")
	return
}

func (c *ApiClient) UserInfo() (statusCode int, response *dto.UserInfoResponse, err error) {
	response = &dto.UserInfoResponse{}
	statusCode, err = c.httpClient.Get(response, c.token, "auth/token/userinfo")
	return
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Admin Endpoints
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *ApiClient) ChainInfo() (statusCode int, response *info2.ChainInfoDTO, err error) {
	response = &info2.ChainInfoDTO{}
	statusCode, err = c.httpClient.Get(response, c.token, "admin/info/chain")
	return
}

func (c *ApiClient) GetCodes() (statusCode int, response *repo.PaginatedResponse[repo.InvitationCodeModel], err error) {
	response = &repo.PaginatedResponse[repo.InvitationCodeModel]{}
	statusCode, err = c.httpClient.Get(response, c.token, "admin/info/codes")
	return
}
func (c *ApiClient) SignIn(body dto.SignInRequest) (statusCode int, response *verifier.Token, err error) {
	response = &verifier.Token{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "admin/accounts/token")
	return
}

func (c *ApiClient) RegisterUser(body dto.RegisterRequest) (statusCode int, response *repo.AccountModel, err error) {
	response = &repo.AccountModel{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "admin/accounts/register")
	return
}

func (c *ApiClient) GetAccount(accountId string) (statusCode int, response *repositories.UserAccount, err error) {
	response = &repositories.UserAccount{}
	statusCode, err = c.httpClient.Get(response, c.token, fmt.Sprintf("admin/accounts/%s", accountId))
	return
}

func (c *ApiClient) GetUsers() (statusCode int, response *repositories.ListUsersAccountResponse, err error) {
	response = &repositories.ListUsersAccountResponse{}
	statusCode, err = c.httpClient.Get(response, c.token, "admin/accounts")
	return
}

func (c *ApiClient) GetUser(accountId string) (statusCode int, response *repositories.UserAccount, err error) {
	response = &repositories.UserAccount{}
	statusCode, err = c.httpClient.Get(response, c.token, fmt.Sprintf("admin/accounts/%s", accountId))
	return
}

func (c *ApiClient) SetEmail(body dto.SetEmailRequest) (statusCode int, response *repositories.UserAccount, err error) {
	response = &repositories.UserAccount{}
	statusCode, err = c.httpClient.Put(body, response, c.token, "admin/accounts/email")
	return
}

func (c *ApiClient) RequestPin(body pin.PinRequestDTO) (statusCode int, response *pin.PinRequestResponse, err error) {
	response = &pin.PinRequestResponse{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/pin")
	return
}

func (c *ApiClient) RedeemPin(body pin.PinRedeemDTO) (statusCode int, response *verifier.Token, err error) {
	response = &verifier.Token{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/pin/redeem")
	return
}

func (c *ApiClient) ListPins() (statusCode int, response []pin.Pin, err error) {
	response = []pin.Pin{}
	statusCode, err = c.httpClient.Get(response, c.token, "admin/pins")
	return
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Auth Endpoints
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *ApiClient) RefreshToken(refreshToken string) (statusCode int, response *verifier.Token, err error) {
	response = &verifier.Token{}
	request := dto.RefreshTokenRequestDTO{RefreshToken: refreshToken}
	statusCode, err = c.httpClient.Post(request, response, c.token, "auth/token/refresh")
	return
}

func (c *ApiClient) VerifyToken(idToken string) (statusCode int, response *verifier.Principal, err error) {
	response = &verifier.Principal{}
	request := dto.VerifyTokenRequestDTO{
		Token: idToken,
	}
	statusCode, err = c.httpClient.Post(request, response, c.token, "auth/token/verify")
	return
}

func (c *ApiClient) Session(body *session.WebsocketMessage) (statusCode int, response *session.WebsocketMessage, err error) {
	response = &session.WebsocketMessage{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/session")
	return
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	SLY Wallets
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *ApiClient) CreateSLYWallet(body *session.WebsocketMessage) (statusCode int, response *session.WebsocketMessage, err error) {
	response = &session.WebsocketMessage{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "auth/session")
	return
}

func (c *ApiClient) SpawnSLYWallet(body *dto.CreateSLYWalletRequest) (statusCode int, response *providers.TransactionTicket, err error) {
	response = &providers.TransactionTicket{}
	statusCode, err = c.httpClient.Post(body, response, c.token, "sly/wallet/spawn")
	return
}

func (c *ApiClient) GetTransactionStatus(hash string) (statusCode int, response *providers.TransactionTicket, err error) {
	response = &providers.TransactionTicket{}
	statusCode, err = c.httpClient.Get(response, c.token, fmt.Sprintf("sly/wallet/receipt/%v", hash))
	return
}
