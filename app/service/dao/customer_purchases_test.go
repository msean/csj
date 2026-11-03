package dao

import (
	"app/service/model"
	"gorm.io/gorm/schema"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestPurchaseSummaryUsesActualModelColumn(t *testing.T) {
	parsed, err := schema.Parse(&model.BatchOrder{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatal(err)
	}
	field := parsed.LookUpField("TotalAmount")
	if field == nil || field.DBName != "amount" {
		t.Fatalf("unexpected amount column: %+v", field)
	}
	source, err := os.ReadFile("customer_purchases.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(source), "SUM("+field.DBName+" + COALESCE(tricycle_fee, 0))") {
		t.Fatal("purchase summary does not use the mapped database amount column")
	}
}
