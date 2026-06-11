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
	// if err = dao.Find(logic.runtime.DB, &_b); err != nil {
	// 	return
	// }
	// 当天只能创建一个批次 测试环境
	// if _b.UID != "" && global.Global.Env() != "test" {
	// 	err = common.BatchDuplicateErr
	// 	return
	// }

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
	if err = utils.First(logic.runtime.DB, &storage, utils.WhereUIDCond(logic.UID)); err != nil {
		logic.runtime.Logger.Error("BatchLogic Update", zap.Error(err))
		return
	}
	switch withBatchGoods {
	case true:
		var oldGoodsList []model.BatchGoods
		if err = logic.runtime.DB.Where("batch_uuid = ?", logic.UID).Find(&oldGoodsList).Error; err != nil {
			logic.runtime.Logger.Error("查询旧批次物品失败", zap.Error(err))
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
			err = logic.runtime.DB.Where("batch_uuid = ? AND goods_uuid IN (?)",
				logic.UID, toDeleteGoodsUUIDs).
				Find(&orderGoodsList).Error
			if err != nil {
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
		dao.BatchDao.Update(logic.runtime.DB, logic.Batch)
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
	var totalMount float64       // 总件数（定装）
	var totalWeight float64      // 总重量（散装）
	var inventoryMount int       // 库存件数
	var inventoryWeight float64  // 库存重量
	var totalSalesAmount float64 // 总销售金额

	for _, goods := range logic.GoodsListRelated {
		if surplus, ok := surplusM[goods.GoodsUUID]; ok {
			goods.BatchGoodsStatFeild = *surplus
			// 统计总销售金额（卖出的总金额）
			totalSalesAmount += surplus.SellTotal

			// 根据货品类型统计总数和库存
			if goods.GoodsTyp == common.GoodsTypeFix {
				// 定装：统计件数
				totalMount += float64(goods.Mount)
				inventoryMount += surplus.Mount
			} else if goods.GoodsTyp == common.GoodsTypeBulk {
				// 散装：统计重量
				totalWeight += goods.Weight
				inventoryWeight += surplus.Weight
			}
		}
	}

	// 6. 填充批次统计数据
	logic.BatchStatFeild = model.BatchStatFeild{
		StatMount:           utils.FloatReserveStr(totalMount, 1),
		StatWeight:          utils.FloatReserveStr(totalWeight, 1),
		StatInventoryMount:  fmt.Sprintf("%d", inventoryMount),
		StatInventoryWeight: utils.FloatReserveStr(inventoryWeight, 1),
		StatSalesAmount:     utils.FloatReserveStr(totalSalesAmount, 1),
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

func (logic *BatchLogic) FromDate(date string) (err error) {
	db := logic.runtime.DB
	var batch model.Batch
	if err = utils.Find(db.Preload("GoodsListRelated"), &batch, utils.WhereSerialNoCond(date), utils.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		return
	}
	if batch.UID == "" {
		return common.ObjectNotExistErr
	}
	logic.Batch = batch
	return logic.SetGoodsFeild()
}

func (logic *BatchLogic) FromUUID(uuid string) (err error) {
	db := logic.runtime.DB
	var batch model.Batch
	if err = utils.Find(db.Preload("GoodsListRelated"), &batch, utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	if batch.UID == "" {
		return common.ObjectNotExistErr
	}
	logic.Batch = batch
	return logic.SetGoodsFeild()
}

func (logic *BatchLogic) FromLatest() (err error) {
	db := logic.runtime.DB
	var batch model.Batch
	if err = utils.First(db.Preload("GoodsListRelated"), &batch, utils.CreatedOrderDescCond(), utils.NewWhereCond("owner_user", logic.OwnerUser)); err != nil {
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
	db := logic.runtime.DB

	_u := NewUser(logic.context)
	u, err := _u.FromUUID(logic.OwnerUser)
	if err != nil {
		logic.runtime.Logger.Error("BatchLogic List FromUUID", zap.Any("req", req), zap.Error(err))
		return
	}

	conds := []utils.Cond{
		req.LimitCond,
		utils.NewWhereCond("owner_user", logic.OwnerUser),
	}

	if req.StartDate != "" {
		startTime, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err != nil {
			logic.runtime.Logger.Error("BatchLogic List startTime", zap.Any("req", req), zap.Error(err))
			return nil, fmt.Errorf("开始日期格式错误: %v", err)
		}
		conds = append(conds, utils.NewCmpCond("created_at", ">=", startTime))
	}

	if req.EndDate != "" {
		endTime, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err != nil {
			logic.runtime.Logger.Error("BatchLogic List endTime", zap.Any("req", req), zap.Error(err))
			return nil, fmt.Errorf("结束日期格式错误: %v", err)
		}
		// 设置为当天的 23:59:59
		endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		conds = append(conds, utils.NewCmpCond("created_at", "<", endTime))
	}

	if req.Status != 0 {
		conds = append(conds, utils.NewWhereCond("status", req.Status))
	}

	conds = append(conds, utils.CreatedOrderDescCond())
	if err = utils.Find(db, &batches, conds...); err != nil {
		logic.runtime.Logger.Error("BatchLogic List FindBatch", zap.Any("req", req), zap.Error(err))
		return
	}

	for _, batch := range batches {
		name := u.Name
		if utils.IsBlankString(name) {
			name = u.Phone
		}
		t := time.Unix(batch.StorageTime, 0)
		rsp = append(rsp, response.BatchListItem{
			Status:   int(batch.Status),
			Time:     t.Format("2006-01-02 15:04:05"),
			UID:      batch.UID,
			Title:    name,
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
