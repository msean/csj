package request

type (
	RegisterParam struct {
		Phone      string `json:"phone"`
		VerifyCode string `json:"verifycode"`
	}
	VerfifyCodeParam struct {
		Phone string `json:"phone"`
		Typ   int    `json:"type"`
	}
	LoginParam struct {
		Phone      string `json:"phone"`
		VerifyCode string `json:"verifycode"`
	}
	UserUpdateParam struct {
		Phone         string `json:"phone"`
		Name          string `json:"name"`
		UID           int64  `json:"uid"`
		UIDCompatible int64
	}
)
