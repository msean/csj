package response

import (
	"app/service/model"
	"app/service/model/request"
	"encoding/json"
	"testing"
)

func TestBatchBasicDetailOmitsStatsAndPreservesGoods(t *testing.T) {
	good := &model.BatchGoods{GoodsUUID: "goods", Price: 12.5, Weight: 26, Mount: 10}
	good.GoodsFeild = model.GoodsFeild{GoodsName: "青菜", GoodsTyp: 1, GoodsWeight: 26}
	good.Surplus = "8"
	good.SerialNo = "batch-serial"
	batch := model.Batch{GoodsListRelated: []*model.BatchGoods{good}, SerialID: 7}
	batch.UID = "batch"
	batch.StatSalesAmount = "25"
	for _, basic := range []bool{false, true} {
		var value any = batch
		if basic {
			value = NewBatchBasicDetail(batch)
		}
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var got map[string]any
		if err = json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"statMount", "statWeight", "statInventoryMount", "statInventoryWeight", "statSalesAmount", "statSellProfit", "statSellWeight", "statSellMount"} {
			if _, ok := got[key]; ok == basic {
				t.Fatalf("basic=%v unexpected presence of %s: %s", basic, key, data)
			}
		}
		row := got["goodsList"].([]any)[0].(map[string]any)
		for _, key := range []string{"sellAmount", "sellAvgPrice", "sellTotalAmount", "sellProfit"} {
			if _, ok := row[key]; ok == basic {
				t.Fatalf("basic=%v unexpected goods field %s", basic, key)
			}
		}
		if row["serialNo"] != "batch-serial" {
			t.Fatalf("serialNo missing or changed for legacy clients: %s", data)
		}
		if row["surplus"] != "8" {
			t.Fatalf("missing stock: %s", data)
		}
		if got["uuid"] != "batch" || row["goodsUUID"] != "goods" || row["price"] != 12.5 || row["weight"] != float64(26) || row["type"] != float64(1) || row["name"] != "青菜" {
			t.Fatalf("lost order fields: %s", data)
		}
	}
	empty, _ := json.Marshal(NewBatchBasicDetail(model.Batch{}))
	var got map[string]any
	json.Unmarshal(empty, &got)
	if rows, ok := got["goodsList"].([]any); !ok || len(rows) != 0 {
		t.Fatalf("empty goods must be array: %s", empty)
	}
}

func TestBatchDetailStatsDefaultRemainsEnabled(t *testing.T) {
	for _, raw := range []string{`{}`, `{"skipStats":false}`, `{"skipStats":true}`} {
		var req request.BatchDetailReq
		if err := json.Unmarshal([]byte(raw), &req); err != nil {
			t.Fatal(err)
		}
		if req.SkipStats != (raw == `{"skipStats":true}`) {
			t.Fatalf("incorrect stats option for %s", raw)
		}
	}
}

func TestBatchGoodsSerialNoAlwaysSerializedForLegacyClients(t *testing.T) {
	data, err := json.Marshal(model.BatchGoods{})
	if err != nil {
		t.Fatal(err)
	}
	var row map[string]any
	if err := json.Unmarshal(data, &row); err != nil {
		t.Fatal(err)
	}
	if value, ok := row["serialNo"]; !ok || value != "" {
		t.Fatalf("old clients require serialNo even when empty: %s", data)
	}
	batch := model.Batch{SerialNo: "parent", GoodsListRelated: []*model.BatchGoods{{}}}
	if got := NewBatchBasicDetail(batch).GoodsList[0].SerialNo; got != "parent" {
		t.Fatalf("missing parent fallback: %s", got)
	}
}
