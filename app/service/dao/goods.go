package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type GoodsDao struct{}

func NewGoodsDao() *GoodsDao {
	return &GoodsDao{}
}

func (dao *GoodsDao) Update(db *gorm.DB, goods model.Goods) error {
	if goods.Typ == model.GoodsTypePack {
		return utils.WhereUIDCond(goods.UID).Cond(db).Updates(&model.Goods{
			Name:   goods.Name,
			Price:  goods.Price,
			Weight: goods.Weight,
		}).Error
	}
	return utils.WhereUIDCond(goods.UID).Cond(db).Updates(&model.Goods{
		Name:  goods.Name,
		Price: goods.Price,
	}).Error
}

func (dao *GoodsDao) NameLike(db *gorm.DB, key string) utils.Cond {
	return utils.NewWhereLikeCond("name", key, utils.LikeTypeBetween)
}
