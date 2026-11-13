package response

type OrderTicketRsp struct {
	TicketBytes string `json:"ticketBytes"`
	BytesLength int    `json:"bytesLength"`
}
