package socket

var (
	InvalidMessageErr       = Error{Code: "invalid_message", Message: "Invalid message."}
	InternalErrorMessageErr = Error{Code: "internal_error", Message: "Something went wrong."}
	UnknownMessageTypeErr   = Error{Code: "unknown_message_type", Message: "Unknown message type."}
)
