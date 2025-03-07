package session

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/api/services/dto"
	"yip/src/config"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

type Controller struct {
	siweService *services.SIWEService
	userService *services.UserService
	MConnector  MConnector
	config      *config.Config
}

func NewController(
	c *config.Config,
	service *services.SIWEService,
	userService *services.UserService,
) Controller {
	return Controller{
		siweService: service,
		userService: userService,
		config:      c,
		MConnector:  InitMConnector(),
	}
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", c.SessionChannel)
	}
}

func (a Controller) SessionChannel(w http.ResponseWriter, r *http.Request) {
	msg := &WebsocketMessage{}
	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		r := createErrorResponse(slyerrors.ErrCodeBadSessionRequest, err.Error(), "", msg.SessionId)
		httpx.RespondWithJSON(w, createJSONErrorResponse(200, r))
		return
	}

	switch msg.MessageType {
	case MessageTypeCreateSessionRequest:
		httpx.RespondWithJSON(w, a.CreateSession(msg))
	case MessageTypeConnectWithAccount:
		httpx.RespondWithJSON(w, a.SetAccount(msg))
	case MessageTypeSubmitSignature:
		httpx.RespondWithJSON(w, a.SubmitSIWE(r.Context(), msg))
	case MessageTypePingToken:
		httpx.RespondWithJSON(w, a.PingResult(msg))
	case MessageTypeCloseSession:
		httpx.RespondWithJSON(w, a.CloseSession(msg))
	default:
		httpx.RespondWithJSON(w, jsonErrorResponse(200, slyerrors.ErrCodeSessionMessageTypeUnknown, "unknown message type", fmt.Sprintf("type %s is not known", msg.MessageType), msg.SessionId))
	}
}

func (a Controller) CreateSession(wm *WebsocketMessage) *httpx.Response {
	payload, err := wm.ParseCreateSessionRequest()
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeBadSessionRequest, err.Error(), "", "")
	}

	var cl *config.Client
	for _, c := range a.config.Clients {
		if c.ID == payload.ClientId {
			cl = &c
		}
	}

	if cl == nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeUnknownClient, "client id does not exist", "", "")
	}

	audiences := a.config.AudiencesByClient(cl.ID)

	if len(audiences) == 0 {
		return jsonErrorResponse(200, slyerrors.ErrCodeAudienceDoesntExist, "no audiences found for client", "", "")
	}

	clients := map[*SessionClient]*SessionClient{}
	clients[newSessionClient(payload.ClientId, &a.MConnector, nil)] = nil

	var s *Session
	if payload.SessionType == SessionTypeAuth {
		s = newAuthSession(&a.MConnector)
		s.clients = clients
		s.AuthFlow.domain = cl.Domain
		s.AuthFlow.audiences = audiences
	} else {
		return jsonErrorResponse(200, slyerrors.ErrCodeBadSessionRequest, "session type does not exist", "", "")
	}

	return httpx.OK(WebsocketMessage{
		MessageType: MessageTypeSessionCreatedResponse,
		SessionId:   s.SessionId.String(),
		Payload: PayloadSessionCreatedResponse{
			SessionId:     s.SessionId.String(),
			SessionType:   s.SessionType,
			ClientId:      payload.ClientId,
			QRCodeContent: getQRCodeContent(a.config.JWT.Issuer, s.SessionId.String(), payload.ClientId, payload.SessionType, a.config.EthConfig.Chain.ID),
		},
	})
}

