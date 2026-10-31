package appmgr

import (
	"encoding/json"
	"github.com/msean/csj/shared/constants"
	model "github.com/msean/csj/shared/model"
	"github.com/msean/csj/shared/sharding"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm/schema"
	"sync"
	"testing"
)

func TestCanonicalModelsPreserveBusinessBehavior(t *testing.T) {
	var base model.BaseModel
	if err := base.BeforeCreate(nil); err != nil {
		t.Fatal(err)
	}
	if _, err := ulid.ParseStrict(base.UID); err != nil {
		t.Fatalf("app IDs must remain ULIDs: %v", err)
	}
	original := base.UID
	_ = base.BeforeCreate(nil)
	if base.UID != original {
		t.Fatal("assigned UID changed")
	}
	user := model.User{Phone: "123"}
	if user.Title() != "123" {
		t.Fatal("user title fallback changed")
	}
	order := model.BatchOrder{}
	order.DefaultSet()
	if order.Status != constants.BatchOrderTemp {
		t.Fatal("order default changed")
	}
	fixed := model.BatchOrderGoods{GoodType: constants.GoodsTypeFix, Price: 12.5, Mount: 3}
	bulk := model.BatchOrderGoods{GoodType: constants.GoodsTypeBulk, Price: 12.5, Weight: 2.3}
	if fixed.Amount() != 37.5 || fixed.SetTotal() != 38 || fixed.Sell() != 3 || bulk.SetTotal() != 29 || bulk.Sell() != 2.3 {
		t.Fatal("financial calculations changed")
	}
	if sharding.GetShardTableName("batch_orders", "alice") != "batch_orders_35" || len(sharding.GetAllShardTableNames("batch_orders")) != 50 {
		t.Fatal("shard routing changed")
	}
	// Both services must use the exact same Go type, not copied persistence structs.
	var canonical model.User
	var backend AppUser = canonical
	canonical = backend
	raw, _ := json.Marshal(canonical)
	if string(raw) == "" {
		t.Fatal("shared serialization failed")
	}
}

func TestCanonicalTableNamesIgnoreAdminPrefix(t *testing.T) {
	for _, entity := range []any{model.User{}, model.Customer{}, model.GoodsCategory{}, model.Goods{}, model.BatchGoods{}, model.BatchOrder{}, model.BatchOrderGoods{}, model.BatchOrderOpr{}, model.BatchOrderPay{}, model.BatchOrderHistory{}, model.MessageCenter{}, model.OfficialMessage{}, model.OfficialMessageRead{}, model.ShareRecord{}, model.Sms{}} {
		parsed, err := schema.Parse(entity, &sync.Map{}, schema.NamingStrategy{TablePrefix: "admin_"})
		if err != nil {
			t.Fatal(err)
		}
		expected := (schema.NamingStrategy{}).TableName(parsed.Name)
		if parsed.Table != expected {
			t.Fatalf("changed app table for %s: %s, expected %s", parsed.Name, parsed.Table, expected)
		}
	}
}
