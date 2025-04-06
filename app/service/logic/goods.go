package logic

import (
	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"app/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	GoodsLogic struct {
		OwnerUser int64
		runtime   *global.RunTime
		context   *gin.Context
	}
	GoodsCategoryLogic struct {
		runtime   *global.RunTime
		context   *gin.Context
		OwnerUser int64
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
func (logic *GoodsCategoryLogic) ListGoodsCategoryByUser(param request.ListGoodsCategoryParam) (gcList []*response.GoodsCategoryRsp, err error) {

	gcList = make([]*response.GoodsCategoryRsp, 0)
	var modelGoodsCategoryList []model.GoodsCategory
	if modelGoodsCategoryList, err = dao.Goods.ListGoodsCatetoryByOwnerUser(logic.runtime.DB, logic.OwnerUser, param); err != nil {
		logic.runtime.Logger.Error("ListGoodsCategoryByUser", zap.Any("param", param), zap.Error(err))
		return
	}
	spew.Dump("ListGoodsCategoryByUser", modelGoodsCategoryList)
	if param.Brief {
		for _, model := range modelGoodsCategoryList {
			gcList = append(gcList, &response.GoodsCategoryRsp{
				GoodsCategory: model,
				Goods:         []response.GoodsDetailRsp{},
			})
		}
		return
	}

	var _goodsList []model.Goods
	_goodsList, err = dao.Goods.ListGoods(logic.runtime.DB, logic.OwnerUser, request.ListGoodsParam{
		OrderBy: "updated_at desc",
	})
	if err != nil {
		logic.runtime.Logger.Error("ListGoodsCategoryByUser", zap.Any("param", param), zap.Error(err))
		return
	}

	goodsCategoriesM := map[int64]*response.GoodsCategoryRsp{}
	goodsCategoriesM[0] = &response.GoodsCategoryRsp{
		GoodsCategory: model.GoodsCategory{
			Name: "未分类",
		},
		Goods: []response.GoodsDetailRsp{},
	}

	for _, model := range modelGoodsCategoryList {
		goodsCategoriesM[model.UID] = &response.GoodsCategoryRsp{
			GoodsCategory: model,
			Goods:         []response.GoodsDetailRsp{},
		}
	}

	for _, goods := range _goodsList {
		if _goodsList, ok := goodsCategoriesM[goods.CategoryID]; ok {
			_goodsList.Goods = append(_goodsList.Goods, response.GoodsDetailRsp{
				Goods: goods,
			})
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
	if gc.UID != 0 {
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
			UID: param.UIDCompatible,
		},
	}
	err = dao.GoodsCategory.Update(logic.runtime.DB, _g)
	return
}

func (logic *GoodsCategoryLogic) Delete(uuid int64) (err error) {
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
	if _goods.UID != 0 {
		return common.GoodsAlreadyExistErr
	}
	return nil
}

func (logic *GoodsLogic) Create(param request.GoodsSaveParam) (goods response.GoodsDetailRsp, err error) {
	_goods := model.Goods{
		CategoryID: param.CategoryIDCompatible,
		Name:       param.Name,
		Typ:        param.Type,
		Price:      param.Price,
		Weight:     param.Weight,
		OwnerUser:  logic.OwnerUser,
	}
	err = utils.GormCreateObj(logic.runtime.DB, &_goods)
	goods = response.GoodsDetailRsp{
		Goods: _goods,
	}
	return
}

func (logic *GoodsLogic) Update(param request.GoodsSaveParam) (goods response.GoodsDetailRsp, err error) {
	_goods := model.Goods{
		CategoryID: param.CategoryIDCompatible,
		Name:       param.Name,
		Typ:        param.Type,
		Price:      param.Price,
		Weight:     param.Weight,
		BaseModel: model.BaseModel{
			UID: param.UIDCompatible,
		},
	}
	err = dao.Goods.Update(logic.runtime.DB, _goods)
	goods = response.GoodsDetailRsp{
		Goods: _goods,
	}
	return
}

func (logic *GoodsLogic) LoadGoods(param request.ListGoodsParam) (goodsList []response.GoodsDetailRsp, err error) {
	var modelGoodsList []model.Goods
	if modelGoodsList, err = dao.Goods.ListGoods(logic.runtime.DB, logic.OwnerUser, param); err != nil {
		logic.runtime.Logger.Error("LoadGoods", zap.Any("param", param), zap.Error(err))
		return
	}
	for _, modelGoods := range modelGoodsList {
		goodsList = append(goodsList, response.GoodsDetailRsp{
			Goods: modelGoods,
		})
	}
	return
}
