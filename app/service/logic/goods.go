package logic

import (
	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	GoodsLogic struct {
		OwnerUser string
		runtime   *global.RunTime
		context   *gin.Context
	}
	GoodsCategoryLogic struct {
		runtime   *global.RunTime
		context   *gin.Context
		OwnerUser string
	}
)

func NewGoodsLogic(context *gin.Context) *GoodsLogic {
	logic := &GoodsLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func NewGoodsCategoryLogic(context *gin.Context) *GoodsCategoryLogic {
	logic := &GoodsCategoryLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

// brief 不需要货品详情
func (logic *GoodsCategoryLogic) ListGoodsCategoryByUser(brief bool, conds ...utils.Cond) (gcList []*model.GoodsCategory, err error) {

	gcList = make([]*model.GoodsCategory, 0)

	conds = append(conds, utils.WhereOwnerUserCond(logic.OwnerUser))
	var _goodCategories []model.GoodsCategory
	if err = utils.GormFind(logic.runtime.DB, &_goodCategories, conds...); err != nil {
		logic.runtime.Logger.Error("ListGoodsCategoryByUser",
			zap.String("uuid", logic.OwnerUser),
			zap.Error(err))
		return
	}

	if brief {
		for _, gc := range _goodCategories {
			gcList = append(gcList, &model.GoodsCategory{
				Name:      gc.Name,
				OwnerUser: gc.OwnerUser,
				Goods:     []model.Goods{},
			})
		}
		return
	}

	var _goodsList []model.Goods
	if err = utils.GormFind(logic.runtime.DB, &_goodsList,
		utils.WhereOwnerUserCond(logic.OwnerUser),
		utils.UpdateOrderDescCond(),
	); err != nil {
		return
	}

	goodsCategoriesM := map[string]*model.GoodsCategory{}
	goodsCategoriesM[""] = &model.GoodsCategory{
		Name:  "未分类",
		Goods: []model.Goods{},
	}

	for _, _goodCategory := range _goodCategories {
		goodsCategoriesM[_goodCategory.UID] = &model.GoodsCategory{
			Name:      _goodCategory.Name,
			OwnerUser: _goodCategory.OwnerUser,
			Goods:     []model.Goods{},
		}
	}

	for _, goods := range _goodsList {
		if _goodsList, ok := goodsCategoriesM[goods.CategoryID]; ok {
			_goodsList.Goods = append(_goodsList.Goods, goods)
		}
	}

	for _, goodsCategory := range goodsCategoriesM {
		gcList = append(gcList, goodsCategory)
	}
	return
}

func (logic *GoodsCategoryLogic) Check(param request.GoodsCategorySaveParam) (err error) {
	var gc model.GoodsCategory
	if err = utils.GormFind(logic.runtime.DB, &gc, utils.WhereOwnerUserCond(logic.OwnerUser), utils.WhereNameCond(param.Name)); err != nil {
		return
	}
	if gc.UID != "" {
		return common.GoodsCategoryAlreadyExistErr
	}
	return nil
}

func (logic *GoodsCategoryLogic) Create(param request.GoodsCategorySaveParam) (_g model.GoodsCategory, err error) {
	_g = model.GoodsCategory{
		Name:      param.Name,
		OwnerUser: logic.OwnerUser,
	}
	err = utils.GormCreateObj(logic.runtime.DB, &_g)
	return
}

func (logic *GoodsCategoryLogic) Update(param request.GoodsCategorySaveParam) (_g model.GoodsCategory, err error) {
	_g = model.GoodsCategory{
		Name:      param.Name,
		OwnerUser: logic.OwnerUser,
		BaseModel: model.BaseModel{
			UID: param.UID,
		},
	}
	err = dao.GoodsCategory.Update(logic.runtime.DB, _g)
	return
}

func (logic *GoodsCategoryLogic) Delete(uuid string) (err error) {
	tx := logic.runtime.DB.Begin()
	if err = dao.GoodsCategory.DeleteGoodsCategory(tx, uuid); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (logic *GoodsLogic) Check(param request.GoodsSaveParam) (err error) {
	var _goods model.Goods
	db := logic.runtime.DB
	err = utils.GormFind(db, &_goods, utils.WhereOwnerUserCond(logic.OwnerUser), utils.WhereNameCond(param.Name))
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if _goods.UID != "" {
		return common.GoodsAlreadyExistErr
	}
	return nil
}

func (logic *GoodsLogic) Create(param request.GoodsSaveParam) (_goods model.Goods, err error) {
	_goods = model.Goods{
		CategoryID: param.CategoryID,
		Name:       param.Name,
		Typ:        param.Type,
		Price:      param.Price,
		Weight:     param.Weight,
		OwnerUser:  logic.OwnerUser,
	}
	err = utils.GormCreateObj(logic.runtime.DB, &_goods)
	return
}

func (logic *GoodsLogic) Update(param request.GoodsSaveParam) (_goods model.Goods, err error) {
	_goods = model.Goods{
		CategoryID: param.CategoryID,
		Name:       param.Name,
		Typ:        param.Type,
		Price:      param.Price,
		Weight:     param.Weight,
		BaseModel: model.BaseModel{
			UID: param.UID,
		},
	}
	err = dao.Goods.Update(logic.runtime.DB, _goods)
	return
}

func (logic *GoodsLogic) LoadGoods(searchKey string, limitCond utils.LimitCond) (goodsList []model.Goods, err error) {
	conds := []utils.Cond{
		utils.WhereOwnerUserCond(logic.OwnerUser),
		limitCond,
	}
	if searchKey != "" {
		conds = append(conds, dao.Goods.NameLike(logic.runtime.DB, searchKey))
	}

	if err = utils.GormFind(logic.runtime.DB, &goodsList, conds...); err != nil {
		logic.runtime.Logger.Error("LoadGoods",
			zap.String("uuid", logic.OwnerUser),
			zap.Any("conditions", conds),
			zap.Error(err))
		return
	}
	return
}
