package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const WS_TAG = "WebSocket"

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fail(w, WS_TAG, fmt.Errorf("failed to upgrade connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	connectionReader(ws)
}

func connectionReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message from client via WS: %v", err)
			return
		}

		log.Printf("Message type: %d. Message: %s", messageType, p)
	}
}
