package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type GoodsCategoryDao struct{}

func NewGoodsCategoryDao() *GoodsCategoryDao {
	return &GoodsCategoryDao{}
}

func (dao *GoodsCategoryDao) Update(db *gorm.DB, goods model.GoodsCategory) error {
	return utils.WhereUIDCond(goods.UID).Cond(db).Updates(&model.GoodsCategory{
		Name: goods.Name,
	}).Error
}

func (dao *GoodsCategoryDao) DeleteGoodsCategory(db *gorm.DB, goodCategoryUUID string) (err error) {
	if goodCategoryUUID == "" {
		return
	}
	var object model.GoodsCategory
	object.UID = goodCategoryUUID
	if err = utils.WhereUIDCond(object.UID).Cond(db).Delete(&object).Error; err != nil {
		return
	}
	return utils.NewWhereCond("category_id", object.UID).Cond(db.Model(&model.Goods{})).Updates(
		map[string]any{
			"category_id": "",
		}).Error
}
