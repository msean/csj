package request

import "app/service/model"

type CustomerParam struct {
	UID   string  `json:"uuid"`
	Name  string  `json:"name"`
	Phone string  `json:"phone"`
	CarNo string  `json:"carNo"`
	Debt  float64 `json:"debt"`
}

type ListCustomerParam struct {
	model.LimitCond
	SearchKey string `json:"searchName"`
}
