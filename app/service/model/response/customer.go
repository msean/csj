package response

import "app/service/model"

type (
	ListCustomerRsp struct {
		model.Customer
		LatestBillDate int `json:"latestBillDate"`
	}
)
