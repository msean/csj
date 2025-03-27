package dao

import (
	"app/service/common"
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type BatchDao struct{}

func NewBatchDao() *BatchDao {
	return &BatchDao{}
}

func (dao *BatchDao) Update(db *gorm.DB, batch model.Batch) (err error) {
	toUpdates := make(map[string]any)
	if batch.StorageTime != 0 {
		toUpdates["storage_time"] = batch.StorageTime
	}
	if batch.Status != 0 {
		toUpdates["status"] = batch.Status
	}
	if len(toUpdates) > 0 {
		return utils.WhereUIDCond(batch.UID).Cond(db).Model(&model.Batch{}).Updates(toUpdates).Error
	}
	return
}

func (logic *BatchDao) FromUUID(db *gorm.DB, uuid int64) (batch model.Batch, err error) {
	if err = utils.GormFind(db, &batch, utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	if batch.UID == 0 {
		err = common.ObjectNotExistErr
		return
	}
	return
}

func (logic *BatchDao) FromDate(db *gorm.DB, ownerUser int64, date string) (batch model.Batch, err error) {
	if err = utils.GormFind(db, &batch, utils.WhereSerialNoCond(date), utils.WhereOwnerUserCond(ownerUser)); err != nil {
		return
	}
	if batch.UID == 0 {
		err = common.ObjectNotExistErr
		return
	}
	return
}

func (logic *BatchDao) FromLatest(db *gorm.DB, ownerUser int64) (batch model.Batch, err error) {
	if err = utils.GormFirst(db, &batch, utils.CreatedOrderDescCond(), utils.NewWhereCond("owner_user", ownerUser)); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			return
		}
	}

	return
}
