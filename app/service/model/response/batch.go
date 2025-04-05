package response

import (
	"app/service/model"
	"app/utils"
	"encoding/json"
)

type (
	BatchRsp struct {
		model.Batch
		Goods []*BatchGoodsRsp `json:"goodsList"`
	}
	BatchGoodsRsp struct {
		model.BatchGoods
		model.GoodsField
		model.SurplusField
	}
)

func (b BatchRsp) MarshalJSON() ([]byte, error) {
	type Alias BatchRsp
	b.UIDCompatible = utils.Violent2String(b.UID)
	return json.Marshal(&struct {
		OwnerUserCompatible string `json:"ownerUser"`
		*Alias
	}{
		OwnerUserCompatible: utils.Violent2String(b.OwnerUser),
		Alias:               (*Alias)(&b),
	})
}

func (b BatchGoodsRsp) MarshalJSON() ([]byte, error) {
	type Alias BatchGoodsRsp
	b.UIDCompatible = utils.Violent2String(b.UID)
	return json.Marshal(&struct {
		OwnerUserCompatible string `json:"ownerUser"`
		BatchUUIDCompatible string `json:"batchUUID"`
		GoodsUUIDCompatible string `json:"goodsUUID"`
		*Alias
	}{
		OwnerUserCompatible: utils.Violent2String(b.OwnerUser),
		BatchUUIDCompatible: utils.Violent2String(b.BatchUUID),
		GoodsUUIDCompatible: utils.Violent2String(b.GoodsUUID),
		Alias:               (*Alias)(&b),
	})
}
