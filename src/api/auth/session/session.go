package session

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

const (
	SessionTypeAuth = "auth_session"
)

var sessionTypes = []string{SessionTypeAuth}

type Session struct {
	connector   *MConnector
	SessionId   uuid.UUID
	clients     map[*SessionClient]*SessionClient
	mutex       *sync.Mutex
	isClosed    bool
	SessionType string
	AuthFlow    *AuthFlow
}

func newAuthSession(connector *MConnector) *Session {
	s := &Session{
		connector:   connector,
		SessionId:   uuid.New(),
		mutex:       &sync.Mutex{},
		SessionType: SessionTypeAuth,
		AuthFlow:    NewAuthFlow(),
	}

	connector.registerNewSession(s)

	return s
}

func newSessionAsync(connector *MConnector, sessionType string) *Session {
	s := &Session{
		connector:   connector,
		SessionId:   uuid.New(),
		mutex:       &sync.Mutex{},
		SessionType: sessionType,
	}

	connector.registerSession <- s

	return s
}

func (s *Session) passMessage(client *SessionClient, msg *WebsocketMessage) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	otherClient, ok := s.clients[client]
	if !ok {
		return fmt.Errorf("other client not connected yet")
	}

	otherClient.sendMessage(msg)

	return nil
}

func (s *Session) isConnected() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.clients) == 2
}

func (s *Session) unregister(client *SessionClient) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, c := range s.clients {
		if c != client {
			s.clients[c] = nil
		}
	}

	delete(s.clients, client)
}

func (s *Session) register(client *SessionClient) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// session can only have 2 participants
	if len(s.clients) == 2 {
		return fmt.Errorf("the session is full")
	}

	// session already has other party
	if len(s.clients) == 1 {
		for _, c := range s.clients {
			if c == client {
				return fmt.Errorf("cannot connect the same client twice")
			}

			s.clients[c] = client
			s.clients[client] = c
		}
	}

	s.clients[client] = nil

	return nil
}

func (s *Session) CreateSessionCreatedResponse(sessionId string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeSessionCreatedResponse,
		SessionId:   s.SessionId.String(),
		Payload: &PayloadSessionCreatedResponse{
			SessionId: sessionId,
		},
	}
}

func (s *Session) CreateSignatureRequest(eoa string, message string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeSignatureRequest,
		SessionId:   s.SessionId.String(),
		Payload: &PayloadSignatureRequest{
			EOA:     eoa,
			Message: message,
		},
	}
}

func (s *Session) CreateAccountsRequest(onlyPrimary bool) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeAccountsRequest,
		SessionId:   s.SessionId.String(),
		Payload:     PayloadAccountsRequest{onlyPrimary},
	}
}

func (s *Session) close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.isClosed = true
}

func jsonErrorResponse(status int, code string, msg string, details string, sessionId string) *httpx.Response {
	wm := &WebsocketMessage{
		MessageType: MessageTypeSessionError,
		SessionId:   sessionId,
		Payload:     &PayloadSessionError{code, msg, details},
	}
	return &httpx.Response{
		Payload:    wm,
		StatusCode: status,
	}
}

func createErrorResponse(code string, msg string, details string, sessionId string) *WebsocketMessage {
	return &WebsocketMessage{
		MessageType: MessageTypeSessionError,
		SessionId:   sessionId,
		Payload:     &PayloadSessionError{code, msg, details},
	}
}

func createJSONErrorResponse(status int, wm *WebsocketMessage) *httpx.Response {
	return &httpx.Response{
		Payload:    wm,
		StatusCode: status,
	}
}

func (s *Session) verifyNotConnectedYet() *httpx.Response {
	if s.isConnected() {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionAlreadyConnected, "session already has connection", "", s.SessionId.String())
	}
	return nil
}
func (s *Session) verifyNotClosedYet() *httpx.Response {
	if s.isClosed {
		return jsonErrorResponse(200, slyerrors.ErrCodeSessionClosed, "session closed", "", s.SessionId.String())
	}
	return nil
}

func (wm *WebsocketMessage) getSessionId() (*uuid.UUID, *httpx.Response) {
	if wm.SessionId == "" {
		return nil, jsonErrorResponse(200, slyerrors.ErrCodeSessionWrongSessionId, "no session id", "", wm.SessionId)
	}
	uu, err := uuid.Parse(wm.SessionId)
	if err != nil {
		return nil, jsonErrorResponse(200, slyerrors.ErrCodeSessionWrongSessionId, err.Error(), "", wm.SessionId)
	}
	return &uu, nil
}

func (mc *MConnector) getSessionFromMessageAndVerifyStatus(wm *WebsocketMessage) (*Session, *httpx.Response) {
	id, errorResponse := wm.getSessionId()
	if errorResponse != nil {
		return nil, errorResponse
	}

	session, err := mc.getSession(*id)
	if err != nil {
		return nil, jsonErrorResponse(200, slyerrors.ErrCodeSessionNotFound, err.Error(), "", "")
	}

	errorResponse = session.verifyNotConnectedYet()
	if errorResponse != nil {
		return nil, errorResponse
	}

	errorResponse = session.verifyNotClosedYet()
	if errorResponse != nil {
		return nil, errorResponse
	}

	return session, nil
}
