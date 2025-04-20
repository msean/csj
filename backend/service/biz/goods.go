package biz

import (
	"github.com/flipped-aurora/gin-vue-admin/server/dao"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz"
	biz_request "github.com/flipped-aurora/gin-vue-admin/server/model/biz/request"
)

type GoodsService struct{}

func (svc *GoodsService) ListGoods(info biz_request.GoodsListParam) (list []*biz.Goods, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&biz.Goods{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OwnerUser != "" {
		db, err = dao.OwnerUserCond(db, info.OwnerUser, "owner_user")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&list).Error
	var ownerUsers []int64
	var categories []int64
	for _, goods := range list {
		ownerUsers = append(ownerUsers, goods.OwnerUser)
		categories = append(categories, goods.CategoryID)
	}
	userMapper := make(map[int64]biz.Users)
	categoryMapper := make(map[int64]biz.GoodsCategory)
	if userMapper, err = dao.MapperByOwnerUser(ownerUsers); err != nil {
		return
	}
	if categoryMapper, err = dao.MapperByCategory(categories); err != nil {
		return
	}

	for _, goods := range list {
		svc.Fill(goods, userMapper[goods.OwnerUser], categoryMapper[goods.CategoryID])
	}
	return
}

func (svc *GoodsService) Fill(goods *biz.Goods, user biz.Users, category biz.GoodsCategory) {
	goods.CategoryName = category.Name
	goods.OwnerUserNamme = user.Name
}
