package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
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
	conds := []model.Cond{
		model.NewOrderCond("serial_no desc"),
		model.NewWhereCond("owner_user", logic.OwnerUser),
	}
	err = model.Find(db, &preBatch, conds...)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	var preBatchGoodsList []model.BatchGoods
	if err = model.Find(db, &preBatchGoodsList, model.WhereOwnerUserCond(logic.OwnerUser), model.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
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
	if err = model.Find(db, &batchOrderedGoodsList, model.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
		return
	}
	for _, batchOrdered := range batchOrderedGoodsList {
		if _, InpreCountM := preCountM[batchOrdered.GoodsUUID]; InpreCountM {
			preCountM[batchOrdered.GoodsUUID].Mount -= batchOrdered.Mount
			preCountM[batchOrdered.GoodsUUID].Weight -= batchOrdered.Weight
			preCountM[batchOrdered.GoodsUUID].SurplusFeild = model.SurplusFeild{
				Mount:  preCountM[batchOrdered.GoodsUUID].Mount,
				Weight: preCountM[batchOrdered.GoodsUUID].Weight,
			}
			preCountM[batchOrdered.GoodsUUID].SurplusFeild.Set()
		}
	}

	goodsM, e := model.BatchGoodsFeildSet(logic.runtime.DB, preGoodsList, logic.OwnerUser)
	if e != nil {
		err = e
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
	if err = model.Find(logic.runtime.DB, &_b, model.WhereSerialNoCond(model.SerioalNo(time.Now()))); err != nil {
		return
	}
	// 当天只能创建一个批次 测试环境
	if _b.UID != "" && global.Global.Env() != "test" {
		err = common.BatchDuplicateErr
		return
	}

	var batchGoods []model.Goods
	var _goodsUIDList []string
	goodsUIDMap := make(map[string]model.Goods)

	for _, goods := range logic.GoodsListRelated {
		_goodsUIDList = append(_goodsUIDList, goods.GoodsUUID)
	}

	if err = model.Find(logic.runtime.DB, &batchGoods, model.InUIDCondFromString(_goodsUIDList)); err != nil {
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

	return model.CreateObj(logic.runtime.DB, &logic.Batch)
}

// withBatchGoods 是否更新批次下的物品信息
func (logic *BatchLogic) Update(withBatchGoods bool) (err error) {
	var storage model.Batch
	if err = model.First(logic.runtime.DB, &storage, model.WhereUIDCond(logic.UID)); err != nil {
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

				goodsM, e := model.BatchGoodsFeildSet(logic.runtime.DB, orderedGoodsUUIDs, logic.OwnerUser)
				if e != nil {
					err = e
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

		if err = model.Find(logic.runtime.DB, &batchGoods, model.InUIDCondFromString(_goodsUIDList)); err != nil {
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
		logic.Batch.Update(logic.runtime.DB)
	}
	return
}

// SetGoodsFeild 设置批次下货品的名称、类型以及剩余信息
func (logic *BatchLogic) SetGoodsFeild() (err error) {
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

	surplusM, e := model.SetSurplusByBatch(logic.runtime.DB, logic.UID)

	if e != nil {
		err = e
		return
	}
	for _, goods := range logic.GoodsListRelated {
		goods.Surplus = surplusM[goods.GoodsUUID].Surplus
		goods.SellAmount = surplusM[goods.GoodsUUID].SellAmount
	}
	return
}

func (logic *BatchLogic) FromDate(date string) (err error) {
	db := logic.runtime.DB
	var batch model.Batch
	if err = model.Find(db.Preload("GoodsListRelated"), &batch, model.WhereSerialNoCond(date), model.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
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
	if err = model.Find(db.Preload("GoodsListRelated"), &batch, model.WhereUIDCond(uuid)); err != nil {
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
	if err = model.First(db.Preload("GoodsListRelated"), &batch, model.CreatedOrderDescCond(), model.NewWhereCond("owner_user", logic.OwnerUser)); err != nil {
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
	if err = model.Find(logic.runtime.DB, &batchGoods, model.WhereUIDCond(uuid)); err != nil {
		return
	}
	if batchGoods.UID == "" {
		err = common.ObjectNotExistErr
		return
	}
	logic.BatchGoods = batchGoods
	var goods []model.Goods
	if err = model.Find(logic.runtime.DB, &goods, model.WhereUIDCond(logic.GoodsUUID)); err != nil {
		return
	}

	if len(goods) > 0 {
		logic.GoodsName = goods[0].Name
		logic.GoodsTyp = goods[0].Typ
	}

	return
}

func (logic *BatchGoodsLogic) Update() (err error) {
	return logic.BatchGoods.Update(logic.runtime.DB)
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

	conds := []model.Cond{
		req.LimitCond,
		model.NewWhereCond("owner_user", logic.OwnerUser),
	}

	if req.StartDate != "" {
		startTime, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err != nil {
			logic.runtime.Logger.Error("BatchLogic List startTime", zap.Any("req", req), zap.Error(err))
			return nil, fmt.Errorf("开始日期格式错误: %v", err)
		}
		conds = append(conds, model.NewCmpCond("created_at", ">=", startTime))
	}

	if req.EndDate != "" {
		endTime, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err != nil {
			logic.runtime.Logger.Error("BatchLogic List endTime", zap.Any("req", req), zap.Error(err))
			return nil, fmt.Errorf("结束日期格式错误: %v", err)
		}
		// 设置为当天的 23:59:59
		endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		conds = append(conds, model.NewCmpCond("created_at", "<", endTime))
	}

	if req.Status != 0 {
		conds = append(conds, model.NewWhereCond("status", req.Status))
	}

	conds = append(conds, model.CreatedOrderDescCond())
	if err = model.Find(db, &batches, conds...); err != nil {
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
			Status: int(batch.Status),
			Time:   t.Format("2006-01-02 15:04:05"),
			UID:    batch.UID,
			Title:  fmt.Sprintf("%s %d", name, batch.SerialID),
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
