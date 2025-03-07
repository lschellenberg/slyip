package session

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
	"yip/src/slyerrors"
)

const (
	CommunicationTypeHTTP1     = "http1"
	CommunicationTypeWebsocket = "websocket"
)

type SessionClient struct {
	ID                string
	session           *Session
	connector         *MConnector
	conn              *websocket.Conn
	mutex             *sync.Mutex
	send              chan []byte
	sendJSON          chan *WebsocketMessage
	CommunicationType string
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func newSessionClient(clientID string, connector *MConnector, conn *websocket.Conn) *SessionClient {
	sc := &SessionClient{
		ID:        clientID,
		connector: connector,
		mutex:     &sync.Mutex{},
	}

	if conn == nil {
		sc.CommunicationType = CommunicationTypeHTTP1
	} else {
		sc.send = make(chan []byte, 256)
		sc.sendJSON = make(chan *WebsocketMessage, 256)
		sc.CommunicationType = CommunicationTypeWebsocket
		return sc
	}
	return sc

}

func newClient(connector *MConnector, conn *websocket.Conn) *SessionClient {

	if conn == nil {
		return &SessionClient{
			connector:         connector,
			mutex:             &sync.Mutex{},
			CommunicationType: CommunicationTypeHTTP1,
		}
	}

	return &SessionClient{
		connector:         connector,
		conn:              conn,
		mutex:             &sync.Mutex{},
		send:              make(chan []byte, 256),
		sendJSON:          make(chan *WebsocketMessage, 256),
		CommunicationType: CommunicationTypeWebsocket,
	}
}

func (c *SessionClient) sendMessage(wm *WebsocketMessage) {
	c.sendJSON <- wm
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *SessionClient) readPump() {
	defer func() { c.closeConnection() }()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		wm, err := c.parseMessage(m)
		if err != nil {
			log.Println("not a proper message", string(m))
			c.sendMessage(createErrorResponse("not a valid mconnector msg", err.Error(), slyerrors.ErrCodeBadSessionRequest, c.getSessionId()))
			continue
		}
		if wm.IsSessionRequest() {
			err = c.handleSessionRequest(wm)
			if err != nil {
				log.Println("not a proper message", string(m))
				c.sendMessage(createErrorResponse("no session with this session id", err.Error(), slyerrors.ErrCodeSessionDifferentMessageTypeExpected, wm.SessionId))
				continue
			}
		}
	}
}

func (c *SessionClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.closeConnection()
	}()
	for {
		select {
		case message, ok := <-c.sendJSON:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println(err)
				}
				return
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Println(err)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *SessionClient) handleSessionRequest(wm *WebsocketMessage) error {
	if c.session != nil {
		return fmt.Errorf("session already established")
	}

	if wm.IsSessionIdEmpty() {
		c.session = newSessionAsync(c.connector, SessionTypeAuth)
		return nil
	}

	uu, err := uuid.Parse(wm.SessionId)
	if err != nil {
		return err
	}

	s, err := c.connector.getSession(uu)
	if err != nil {
		return err
	}

	c.session = s
	return nil
}

func (c *SessionClient) getSessionId() string {
	if c.session != nil {
		return c.session.SessionId.String()
	}
	return ""
}

func (c *SessionClient) closeConnection() {
	c.connector.unregister <- c
	if c.session != nil {
		c.session.unregister(c)
	}
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (c *SessionClient) parseMessage(rawMessage []byte) (*WebsocketMessage, error) {
	msg := &WebsocketMessage{}
	err := json.Unmarshal(rawMessage, msg)
	if err != nil {
		return nil, err
	}
	return msg, err
}
