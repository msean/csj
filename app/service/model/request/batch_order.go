package request

import (
	"app/utils"
	"encoding/json"
	"strconv"
)

type (
	BatchOrderGoodsParam struct {
		GoodsUUID           string `json:"goodsUUID"`
		GoodsUUIDCompatible int64
		Price               float64 `json:"price"`
		Weight              float64 `json:"weight"`
		Mount               int32   `json:"mount"`
		SerialNo            string  `json:"serialNo"`
	}
	CreateTempBatchOrderParam struct {
		BatchUUID              string `json:"batchUUID"`
		BatchUUIDCompatible    int64
		CustomerUUID           string `json:"customerUUID"`
		CustomerUUIDCompatible int64
		GoodsList              []BatchOrderGoodsParam `json:"goodsList"`
	}
	BatchOrderDetailParam struct {
		UUID           string `json:"uuid"`
		UUIDCompatible int64
	}
	CreateBatchOrderParam struct {
		BatchUUID              string `json:"batchUUID"`
		BatchUUIDCompatible    int64
		CustomerUUID           string `json:"customerUUID"`
		CustomerUUIDCompatible int64
		GoodsList              []BatchOrderGoodsParam `json:"goodsList"`
		FPayAmount             float64                `json:"payAmount"` // 总计
		PayType                int32                  `json:"payType"`   // 支付方式
	}
	UpdateBatchOrderParam struct {
		BatchOrderUUID           string `json:"uuid"`
		BatchOrderUUIDCompatible int64
		CustomerUUID             string `json:"customerUUID"`
		CustomerUUIDCompatible   int64
		GoodsList                []BatchOrderGoodsParam `json:"goodsList"`
	}
	UpdateBatchOrderStatusParam struct {
		BatchOrderUUID           string `json:"uuid"`
		BatchOrderUUIDCompatible int64
		Status                   int32 `json:"status"`
	}
	GoodsLatestOrderParam struct {
		GoodsUUIDList []string `json:"goodsUUIDList"`
	}
	ListBatchOrderParam struct {
		utils.LimitCond
		Status             int32  `json:"status"`
		UserUUID           string `json:"userUUID"`
		UserUUIDCompatible int64
		StartTime          int64 `json:"startTime"`
		EndTime            int64 `json:"endTime"`
	}
	ShareBatchOrderParam struct {
		UUID           string `json:"uuid"`
		UUIDCompatible int64
	}
)

func (param *BatchOrderGoodsParam) UnmarshalJSON(data []byte) error {
	type Alias BatchOrderGoodsParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.GoodsUUID != "" {
		if uid, err := strconv.ParseInt(param.GoodsUUID, 10, 64); err == nil {
			param.GoodsUUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *CreateTempBatchOrderParam) UnmarshalJSON(data []byte) error {
	type Alias CreateTempBatchOrderParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.BatchUUID != "" {
		if uid, err := strconv.ParseInt(param.BatchUUID, 10, 64); err == nil {
			param.BatchUUIDCompatible = uid
		} else {
			return err
		}
	}

	return nil
}

func (param *BatchOrderDetailParam) UnmarshalJSON(data []byte) error {
	type Alias BatchOrderDetailParam
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

func (param *CreateBatchOrderParam) UnmarshalJSON(data []byte) error {
	type Alias CreateBatchOrderParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.BatchUUID != "" {
		if uid, err := strconv.ParseInt(param.BatchUUID, 10, 64); err == nil {
			param.BatchUUIDCompatible = uid
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

func (param *UpdateBatchOrderParam) UnmarshalJSON(data []byte) error {
	type Alias UpdateBatchOrderParam
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

func (param *UpdateBatchOrderStatusParam) UnmarshalJSON(data []byte) error {
	type Alias UpdateBatchOrderStatusParam
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
	return nil
}

func (param *ListBatchOrderParam) UnmarshalJSON(data []byte) error {
	type Alias ListBatchOrderParam
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(param),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if param.UserUUID != "" {
		if uid, err := strconv.ParseInt(param.UserUUID, 10, 64); err == nil {
			param.UserUUIDCompatible = uid
		} else {
			return err
		}
	}

	if param.UserUUID != "" {
		if uid, err := strconv.ParseInt(param.UserUUID, 10, 64); err == nil {
			param.UserUUIDCompatible = uid
		} else {
			return err
		}
	}
	return nil
}

func (param *ShareBatchOrderParam) UnmarshalJSON(data []byte) error {
	type Alias ShareBatchOrderParam
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
