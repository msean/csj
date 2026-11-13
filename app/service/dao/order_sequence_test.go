package dao

import (
	"app/service/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestOrderGoodsSequenceSavedAcrossInsertChunks(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "user:password@tcp(localhost:3306)/test", SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, DryRun: true, SkipDefaultTransaction: true})
	if err != nil {
		t.Fatal(err)
	}
	// Multiple rows for one goods UUID remain separate and retain submitted order.
	goods := make([]*model.BatchOrderGoods, 103)
	for index := range goods {
		goods[index] = &model.BatchOrderGoods{GoodsUUID: "same-goods", SortNo: 999}
	}
	if err := OrderDao.CreateGoodsBatch(db, "owner", goods); err != nil {
		t.Fatal(err)
	}
	for index, good := range goods {
		if good.SortNo != index+1 {
			t.Fatalf("row %d saved with sortNo %d", index, good.SortNo)
		}
	}
	// A changed cart order is numbered afresh on the edit save path.
	reordered := []*model.BatchOrderGoods{goods[102], goods[0]}
	if err := OrderDao.CreateGoodsBatch(db, "owner", reordered); err != nil {
		t.Fatal(err)
	}
	if reordered[0].SortNo != 1 || reordered[1].SortNo != 2 {
		t.Fatal("edit did not preserve new array order")
	}
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return OrderGoodsInSequence(tx.Table(ShardOrderGoodsTable("owner"))).Find(&[]model.BatchOrderGoods{})
	})
	if !strings.Contains(sql, "ORDER BY sort_no ASC, created_at ASC, uid ASC") {
		t.Fatalf("unstable read ordering: %s", sql)
	}
}
