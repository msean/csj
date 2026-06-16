package request

type (
	LoginReq struct {
		Phone      string `json:"phone"`
		VerifyCode string `json:"verifyCode"`
	}
	SendVerifyCodeReq struct {
		Phone string `json:"phone"`
		Typ   int    `json:"type"`
	}
)
