package dao

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"

	"gorm.io/gorm"
)

type goodsDao struct{}

func newGoodsDao() *goodsDao {
	return &goodsDao{}
}

func (dao *goodsDao) Update(db *gorm.DB, goods model.Goods) (err error) {
	if goods.Typ == common.GoodsTypeBulk {
		goods.Weight = 0
	}
	return utils.WhereUIDCond(goods.UID).Cond(db).Updates(&model.Goods{
		Name:   goods.Name,
		Price:  goods.Price,
		Weight: goods.Weight,
	}).Error
}

func (dao *goodsDao) ResetCategory(db *gorm.DB, goodsCategory model.GoodsCategory) error {
	return utils.WhereUIDCond(goodsCategory.UID).Cond(db).Delete(&goodsCategory).Error
}

func (dao *goodsDao) DeleteCategory(db *gorm.DB, goodsCategory model.GoodsCategory) (err error) {
	tx := db.Begin()
	if err = utils.WhereUIDCond(goodsCategory.UID).Cond(tx).Delete(&goodsCategory).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = utils.NewWhereCond("category_id", goodsCategory.UID).Cond(db.Model(&model.Goods{})).Updates(
		map[string]any{
			"category_id": "",
		}).Error; err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (dao *goodsDao) UpdateGoodsCategory(db *gorm.DB, goodsCategory model.GoodsCategory) error {
	return utils.WhereUIDCond(goodsCategory.UID).Cond(db).Updates(&model.GoodsCategory{
		Name: goodsCategory.Name,
	}).Error
}

func (dao *goodsDao) List(db *gorm.DB, ownerUser string, conditions request.GoodsListReq) (goodsList []model.Goods, err error) {
	conds := []utils.Cond{
		utils.WhereOwnerUserCond(ownerUser),
	}
	if !conditions.LoadAll {
		conds = append(conds, conditions.LimitCond)
	}
	if conditions.SearchKey != "" {
		conds = append(conds, utils.NewWhereLikeCond("name", conditions.SearchKey, utils.LikeTypeBetween))
	}

	err = utils.Find(db, &goodsList, conds...)
	return
}
