package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/cache"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	BatchLogic struct {
		context *gin.Context
		runtime *global.RunTime
		model.Batch
	}
	BatchGoodsLogic struct {
		model.BatchGoods
		context *gin.Context
		runtime *global.RunTime
	}
)

func NewBatchLogic(context *gin.Context) *BatchLogic {
	logic := &BatchLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func NewBatchGoodsLogic(context *gin.Context) *BatchGoodsLogic {
	logic := &BatchGoodsLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *BatchLogic) CalSurplus() (surplusGoodsList []BatchGoodsLogic, err error) {
	db := logic.runtime.DB
	// 找到上一批
	var preBatch model.Batch
	conds := []utils.Cond{
		utils.NewOrderCond("serial_no desc"),
		utils.NewWhereCond("owner_user", logic.OwnerUser),
	}
	err = utils.Find(db, &preBatch, conds...)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	var preBatchGoodsList []model.BatchGoods
	if err = utils.Find(db, &preBatchGoodsList, utils.WhereOwnerUserCond(logic.OwnerUser), utils.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
		return
	}
	// key 是goodsUUID
	preGoodsList := []string{}
	preCountM := make(map[string]*BatchGoodsLogic)
	for _, preBatchGoods := range preBatchGoodsList {
		preCountM[preBatchGoods.GoodsUUID] = &BatchGoodsLogic{
			BatchGoods: preBatchGoods,
		}
		preGoodsList = append(preGoodsList, preBatchGoods.GoodsUUID)
	}

	// 上一批开单的
	var batchOrderedGoodsList []model.BatchOrderGoods
	if err = utils.Find(db, &batchOrderedGoodsList, utils.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
		return
	}
	for _, batchOrdered := range batchOrderedGoodsList {
		if _, InpreCountM := preCountM[batchOrdered.GoodsUUID]; InpreCountM {
			preCountM[batchOrdered.GoodsUUID].Mount -= batchOrdered.Mount
			preCountM[batchOrdered.GoodsUUID].Weight -= batchOrdered.Weight
			preCountM[batchOrdered.GoodsUUID].BatchGoodsStatFeild = model.BatchGoodsStatFeild{
				Mount:  preCountM[batchOrdered.GoodsUUID].Mount,
				Weight: preCountM[batchOrdered.GoodsUUID].Weight,
			}
			preCountM[batchOrdered.GoodsUUID].BatchGoodsStatFeild.Set()
		}
	}

	goodsM := make(map[string]model.GoodsFeild, 0)
	if goodsM, err = cache.GoodsCache.BatchGoodsFeildSet(preGoodsList, logic.OwnerUser); err != nil {
		return
	}

	for goodsUUID, batchGoods := range preCountM {
		batchGoods.GoodsName = goodsM[goodsUUID].GoodsName
		batchGoods.GoodsTyp = goodsM[goodsUUID].GoodsTyp
		surplusGoodsList = append(surplusGoodsList, *batchGoods)
	}
	return
}

func (logic *BatchLogic) Create() (err error) {
	var _b model.Batch
	var batchGoods []model.Goods
	var _goodsUIDList []string
	goodsUIDMap := make(map[string]model.Goods)

	for _, goods := range logic.GoodsListRelated {
		_goodsUIDList = append(_goodsUIDList, goods.GoodsUUID)
	}

	if err = utils.Find(logic.runtime.DB, &batchGoods, utils.InUIDCondFromString(_goodsUIDList)); err != nil {
		logic.runtime.Logger.Error("BatchLogic Create", zap.Any("_b", _b))
		return
	}

	logic.runtime.Logger.Debug("BatchLogic Create", zap.Any("batchGoods", batchGoods))
	for _, goods := range batchGoods {
		goodsUIDMap[goods.UID] = goods
	}
	logic.Default()
	for _, goods := range logic.GoodsListRelated {
		goods.SerialNo = logic.SerialNo
		goods.OwnerUser = logic.OwnerUser
		goods.GoodType = int(goodsUIDMap[goods.GoodsUUID].Typ)

	}

	return utils.CreateObj(logic.runtime.DB, &logic.Batch)
}

// withBatchGoods 是否更新批次下的物品信息
func (logic *BatchLogic) Update(withBatchGoods bool) (err error) {
	var storage model.Batch
	if storage, err = dao.BatchDao.FromUUID(logic.runtime.DB, logic.UID, false); err != nil {
		logic.runtime.Logger.Error("BatchLogic Update", zap.Error(err))
		return
	}
	logic.SerialID = storage.SerialID
	switch withBatchGoods {
	case true:
		var oldGoodsList []model.BatchGoods
		if oldGoodsList, err = dao.BatchGoodsDao.ListByBatchUUID(logic.runtime.DB, logic.UID); err != nil {
			logic.runtime.Logger.Error("BatchLogic Update", zap.Any("object", logic.Batch), zap.Error(err))
			return
		}
		// todo 开了单的不能删除
		// 2. 构建新 GoodsUUID 的 Map
		newGoodsUUIDMap := make(map[string]bool)
		for _, goods := range logic.GoodsListRelated {
			newGoodsUUIDMap[goods.GoodsUUID] = true
		}

		// 3. 找出要删除的 GoodsUUID（旧有但新无）
		var toDeleteGoodsUUIDs []string
		for _, oldGoods := range oldGoodsList {
			if !newGoodsUUIDMap[oldGoods.GoodsUUID] {
				toDeleteGoodsUUIDs = append(toDeleteGoodsUUIDs, oldGoods.GoodsUUID)
			}
		}
		// 开单检测
		if len(toDeleteGoodsUUIDs) > 0 {
			var orderGoodsList []model.BatchOrderGoods
			if orderGoodsList, err = dao.OrderDao.ListByBatchUUIDIn(logic.runtime.DB, logic.OwnerUser, logic.UID, toDeleteGoodsUUIDs); err != nil {
				logic.runtime.Logger.Error("查询开单记录失败", zap.Error(err))
				return
			}
			// 如果存在开单记录，则不允许删除
			if len(orderGoodsList) > 0 {
				// 收集已开单的 GoodsUUID
				orderedGoodsMap := make(map[string]bool)
				for _, og := range orderGoodsList {
					orderedGoodsMap[og.GoodsUUID] = true
				}

				var orderedGoodsNames []string
				var orderedGoodsUUIDs []string
				for _, goodsUUID := range toDeleteGoodsUUIDs {
					if orderedGoodsMap[goodsUUID] {
						// 这里可以根据 GoodsUUID 查询具体的货品名称
						orderedGoodsUUIDs = append(orderedGoodsUUIDs, goodsUUID)
					}
				}

				goodsM := make(map[string]model.GoodsFeild, 0)
				if goodsM, err = cache.GoodsCache.BatchGoodsFeildSet(orderedGoodsUUIDs, logic.OwnerUser); err != nil {
					return
				}

				for _, goods := range goodsM {
					orderedGoodsNames = append(orderedGoodsNames, goods.GoodsName)
				}

				return fmt.Errorf("以下货品已开单，不能删除: %v", orderedGoodsNames)
			}
		}

		tx := logic.runtime.DB.Begin()
		if err = tx.Delete(&model.BatchGoods{}, "batch_uuid=?", logic.UID).Error; err != nil {
			tx.Rollback()
			return err
		}

		var _goodsUIDList []string
		var batchGoods []model.Goods
		goodsUIDMap := make(map[string]model.Goods)

		for _, goods := range logic.GoodsListRelated {
			_goodsUIDList = append(_goodsUIDList, goods.GoodsUUID)
		}

		if err = utils.Find(logic.runtime.DB, &batchGoods, utils.InUIDCondFromString(_goodsUIDList)); err != nil {
			logic.runtime.Logger.Error("BatchLogic Create", zap.Any("_b", logic.Batch))
			return
		}

		for _, goods := range batchGoods {
			goodsUIDMap[goods.UID] = goods
		}
		logic.runtime.Logger.Debug("BatchLogic goodsUIDMap", zap.Any("goodsUIDMap", goodsUIDMap))
		for _, goods := range logic.GoodsListRelated {
			goods.SerialNo = logic.SerialNo
			goods.OwnerUser = logic.OwnerUser
			goods.BatchUUID = logic.UID
			goods.GoodType = int(goodsUIDMap[goods.GoodsUUID].Typ)
		}
		if len(logic.GoodsListRelated) > 0 {
			if err = tx.Create(&logic.GoodsListRelated).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		logic.SetGoodsFeild()

	case false:
		// logic.Batch.Update(logic.runtime.DB)
		err = dao.BatchDao.Update(logic.runtime.DB, logic.Batch)
	}
	return
}

// SetGoodsFeild 设置批次下货品的名称、类型、剩余信息以及统计批次总体数据
func (logic *BatchLogic) SetGoodsFeild() (err error) {
	// 1. 收集货品UUID列表
	var goodsUUIDList []string
	for _, goods := range logic.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM := make(map[string]model.GoodsFeild, 0)
	if goodsM, err = cache.GoodsCache.BatchGoodsFeildSet(goodsUUIDList, logic.OwnerUser); err != nil {
		return
	}

	// 3. 为每个货品设置名称和类型
	for _, goods := range logic.GoodsListRelated {
		goods.GoodsFeild = goodsM[goods.GoodsUUID]
	}

	// 4. 获取批次剩余库存和销售数据
	surplusM, e := BatchGoodsStat(logic.runtime.DB, logic.UID)
	if e != nil {
		err = e
		return
	}

	// 5. 为每个货品设置剩余库存和销售数量，同时统计总体数据
	var sellMount int            // 卖出总件数
	var sellWeight float64       // 卖出总重量
	var inventoryMount int       // 库存件数
	var inventoryWeight float64  // 库存重量
	var totalSalesAmount float64 // 总销售金额
	var totalProfit float64      // 总盈利

	for _, goods := range logic.GoodsListRelated {
		if surplus, ok := surplusM[goods.GoodsUUID]; ok {
			goods.BatchGoodsStatFeild = *surplus
			// 统计总销售金额（卖出的总金额）
			totalSalesAmount += surplus.SellTotal

			// 计算盈利：销售收入 - 入库成本
			var salesRevenue float64 // 销售收入
			var costAmount float64   // 入库成本

			if goods.GoodsTyp == common.GoodsTypeFix {
				// 定装：统计件数
				inventoryMount += surplus.Mount
				sellMount += surplus.SellMount

				// 销售收入 = 卖出件数 × 卖出价格（均价）
				salesRevenue = surplus.SellTotal
				// 入库成本 = 入库件数 × 入库价格
				costAmount = float64(goods.Mount) * goods.Price
			} else if goods.GoodsTyp == common.GoodsTypeBulk {
				// 散装：统计重量
				sellWeight += surplus.SellWeight
				inventoryWeight += surplus.Weight

				// 销售收入 = 卖出重量 × 卖出价格
				salesRevenue = surplus.SellTotal
				// 入库成本 = 入库重量 × 入库价格
				costAmount = goods.Weight * goods.Price
			}

			// 累加盈利
			totalProfit += (salesRevenue - costAmount)
		}
	}

	// 6. 填充批次统计数据
	logic.BatchStatFeild = model.BatchStatFeild{
		StatMount:           fmt.Sprintf("%d", sellMount),
		StatWeight:          utils.FloatReserveStr(sellWeight, 1),
		StatInventoryMount:  fmt.Sprintf("%d", inventoryMount),
		StatInventoryWeight: utils.FloatReserveStr(inventoryWeight, 1),
		StatSalesAmount:     utils.FloatReserveStr(totalSalesAmount, 1),
		StatSellProfit:      utils.FloatReserve(totalProfit, 1),
	}

	return nil
}

// // SetGoodsFeild 设置批次下货品的名称、类型以及剩余信息
// func (logic *BatchLogic) SetGoodsFeild() (err error) {
// 	var goodsUUIDList []string
// 	for _, goods := range logic.GoodsListRelated {
// 		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
// 	}

// 	goodsM, e := model.BatchGoodsFeildSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
// 	if e != nil {
// 		logic.runtime.Logger.Error("BatchLogic SetGoodsFeild", zap.Error(e))
// 		err = e
// 		return
// 	}

// 	for _, goods := range logic.GoodsListRelated {
// 		goods.GoodsName = goodsM[goods.GoodsUUID].GoodsName
// 		goods.GoodsTyp = goodsM[goods.GoodsUUID].GoodsTyp
// 	}

// 	surplusM, e := model.SetSurplusByBatch(logic.runtime.DB, logic.UID)

// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	for _, goods := range logic.GoodsListRelated {
// 		goods.Surplus = surplusM[goods.GoodsUUID].Surplus
// 		goods.SellAmount = surplusM[goods.GoodsUUID].SellAmount
// 	}
// 	return
// }

// func (logic *BatchLogic) FromDate(date string) (err error) {
// 	db := logic.runtime.DB
// 	var batch model.Batch
// 	if err = utils.Find(db.Preload("GoodsListRelated"), &batch, utils.WhereSerialNoCond(date), utils.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
// 		return
// 	}
// 	if batch.UID == "" {
// 		return common.ObjectNotExistErr
// 	}
// 	logic.Batch = batch
// 	return logic.SetGoodsFeild()
// }

func (logic *BatchLogic) FromUUID(uuid string) (err error) {
	var batch model.Batch
	if batch, err = dao.BatchDao.FromUUID(logic.runtime.DB, uuid, true); err != nil {
		logic.runtime.Logger.Error("BatchLogic FromUUID", zap.Any("uuid", uuid), zap.Error(err))
		return
	}
	logic.Batch = batch
	return logic.SetGoodsFeild()
}

func (logic *BatchLogic) FromLatest() (err error) {
	var batch model.Batch
	if batch, err = dao.BatchDao.FromLatest(logic.runtime.DB, logic.OwnerUser, true); err != nil {
		if err == gorm.ErrRecordNotFound {
			batch.GoodsListRelated = []*model.BatchGoods{}
			err = nil
		} else {
			return
		}
	}

	logic.Batch = batch
	return logic.SetGoodsFeild()
}

func (logic *BatchGoodsLogic) FromUUID(uuid string) (err error) {
	var batchGoods model.BatchGoods
	if err = utils.Find(logic.runtime.DB, &batchGoods, utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	if batchGoods.UID == "" {
		err = common.ObjectNotExistErr
		return
	}
	logic.BatchGoods = batchGoods
	var goods []model.Goods
	if err = utils.Find(logic.runtime.DB, &goods, utils.WhereUIDCond(logic.GoodsUUID)); err != nil {
		return
	}

	if len(goods) > 0 {
		logic.GoodsName = goods[0].Name
		logic.GoodsTyp = goods[0].Typ
	}

	return
}

func (logic *BatchGoodsLogic) Update() (err error) {
	return dao.BatchDao.UpdateBatchGoods(logic.runtime.DB, logic.BatchGoods)
}

func (logic *BatchLogic) List(req request.BatchListReq) (rsp []response.BatchListItem, err error) {
	var batches []model.Batch
	var userModel model.User
	if userModel, err = dao.UserDao.FromUUID(logic.runtime.DB, logic.OwnerUser); err != nil {
		return
	}

	if batches, err = dao.BatchDao.List(logic.runtime.DB, logic.OwnerUser, req); err != nil {
		return
	}
	for _, batch := range batches {
		t := time.Unix(batch.StorageTime, 0)
		rsp = append(rsp, response.BatchListItem{
			Status:   int(batch.Status),
			Time:     t.Format("2006-01-02 15:04:05"),
			UID:      batch.UID,
			Title:    userModel.Title(),
			SerialID: batch.SerialID,
		})
	}
	return
}

func (logic *BatchLogic) Precreate() (rsp response.PrecreateRsp, err error) {
	var maxID int
	db := logic.runtime.DB
	if maxID, err = logic.Batch.GenerateSerialID(db); err != nil {
		return
	}
	rsp.SerialID = maxID
	return
}
