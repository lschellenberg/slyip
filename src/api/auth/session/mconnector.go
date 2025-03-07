package session

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type MConnector struct {
	upgrader          websocket.Upgrader
	sessions          map[uuid.UUID]*Session
	clients           map[*SessionClient]bool
	register          chan *SessionClient
	unregister        chan *SessionClient
	registerSession   chan *Session
	unregisterSession chan *Session
	lock              *sync.Mutex
}

func InitMConnector() MConnector {
	return MConnector{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		sessions:          make(map[uuid.UUID]*Session),
		register:          make(chan *SessionClient),
		unregister:        make(chan *SessionClient),
		registerSession:   make(chan *Session),
		unregisterSession: make(chan *Session),
		lock:              &sync.Mutex{},
	}
}

func (h *MConnector) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.sendJSON)
			}
		case session := <-h.registerSession:
			h.registerNewSession(session)
		case session := <-h.unregisterSession:
			h.lock.Lock()
			if _, ok := h.sessions[session.SessionId]; ok {
				delete(h.sessions, session.SessionId)
				session.close()
			}
			h.lock.Unlock()
		}
	}
}

func (mc *MConnector) registerNewSession(s *Session) {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.sessions[s.SessionId] = s
}

func (mc *MConnector) getSession(sessionId uuid.UUID) (*Session, error) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	s, ok := mc.sessions[sessionId]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return s, nil
}

func (mc *MConnector) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := mc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(mc, conn)
	mc.register <- client

	go client.writePump()
	go client.readPump()
}
