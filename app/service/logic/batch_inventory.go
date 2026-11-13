package logic

import (
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"fmt"
)

func (logic *BatchLogic) setInventoryOnly() error {
	if len(logic.GoodsListRelated) == 0 {
		return nil
	}
	var usage []dao.BatchInventoryUsage
	if err := dao.BatchInventoryUsageQuery(logic.runtime.DB, logic.OwnerUser, logic.UID).Find(&usage).Error; err != nil {
		return err
	}
	applyBatchInventory(logic.GoodsListRelated, usage)
	return nil
}

func applyBatchInventory(goods []*model.BatchGoods, usage []dao.BatchInventoryUsage) {
	sold := make(map[string]dao.BatchInventoryUsage, len(usage))
	for _, item := range usage {
		sold[item.GoodsUUID] = item
	}
	for _, good := range goods {
		used := sold[good.GoodsUUID]
		if good.GoodsTyp == common.GoodsTypeBulk {
			good.Surplus = fmt.Sprintf("%.2f", good.Weight-used.Weight)
		} else {
			good.Surplus = fmt.Sprintf("%d", good.Mount-used.Mount)
		}
	}
}
