package response

type RegisterRsp struct {
	UUID  string `json:"uuid"`
	Token string `json:"token"`
}
