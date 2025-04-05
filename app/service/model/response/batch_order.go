package response

import (
	"app/service/model"
	"app/utils"
	"encoding/json"
)

type (
	BatchOrderRsp struct {
		model.BatchOrder
		model.CustomerField
		History model.BatchOrderHistory `json:"history"`
	}
)

func (b BatchOrderRsp) MarshalJSON() ([]byte, error) {
	type Alias BatchOrderRsp
	b.UIDCompatible = utils.Violent2String(b.UID)
	return json.Marshal(&struct {
		OwnerUserCompatible string `json:"ownerUser"`
		BatchUUIDCompatible string `json:"batchUUID"`
		CustomerUUID        string `json:"customerUUID"`
		*Alias
	}{
		OwnerUserCompatible: utils.Violent2String(b.OwnerUser),
		BatchUUIDCompatible: utils.Violent2String(b.BatchUUID),
		CustomerUUID:        utils.Violent2String(b.UserUUID),
		Alias:               (*Alias)(&b),
	})
}
