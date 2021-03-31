package socket

import "fmt"

var (
	InvalidMessageErr       = Error{Code: "invalid_message", Message: "Invalid message."}
	InternalErrorMessageErr = Error{Code: "internal_error", Message: "Something went wrong."}
	UnknownMessageTypeErr   = Error{Code: "unknown_message_type", Message: "Unknown message type."}
)

func loggingError(typeTag string, msgID string, err error) error {
	return fmt.Errorf("[%s][%s] %v", typeTag, msgID, err)
}
