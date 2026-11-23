package logic

import (
	"app/service/model"
	"app/service/model/response"
	"testing"
	"time"
)

func TestRepaymentAllocationDecoding(t *testing.T) {
	for _, raw := range []string{"", "null", "broken", `{"unexpected":true}`} {
		if len(decodeRepaymentAllocations(raw)) != 0 {
			t.Fatalf("invalid payload accepted: %s", raw)
		}
	}
	for _, raw := range []string{`[{"type":2,"orderUUID":"order","amount":"12.50"}]`, `{"details":[{"type":2,"orderUUID":"order","amount":12.5}]}`} {
		rows := decodeRepaymentAllocations(raw)
		if len(rows) != 1 || rows[0].Amount != 12.5 || rows[0].OrderUUID != "order" {
			t.Fatalf("bad decode: %+v", rows)
		}
	}
}
func TestRepaymentOrderEnrichmentPreservesAllocation(t *testing.T) {
	rows := []response.RepaymentAllocation{{Type: 2, OrderUUID: "one", Amount: 20}, {Type: 1, Amount: 5}, {Type: 2, OrderUUID: "missing", Amount: 3}}
	order := model.BatchOrder{TotalAmount: 100.5, TricycleFee: 2, OrderTime: 1700000000, Status: 1}
	order.UID = "one"
	order.TicketNo = "01234567"
	enrichRepaymentAllocations(rows, []model.BatchOrder{order})
	if rows[0].TicketNo != "01234567" || rows[0].Amount != 20 || !rows[0].OrderAvailable || rows[0].OrderTotalAmount == nil || *rows[0].OrderTotalAmount != 102.5 {
		t.Fatalf("bad enrichment: %+v", rows[0])
	}
	if rows[0].OrderTime == nil || !rows[0].OrderTime.Equal(time.Unix(1700000000, 0)) {
		t.Fatal("order date missing")
	}
	for _, row := range rows[1:] {
		if row.OrderAvailable || row.OrderTime != nil || row.OrderTotalAmount != nil {
			t.Fatal("invented metadata for debt/missing order")
		}
	}
}

func TestRepaymentTicketSurvivesDeletedOrder(t *testing.T) {
	pay := model.BatchOrderPay{BatchOrderUUID: "one", Amount: 20, PayDetails: "{}"}
	rows := paymentAllocationRows(&pay)
	if len(rows) != 1 || rows[0].OrderUUID != "one" {
		t.Fatalf("missing allocation for linked order: %+v", rows)
	}
	deleted := time.Now()
	order := model.BatchOrder{TicketNo: "01234567"}
	order.UID = "one"
	order.DeletedAt = &deleted
	enrichRepaymentAllocations(rows, []model.BatchOrder{order})
	if rows[0].TicketNo != "01234567" || rows[0].OrderAvailable {
		t.Fatalf("bad allocation for deleted order: %+v", rows[0])
	}
}
