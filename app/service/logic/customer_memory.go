package logic

import (
	"app/service/dao"
	"app/service/model"
	"fmt"
)

type CustomerMemoryRsp struct {
	OrderUUID string                   `json:"orderUUID"`
	GoodsList []*model.BatchOrderGoods `json:"goodsList"`
}

// CustomerMemory scopes both customer and sharded orders to the authenticated merchant.
func (logic *OrderLogic) CustomerMemory(customerUUID string) (rsp CustomerMemoryRsp, err error) {
	rsp.GoodsList = make([]*model.BatchOrderGoods, 0)
	var customer model.Customer
	if err = logic.runtime.DB.Where("uid = ? AND owner_user = ? AND deleted_at IS NULL", customerUUID, logic.OwnerUser).First(&customer).Error; err != nil {
		return
	}
	if customer.Name == "现金客户" {
		return rsp, fmt.Errorf("现金客户不支持客户记忆开单")
	}
	rows, err := dao.CustomerPurchaseList(logic.runtime.DB, logic.OwnerUser, customerUUID, 0, 0, false, 0, 1)
	if err != nil || len(rows) == 0 {
		return rsp, err
	}
	logic.BatchFeilds(rows)
	rsp.OrderUUID = rows[0].UID
	rsp.GoodsList = rows[0].GoodsListRelated
	for _, good := range rsp.GoodsList {
		good.Mount = 0
		good.Weight = 0
		good.Total = 0
	}
	return rsp, nil
}
