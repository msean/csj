package request

type OrderTicketReq struct {
	OrderUUID  string `json:"orderUUID" binding:"required,max=64"`
	PaperWidth int    `json:"paperWidth" binding:"required,oneof=58 80"`
	PrintMode  string `json:"printMode" binding:"required,oneof=image"`
	FeedLines  *int   `json:"feedLines" binding:"required,oneof=0 5"`
}
