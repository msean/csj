package model

import (
	"gorm.io/gorm"
)

const (
	GoodsTypePack   = 1 // 定装
	GoodsTypeWeight = 2 // 散装
)

type GoodsCategory struct {
	BaseModel
	OwnerUser string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
	Name      string  `gorm:"column:name;comment:客户名字" json:"name"`
	Goods     []Goods `gorm:"-" json:"goodsList"`
}

type Goods struct {
	BaseModel
	CategoryID string  `gorm:"column:category_id;comment:所属分类" json:"categoryID"`
	OwnerUser  string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
	Name       string  `gorm:"column:name;comment:客户名字" json:"name"`
	Typ        int32   `gorm:"column:type;comment:类型 1 定装 2散装" json:"type"`
	Price      float32 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
	Weight     float32 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
	Status     int     `gorm:"column:status;comment:状态" json:"status"`
}

func (g *Goods) Update(db *gorm.DB) error {
	if g.Typ == GoodsTypePack {
		return WhereUIDCond(g.UID).Cond(db).Updates(&Goods{
			Name:   g.Name,
			Price:  g.Price,
			Weight: g.Weight,
		}).Error
	}
	return WhereUIDCond(g.UID).Cond(db).Updates(&Goods{
		Name:  g.Name,
		Price: g.Price,
	}).Error
}

func (g *GoodsCategory) Update(db *gorm.DB) error {
	return WhereUIDCond(g.UID).Cond(db).Updates(&GoodsCategory{
		Name: g.Name,
	}).Error
}

func DeleteGoodsCategory(db *gorm.DB, goodCategoryUUID string) (err error) {
	if goodCategoryUUID == "" {
		return
	}
	var object GoodsCategory
	object.UID = goodCategoryUUID
	if err = WhereUIDCond(object.UID).Cond(db).Delete(&object).Error; err != nil {
		return
	}
	return NewWhereCond("category_id", object.UID).Cond(db.Model(&Goods{})).Updates(
		map[string]any{
			"category_id": "",
		}).Error
}

func (g *Goods) NameLike(db *gorm.DB, key string) Cond {
	return NewWhereLikeCond("name", key, LikeTypeBetween)
}
