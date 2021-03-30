package handler

import (
	"fmt"
	"net/http"

	"github.com/sergeikus/go-rest-template/pkg/socket"
)

func (api *API) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	socket.SocketUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := socket.SocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fail(w, socket.WS_TAG, fmt.Errorf("failed to upgrade connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	socket.ConnectionReader(ws)
}
