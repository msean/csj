package logic

import (
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"fmt"
)

type CustomerPurchasesRsp struct {
	Items  []*model.BatchOrder `json:"items"`
	Total  int64               `json:"total"`
	Amount float64             `json:"amount"`
}

func (logic *OrderLogic) CustomerPurchases(req request.CustomerPurchasesReq) (rsp CustomerPurchasesRsp, err error) {
	if req.EndTime > 0 && req.StartTime >= req.EndTime {
		return rsp, fmt.Errorf("时间范围无效")
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageCount < 1 {
		req.PageCount = 10
	}
	if req.PageCount > 100 {
		req.PageCount = 100
	}
	var customer model.Customer
	if err = logic.runtime.DB.Where("uid = ? AND owner_user = ? AND deleted_at IS NULL", req.UserUUID, logic.OwnerUser).First(&customer).Error; err != nil {
		return
	}
	rsp.Total, rsp.Amount, err = dao.CustomerPurchaseSummary(logic.runtime.DB, logic.OwnerUser, req.UserUUID, req.StartTime, req.EndTime, req.CreditOnly)
	if err != nil {
		return
	}
	rsp.Items, err = dao.CustomerPurchaseList(logic.runtime.DB, logic.OwnerUser, req.UserUUID, req.StartTime, req.EndTime, req.CreditOnly, (req.Page-1)*req.PageCount, req.PageCount)
	if err == nil {
		logic.BatchFeilds(rsp.Items)
	}
	return
}
