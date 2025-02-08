package request

import "app/utils"

type CustomerParam struct {
	UID   string  `json:"uuid"`
	Name  string  `json:"name"`
	Phone string  `json:"phone"`
	CarNo string  `json:"carNo"`
	Debt  float64 `json:"debt"`
}

type ListCustomerParam struct {
	utils.LimitCond
	SearchKey string `json:"searchName"`
}
