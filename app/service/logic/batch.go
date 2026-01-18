package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	BatchLogic struct {
		context *gin.Context
		runtime *global.RunTime
		model.Batch
		Date string `json:"date"`
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

func (logic *BatchLogic) CalSurplus(ownerUser string) (surplusGoodsList []BatchGoodsLogic, err error) {
	db := logic.runtime.DB
	// 找到上一批
	var preBatch model.Batch
	conds := []model.Cond{
		model.NewOrderCond("serial_no desc"),
		model.NewWhereCond("owner_user", ownerUser),
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
	if err = model.Find(db, &preBatchGoodsList, model.WhereOwnerUserCond(ownerUser), model.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
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
	// 当天只能创建一个批次
	if _b.UID != "" {
		err = common.BatchDuplicateErr
		return
	}

	logic.Default()
	for _, goods := range logic.GoodsListRelated {
		goods.SerialNo = logic.SerialNo
		goods.OwnerUser = logic.OwnerUser
	}

	return model.CreateObj(logic.runtime.DB, &logic.Batch)
}

// withBatchGoods 是否更新批次下的物品信息
func (logic *BatchLogic) Update(withBatchGoods bool) (err error) {
	switch withBatchGoods {
	case true:
		tx := logic.runtime.DB.Begin()
		if err = tx.Delete(&model.BatchGoods{}, "batch_uuid=?", logic.UID).Error; err != nil {
			tx.Rollback()
			return err
		}
		for _, goods := range logic.GoodsListRelated {
			goods.SerialNo = logic.SerialNo
			goods.OwnerUser = logic.OwnerUser
		}
		if err = tx.Save(&logic.Batch).Error; err != nil {
			tx.Rollback()
			return
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
	}
	return
}

func (logic *BatchLogic) FromDate() (err error) {
	db := logic.runtime.DB
	var batch model.Batch
	if err = model.Find(db.Preload("GoodsListRelated"), &batch, model.WhereSerialNoCond(logic.Date), model.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		return
	}
	if batch.UID == "" {
		return common.ObjectNotExistErr
	}
	logic.Batch = batch
	return logic.SetGoodsFeild()
}

func (logic *BatchLogic) FromUUID() (err error) {
	db := logic.runtime.DB
	var batch model.Batch
	if err = model.Find(db.Preload("GoodsListRelated"), &batch, model.WhereUIDCond(logic.OwnerUser)); err != nil {
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

func (logic *BatchGoodsLogic) FromUUID() (err error) {
	var batchGoods model.BatchGoods
	if err = model.Find(logic.runtime.DB, &batchGoods, model.WhereUIDCond(logic.UID)); err != nil {
		return
	}
	if batchGoods.UID == "" {
		err = common.ObjectNotExistErr
		return
	}
	logic.BatchGoods = batchGoods
	var goods []model.Goods
	if err = model.Find(logic.runtime.DB, &goods, model.WhereUIDCond(logic.UID)); err != nil {
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
