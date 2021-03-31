package socket

func handleStatus(in Inbound) Outbound {
	return Outbound{
		ID:    in.ID,
		Error: nil,
		Data:  "{\"status\":\"ok\"}",
	}
}
