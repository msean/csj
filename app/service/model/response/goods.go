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
		Goods []GoodsDetailRsp `gorm:"-" json:"goodsList"`
	}
)

// 自定义 JSON 序列化逻辑
func (g GoodsDetailRsp) MarshalJSON() ([]byte, error) {
	type Alias GoodsDetailRsp
	g.UIDCompatible = utils.Violent2String(g.UID)
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

// 自定义 JSON 序列化逻辑
func (g GoodsCategoryRsp) MarshalJSON() ([]byte, error) {
	type Alias GoodsCategoryRsp
	g.UIDCompatible = utils.Violent2String(g.UID)
	return json.Marshal(&struct {
		Alias
		OwnerUserCompatible string `json:"ownerUser"`
	}{
		Alias:               Alias(g),
		OwnerUserCompatible: utils.Violent2String(g.OwnerUser),
	})
}
