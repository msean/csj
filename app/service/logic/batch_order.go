package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		Price     float64 `json:"price"`
	}
)

func NewBatchOrderLogic(context *gin.Context) *BatchOrderLogic {
	logic := &BatchOrderLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	fmt.Println(">>>>>>>>>>>>>logic.OwnerUser", logic.OwnerUser)
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

	// 在创建前加这一行
	if err = logic.FillGoodsWeightWithTotal(logic.runtime.DB); err != nil {
		return
	}

	// 创建订单
	if err = model.CreateObj(logic.runtime.DB, &logic.BatchOrder); err != nil {
		global.Global.Logger.Error("BatchOrderLogic TempCreate CreateObj", zap.Any("logic.BatchOrder", logic.BatchOrder), zap.Error(err))
		return
	}

	logic.SetFeilds()
	return
}

func (logic *BatchOrderLogic) FillGoodsWeightWithTotal(db *gorm.DB) error {
	list := logic.GoodsListRelated

	if len(list) == 0 {
		return nil
	}

	// 1. 收集 goodsUUID
	var goodsUUIDs []string
	for _, item := range list {
		if item.GoodsUUID != "" {
			goodsUUIDs = append(goodsUUIDs, item.GoodsUUID)
			item.BatchUUID = logic.BatchUUID
			item.OwnerUser = logic.OwnerUser
			item.UserUUID = logic.UserUUID
		}
	}

	// 2. 一次查库
	var goodsList []model.Goods
	if err := db.
		Where("uid IN ?", goodsUUIDs).
		Find(&goodsList).Error; err != nil {
		return err
	}

	// 3. 转 map
	goodsMap := make(map[string]model.Goods)
	for _, g := range goodsList {
		goodsMap[g.UID] = g
	}

	var totalAmount float64

	// 4. 计算 weight
	for _, item := range list {
		goods, ok := goodsMap[item.GoodsUUID]
		if !ok {
			continue
		}

		item.GoodsTyp = goods.Typ
		var subTotal float64
		// 定装：weight = mount * goods.weight
		if goods.Typ == common.GoodsTypeFix {
			item.Weight = utils.FloatReserve(float64(item.Mount)*goods.Weight, 1)
			subTotal = utils.FloatReserve(float64(item.Mount)*item.Price, 0)
		}

		// 将散装设置成0
		if goods.Typ == common.GoodsTypeBulk {
			item.Mount = 0
			subTotal = utils.FloatReserve(item.Price*item.Weight, 0)
		}
		item.Total = subTotal
		totalAmount += subTotal
	}
	logic.TotalAmount = totalAmount

	return nil
}

// 下单
func (logic *BatchOrderLogic) Create(tx *gorm.DB) (err error) {
	if tx == nil {
		tx = logic.runtime.DB.Begin()
		defer tx.Commit()
	}
	logic.BatchOrder.Shared = common.BatchOrderUnshare

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
	case common.BatchOrderTemp:
		logic.Record(true, model.HistoryStepCredit, model.PayFeild{})
	case common.BatchOrderCancel, common.BatchOrderRefund:
		logic.Record(true, model.HistoryStepCrash, model.PayFeild{})
	}
	return
}

// 更新单次
func (logic *BatchOrderLogic) Update(tx *gorm.DB) (err error) {
	if err = tx.Delete(&model.BatchOrderGoods{}, "batch_order_uuid=?", logic.UID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Delete(&model.BatchOrderPay{}, "batch_order_uuid=?", logic.UID).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
		goods.BatchOrderUID = logic.UID
	}
	if err = tx.Omit("created_at").Save(&logic.BatchOrder).Error; err != nil {
		tx.Rollback()
		return
	}

	// logic.SetFeilds()
	// logic.Record(false, model.HistoryStepOrderFix, model.PayFeild{})
	return
}

