package logic

import (
	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	BatchLogic struct {
		context   *gin.Context
		runtime   *global.RunTime
		OwnerUser int64
	}
	BatchGoodsLogic struct {
		OwnerUser int64
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

func (logic *BatchLogic) CalSurplus() (surplusGoodsList []response.BatchGoodsRsp, err error) {
	db := logic.runtime.DB

	var preBatch model.Batch
	conditions := []utils.Cond{
		utils.NewOrderCond("serial_no desc"),
		utils.NewWhereCond("owner_user", logic.OwnerUser),
	}
	err = utils.GormFind(db, &preBatch, conditions...)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	var preBatchGoodsList []response.BatchGoodsRsp
	if err = utils.GormFind(db, &preBatchGoodsList, utils.WhereOwnerUserCond(logic.OwnerUser), utils.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
		return
	}

	preGoodsList := []int64{}
	preCountM := make(map[int64]*response.BatchGoodsRsp)
	for _, preBatchGoods := range preBatchGoodsList {
		preCountM[preBatchGoods.GoodsUUID] = &preBatchGoods
		preGoodsList = append(preGoodsList, preBatchGoods.GoodsUUID)
	}

	// 上一批开单的
	var batchOrderedGoodsList []model.BatchOrderGoods
	if err = utils.GormFind(db, &batchOrderedGoodsList, utils.NewWhereCond("batch_uuid", preBatch.UID)); err != nil {
		return
	}
	for _, batchOrdered := range batchOrderedGoodsList {
		if _, InpreCountM := preCountM[batchOrdered.GoodsUUID]; InpreCountM {
			preCountM[batchOrdered.GoodsUUID].BatchGoods.Mount -= batchOrdered.Mount
			preCountM[batchOrdered.GoodsUUID].BatchGoods.Weight -= batchOrdered.Weight
			preCountM[batchOrdered.GoodsUUID].SurplusField = model.SurplusField{
				Mount:  preCountM[batchOrdered.GoodsUUID].BatchGoods.Mount,
				Weight: preCountM[batchOrdered.GoodsUUID].BatchGoods.Weight,
			}
			preCountM[batchOrdered.GoodsUUID].SurplusField.Set()
		}
	}

	goodsM, e := dao.BatchGoodsFieldSet(logic.runtime.DB, preGoodsList, logic.OwnerUser)
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

func (logic *BatchLogic) Create(param request.CreateBatchParam) (err error) {
	var checkBatch model.Batch
	serialNo := model.SerialNo(time.Now())
	if err = utils.GormFind(logic.runtime.DB, &checkBatch, utils.WhereSerialNoCond(serialNo)); err != nil {
		return
	}
	// 当天只能创建一个批次
	if checkBatch.UID != 0 {
		err = common.BatchDuplicateErr
		return
	}

	batchModel := model.Batch{
		OwnerUser:   logic.OwnerUser,
		SerialNo:    serialNo,
		StorageTime: param.StorageTime,
		Status:      model.BatchStatusOnSelling,
	}

	var batchGoods []model.BatchGoods
	for _, goods := range param.Goods {
		batchGoods = append(batchGoods, model.BatchGoods{
			OwnerUser: logic.OwnerUser,
			SerialNo:  serialNo,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
		})
	}
	db := logic.runtime.DB
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&batchModel).Error; err != nil {
			return err
		}

		if err := tx.Create(&batchGoods).Error; err != nil {
			return err
		}

		return nil
	})
	return
}

