package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	BatchOrderLogic struct {
		context *gin.Context
		runtime *global.RunTime
		// model.BatchOrder
		OwnerUser string
		History   model.BatchOrderHistory `json:"history"`
	}
	BatchOrderGoodsLogic struct {
		context *gin.Context
		runtime *global.RunTime
		*model.BatchOrderGoods
	}
	BatchOrderGoodsOrder struct {
		GoodsUUID string  `json:"goodsUUID"`
		Price     float64 `json:"price"`
	}
)

func NewBatchOrderLogic(context *gin.Context) *BatchOrderLogic {
	logic := &BatchOrderLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func NewBatchOrderGoodsLogic(context *gin.Context) *BatchOrderGoodsLogic {
	logic := &BatchOrderGoodsLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

// 码单
func (logic *BatchOrderLogic) TempCreate(param request.CreateTempBatchOrderParam) (batchOrder model.BatchOrder, err error) {
	if param.BatchUUID == "" {
		err = common.BatchUUIDRequireErr
		return
	}
	if len(param.GoodsList) == 0 {
		err = common.BatchOrderGoodsRequireErr
		return
	}

	batchOrder = model.BatchOrder{
		BatchUUID: param.BatchUUID,
		OwnerUser: logic.OwnerUser,
		UserUUID:  param.CustomerUUID,
		Shared:    model.BatchOrderUnshare,
		Status:    model.BatchOrderTemp,
	}

	for _, goods := range param.GoodsList {
		batchOrder.GoodsListRelated = append(batchOrder.GoodsListRelated, &model.BatchOrderGoods{
			OwnerUser: logic.OwnerUser,
			BatchUUID: batchOrder.BatchUUID,
			UserUUID:  batchOrder.UserUUID,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
			SerialNo:  goods.SerialNo,
		})
	}

	batchOrder.TotalAmount = batchOrder.SetTotalAmount()
	batchOrder.CreditAmount = batchOrder.TotalAmount

	if err = model.CreateObj(logic.runtime.DB, &batchOrder); err != nil {
		return
	}
	logic.SetFeilds(batchOrder)
	return
}

// 下单
func (logic *BatchOrderLogic) Create(tx *gorm.DB, param request.CreateBatchOrderParam) (batchOrder model.BatchOrder, err error) {
	if param.BatchUUID == "" {
		err = common.BatchUUIDRequireErr
		return
	}

	if len(param.GoodsList) == 0 {
		err = common.BatchOrderGoodsRequireErr
		return
	}
	batchOrder = model.BatchOrder{
		BatchUUID: param.BatchUUID,
		OwnerUser: logic.OwnerUser,
		UserUUID:  param.CustomerUUID,
		Shared:    model.BatchOrderUnshare,
	}

	for _, goods := range param.GoodsList {
		batchOrder.GoodsListRelated = append(batchOrder.GoodsListRelated, &model.BatchOrderGoods{
			OwnerUser: logic.OwnerUser,
			BatchUUID: batchOrder.BatchUUID,
			UserUUID:  batchOrder.UserUUID,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
			SerialNo:  goods.SerialNo,
		})
	}

	batchOrder.TotalAmount = batchOrder.SetTotalAmount()
	batchOrder.CreditAmount = batchOrder.TotalAmount - param.FPayAmount
	if common.FloatEqual(batchOrder.TotalAmount, param.FPayAmount) || common.FloatGreat(param.FPayAmount, batchOrder.TotalAmount) {
		batchOrder.Status = model.BatchOrderFinish
	} else {
		batchOrder.Status = model.BatchOrderedCredit
	}

	if err = model.CreateObj(tx, &batchOrder); err != nil {
		return
	}
	logic.SetFeilds(batchOrder)
	return
}

func (logic *BatchOrderLogic) Shared() (err error) {
	if err = logic.BatchOrder.UpdateShare(logic.runtime.DB); err != nil {
		return
	}
	logic.Record(true, model.HistoryStepOrderShare, model.PayFeild{})
	return
}

func (logic *BatchOrderLogic) UpdateStatus(param request.UpdateBatchOrderStatusParam) (batchOrder model.BatchOrder, err error) {
	if batchOrder, err = logic.FromUUID(param.BatchOrderUUID); err != nil {
		return
	}
	if err = batchOrder.UpdateStatus(logic.runtime.DB, param.Status); err != nil {
		return
	}
	switch param.Status {
	case model.BatchOrderedCredit:
		logic.Record(batchOrder, true, model.HistoryStepCredit, model.PayFeild{})
	case model.BatchOrderCancel, model.BatchOrderRefund:
		logic.Record(batchOrder, true, model.HistoryStepCrash, model.PayFeild{})
	}
	return
}

func (logic *BatchOrderLogic) Update(param request.UpdateBatchOrderParam) (batchOrder model.BatchOrder, err error) {
	if batchOrder, err = logic.FromUUID(param.BatchOrderUUID); err != nil {
		return
	}
	tx := logic.runtime.DB.Begin()
	if err = tx.Delete(&model.BatchOrderGoods{}, "batch_order_uuid=?", param.BatchOrderUUID).Error; err != nil {
		tx.Rollback()
		return
	}
	batchOrder.UserUUID = param.CustomerUUID
	for _, goods := range param.GoodsList {
		batchOrder.GoodsListRelated = append(batchOrder.GoodsListRelated, &model.BatchOrderGoods{
			OwnerUser:     logic.OwnerUser,
			BatchUUID:     batchOrder.BatchUUID,
			UserUUID:      batchOrder.UserUUID,
			BatchOrderUID: batchOrder.UID,
			Price:         goods.Price,
			Mount:         goods.Mount,
			Weight:        goods.Weight,
			SerialNo:      goods.SerialNo,
		})

	}
	if err = tx.Save(&batchOrder).Error; err != nil {
		tx.Rollback()
		return
	}

	logic.SetFeilds(batchOrder)
	logic.Record(batchOrder, false, model.HistoryStepOrderFix, model.PayFeild{})
	return
}

func (logic *BatchOrderLogic) FromUUID(uuid string) (batchOrder model.BatchOrder, err error) {
	if err = model.Find(logic.runtime.DB.Preload("GoodsListRelated"), &batchOrder); err != nil {
		return
	}
	logic.LoadHistory()
	logic.SetFeilds(batchOrder)
	return
}

func (logic *BatchOrderLogic) FindLatestGoods(goodsUUIDList []string) (goodsOrderList []BatchOrderGoodsOrder, err error) {
	for _, goodsUUID := range goodsUUIDList {
		conds := []model.Cond{
			model.NewWhereCond("goods_uuid", goodsUUID),
			model.CreatedOrderAscCond(),
		}
		var goodsOrder model.BatchOrderGoods
		var _price float64
		if err = model.First(logic.runtime.DB, &goodsOrder, conds...); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = nil
			} else {
				return
			}
		}
		_price = goodsOrder.Price
		goodsOrderList = append(goodsOrderList, BatchOrderGoodsOrder{
			GoodsUUID: goodsUUID,
			Price:     _price,
		})
	}

	return
}