func (a Controller) SetAccount(wm *WebsocketMessage) *httpx.Response {
	session, response := a.MConnector.getSessionFromMessageAndVerifyStatus(wm)
	if response != nil {
		return response
	}

	if session.AuthFlow == nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionWrongSessionType, "flow not initiated by party", "", session.SessionId.String())
	}

	if !session.AuthFlow.isCreatedState() {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionDifferentMessageTypeExpected, "not expecting this message", "", session.SessionId.String())
	}

	payload, err := wm.ParseAccountsResponse()
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeBadSessionRequest, err.Error(), "", session.SessionId.String())
	}

	if a.config.EthConfig.Chain.ID != payload.ChainID {
		return jsonErrorResponse(200, slyerrors.ErrCodeWrongChainId, "wrong chain id", "", session.SessionId.String())
	}

	session.AuthFlow.setPayload(payload)

	r, err := a.siweService.Challenge(&dto.ChallengeRequestDTO{
		Address: session.AuthFlow.eoa,
		ChainId: payload.ChainID,
		Domain:  session.AuthFlow.domain,
	})
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionCantCreateSIWEMessage, err.Error(), "", session.SessionId.String())
	}

	wmResponse := wm.response()
	wmResponse.Payload = r

	return httpx.OK(wmResponse)
}

func (a Controller) SubmitSIWE(ctx context.Context, wm *WebsocketMessage) *httpx.Response {
	session, response := a.MConnector.getSessionFromMessageAndVerifyStatus(wm)
	if response != nil {
		return response
	}

	if session.AuthFlow == nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionWrongSessionType, "flow not initiated by party", "", session.SessionId.String())
	}

	if !session.AuthFlow.isConnectedState() {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionDifferentMessageTypeExpected, "not expecting this message", "", session.SessionId.String())
	}

	payload, err := wm.ParseSubmitRequest()
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeBadSessionRequest, err.Error(), "", session.SessionId.String())
	}

	verificationResult, err := a.siweService.Verify(&dto.SubmitRequestDTO{
		Message:   payload.Message,
		Signature: payload.Signature,
		Audience:  payload.Audience,
	})

	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeWrongSignature, err.Error(), "", session.SessionId.String())
	}

	if session.AuthFlow.slyWalletAddress != "" {
		ma := common.HexToAddress(session.AuthFlow.slyWalletAddress)
		eoa := common.HexToAddress(session.AuthFlow.eoa)

		result := a.siweService.Authenticate(eoa, ma)

		if !result.IsAuthenticated {
			return jsonErrorResponse(result.StatusCode, result.ErrorCode, result.ErrorMessage, result.ErrorDetails, session.SessionId.String())
		}
	}

	ecdsa, err := a.siweService.GetOrCreateAccount(ctx, session.AuthFlow.eoa)
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeCantCreateOrGetAccount, "error getting or creating account", err.Error(), session.SessionId.String())
	}

	uu, err := uuid.Parse(ecdsa.AccountId)
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeParsingUUID, "error parsing uuid of account", err.Error(), session.SessionId.String())
	}

	account, err := a.userService.GetAccountById(ctx, uu)
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeCantCreateOrGetAccount, "error parsing retrieving account", err.Error(), session.SessionId.String())
	}

	session.AuthFlow.setVerified(account.ID)

	return httpx.OK(CreateVerificationResponse(session.SessionId.String(), verificationResult))
}

func (a Controller) PingResult(wm *WebsocketMessage) *httpx.Response {
	session, response := a.MConnector.getSessionFromMessageAndVerifyStatus(wm)
	if response != nil {
		return response
	}

	if session.AuthFlow == nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionWrongSessionType, "flow not initiated by party", "", session.SessionId.String())
	}

	if !session.AuthFlow.isVerified() {
		return httpx.OK(CreatePingResponse(session.SessionId.String(), FlowStatePending, nil))
	}

	info := session.AuthFlow
	token, err := a.siweService.CreateToken(info.audiences, info.accountId, info.eoa, info.slyWalletAddress, verifier.RoleBasic)
	if err != nil {
		return jsonErrorResponse(200, slyerrors.ErrCodeCantCreateToken, "cant create token", err.Error(), session.SessionId.String())
	}

	return httpx.OK(CreatePingResponse(session.SessionId.String(), FlowStateSuccess, token))
}

func (a Controller) CloseSession(wm *WebsocketMessage) *httpx.Response {
	session, response := a.MConnector.getSessionFromMessageAndVerifyStatus(wm)
	if response != nil {
		return response
	}

	session.close()

	return httpx.OK(CreateCloseResponse(session.SessionId.String()))
}
