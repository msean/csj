package dao

import (
	"app/pkg/utils"
	"app/service/model"
	"app/service/model/request"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type batchDao struct{}

func newbatchDao() *batchDao {
	return &batchDao{}
}

func (dao *batchDao) Update(db *gorm.DB, object model.Batch) (err error) {
	toUpdates := make(map[string]any)
	if object.StorageTime != 0 {
		toUpdates["storage_time"] = object.StorageTime
	}
	if object.Status != 0 {
		toUpdates["status"] = object.Status
	}
	if len(toUpdates) > 0 {
		return utils.WhereUIDCond(object.UID).Cond(db).Model(&model.Batch{}).Updates(toUpdates).Error
	}
	return
}

func (dao *batchDao) UpdateBatchGoods(db *gorm.DB, object model.BatchGoods) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchGoods{
		Price:  object.Price,
		Weight: object.Weight,
		Mount:  object.Mount,
	}).Error
}

func (dao *batchDao) List(db *gorm.DB, ownerUser string, conditions request.BatchListReq) (batches []model.Batch, err error) {
	conds := []utils.Cond{
		conditions.LimitCond,
		utils.NewWhereCond("owner_user", ownerUser),
	}

	if conditions.StartDate != "" {
		startTime, err := time.ParseInLocation("2006-01-02", conditions.StartDate, time.Local)
		if err != nil {
			return nil, fmt.Errorf("开始日期格式错误: %v", err)
		}
		conds = append(conds, utils.NewCmpCond("created_at", ">=", startTime))
	}

	if conditions.EndDate != "" {
		endTime, err := time.ParseInLocation("2006-01-02", conditions.EndDate, time.Local)
		if err != nil {
			return nil, fmt.Errorf("结束日期格式错误: %v", err)
		}
		// 设置为当天的 23:59:59
		endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		conds = append(conds, utils.NewCmpCond("created_at", "<", endTime))
	}

	if conditions.Status != 0 {
		conds = append(conds, utils.NewWhereCond("status", conditions.Status))
	}

	conds = append(conds, utils.CreatedOrderDescCond())
	utils.Find(db, &batches, conds...)
	return
}

func (dao *batchDao) FromUUID(db *gorm.DB, uuid string, withGoods bool) (batch model.Batch, err error) {
	if withGoods {
		db = db.Preload("GoodsListRelated")
	}
	err = utils.First(db, &batch, utils.WhereUIDCond(uuid))
	return
}

func (dao *batchDao) FromLatest(db *gorm.DB, ownerUser string, withGoods bool) (batch model.Batch, err error) {
	if withGoods {
		db = db.Preload("GoodsListRelated")
	}
	err = utils.First(db, &batch, utils.CreatedOrderDescCond(), utils.WhereOwnerUserCond(ownerUser))
	return
}
