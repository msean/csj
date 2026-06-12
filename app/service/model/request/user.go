package request

type (
	SendVerifyCodeReq struct {
		Phone string `json:"phone"`
		Typ   int    `json:"type"`
	}
)
