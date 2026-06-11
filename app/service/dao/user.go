package dao

import (
	"app/pkg/utils"
	"app/service/model"

	"gorm.io/gorm"
)

type userDao struct{}

func newUserDao() *userDao {
	return &userDao{}
}

func (dao *userDao) Update(db *gorm.DB, user model.User) (err error) {
	toUpdate := make(map[string]any)
	if user.Name != "" {
		toUpdate["name"] = user.Name
	}
	if user.Phone != "" {
		toUpdate["phone"] = user.Phone
	}
	return utils.WhereUIDCond(user.UID).Cond(db.Model(&model.User{})).Updates(toUpdate).Error
}

func (dao *userDao) WherePhoneCond(user model.User) (w utils.WhereCond) {
	return utils.NewWhereCond("phone", user.Phone)
}
