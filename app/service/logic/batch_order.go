package logic

import (
	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
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
func (logic *BatchOrderLogic) TempCreate(param request.CreateTempBatchOrderParam) (err error) {
	if param.BatchUUID == "" {
		err = common.BatchUUIDRequireErr
		return
	}
	if len(param.GoodsList) == 0 {
		err = common.BatchOrderGoodsRequireErr
		return
	}

	batchOrder := model.BatchOrder{
		BatchUUID: param.BatchUUIDCompatible,
		OwnerUser: logic.OwnerUser,
		UserUUID:  param.CustomerUUIDCompatible,
		Shared:    model.BatchOrderUnshare,
		Status:    model.BatchOrderTemp,
	}

	var batchOrderGoodsList []model.BatchOrderGoods
	for _, goods := range param.GoodsList {
		batchOrderGoodsList = append(batchOrderGoodsList, model.BatchOrderGoods{
			OwnerUser: logic.OwnerUser,
			BatchUUID: batchOrder.BatchUUID,
			UserUUID:  batchOrder.UserUUID,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
			SerialNo:  goods.SerialNo,
		})
	}

	batchOrder.TotalAmount = batchOrder.SetTotalAmount(batchOrderGoodsList)
	batchOrder.CreditAmount = batchOrder.TotalAmount

	db := logic.runtime.DB
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&batchOrder).Error; err != nil {
			return err
		}

		if err := tx.Create(&batchOrderGoodsList).Error; err != nil {
			return err
		}
		return nil
	})
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

	if err = utils.GormCreateObj(tx, &batchOrder); err != nil {
		return
	}
	logic.SetFeilds(&batchOrder)
	return
}

func (logic *BatchOrderLogic) Shared(uuid string) (batchOrder model.BatchOrder, err error) {
	if batchOrder, err = logic.FromUUID(uuid); err != nil {
		return
	}
	if err = dao.BatchOrder.UpdateShare(logic.runtime.DB, uuid, model.BatchOrderShared); err != nil {
		return
	}
	logic.Record(batchOrder, true, model.HistoryStepOrderShare, model.PayField{})
	return
}

func (logic *BatchOrderLogic) UpdateStatus(param request.UpdateBatchOrderStatusParam) (batchOrder model.BatchOrder, err error) {
	if batchOrder, err = logic.FromUUID(param.BatchOrderUUID); err != nil {
		return
	}
	if err = dao.BatchOrder.UpdateStatus(logic.runtime.DB, param.BatchOrderUUID, param.Status); err != nil {
		return
	}
	switch param.Status {
	case model.BatchOrderedCredit:
		logic.Record(batchOrder, true, model.HistoryStepCredit, model.PayField{})
	case model.BatchOrderCancel, model.BatchOrderRefund:
		logic.Record(batchOrder, true, model.HistoryStepCrash, model.PayField{})
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

	logic.SetFeilds(&batchOrder)
	logic.Record(batchOrder, false, model.HistoryStepOrderFix, model.PayField{})
	return
}

func (logic *BatchOrderLogic) FromUUID(uuid string) (batchOrder model.BatchOrder, err error) {
	if err = utils.GormFind(logic.runtime.DB.Preload("GoodsListRelated"), &batchOrder); err != nil {
		return
	}
	logic.LoadHistory(&batchOrder)
	logic.SetFeilds(&batchOrder)
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

func (logic *BatchOrderLogic) SetGoodsFeild(batchOrder *model.BatchOrder) (err error) {
	var goodsUUIDList []string
	for _, goods := range batchOrder.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM, e := model.BatchGoodsFieldSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
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

func (logic *BatchOrderLogic) SetCustomerField(batchOrder *model.BatchOrder) (err error) {
	batchOrder.CustomerField, err = model.CustomerFieldSet(logic.runtime.DB, batchOrder.UserUUID, logic.OwnerUser)
	return
}

func (logic *BatchOrderLogic) List(param request.ListBatchOrderParam) (orderList []*model.BatchOrder, err error) {

	conds := []utils.Cond{
		utils.DefaultSetLimitCond(param.LimitCond),
		utils.NewWhereCond("owner_user", logic.OwnerUser),
	}
	if param.UserUUID != "" {
		conds = append(conds, utils.NewWhereCond("user_uuid", param.UserUUID))
	}

	if param.StartTime != 0 {
		conds = append(conds, utils.NewCmpCond("created_at", ">=", time.Unix(param.StartTime, 0)))
	}
	if param.EndTime != 0 {
		conds = append(conds, utils.NewCmpCond("created_at", "<=", time.Unix(param.EndTime, 0)))
	}
	if param.Status != 0 {
		conds = append(conds, utils.NewWhereCond("status", param.Status))
	}
	conds = append(conds, utils.CreatedOrderDescCond())
	if err = utils.GormFind(logic.runtime.DB.Preload("GoodsListRelated"), &orderList, conds...); err != nil {
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
		_cm, _ := model.BatchCustomerFieldSet(logic.runtime.DB, _cl, logic.OwnerUser)

		for _, o := range orderList {
			o.CustomerField = _cm[o.UserUUID]
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

func (logic *BatchOrderLogic) SetFeilds(batchOrder *model.BatchOrder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		dao.SetCustomerField(batchOrder)
		wg.Done()
	}()
	go func() {
		logic.SetGoodsFeild(batchOrder)
		wg.Done()
	}()
	wg.Wait()
}

func (logic *BatchOrderLogic) Record(batchOrder model.BatchOrder, batchOrderGoods []model.BatchOrderGoods, stepType int32, pay model.PayField) {
	go func() {
		// if load {
		// 	if err := utils.GormFind(logic.runtime.DB, &batchOrder); err != nil {
		// 		logic.runtime.Logger.Error(fmt.Sprintf("BatchOrderLogic Record: %s", err))
		// 		return
		// 	}

		// 	logic.SetFeilds(&batchOrder)
		// }
		dao.BatchOrder.Record(logic.runtime.DB, batchOrder, batchOrderGoods, stepType, pay)
	}()
	return
}

func (logic *BatchOrderLogic) LoadHistory(batchOrder *model.BatchOrder) (err error) {
	var history model.BatchOrderHistory
	if err = utils.GormFirst(logic.runtime.DB, &history, utils.NewWhereCond("batch_order_uuid", batchOrder.UID)); err != nil {
		global.Global.Logger.Error(fmt.Sprintf("BatchOrderLogic LoadHistory err", err))
		return
	}
	if history.UID != "" {
		batchOrder.History = history
	}
	return
}
