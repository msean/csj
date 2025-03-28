package request

import (
	"encoding/json"
	"strconv"
)

type (
	CreateBatchOrderPayParam struct {
		BatchOrderUUID           string  `json:"batchOrderUUID"`
		CustomerUUID             string  `json:"customerUUID"`
		Amount                   float64 `json:"amount"`
		PayType                  int32   `json:"payType"`
		BatchOrderUUIDCompatible int64
		CustomerUUIDCompatible   int64
	}
	UpdateBatchOrderPayParam struct {
		BatchOrderPayUUID           string  `json:"uuid"`
		Amount                      float64 `json:"amount"`
		PayType                     int32   `json:"payType"`
		BatchOrderPayUUIDCompatible int64
	}
	GetBatchOrderPayParam struct {
		BatchOrderPayUUID           string `json:"uuid"`
		BatchOrderPayUUIDCompatible int64
	}
)

func (param *CreateBatchOrderPayParam) UnmarshalJSON(data []byte) error {
	type Alias CreateBatchOrderPayParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.BatchOrderUUID != "" {
		if uid, err := strconv.ParseInt(param.BatchOrderUUID, 10, 64); err == nil {
			param.BatchOrderUUIDCompatible = uid
		} else {
			return err
		}
	}

	if param.CustomerUUID != "" {
		if uid, err := strconv.ParseInt(param.CustomerUUID, 10, 64); err == nil {
			param.CustomerUUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *UpdateBatchOrderPayParam) UnmarshalJSON(data []byte) error {
	type Alias UpdateBatchOrderPayParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.BatchOrderPayUUID != "" {
		if uid, err := strconv.ParseInt(param.BatchOrderPayUUID, 10, 64); err == nil {
			param.BatchOrderPayUUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *GetBatchOrderPayParam) UnmarshalJSON(data []byte) error {
	type Alias GetBatchOrderPayParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.BatchOrderPayUUID != "" {
		if uid, err := strconv.ParseInt(param.BatchOrderPayUUID, 10, 64); err == nil {
			param.BatchOrderPayUUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}
