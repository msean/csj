package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	BatchLogic struct {
		context   *gin.Context
		runtime   *global.RunTime
		OwnerUser string
		// model.Batch
		// Date string `json:"date"`
	}
	BatchGoodsLogic struct {
		// model.BatchGoods
		OwnerUser string
		context   *gin.Context
		runtime   *global.RunTime
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

func (logic *BatchLogic) CalSurplus() (surplusGoodsList []model.BatchGoods, err error) {
	db := logic.runtime.DB

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

	preGoodsList := []string{}
	preCountM := make(map[string]*model.BatchGoods)
	for _, preBatchGoods := range preBatchGoodsList {
		preCountM[preBatchGoods.GoodsUUID] = &preBatchGoods
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

func (logic *BatchLogic) Create(param request.CreateBatchParam) (batch model.Batch, err error) {
	var checkBatch model.Batch
	serialNo := model.SerioalNo(time.Now())
	if err = model.Find(logic.runtime.DB, &checkBatch, model.WhereSerialNoCond(serialNo)); err != nil {
		return
	}
	// 当天只能创建一个批次
	if checkBatch.UID != "" {
		err = common.BatchDuplicateErr
		return
	}

	batch = model.Batch{
		OwnerUser:   logic.OwnerUser,
		SerialNo:    serialNo,
		StorageTime: param.StroageTime,
		Status:      model.BatchStatusOnSellering,
	}

	for _, goods := range param.Goods {
		goods := model.BatchGoods{
			OwnerUser: logic.OwnerUser,
			SerialNo:  serialNo,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
		}
		batch.GoodsListRelated = append(batch.GoodsListRelated, &goods)
	}

	err = model.CreateObj(logic.runtime.DB, &batch)
	return
}

func (logic *BatchLogic) Update(param request.UpdateBatchParam) (batch model.Batch, err error) {
	if batch, err = logic.FromUUID(param.UUID, false); err != nil {
		return
	}
	tx := logic.runtime.DB.Begin()
	if err = tx.Delete(&model.BatchGoods{}, "batch_uuid=?", param.UUID).Error; err != nil {
		tx.Rollback()
		return
	}
	for _, goods := range param.Goods {
		batch.GoodsListRelated = append(batch.GoodsListRelated, &model.BatchGoods{
			SerialNo:  batch.SerialNo,
			OwnerUser: batch.OwnerUser,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
		})
	}
	if err = tx.Save(&batch).Error; err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	logic.SetGoodsFeild(batch)
	return
}

func (logic *BatchLogic) UpdateStatus(param request.UpdateBatchStatusParam) (err error) {
	var batch model.Batch
	if batch, err = logic.FromUUID(param.UUID, false); err != nil {
		return
	}
	if err = model.WhereUIDCond(batch.UID).Cond(logic.runtime.DB).Model(&model.Batch{}).Updates(map[string]any{
		"status": param.Status,
	}).Error; err != nil {
		return
	}
	err = logic.SetGoodsFeild(batch)
	return
}

// SetGoodsFeild 设置批次下货品的名称、类型以及剩余信息
func (logic *BatchLogic) SetGoodsFeild(batch model.Batch) (err error) {
	var goodsUUIDList []string
	for _, goods := range batch.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM, e := model.BatchGoodsFeildSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
	if e != nil {
		err = e
		return
	}

	for _, goods := range batch.GoodsListRelated {
		goods.GoodsName = goodsM[goods.GoodsUUID].GoodsName
		goods.GoodsTyp = goodsM[goods.GoodsUUID].GoodsTyp
	}

	surplusM, e := model.SetSurplusByBatch(logic.runtime.DB, batch.UID)
	if e != nil {
		err = e
		return
	}
	for _, goods := range batch.GoodsListRelated {
		goods.Surplus = surplusM[goods.GoodsUUID].Surplus
	}
	return
}

func (logic *BatchLogic) FromDate(date string) (batch model.Batch, err error) {
	db := logic.runtime.DB
	if err = model.Find(db.Preload("GoodsListRelated"), &batch, model.WhereSerialNoCond(date), model.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		return
	}
	if batch.UID == "" {
		err = common.ObjectNotExistErr
		return
	}
	err = logic.SetGoodsFeild(batch)
	return
}

func (logic *BatchLogic) FromUUID(uuid string, withGoods bool) (batch model.Batch, err error) {
	db := logic.runtime.DB
	if withGoods {
		db = db.Preload("GoodsListRelated")
	}
	if err = model.Find(db, &batch,
		model.WhereUIDCond(uuid),
		model.WhereOwnerUserCond(logic.OwnerUser),
	); err != nil {
		return
	}
	if batch.UID == "" {
		err = common.ObjectNotExistErr
		return
	}
	err = logic.SetGoodsFeild(batch)
	return
}

func (logic *BatchLogic) FromLatest() (batch model.Batch, err error) {
	db := logic.runtime.DB
	if err = model.First(db.Preload("GoodsListRelated"), &batch, model.CreatedOrderDescCond(), model.NewWhereCond("owner_user", logic.OwnerUser)); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			return
		}
	}

	err = logic.SetGoodsFeild(batch)
	return
}

func (logic *BatchGoodsLogic) FromUUID(uuid string) (batchGoods model.BatchGoods, err error) {
	if err = model.Find(logic.runtime.DB, &batchGoods, model.WhereUIDCond(uuid), model.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		return
	}
	if batchGoods.UID == "" {
		err = common.ObjectNotExistErr
		return
	}
	var goods []model.Goods
	if err = model.Find(logic.runtime.DB, &goods, model.WhereUIDCond(batchGoods.GoodsUUID)); err != nil {
		return
	}

	if len(goods) > 0 {
		batchGoods.GoodsName = goods[0].Name
		batchGoods.GoodsTyp = goods[0].Typ
	}

	return
}

func (logic *BatchGoodsLogic) Update(param request.UpdateBatchGoodsParam) (batchGoods model.BatchGoods, err error) {
	batchGoods.UID = param.UUID
	batchGoods.Price = param.Price
	batchGoods.Mount = param.Mount
	batchGoods.Weight = param.Weight
	err = batchGoods.Update(logic.runtime.DB)
	return
}
