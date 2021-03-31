package socket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const WS_TAG = "WebSocket"

// SocketUpgrader is a websocket configuration.
var SocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	TypeStatus      = "status"
	TypeCreateArray = "create-array"
)

var functions = map[string]func(in Inbound) (Outbound, error){
	TypeStatus:      handleStatus,
	TypeCreateArray: handleCreateArray,
}

// ConnectionReader is a main function which handles websocket messaging
// and data transfer.
func ConnectionReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if messageType >= websocket.CloseNormalClosure &&
			messageType <= websocket.CloseTLSHandshake || messageType == websocket.CloseMessage {
			log.Printf("[%s] Received close code", WS_TAG)
			return
		}
		if err != nil {
			fail(conn, Inbound{}, fmt.Errorf("failed to read message from client via WS: %v", err), InvalidMessageErr)
			return
		}

		var in Inbound
		if err := json.Unmarshal(p, &in); err != nil {
			fail(conn, in, fmt.Errorf("failed to unmarshal message: %v", err), InvalidMessageErr)
			continue
		}

		function, exist := functions[in.Type]
		if !exist {
			fail(conn, in, fmt.Errorf("unknown message type: '%s'", in.Type), UnknownMessageTypeErr)
			continue
		}

		out, err := function(in)
		if err != nil {
			// Here we don't use fail as outbound message is already
			// formatted correctly
			log.Printf("[%s] %v", WS_TAG, err)
		}

		outBytes, err := json.Marshal(&out)
		if err != nil {
			fail(conn, in, fmt.Errorf("failed to marshal outbound message: %v", err), InternalErrorMessageErr)
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, outBytes); err != nil {
			fail(conn, in, fmt.Errorf("failed to write a response message: %v", err), InternalErrorMessageErr)
			continue
		}
	}
}

func fail(conn *websocket.Conn, in Inbound, logError error, msgError Error) {
	log.Printf("[%s] %v", WS_TAG, logError)
	out := Outbound{
		ID:    in.ID,
		Error: &msgError,
	}
	outBytes, logError := json.Marshal(&out)
	if logError != nil {
		log.Printf("[%s] failed to marshal oubound error response: %v", WS_TAG, logError)
		return
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, outBytes); err != nil {
		log.Printf("[%s] failed to write error response: %v", WS_TAG, err)
	}
}
