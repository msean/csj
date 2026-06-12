package request

import "app/pkg/utils"

type (
	GoodsCategoryListReq struct {
		Brief bool `json:"brief"`
		utils.LimitCond
	}
	GoodsListReq struct {
		SearchKey string `json:"searchName"`
		utils.LimitCond
		LoadAll bool `json:"loadALl"`
	}
)
