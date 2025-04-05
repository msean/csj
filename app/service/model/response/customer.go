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
	b.UIDCompatible = utils.Violent2String(b.UID)
	type Alias ListCustomerRsp
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&b),
	})
}
