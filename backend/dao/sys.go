package dao

import (
	"github.com/msean/csj/backend/model/system"
	"gorm.io/gorm"
)

type sysDao struct{}

func newsysDao() *sysDao {
	return &sysDao{}
}

func (dao *sysDao) NameMapperFromIDList(db *gorm.DB, userID []int64) (mapper map[int64]string, err error) {
	var models []system.SysUser
	mapper = make(map[int64]string)
	if err = db.Find(&models, "id in (?)", userID).Error; err != nil {
		return
	}
	for _, model := range models {
		mapper[int64(model.ID)] = model.Username
	}
	return
}

func (dao *sysDao) AllUser(db *gorm.DB) (models []system.SysUser, err error) {
	err = db.Find(&models).Error
	return
}
