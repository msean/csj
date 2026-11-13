package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestBatchInventoryQueryScopesSalesAndExcludesRevenue(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "user:password@tcp(localhost:3306)/test", SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return BatchInventoryUsageQuery(tx, "owner", "batch").Find(&[]BatchInventoryUsage{})
	})
	for _, want := range []string{"SUM(bog.mount)", "SUM(bog.weight)", "GROUP BY", "bog.owner_user = 'owner'", "bo.owner_user = 'owner'", "bog.batch_uuid = 'batch'", "bo.status IN (1,3)", "bog.deleted_at IS NULL", "bo.deleted_at IS NULL"} {
		if !strings.Contains(sql, want) {
			t.Fatalf("missing %s in %s", want, sql)
		}
	}
	if strings.Contains(sql, "price") || strings.Contains(sql, "total") {
		t.Fatalf("inventory query reads monetary fields: %s", sql)
	}
}
