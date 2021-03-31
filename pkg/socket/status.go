package socket

func handleStatus(in Inbound) (Outbound, error) {
	return Outbound{
		ID:    in.ID,
		Error: nil,
		Data:  "{\"status\":\"ok\"}",
	}, nil
}