func (logic *BatchOrderLogic) FromUUID(uuid string) (err error) {
	if err = model.Find(logic.runtime.DB.Preload("GoodsListRelated"), &logic.BatchOrder, model.WhereUIDCond(uuid)); err != nil {
		return
	}
	logic.LoadHistory()
	logic.SetFeilds()
	// logic.SetTotal()
	return
}

func (logic *BatchOrderLogic) LoadSingle(uuid string) (batchOrder model.BatchOrder, err error) {
	err = model.Find(logic.runtime.DB, &batchOrder, model.WhereUIDCond(uuid))
	// logic.SetTotal()
	return
}

func (logic *BatchOrderLogic) FindLatestGoods(goodsUUIDList []string) (goodsOrderList []BatchOrderGoodsOrder, err error) {
	for _, goodsUUID := range goodsUUIDList {
		conds := []model.Cond{
			model.NewWhereCond("goods_uuid", goodsUUID),
			model.CreatedOrderAscCond(),
		}
		// todo BatchOrderGoods需要加索引
		var goodsOrder model.BatchOrderGoods
		var _price float64
		if err = model.First(logic.runtime.DB, &goodsOrder, conds...); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var goods model.Goods
				if err = model.First(logic.runtime.DB, &goods, model.NewWhereCond("uid", goodsUUID)); err != nil {
					return
				}
				_price = goods.Price
			} else {
				return
			}
		} else {
			_price = goodsOrder.Price
		}
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
		goods.GoodsWeight = goodsM[goods.GoodsUUID].GoodsWeight
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
	} else {
		// conds = append(conds, model.NewInCondFromInt("status", common.ExCludeTempBatchOrder))
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
				g.SetTotal()
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

// func (logic *BatchOrderLogic) SetTotal() {
// 	for _, g := range logic.GoodsListRelated {
// 		logic.TotalAmount += g.SetTotal()
// 	}
// }

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
}

func (logic *BatchOrderLogic) LoadHistory() (err error) {
	var history model.BatchOrderHistory
	if err = model.First(logic.runtime.DB, &history, model.NewWhereCond("batch_order_uuid", logic.UID)); err != nil {
		logic.runtime.Logger.Error("BatchOrderLogic LoadHistory", zap.Error(err))
		return
	}
	if history.UID != "" {
		logic.History = history
	}
	return
}

func (logic *BatchOrderLogic) GoodsList(req request.BatchGoodsListReq) (rsp []*response.BatchGoodsGroupItem, err error) {
	rsp = make([]*response.BatchGoodsGroupItem, 0)
	var storages []model.BatchOrderGoods
	db := logic.runtime.DB

	conds := []model.Cond{
		model.NewWhereCond("batch_uuid", req.BatchUUID),
		model.NewWhereCond("goods_uuid", req.GoodsUUID),
	}

	conds = append(conds, model.CreatedOrderDescCond())
	if err = model.Find(db, &storages, conds...); err != nil {
		logic.runtime.Logger.Error("BatchLogic List FindBatch", zap.Any("req", req), zap.Error(err))
		return
	}

	var customerUUIDList []string
	for _, storage := range storages {
		customerUUIDList = append(customerUUIDList, storage.UserUUID)
	}

	_cm, _ := model.BatchCustomerFeildSet(logic.runtime.DB, customerUUIDList, logic.OwnerUser)

	for _, storage := range storages {
		rsp = append(rsp, &response.BatchGoodsGroupItem{
			SellTime:     storage.CreatedAt.Format("2006-01-02 15:04:05"),
			SellPrice:    utils.FloatReserveStr(storage.Price, 1), //
			SellAmount:   utils.FloatReserveStr(storage.Amount(), 1),
			CustomerName: _cm[storage.UserUUID].CustomerName,
			SellTotall:   utils.FloatReserveStr(float64(storage.Price)*storage.Amount(), 1),
			GoodType:     storage.GType(),
		})
	}

	return
}
