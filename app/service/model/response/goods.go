package response

import (
	"app/service/model"
	"app/utils"
	"encoding/json"
)

type (
	GoodsDetailRsp struct {
		model.Goods
	}
)

type (
	GoodsCategoryRsp struct {
		model.GoodsCategory
	}
)

// 自定义 JSON 序列化逻辑
func (g GoodsDetailRsp) MarshalJSON() ([]byte, error) {
	type Alias GoodsDetailRsp
	return json.Marshal(&struct {
		Alias
		OwnerUserCompatible  string `json:"ownerUser"`
		CategoryIDCompatible string `json:"categoryID"`
	}{
		Alias:                Alias(g),
		OwnerUserCompatible:  utils.Violent2String(g.OwnerUser),
		CategoryIDCompatible: utils.Violent2String(g.CategoryID),
	})
}
