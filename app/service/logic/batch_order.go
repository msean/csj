package logic

import (
	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"app/utils"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	BatchOrderLogic struct {
		context   *gin.Context
		runtime   *global.RunTime
		OwnerUser int64
	}
	BatchOrderGoodsLogic struct {
		context   *gin.Context
		runtime   *global.RunTime
		OwnerUser int64
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
func (logic *BatchOrderLogic) TempCreate(param request.CreateTempBatchOrderParam) (batchOrder response.BatchOrderRsp, err error) {
	if param.BatchUUID == "" {
		err = common.BatchUUIDRequireErr
		return
	}
	if len(param.GoodsList) == 0 {
		err = common.BatchOrderGoodsRequireErr
		return
	}

	batchOrder.BatchOrder = model.BatchOrder{
		BatchUUID: param.BatchUUIDCompatible,
		OwnerUser: logic.OwnerUser,
		UserUUID:  param.CustomerUUIDCompatible,
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
			GoodsUUID: goods.GoodsUUIDCompatible,
		})
	}

	batchOrder.TotalAmount = batchOrder.SetTotalAmount()
	batchOrder.CreditAmount = batchOrder.TotalAmount

	logic.runtime.DB.Create(&batchOrder.BatchOrder)

	logic.SetField(&batchOrder)
	return
}

// 下单
func (logic *BatchOrderLogic) Create(
	tx *gorm.DB,
	param request.CreateBatchOrderParam,
) (batchOrder response.BatchOrderRsp, err error) {
	if param.BatchUUIDCompatible == 0 {
		err = common.BatchUUIDRequireErr
		return
	}

	if len(param.GoodsList) == 0 {
		err = common.BatchOrderGoodsRequireErr
		return
	}
	batchOrder.BatchOrder = model.BatchOrder{
		BatchUUID: param.BatchUUIDCompatible,
		OwnerUser: logic.OwnerUser,
		UserUUID:  param.CustomerUUIDCompatible,
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
			GoodsUUID: goods.GoodsUUIDCompatible,
		})
	}

	batchOrder.TotalAmount = batchOrder.SetTotalAmount()
	batchOrder.CreditAmount = batchOrder.TotalAmount - param.FPayAmount
	if common.FloatEqual(batchOrder.TotalAmount, param.FPayAmount) || common.FloatGreat(param.FPayAmount, batchOrder.TotalAmount) {
		batchOrder.Status = model.BatchOrderFinish
	} else {
		batchOrder.Status = model.BatchOrderedCredit
	}

	if err = utils.GormCreateObj(tx, &batchOrder.BatchOrder); err != nil {
		return
	}

	logic.SetField(&batchOrder)
	return
}

func (logic *BatchOrderLogic) Shared(uuid int64) (err error) {
	var batchOrder response.BatchOrderRsp
	if batchOrder, err = logic.FromUUID(uuid); err != nil {
		return
	}
	if err = dao.BatchOrder.UpdateShare(logic.runtime.DB, uuid, model.BatchOrderShared); err != nil {
		return
	}
	logic.Record(&batchOrder, true, model.HistoryStepOrderShare, model.PayField{})
	return
}

func (logic *BatchOrderLogic) UpdateStatus(param request.UpdateBatchOrderStatusParam) (batchOrder response.BatchOrderRsp, err error) {
	if batchOrder, err = logic.FromUUID(param.BatchOrderUUIDCompatible); err != nil {
		return
	}
	if err = dao.BatchOrder.UpdateStatus(logic.runtime.DB, param.BatchOrderUUIDCompatible, param.Status); err != nil {
		return
	}
	switch param.Status {
	case model.BatchOrderedCredit:
		logic.Record(&batchOrder, true, model.HistoryStepCredit, model.PayField{})
	case model.BatchOrderCancel, model.BatchOrderRefund:
		logic.Record(&batchOrder, true, model.HistoryStepCrash, model.PayField{})
	}
	return
}

func (logic *BatchOrderLogic) Update(param request.UpdateBatchOrderParam) (batchOrder response.BatchOrderRsp, err error) {
	if batchOrder, err = logic.FromUUID(param.BatchOrderUUIDCompatible); err != nil {
		return
	}
	tx := logic.runtime.DB.Begin()
	if err = tx.Unscoped().Delete(&model.BatchOrderGoods{}, "batch_order_uuid=?", param.BatchOrderUUID).Error; err != nil {
		tx.Rollback()
		return
	}
	batchOrder.UserUUID = param.CustomerUUIDCompatible
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
	if err = tx.Save(&batchOrder.BatchOrder).Error; err != nil {
		tx.Rollback()
		return
	}

	logic.Record(&batchOrder, false, model.HistoryStepOrderFix, model.PayField{})
	return
}

func (logic *BatchOrderLogic) FromUUID(uuid int64) (batchOrder response.BatchOrderRsp, err error) {

	if err = utils.GormFind(logic.runtime.DB.Preload("GoodsListRelated"),
		&batchOrder.BatchOrder,
		utils.WhereUIDCond(uuid),
	); err != nil {
		return
	}

	logic.LoadHistory(&batchOrder)
	logic.SetField(&batchOrder)
	return
}

