package appmgr

import (
	"github.com/glebarez/sqlite"
	"github.com/msean/csj/backend/global"
	appReq "github.com/msean/csj/backend/model/appmgr/request"
	commonReq "github.com/msean/csj/backend/model/common/request"
	model "github.com/msean/csj/shared/model"
	"github.com/msean/csj/shared/sharding"
	"gorm.io/gorm"
	"testing"
	"time"
)

func newAppDataTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}
func TestAppDataScopesAndFilters(t *testing.T) {
	db := newAppDataTestDB(t)
	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		t.Fatal(err)
	}
	deleted := time.Now()
	customers := []model.Customer{
		{BaseModel: model.BaseModel{UID: "a1"}, OwnerUser: "alice", Name: "甲", Phone: "100"},
		{BaseModel: model.BaseModel{UID: "a2"}, OwnerUser: "alice", Name: "乙", Phone: "200"},
		{BaseModel: model.BaseModel{UID: "b1"}, OwnerUser: "bob", Name: "甲", Phone: "100"},
		{BaseModel: model.BaseModel{UID: "deleted", DeletedAt: &deleted}, OwnerUser: "alice", Name: "甲", Phone: "100"},
	}
	if err := db.Create(&customers).Error; err != nil {
		t.Fatal(err)
	}
	service := AppDataService{}
	result, total, err := service.list(db, "customers", appReq.AppDataSearch{OwnerUser: "alice", Keyword: "100"})
	if err != nil || total != 1 {
		t.Fatalf("owner/OR search leaked records: %v %d", err, total)
	}
	if got := result.([]model.Customer); len(got) != 1 || got[0].UID != "a1" {
		t.Fatalf("unexpected rows: %+v", got)
	}
	result, total, err = service.list(db, "customers", appReq.AppDataSearch{OwnerUser: "alice", PageInfo: commonReq.PageInfo{Page: 1, PageSize: 1}})
	if err != nil || total != 2 || len(result.([]model.Customer)) != 1 {
		t.Fatalf("bad paging: %v %d", err, total)
	}
	if _, err = service.detail(db, "customers", "a1", "bob"); err != gorm.ErrRecordNotFound {
		t.Fatalf("cross-owner detail allowed: %v", err)
	}
	for _, resource := range []string{"orders", "payments", "unknown", "customers; DROP TABLE users"} {
		if _, _, err = service.list(db, resource, appReq.AppDataSearch{}); err == nil {
			t.Fatalf("unsafe resource accepted: %s", resource)
		}
	}
}
func TestShardedOrderDetails(t *testing.T) {
	db := newAppDataTestDB(t)
	owner := "alice"
	ordersTable := sharding.GetShardTableName("batch_orders", owner)
	goodsTable := sharding.GetShardTableName("batch_order_goods", owner)
	if err := db.Table(ordersTable).AutoMigrate(&model.BatchOrder{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("ALTER TABLE " + ordersTable + " ADD COLUMN id integer NOT NULL DEFAULT 0").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table(goodsTable).AutoMigrate(&model.BatchOrderGoods{}); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.BatchOrderHistory{}); err != nil {
		t.Fatal(err)
	}
	orders := []model.BatchOrder{{BaseModel: model.BaseModel{UID: "order-a"}, OwnerUser: owner, Status: 1, UserUUID: "customer-a"}, {BaseModel: model.BaseModel{UID: "order-b"}, OwnerUser: "bob", Status: 3}}
	if err := db.Table(ordersTable).Omit("GoodsListRelated").Create(&orders).Error; err != nil {
		t.Fatal(err)
	}
	goods := []model.BatchOrderGoods{{BaseModel: model.BaseModel{UID: "good-a"}, OwnerUser: owner, BatchOrderUID: "order-a"}, {BaseModel: model.BaseModel{UID: "foreign"}, OwnerUser: "bob", BatchOrderUID: "order-a"}}
	if err := db.Table(goodsTable).Create(&goods).Error; err != nil {
		t.Fatal(err)
	}
	status := 1
	service := AppDataService{}
	result, total, err := service.list(db, "orders", appReq.AppDataSearch{OwnerUser: owner, Status: &status, CustomerUUID: "customer-a"})
	if err != nil || total != 1 || len(result.([]model.BatchOrder)) != 1 {
		t.Fatalf("wrong shard/filter: %d %v", total, err)
	}
	detail, err := service.detail(db, "orders", "order-a", owner)
	if err != nil {
		t.Fatal(err)
	}
	if got := detail.Related["goodsList"].([]model.BatchOrderGoods); len(got) != 1 || got[0].UID != "good-a" {
		t.Fatalf("related records leaked: %+v", got)
	}
}
func TestRechargePendingAndSoftDelete(t *testing.T) {
	db := newAppDataTestDB(t)
	if err := db.AutoMigrate(&model.RechargeOrder{}, &model.AppVersion{}); err != nil {
		t.Fatal(err)
	}
	old := global.GVA_MYSQL
	global.GVA_MYSQL = db
	t.Cleanup(func() { global.GVA_MYSQL = old })
	orders := []model.RechargeOrder{{OrderNo: "pending", PayStatus: 0}, {OrderNo: "paid", PayStatus: 1}}
	if err := db.Create(&orders).Error; err != nil {
		t.Fatal(err)
	}
	pending := 0
	rows, total, err := (&RechargeOrderService{}).GetRechargeOrderInfoList(appReq.RechargeOrderSearch{PayStatus: &pending})
	if err != nil || total != 1 || len(rows) != 1 || rows[0].OrderNo != "pending" {
		t.Fatalf("status 0 filter failed: %+v %d %v", rows, total, err)
	}
	version := model.AppVersion{Version: "1.0", DownloadAddr: "https://example.com/app"}
	service := AppVersionService{}
	if err = service.CreateAppVersion(&version); err != nil {
		t.Fatal(err)
	}
	if err = service.DeleteAppVersion(version.UID); err != nil {
		t.Fatal(err)
	}
	var persisted model.AppVersion
	if err = db.Where("uid = ?", version.UID).First(&persisted).Error; err != nil || persisted.DeletedAt == nil {
		t.Fatalf("delete must retain row: %v", err)
	}
	_, total, err = service.GetAppVersionInfoList(appReq.AppVersionSearch{})
	if err != nil || total != 0 {
		t.Fatalf("deleted visible: %d %v", total, err)
	}
}
