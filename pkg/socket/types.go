package socket

import "encoding/json"

// Inbound is a data that comes from the client.
type Inbound struct {
	// Message unique ID can help client
	ID string `json:"id"`
	// Internal type of message
	Type string `json:"type"`
	// Data which differs from type to type
	Data json.RawMessage `json:"data"`
}

// Outbond defines data that will be sent back to the client.
type Outbound struct {
	ID    string `json:"id"`
	Error string `json:"error"`
	Data  string `json:"data"`
}

// Error is a wrapper for any errors happend inside the server,
// where Code is a unique error code for message to be translatable
// and Message is a human readable message.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