func (logic *BatchOrderLogic) FindLatestGoods(goodsUUIDList []string) (goodsOrderList []BatchOrderGoodsOrder, err error) {
	for _, goodsUUID := range goodsUUIDList {
		conds := []utils.Cond{
			utils.NewWhereCond("goods_uuid", goodsUUID),
			utils.CreatedOrderAscCond(),
		}
		var goodsOrder model.BatchOrderGoods
		var _price float64
		if err = utils.GormFirst(logic.runtime.DB, &goodsOrder, conds...); err != nil {
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

func (logic *BatchOrderLogic) SetGoodsFeild(batchOrder *response.BatchOrderRsp) (err error) {
	var goodsUUIDList []int64
	for _, goods := range batchOrder.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM, e := dao.BatchGoodsFieldSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
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

func (logic *BatchOrderLogic) SetCustomerField(batchOrder *response.BatchOrderRsp) (err error) {
	batchOrder.CustomerField, err = dao.CustomerFieldSet(logic.runtime.DB, batchOrder.UserUUID, logic.OwnerUser)
	return
}

func (logic *BatchOrderLogic) List(param request.ListBatchOrderParam) (orderList []*response.BatchOrderRsp, err error) {
	var modelOrderList []model.BatchOrder

	conditions := []utils.Cond{
		utils.DefaultSetLimitCond(param.LimitCond),
		utils.NewWhereCond("owner_user", logic.OwnerUser),
	}
	if param.UserUUID != "" {
		conditions = append(conditions, utils.NewWhereCond("user_uuid", param.UserUUID))
	}

	if param.StartTime != 0 {
		conditions = append(conditions, utils.NewCmpCond("created_at", ">=", time.Unix(param.StartTime, 0)))
	}
	if param.EndTime != 0 {
		conditions = append(conditions, utils.NewCmpCond("created_at", "<=", time.Unix(param.EndTime, 0)))
	}
	if param.Status != 0 {
		conditions = append(conditions, utils.NewWhereCond("status", param.Status))
	}
	conditions = append(conditions, utils.CreatedOrderDescCond())
	if err = utils.GormFind(logic.runtime.DB.Preload("GoodsListRelated"), &modelOrderList, conditions...); err != nil {
		return
	}

	for _, modelBatchOrder := range modelOrderList {
		orderList = append(orderList, &response.BatchOrderRsp{
			BatchOrder: modelBatchOrder,
		})
	}

	logic.BatchField(orderList)
	return
}

func (logic *BatchOrderLogic) BatchField(orderList []*response.BatchOrderRsp) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var _cl []int64
		for _, o := range orderList {
			_cl = append(_cl, o.UserUUID)
		}
		_cm, _ := dao.BatchCustomerFieldSet(logic.runtime.DB, _cl, logic.OwnerUser)

		for _, o := range orderList {
			o.CustomerField = _cm[o.UserUUID]
		}
		wg.Done()
	}()
	go func() {
		var _gl []int64
		for _, o := range orderList {
			for _, g := range o.GoodsListRelated {
				_gl = append(_gl, g.GoodsUUID)
			}
		}
		_gm, _ := dao.BatchGoodsFieldSet(logic.runtime.DB, _gl, logic.OwnerUser)
		for _, o := range orderList {
			for _, g := range o.GoodsListRelated {
				g.GoodsField = _gm[g.GoodsUUID]
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

// func (logic *BatchOrderLogic) BatchSetFeilds(orders []model.BatchOrder) {
// 	wg := sync.WaitGroup{}
// 	wg.Add(2)
// 	go func() {
// 		logic.SetCustomerFeild()
// 		wg.Done()
// 	}()
// 	go func() {
// 		logic.SetGoodsFeild()
// 		wg.Done()
// 	}()
// 	wg.Wait()
// }

func (logic *BatchOrderLogic) SetField(batchOrder *response.BatchOrderRsp) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		logic.SetCustomerField(batchOrder)
		wg.Done()
	}()
	go func() {
		logic.SetGoodsFeild(batchOrder)
		wg.Done()
	}()
	wg.Wait()
}

func (logic *BatchOrderLogic) Record(batchOrder *response.BatchOrderRsp, load bool, stepType int32, pay model.PayField) {
	func() {
		if load {
			if err := utils.GormFind(logic.runtime.DB.Preload("GoodsListRelated"), &batchOrder.BatchOrder); err != nil {
				logic.runtime.Logger.Error(fmt.Sprintf("BatchOrderLogic Record: %s", err))
				return
			}

			logic.SetField(batchOrder)
		}
		dao.BatchOrder.Record(logic.runtime.DB, *batchOrder, stepType, pay)
	}()
}

func (logic *BatchOrderLogic) LoadHistory(batchOrder *response.BatchOrderRsp) (err error) {
	var history model.BatchOrderHistory
	if err = utils.GormFirst(logic.runtime.DB, &history, utils.NewWhereCond("batch_order_uuid", batchOrder.UID)); err != nil {
		global.Global.Logger.Error(fmt.Sprintf("BatchOrderLogic LoadHistory err", err))
		return
	}
	if history.UID != 0 {
		batchOrder.History = history
	}
	return
}
