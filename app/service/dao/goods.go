package dao

import (
	"app/service/model"
	"app/service/model/request"
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

func (dao *GoodsDao) ListGoods(
	db *gorm.DB,
	ownerUser int64,
	param request.ListGoodsParam,
) (modelGoodsList []model.Goods, err error) {

	conds := []utils.Cond{
		utils.WhereOwnerUserCond(ownerUser),
		param.LimitCond,
	}
	if param.SearchKey != "" {
		conds = append(conds, dao.NameLike(db, param.SearchKey))
	}
	if param.OrderBy != "" {
		conds = append(conds, utils.NewOrderCond(param.OrderBy))
	}

	err = utils.GormFind(db, &modelGoodsList, conds...)
	return
}

func (dao *GoodsDao) ListGoodsCatetoryByOwnerUser(
	db *gorm.DB,
	ownerUser int64,
	param request.ListGoodsCategoryParam,
) (modelGoodsList []model.GoodsCategory, err error) {
	conds := []utils.Cond{param.LimitCond, utils.WhereOwnerUserCond(ownerUser)}
	err = utils.GormFind(db, &modelGoodsList, conds...)
	return
}
