package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"
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
		model.BatchOrder
		History model.BatchOrderHistory `json:"history"`
	}
	BatchOrderGoodsLogic struct {
		context *gin.Context
		runtime *global.RunTime
		*model.BatchOrderGoods
	}
	BatchOrderGoodsOrder struct {
		GoodsUUID string  `json:"goodsUUID"`
		Price     float32 `json:"price"`
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
func (logic *BatchOrderLogic) TempCreate() (err error) {
	logic.BatchOrder.DefaultSet()

	if logic.BatchUUID == "" {
		return common.BatchUUIDRequireErr
	}
	if len(logic.GoodsListRelated) == 0 {
		return common.BatchOrderGoodsRequireErr
	}

	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
	}

	if err = model.CreateObj(logic.runtime.DB, &logic.BatchOrder); err != nil {
		return
	}
	logic.SetFeilds()
	return
}

// 下单
func (logic *BatchOrderLogic) Create(tx *gorm.DB) (err error) {
	if tx == nil {
		tx = logic.runtime.DB.Begin()
		defer tx.Commit()
	}
	logic.BatchOrder.Shared = model.BatchOrderUnshare

	if logic.BatchUUID == "" {
		return common.BatchUUIDRequireErr
	}
	if len(logic.GoodsListRelated) == 0 {
		return common.BatchOrderGoodsRequireErr
	}

	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
	}

	if err = model.CreateObj(tx, &logic.BatchOrder); err != nil {
		tx.Rollback()
		return
	}
	logic.SetFeilds()
	return
}

func (logic *BatchOrderLogic) Shared() (err error) {
	if err = logic.BatchOrder.UpdateShare(logic.runtime.DB); err != nil {
		return
	}
	logic.Record(true, model.HistoryStepOrderShare, model.PayFeild{})
	return
}

func (logic *BatchOrderLogic) UpdateStatus() (err error) {
	if err = logic.BatchOrder.UpdateStatus(logic.runtime.DB, logic.Status); err != nil {
		return
	}
	switch logic.Status {
	case model.BatchOrderedCredit:
		logic.Record(true, model.HistoryStepCredit, model.PayFeild{})
	case model.BatchOrderCancel, model.BatchOrderRefund:
		logic.Record(true, model.HistoryStepCrash, model.PayFeild{})
	}
	return
}

func (logic *BatchOrderLogic) Update() (err error) {
	tx := logic.runtime.DB.Begin()
	if err = tx.Delete(&model.BatchOrderGoods{}, "batch_order_uuid=?", logic.UID).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
		goods.BatchOrderUID = logic.UID
	}
	if err = tx.Save(&logic.BatchOrder).Error; err != nil {
		tx.Rollback()
		return
	}

	logic.SetFeilds()
	logic.Record(false, model.HistoryStepOrderFix, model.PayFeild{})
	return
}

func (logic *BatchOrderLogic) FromUUID() (err error) {
	if err = model.Find(logic.runtime.DB.Preload("GoodsListRelated"), &logic.BatchOrder); err != nil {
		return
	}
	logic.LoadHistory()
	logic.SetFeilds()
	return
}

func (logic *BatchOrderLogic) FindLatestGoods(goodsUUIDList []string) (goodsOrderList []BatchOrderGoodsOrder, err error) {
	for _, goodsUUID := range goodsUUIDList {
		conds := []model.Cond{
			model.NewWhereCond("goods_uuid", goodsUUID),
			model.CreatedOrderAscCond(),
		}
		var goodsOrder model.BatchOrderGoods
		var _price float32
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

func (logic *BatchOrderLogic) SetGoodsFeild() (err error) {
	var goodsUUIDList []string
	for _, goods := range logic.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM, e := model.BatchGoodsFeildSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
	if e != nil {
		err = e
		return
	}

	for _, goods := range logic.GoodsListRelated {
		goods.GoodsName = goodsM[goods.GoodsUUID].GoodsName
		goods.GoodsTyp = goodsM[goods.GoodsUUID].GoodsTyp
	}

	return
}

func (logic *BatchOrderLogic) SetCustomerFeild() (err error) {
	logic.CustomerFeild, err = model.CustomerFeildSet(logic.runtime.DB, logic.UserUUID, logic.OwnerUser)
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

func (logic *BatchOrderLogic) SetFeilds() {
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

func (logic *BatchOrderLogic) Record(loadself bool, stepType int32, pay model.PayFeild) {
	go func() {
		if loadself {
			if err := model.Find(logic.runtime.DB.Preload("GoodsListRelated"), &logic.BatchOrder); err != nil {
				logic.runtime.Logger.Error(fmt.Sprintf("BatchOrderLogic Record: %s", err))
				return
			}
			logic.SetFeilds()
		}
		logic.BatchOrder.Record(logic.runtime.DB, stepType, pay)
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
