package session

import (
	"encoding/json"
	"fmt"
	"yip/src/api/auth/verifier"
	"yip/src/api/services/dto"
	"yip/src/slyerrors"
	"yip/src/utils"
)

const (
	MessageTypeCreateSessionRequest   = "create_session"
	MessageTypeSessionCreatedResponse = "session_created"
	MessageTypeSessionError           = "session_error"
	MessageTypeSignatureRequest       = "eth_sign"
	MessageTypeSubmitSignature        = "eth_sign_response"
	MessageTypeVerificationResponse   = "eth_sign_verification_response"
	MessageTypeAccountsRequest        = "eth_accounts"
	MessageTypeConnectWithAccount     = "connect_with_account"
	MessageTypePingToken              = "ping_token"
	MessageTypePingTokenResponse      = "ping_token_response"
	MessageTypeCloseSession           = "session_close"
	MessageTypeCloseSessionResponse   = "session_close_response"
)

type WebsocketMessage struct {
	MessageType string      `json:"messageType"`
	SessionId   string      `json:"sessionId"`
	Payload     interface{} `json:"payload"`
}

func (wm *WebsocketMessage) response() *WebsocketMessage {
	mt := ""

	switch wm.MessageType {
	case MessageTypeConnectWithAccount:
		mt = MessageTypeSignatureRequest
	}
	return &WebsocketMessage{
		MessageType: mt,
		SessionId:   wm.SessionId,
		Payload:     nil,
	}
}
func (wm *WebsocketMessage) IsSessionRequest() bool {
	return wm.MessageType == MessageTypeCreateSessionRequest
}

func (wm *WebsocketMessage) IsSessionIdEmpty() bool {
	return wm.SessionId == ""
}

type PayloadCreateSessionRequest struct {
	ClientId    string `json:"clientId"`
	SessionType string `json:"sessionType"`
}

func CreateSessionMessage(clientId string, sessionType string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeCreateSessionRequest,
		SessionId:   "",
		Payload: PayloadCreateSessionRequest{
			ClientId:    clientId,
			SessionType: sessionType,
		},
	}
}

func (wm *WebsocketMessage) ParseCreateSessionRequest() (*PayloadCreateSessionRequest, error) {
	payload := PayloadCreateSessionRequest{}

	b, err := json.Marshal(wm.Payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}

	err = json.Unmarshal(b, &payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}

	err = slyerrors.NewValidation("400").
		ValidateNotEmpty("clientId", payload.ClientId).
		ValidateInList("sessionType", payload.SessionType, sessionTypes).
		Error()

	if err != nil {
		return nil, err
	}

	return &payload, nil
}

type PayloadSessionCreatedResponse struct {
	SessionId     string `json:"sessionId"`
	SessionType   string `json:"sessionType"`
	ClientId      string `json:"clientId"`
	QRCodeContent string `json:"qrCodeContent"`
}

type PayloadSignatureRequest struct {
	EOA     string `json:"eoa"`
	Message string `json:"message"`
}

type PayloadSignatureResponse = dto.SubmitRequestDTO

func CreateSubmitSignature(sessionId string, message string, signature string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeSubmitSignature,
		SessionId:   sessionId,
		Payload: PayloadSignatureResponse{
			Message:   message,
			Signature: signature,
		},
	}
}

func (wm *WebsocketMessage) ParseSubmitRequest() (*PayloadSignatureResponse, error) {
	payload := PayloadSignatureResponse{}

	err := utils.MapToStruct(wm.Payload, &payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}

	err = slyerrors.NewValidation("400").
		ValidateNotEmpty("message", payload.Message).
		ValidateNotEmpty("signature", payload.Signature).
		Error()

	if err != nil {
		return nil, err
	}

	return &payload, nil
}

type PayloadVerificationResponse = dto.VerifyResponse

func CreateVerificationResponse(sessionId string, response *dto.VerifyResponse) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeVerificationResponse,
		SessionId:   sessionId,
		Payload:     response,
	}
}

type PayloadAccountsRequest struct {
	OnlyPrimary bool `json:"onlyPrimary"`
}

type PayloadAccountsResponse struct {
	EOA              string `json:"eoa"`
	SLYWalletAddress string `json:"slyWalletAddress"`
	ChainID          string `json:"chainId"`
}

func CreateAccountResponse(sessionId string, eoa string, slyWalletAddress string, chainId string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeConnectWithAccount,
		SessionId:   sessionId,
		Payload: PayloadAccountsResponse{
			EOA:              eoa,
			SLYWalletAddress: slyWalletAddress,
			ChainID:          chainId,
		},
	}
}

func (wm *WebsocketMessage) ParseAccountsResponse() (*PayloadAccountsResponse, error) {
	payload := PayloadAccountsResponse{}

	b, err := json.Marshal(wm.Payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}

	err = slyerrors.NewValidation("400").
		ValidateEthAddress("eoa", payload.EOA).
		ValidateNotEmpty("chainId", payload.ChainID).
		Error()

	if err != nil {
		return nil, err
	}

	return &payload, nil
}

type PayloadSessionError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Details string `json:"details"`
}

func (wm *WebsocketMessage) ParseSessionError() (*PayloadSessionError, error) {
	payload := PayloadSessionError{}

	err := utils.MapToStruct(wm.Payload, &payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}

	if err != nil {
		return nil, err
	}

	return &payload, nil
}

type PayloadPingTokenRequest struct {
}

type PayloadPingTokenResponse struct {
	AuthState string          `json:"authState"`
	Token     *verifier.Token `json:"token"`
}

func CreatePingRequest(sessionId string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypePingToken,
		SessionId:   sessionId,
		Payload:     PayloadPingTokenRequest{},
	}
}

func CreatePingResponse(sessionId string, state string, token *verifier.Token) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypePingTokenResponse,
		SessionId:   sessionId,
		Payload: PayloadPingTokenResponse{
			AuthState: state,
			Token:     token,
		},
	}
}

type PayloadSessionClosed struct {
}

func CreateCloseResponse(sessionId string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeCloseSessionResponse,
		SessionId:   sessionId,
		Payload:     PayloadSessionClosed{},
	}
}

type PayloadSessionClose struct {
}

func CreateCloseSessionRequest(sessionId string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeCloseSession,
		SessionId:   sessionId,
		Payload:     PayloadSessionClose{},
	}
}

func (wm *WebsocketMessage) ParsePingResponse() (*PayloadPingTokenResponse, error) {
	payload := &PayloadPingTokenResponse{}

	err := utils.MapToStruct(wm.Payload, payload)
	if err != nil {
		return nil, fmt.Errorf("payload is not according to message type")
	}

	return payload, nil
}

func getQRCodeContent(uri string, sessionId string, clientId string, sessionType string, chainId string) string {
	return fmt.Sprintf("%s/%s?sid=%s&cid=%s&flow=%s&chainId=%s", uri, "api/v1/auth/session", sessionId, clientId, sessionType, chainId)
}