func (logic *BatchLogic) Update(param request.UpdateBatchParam) (err error) {
	var batch model.Batch
	if batch, err = dao.Batch.FromUUID(logic.runtime.DB, param.UUIDCompatible); err != nil {
		return
	}
	tx := logic.runtime.DB.Begin()
	if err = tx.Unscoped().Delete(&model.BatchGoods{}, "batch_uuid=?", param.UUIDCompatible).Error; err != nil {
		tx.Rollback()
		return
	}

	batch.StorageTime = param.StorageTime

	var batchGoods []model.BatchGoods
	for _, goods := range param.Goods {
		batchGoods = append(batchGoods, model.BatchGoods{
			SerialNo:  batch.SerialNo,
			OwnerUser: batch.OwnerUser,
			Price:     goods.Price,
			Weight:    goods.Weight,
			Mount:     goods.Mount,
		})
	}

	if err = dao.Batch.Update(tx, batch); err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Create(&batchGoods).Error; err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (logic *BatchLogic) Detail(param request.BatchDetailParam) (rsp response.BatchRsp, err error) {
	var batch model.Batch
	if param.UUID != "" {
		if batch, err = dao.Batch.FromUUID(logic.runtime.DB, param.UUIDCompatible); err != nil {
			return
		}
	} else {
		if param.Date != "" {
			if batch, err = dao.Batch.FromDate(logic.runtime.DB, logic.OwnerUser, param.Date); err != nil {
				return
			}
		} else {
			if batch, err = dao.Batch.FromLatest(logic.runtime.DB, logic.OwnerUser); err != nil {
				return
			}
		}
	}

	rsp.Batch = batch
	var batchGoods []model.BatchGoods
	batchGoods, err = dao.BatchGoods.FromBatchUUID(logic.runtime.DB, batch.UID)
	for _, goods := range batchGoods {
		rsp.Goods = append(rsp.Goods, &response.BatchGoodsRsp{
			BatchGoods: goods,
		})
	}

	logic.SetGoodsField(rsp)
	return
}

func (logic *BatchLogic) UpdateStatus(param request.UpdateBatchStatusParam) (err error) {
	var batch model.Batch
	if batch, err = dao.Batch.FromUUID(logic.runtime.DB, param.UUIDCompatible); err != nil {
		return
	}
	batch.Status = int32(param.Status)
	err = dao.Batch.Update(logic.runtime.DB.Begin(), batch)
	return
}

// SetGoodsField 设置批次下货品的名称、类型以及剩余信息
func (logic *BatchLogic) SetGoodsField(batch response.BatchRsp) (err error) {
	var goodsUUIDList []int64
	for _, goods := range batch.Goods {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM, e := dao.BatchGoodsFieldSet(logic.runtime.DB, goodsUUIDList, logic.OwnerUser)
	if e != nil {
		err = e
		return
	}

	for _, goods := range batch.Goods {
		goods.GoodsName = goodsM[goods.GoodsUUID].GoodsName
		goods.GoodsTyp = goodsM[goods.GoodsUUID].GoodsTyp
	}

	surplusM, e := dao.SetSurplusByBatch(logic.runtime.DB, batch.UID)
	if e != nil {
		err = e
		return
	}
	for _, goods := range batch.Goods {
		goods.Surplus = surplusM[goods.GoodsUUID].Surplus
	}
	return
}

func (logic *BatchGoodsLogic) FromUUID(uuid int64) (batchGoods response.BatchGoodsRsp, err error) {
	if err = utils.GormFind(logic.runtime.DB, &batchGoods, utils.WhereUIDCond(uuid), utils.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		return
	}
	if batchGoods.UID == 0 {
		err = common.ObjectNotExistErr
		return
	}
	var goods []model.Goods
	if err = utils.GormFind(logic.runtime.DB, &goods, utils.WhereUIDCond(batchGoods.GoodsUUID)); err != nil {
		return
	}

	if len(goods) > 0 {
		batchGoods.GoodsName = goods[0].Name
		batchGoods.GoodsTyp = goods[0].Typ
	}

	return
}

func (logic *BatchGoodsLogic) Update(param request.UpdateBatchGoodsParam) (batchGoods response.BatchGoodsRsp, err error) {
	batchGoods.UID = param.UUIDCompatible
	batchGoods.Price = param.Price
	batchGoods.BatchGoods.Mount = param.Mount
	batchGoods.BatchGoods.Weight = param.Weight
	err = dao.BatchGoods.Update(logic.runtime.DB, batchGoods.BatchGoods)
	return
}
