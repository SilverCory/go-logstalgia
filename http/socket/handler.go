package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Handler struct {
	clients     map[*websocket.Conn]bool
	clientsLock *sync.Mutex
}

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"log"},
}

func New() *Handler {
	handler := &Handler{
		clients:     make(map[*websocket.Conn]bool),
		clientsLock: &sync.Mutex{},
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			handler.clientsLock.Lock()
			for k, v := range handler.clients {
				if v {
					if err := k.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						delete(handler.clients, k)
					}
				}
			}
			handler.clientsLock.Unlock()
		}
	}()

	return handler
}

func (h *Handler) BroadcastJSON(i interface{}) {
	h.clientsLock.Lock()
	for k, v := range h.clients {
		if v {
			if err := k.WriteJSON(i); err != nil {
				log.Println("Error in broadcast! ", err)
				delete(h.clients, k)
			}
		}
	}
	h.clientsLock.Unlock()
}

func (h *Handler) UpgradeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("An error occurred on upgrade: ", err)
		return
	}

	conn.SetReadLimit(512)
	conn.EnableWriteCompression(true)
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	conn.SetPingHandler(func(data string) error {
		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return err
		}
		return conn.WriteMessage(websocket.PongMessage, []byte(data))
	})

	h.clientsLock.Lock()
	h.clients[conn] = true
	h.clientsLock.Unlock()
}
