package request

import (
	"app/utils"
	"encoding/json"
	"strconv"
)

type (
	GoodsCategorySaveParam struct {
		Name          string `json:"name"`
		UID           string `json:"uuid"`
		UIDCompatible int64
	}
	GoodsSaveParam struct {
		UID                  string  `json:"uuid"`
		CategoryID           string  `json:"categoryID"`
		Name                 string  `json:"name"`
		Type                 int32   `json:"type"`
		Price                float32 `json:"price"`
		Weight               float32 `json:"weight"`
		UIDCompatible        int64
		CategoryIDCompatible int64
	}
	ListGoodsParam struct {
		SearchKey string `json:"searchName"`
		utils.LimitCond
		OrderBy string `json:"orderBy"`
	}
	ListGoodsCategoryParam struct {
		Brief bool `json:"brief"`
		utils.LimitCond
	}
	DeleteGoodsCategoryParam struct {
		UID           string `json:"uuid"`
		UIDCompatible int64
	}
)

func (param *GoodsCategorySaveParam) UnmarshalJSON(data []byte) error {
	type Alias GoodsCategorySaveParam
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

func (param *GoodsSaveParam) UnmarshalJSON(data []byte) error {
	type Alias GoodsSaveParam
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

	if param.CategoryID != "" {
		if uid, err := strconv.ParseInt(param.CategoryID, 10, 64); err == nil {
			param.CategoryIDCompatible = uid
		} else {
			return err
		}
	}
	return nil
}

func (param *DeleteGoodsCategoryParam) UnmarshalJSON(data []byte) error {
	type Alias DeleteGoodsCategoryParam
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
