package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func httpToWS(t *testing.T, u string) string {
	wsURL, err := url.Parse(u)
	require.NoError(t, err, "URL parsing failed: %v", err)

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}

func testWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	s := httptest.NewServer(h)
	wsURL := httpToWS(t, s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to dial WS url: %v", err)
	return s, ws
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg []byte) {
	require.NoError(t, ws.WriteMessage(websocket.BinaryMessage, msg))
}

func receiveMessage(t *testing.T, ws *websocket.Conn) []byte {
	_, m, err := ws.ReadMessage()
	require.NoError(t, err, "Failed to read message: %v", err)

	return m
}

type wsEndopointHandler struct {
}

func (weh wsEndopointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	SocketUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := SocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer ws.Close()

	ConnectionReader(ws)
}

func marshal(t *testing.T, in interface{}) []byte {
	b, err := json.Marshal(&in)
	require.NoError(t, err, "failed to marshal: %v", err)
	return b
}

func TestSocketHandler(t *testing.T) {
	tt := []struct {
		name             string
		inbound          Inbound
		expectedResponse Outbound
	}{
		{
			name: "Empty message",
			inbound: Inbound{
				ID: "1",
			},
			expectedResponse: Outbound{
				ID:    "1",
				Error: &UnknownMessageTypeErr,
			},
		},
		{
			name: "Status check",
			inbound: Inbound{
				ID:   "1",
				Type: TypeStatus,
			},
			expectedResponse: Outbound{
				ID:   "1",
				Data: "{\"status\":\"ok\"}",
			},
		},
		{
			name: "Create array check",
			inbound: Inbound{
				ID:   "1",
				Type: TypeCreateArray,
				Data: marshal(t,
					CreateArrayRequest{
						Number: 4,
					}),
			},
			expectedResponse: Outbound{
				ID:   "1",
				Data: string(marshal(t, CreateArrayResponse{Numbers: []int{0, 1, 2, 3}})),
			},
		},
		{
			name: "Create array negative number",
			inbound: Inbound{
				ID:   "1",
				Type: TypeCreateArray,
				Data: marshal(t, CreateArrayRequest{Number: -10}),
			},
			expectedResponse: Outbound{
				ID:    "1",
				Error: &CreateArrayNegativeNumberErr,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, ws := testWSServer(t, wsEndopointHandler{})
			defer s.Close()
			defer ws.Close()

			msgBytes, err := json.Marshal(&tc.inbound)
			require.NoError(t, err, "failed to marshal inbound message: %v", err)

			sendMessage(t, ws, msgBytes)
			response := receiveMessage(t, ws)

			require.Equal(t, string(marshal(t, tc.expectedResponse)), string(response))
		})
	}
}
