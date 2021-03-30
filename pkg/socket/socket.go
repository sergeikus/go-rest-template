package socket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const WS_TAG = "WebSocket"

var SocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectionReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if messageType >= websocket.CloseNormalClosure &&
			messageType <= websocket.CloseTLSHandshake || messageType == websocket.CloseMessage {
			log.Printf("[%s] Received close code", WS_TAG)
			return
		}
		if err != nil {
			fail(conn, Inbound{}, fmt.Errorf("failed to read message from client via WS: %v", err))
			return
		}

		var in Inbound
		if err := json.Unmarshal(p, &in); err != nil {
			fail(conn, in, fmt.Errorf("failed to unmarshal message: %v", err))
			continue
		}
		log.Printf("Marshalled inbound: %#v", in)

		out := Outbound{
			ID: in.ID,
		}
		outBytes, err := json.Marshal(&out)
		if err != nil {
			fail(conn, in, fmt.Errorf("failed to marshal outbound message: %v", err))
			continue
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, outBytes); err != nil {
			fail(conn, in, fmt.Errorf("failed to write a response message: %v", err))
			continue
		}
	}
}

func fail(conn *websocket.Conn, in Inbound, err error) {
	log.Printf("[%s] %v", WS_TAG, err)
	out := Outbound{
		ID:    in.ID,
		Error: err.Error(),
	}
	outBytes, err := json.Marshal(&out)
	if err != nil {
		log.Printf("[%s] failed to marshal oubound error response: %v", WS_TAG, err)
		return
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, outBytes); err != nil {
		log.Printf("[%s] failed to write error response: %v", WS_TAG, err)
	}
}
