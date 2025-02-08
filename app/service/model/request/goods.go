package request

import "app/utils"

type (
	GoodsCategorySaveParam struct {
		Name string `json:"name"`
		UID  string `json:"uuid"`
	}
	GoodsSaveParam struct {
		UID        string  `json:"uuid"`
		CategoryID string  `json:"categoryID"`
		Name       string  `json:"name"`
		Type       int32   `json:"type"`
		Price      float32 `json:"price"`
		Weight     float32 `json:"weight"`
	}
	ListGoodsParam struct {
		SearchKey string `json:"searchName"`
		utils.LimitCond
	}
	ListGoodsCategoryParam struct {
		Brief bool `json:"brief"`
		utils.LimitCond
	}
	DeleteGoodsCategoryParam struct {
		UUID string `json:"uuid"`
	}
)
