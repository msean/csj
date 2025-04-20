package dao

import (
	"fmt"
	"unicode"

	"gorm.io/gorm"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz"
)

func IsAllDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0 // 如果字符串为空，返回 false
}

func MapperByOwnerUser(ownerUsers []int64) (userMapper map[int64]biz.Users, err error) {
	userMapper = make(map[int64]biz.Users)
	db := global.GVA_DB.Model(&biz.Users{})
	var userss []biz.Users
	if err = db.Where("uid in (?)", ownerUsers).Find(&userss).Error; err != nil {
		return
	}
	for _, user := range userss {
		userMapper[user.Uid] = user
	}
	return
}

func MapperByCategory(categoryIDList []int64) (mapper map[int64]biz.GoodsCategory, err error) {
	mapper = make(map[int64]biz.GoodsCategory)
	db := global.GVA_DB.Model(&biz.GoodsCategory{})
	var categories []biz.GoodsCategory
	if err = db.Where("uid in (?)", categoryIDList).Find(&categories).Error; err != nil {
		return
	}
	for _, category := range categories {
		mapper[category.UID] = category
	}
	return
}

func OwnerUserCond(db *gorm.DB, value string, ownerUserModelField string) (out *gorm.DB, err error) {
	var user biz.Users
	if IsAllDigits(value) {
		if err = db.Where("uid = ?", value).First(&user).Error; err != nil {
			return
		}
		if user.Uid == 0 {
			if err = db.Where(fmt.Sprintf("name like %s%s%s", "%", value, "%"), value).First(&user).Error; err != nil {
				return
			}
		}
	} else {
		if err = db.Where(fmt.Sprintf("name like %s%s%s", "%", value, "%"), value).First(&user).Error; err != nil {
			return
		}
	}
	if user.Uid != 0 {
		out = db.Where(fmt.Sprintf("%s = ?", ownerUserModelField), user.Uid)
	}
	return
}
