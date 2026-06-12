package dao

import (
	"app/pkg/utils"
	"app/service/common"
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
	return
}

func (dao *userDao) FindByPhone(db *gorm.DB, phone string) (user model.User, err error) {
	err = utils.Find(db, &user, utils.NewWhereCond("phone", user.Phone))
	return
}

func (dao *userDao) FromUUID(db *gorm.DB, phone string) (user model.User, err error) {
	if err = utils.Find(db, &user, utils.NewWhereCond("phone", user.Phone)); err != nil {
		return
	}
	if user.UID == "" {
		err = common.UnRegisterErr
	}
	return
}