func (logic *BatchOrderLogic) SetGoodsFeild(batchOrder model.BatchOrder) (err error) {
	var goodsUUIDList []string
	for _, goods := range batchOrder.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM, e := model.BatchGoodsFeildSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
	if e != nil {
		err = e
		return
	}

	for _, goods := range batchOrder.GoodsListRelated {
		goods.GoodsName = goodsM[goods.GoodsUUID].GoodsName
		goods.GoodsTyp = goodsM[goods.GoodsUUID].GoodsTyp
	}

	return
}

func (logic *BatchOrderLogic) SetCustomerFeild(batchOrder model.BatchOrder) (err error) {
	batchOrder.CustomerFeild, err = model.CustomerFeildSet(logic.runtime.DB, batchOrder.UserUUID, logic.OwnerUser)
	return
}

func (logic *BatchOrderLogic) List(userUUID string, startTime, endTime int64, status int32, limitCond model.LimitCond) (orderList []*model.BatchOrder, err error) {

	conds := []model.Cond{
		limitCond,
		model.NewWhereCond("owner_user", logic.OwnerUser),
	}
	if userUUID != "" {
		conds = append(conds, model.NewWhereCond("user_uuid", userUUID))
	}

	if startTime != 0 {
		conds = append(conds, model.NewCmpCond("created_at", ">=", time.Unix(startTime, 0)))
	}
	if endTime != 0 {
		conds = append(conds, model.NewCmpCond("created_at", "<=", time.Unix(endTime, 0)))
	}
	if status != 0 {
		conds = append(conds, model.NewWhereCond("status", status))
	}
	conds = append(conds, model.CreatedOrderDescCond())
	if err = model.Find(logic.runtime.DB.Preload("GoodsListRelated"), &orderList, conds...); err != nil {
		return
	}

	logic.BatchFeilds(orderList)
	return
}

func (logic *BatchOrderLogic) BatchFeilds(orderList []*model.BatchOrder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var _cl []string
		for _, o := range orderList {
			_cl = append(_cl, o.UserUUID)
		}
		_cm, _ := model.BatchCustomerFeildSet(logic.runtime.DB, _cl, logic.OwnerUser)

		for _, o := range orderList {
			o.CustomerFeild = _cm[o.UserUUID]
		}
		wg.Done()
	}()
	go func() {
		var _gl []string
		for _, o := range orderList {
			for _, g := range o.GoodsListRelated {
				_gl = append(_gl, g.GoodsUUID)
			}
		}
		_gm, _ := model.BatchGoodsFeildSet(logic.runtime.DB, _gl, logic.OwnerUser)
		for _, o := range orderList {
			for _, g := range o.GoodsListRelated {
				g.GoodsFeild = _gm[g.GoodsUUID]
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

func (logic *BatchOrderLogic) BatchSetFeilds(orders []model.BatchOrder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		logic.SetCustomerFeild()
		wg.Done()
	}()
	go func() {
		logic.SetGoodsFeild()
		wg.Done()
	}()
	wg.Wait()
}

func (logic *BatchOrderLogic) SetFeilds(batchOrder model.BatchOrder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		logic.SetCustomerFeild(batchOrder)
		wg.Done()
	}()
	go func() {
		logic.SetGoodsFeild(batchOrder)
		wg.Done()
	}()
	wg.Wait()
}

func (logic *BatchOrderLogic) Record(batchOrder model.BatchOrder, loadself bool, stepType int32, pay model.PayFeild) {
	go func() {
		if loadself {
			if err := model.Find(logic.runtime.DB.Preload("GoodsListRelated"), &batchOrder); err != nil {
				logic.runtime.Logger.Error(fmt.Sprintf("BatchOrderLogic Record: %s", err))
				return
			}
			logic.SetFeilds(batchOrder)
		}
		batchOrder.Record(logic.runtime.DB, stepType, pay)
	}()
	return
}

func (logic *BatchOrderLogic) LoadHistory() (err error) {
	var history model.BatchOrderHistory
	if err = model.First(logic.runtime.DB, &history, model.NewWhereCond("batch_order_uuid", logic.UID)); err != nil {
		global.Global.Logger.Error(fmt.Sprintf("BatchOrderLogic LoadHistory err", err))
		return
	}
	if history.UID != "" {
		logic.History = history
	}
	return
}
