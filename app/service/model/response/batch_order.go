package response

import "app/service/model"

type (
	BatchOrderRsp struct {
		model.BatchOrder
		model.CustomerField
		History model.BatchOrderHistory `json:"history"`
	}
)
