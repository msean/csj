package response

import (
	"app/service/model"
	"encoding/json"
	"testing"
)

func TestCustomerRspModel(ts *testing.T) {
	_model := ListCustomerRsp{
		Customer: model.Customer{
			Name:  "测试",
			Phone: "1234567890",
			Debt:  100.50,
			Addr:  "测试地址",
			CarNo: "测试车牌",
		},
		LatestBillDate: 20230405,
	}

	_model.UID = 1011
	_bs, err := json.Marshal(_model)

	if err != nil {
		ts.Fatal("Failed to marshal:", err)
	}

	ts.Log(string(_bs))
}
