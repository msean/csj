package dao

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"
	"time"

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
	if user.ContactPhones != "" {
		toUpdate["contact_phones"] = user.ContactPhones
	}
	if user.Province != "" {
		toUpdate["province"] = user.Province
	}
	if user.MarketName != "" {
		toUpdate["market_name"] = user.MarketName
	}
	if user.MarketAddress != "" {
		toUpdate["market_address"] = user.MarketAddress
	}
	if user.StallName != "" {
		toUpdate["stall_name"] = user.StallName
	}
	return utils.WhereUIDCond(user.UID).Cond(db.Model(&model.User{})).Updates(toUpdate).Error
}

func (dao *userDao) FindByPhone(db *gorm.DB, phone string) (user model.User, err error) {
	err = utils.Find(db, &user, utils.NewWhereCond("phone", phone))
	return
}

func (dao *userDao) FromUUID(db *gorm.DB, uuid string) (user model.User, err error) {
	if err = utils.Find(db, &user, utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	if user.UID == "" {
		err = common.UnRegisterErr
	}
	return
}

// UpdateVipDeadline 更新用户会员截止时间
func (dao *userDao) UpdateVipDeadline(db *gorm.DB, userUUID string, deadline time.Time) error {
	return db.Model(&model.User{}).Where("uid = ?", userUUID).
		Update("vip_deadline", deadline).Error
}
