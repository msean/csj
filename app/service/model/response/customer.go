package response

import (
	"app/service/model"
	"app/utils"
	"encoding/json"
)

type (
	ListCustomerRsp struct {
		model.Customer
		LatestBillDate int `json:"latestBillDate"`
	}
)

func (b ListCustomerRsp) MarshalJSON() ([]byte, error) {
	type Alias ListCustomerRsp
	return json.Marshal(&struct {
		OwnerUserCompatible string `json:"ownerUser"`
		*Alias
	}{
		OwnerUserCompatible: utils.Violent2String(b.OwnerUser),
		Alias:               (*Alias)(&b),
	})
}
