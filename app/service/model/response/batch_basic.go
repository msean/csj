package response

import "app/service/model"

// BatchBasicDetail includes stock but excludes sales and profit statistics.
type BatchBasicDetail struct {
	model.BaseModel
	OwnerUser   string            `json:"ownerUser"`
	SerialNo    string            `json:"serialNo"`
	StorageTime int64             `json:"storageTime"`
	Status      int32             `json:"status"`
	SerialID    int               `json:"serialID"`
	GoodsList   []BatchBasicGoods `json:"goodsList"`
}

type BatchBasicGoods struct {
	Surplus string `json:"surplus"`
	model.BaseModel
	model.GoodsFeild
	OwnerUser string  `json:"ownerUser"`
	BatchUUID string  `json:"batchUUID"`
	SerialNo  string  `json:"serialNo"`
	SortNo    int     `json:"sortNo"`
	GoodsUUID string  `json:"goodsUUID"`
	GoodType  int     `json:"goods_type"`
	Price     float64 `json:"price"`
	Weight    float64 `json:"weight"`
	Mount     int     `json:"mount"`
}

func NewBatchBasicDetail(batch model.Batch) BatchBasicDetail {
	rsp := BatchBasicDetail{BaseModel: batch.BaseModel, OwnerUser: batch.OwnerUser, SerialNo: batch.SerialNo, StorageTime: batch.StorageTime, Status: batch.Status, SerialID: batch.SerialID, GoodsList: make([]BatchBasicGoods, 0, len(batch.GoodsListRelated))}
	for _, good := range batch.GoodsListRelated {
		serialNo := good.SerialNo
		if serialNo == "" {
			serialNo = batch.SerialNo
		}
		rsp.GoodsList = append(rsp.GoodsList, BatchBasicGoods{Surplus: good.Surplus, BaseModel: good.BaseModel, GoodsFeild: good.GoodsFeild, OwnerUser: good.OwnerUser, BatchUUID: good.BatchUUID, SerialNo: serialNo, SortNo: good.SortNo, GoodsUUID: good.GoodsUUID, GoodType: good.GoodType, Price: good.Price, Weight: good.Weight, Mount: good.Mount})
	}
	return rsp
}
