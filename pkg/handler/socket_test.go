package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func httpToWS(t *testing.T, u string) string {
	t.Helper()

	wsURL, err := url.Parse(u)
	require.NoError(err)

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}

func testWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	s := httptest.NewServer(h)
	wsURL := httpToWS(t, s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError("Failed to dial WS url: %v", err)
	return s, ws
}

func TestSocketHandler(t *testing.T) {
	tt := []struct {
		name             string
		message          string
		expectedResponse string
	}{
		{
			name:             "Empty message",
			message:          "",
			expectedResponse: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			s, ws := testWSServer(wsEndpoint())
		})
	}
}
