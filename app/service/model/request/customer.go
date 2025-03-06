package request

import (
	"app/utils"
	"encoding/json"
	"strconv"
)

type CustomerParam struct {
	UID           string  `json:"uuid"`
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	CarNo         string  `json:"carNo"`
	Debt          float64 `json:"debt"`
	UIDCompatible int64
}

type ListCustomerParam struct {
	utils.LimitCond
	SearchKey string `json:"searchName"`
}

func (param *CustomerParam) UnmarshalJSON(data []byte) error {
	type Alias CustomerParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UID != "" {
		if uid, err := strconv.ParseInt(param.UID, 10, 64); err == nil {
			param.UIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}
