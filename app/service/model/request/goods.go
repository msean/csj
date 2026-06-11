package request

import "app/pkg/utils"

type GoodsListReq struct {
	SearchKey string `json:"searchName"`
	utils.LimitCond
	LoadAll bool `json:"loadALl"`
}
