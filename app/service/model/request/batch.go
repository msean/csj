package request

import (
	"encoding/json"
	"strconv"
)

type (
	SaveBatchGoods struct {
		UUID           string `json:"goodsUUID"`
		UUIDCompatible int64
		Price          float64 `json:"price"`
		Weight         float64 `json:"weight"`
		Mount          int32   `json:"mount"`
	}
	CreateBatchParam struct {
		StorageTime int64            `json:"storageTime"`
		Goods       []SaveBatchGoods `json:"goodsList"`
	}
	BatchDetailParam struct {
		UUID           string `json:"uuid"`
		Date           string `json:"date"`
		UUIDCompatible int64
	}
	UpdateBatchParam struct {
		StorageTime    int64  `json:"storageTime"`
		UUID           string `json:"uuid"`
		UUIDCompatible int64
		Goods          []SaveBatchGoods `json:"goodsList"`
	}
	UpdateBatchStatusParam struct {
		UUID           string `json:"uuid"`
		Status         int32  `json:"status"`
		UUIDCompatible int64
	}
	FindBatchGoodsParam struct {
		UUID           string `json:"uuid"`
		UUIDCompatible int64
	}
	UpdateBatchGoodsParam struct {
		UUID           string `json:"uuid"`
		UUIDCompatible int64
		Price          float64 `json:"price"`
		Weight         float64 `json:"weight"`
		Mount          int32   `json:"mount"`
	}
)

func (param *SaveBatchGoods) UnmarshalJSON(data []byte) error {
	type Alias SaveBatchGoods
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UUID != "" {
		if uid, err := strconv.ParseInt(param.UUID, 10, 64); err == nil {
			param.UUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *BatchDetailParam) UnmarshalJSON(data []byte) error {
	type Alias BatchDetailParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UUID != "" {
		if uid, err := strconv.ParseInt(param.UUID, 10, 64); err == nil {
			param.UUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *UpdateBatchParam) UnmarshalJSON(data []byte) error {
	type Alias UpdateBatchParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UUID != "" {
		if uid, err := strconv.ParseInt(param.UUID, 10, 64); err == nil {
			param.UUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *UpdateBatchStatusParam) UnmarshalJSON(data []byte) error {
	type Alias UpdateBatchStatusParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UUID != "" {
		if uid, err := strconv.ParseInt(param.UUID, 10, 64); err == nil {
			param.UUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *FindBatchGoodsParam) UnmarshalJSON(data []byte) error {
	type Alias FindBatchGoodsParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UUID != "" {
		if uid, err := strconv.ParseInt(param.UUID, 10, 64); err == nil {
			param.UUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *UpdateBatchGoodsParam) UnmarshalJSON(data []byte) error {
	type Alias UpdateBatchGoodsParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UUID != "" {
		if uid, err := strconv.ParseInt(param.UUID, 10, 64); err == nil {
			param.UUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}
