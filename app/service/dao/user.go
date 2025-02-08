package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (dao *UserDao) Update(db *gorm.DB, user model.User) (err error) {
	toUpdate := make(map[string]any)
	if user.Name != "" {
		toUpdate["name"] = user.Name
	}
	if user.Phone != "" {
		toUpdate["phone"] = user.Phone
	}
	return utils.WhereUIDCond(user.UID).Cond(db.Model(&model.User{})).Updates(toUpdate).Error
}

func (dao *UserDao) WherePhoneCond(user model.User) (w utils.WhereCond) {
	return utils.NewWhereCond("phone", user.Phone)
}
