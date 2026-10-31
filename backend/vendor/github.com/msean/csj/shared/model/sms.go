package model

type Sms struct {
	BaseModel
	Phone string
	Code  string
}
