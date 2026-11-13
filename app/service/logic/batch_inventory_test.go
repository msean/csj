package logic

import (
	"app/service/dao"
	"app/service/model"
	"testing"
)

func TestBatchInventoryUnitsAndUnsoldGoods(t *testing.T) {
	fixed := &model.BatchGoods{GoodsUUID: "fixed", Mount: 10, Weight: 26}
	fixed.GoodsTyp = 1
	bulk := &model.BatchGoods{GoodsUUID: "bulk", Weight: 12.75}
	bulk.GoodsTyp = 2
	unsold := &model.BatchGoods{GoodsUUID: "unsold", Mount: 4}
	unsold.GoodsTyp = 1
	oversold := &model.BatchGoods{GoodsUUID: "oversold", Weight: 2}
	oversold.GoodsTyp = 2
	applyBatchInventory([]*model.BatchGoods{fixed, bulk, unsold, oversold}, []dao.BatchInventoryUsage{
		{GoodsUUID: "fixed", Mount: 3, Weight: 78},
		{GoodsUUID: "bulk", Weight: 2.5},
		{GoodsUUID: "oversold", Weight: 3.5},
	})
	if fixed.Surplus != "7" || bulk.Surplus != "10.25" || unsold.Surplus != "4" || oversold.Surplus != "-1.50" {
		t.Fatalf("incorrect units or subtraction: %s %s %s %s", fixed.Surplus, bulk.Surplus, unsold.Surplus, oversold.Surplus)
	}
	if fixed.SellAmount != "" || bulk.SellTotalAmount != "" {
		t.Fatal("inventory-only mode populated sales statistics")
	}
}
