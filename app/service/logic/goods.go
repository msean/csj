package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	GoodsLogic struct {
		model.Goods
		runtime *global.RunTime
		context *gin.Context
	}
	GoodsCategoryLogic struct {
		model.GoodsCategory
		GoodsList []GoodsLogic `json:"goodsList"`
		runtime   *global.RunTime
		context   *gin.Context
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
func (logic *GoodsCategoryLogic) ListGoodsCategoryByUser(brief bool, conds ...utils.Cond) (gcList []*GoodsCategoryLogic, err error) {
	gcList = make([]*GoodsCategoryLogic, 0)

	conds = append(conds, utils.WhereOwnerUserCond(logic.OwnerUser))
	var _goodCategories []model.GoodsCategory
	if err = utils.Find(logic.runtime.DB, &_goodCategories, conds...); err != nil {
		return
	}

	if brief {
		for _, gc := range _goodCategories {
			gcList = append(gcList, &GoodsCategoryLogic{
				GoodsCategory: gc,
				GoodsList:     []GoodsLogic{},
			})
		}
		return
	}

	var _goodsList []model.Goods
	if err = utils.Find(logic.runtime.DB, &_goodsList, utils.WhereOwnerUserCond(logic.OwnerUser), utils.UpdateOrderDescCond()); err != nil {
		return
	}

	goodsCategoriesM := map[string]*GoodsCategoryLogic{}
	goodsCategoriesM[""] = &GoodsCategoryLogic{
		GoodsCategory: model.GoodsCategory{
			Name: "未分类",
		},
		GoodsList: []GoodsLogic{},
	}

	for _, _goodCategory := range _goodCategories {
		goodsCategoriesM[_goodCategory.UID] = &GoodsCategoryLogic{
			GoodsCategory: _goodCategory,
			GoodsList:     []GoodsLogic{},
		}
	}

	for _, goods := range _goodsList {
		if _goodsList, ok := goodsCategoriesM[goods.CategoryID]; ok {
			_goodsList.GoodsList = append(_goodsList.GoodsList, GoodsLogic{
				Goods: goods,
			})
		}
	}

	for _, goodsCategory := range goodsCategoriesM {
		gcList = append(gcList, goodsCategory)
	}
	return
}

func (logic *GoodsCategoryLogic) Check() (err error) {
	var gc model.GoodsCategory
	if err = utils.Find(logic.runtime.DB, &gc, utils.WhereOwnerUserCond(logic.OwnerUser), utils.WhereNameCond(logic.Name)); err != nil {
		return
	}
	if gc.UID != "" {
		return common.GoodsCategoryAlreadyExistErr
	}
	return nil
}

func (logic *GoodsCategoryLogic) Create() (err error) {
	return utils.CreateObj(logic.runtime.DB, &logic.GoodsCategory)
}

func (logic *GoodsCategoryLogic) Update() (err error) {
	return dao.GoodsDao.UpdateGoodsCategory(logic.runtime.DB, logic.GoodsCategory)
}

func (logic *GoodsCategoryLogic) Delete() (err error) {
	return dao.GoodsDao.DeleteCategory(logic.runtime.DB, logic.GoodsCategory)
}

func (logic *GoodsLogic) Check() (err error) {
	var _goods model.Goods
	db := logic.runtime.DB
	err = utils.Find(db, &_goods, utils.WhereOwnerUserCond(logic.OwnerUser), utils.WhereNameCond(logic.Name))
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if _goods.UID != "" {
		return common.GoodsAlreadyExistErr
	}
	return nil
}

func (logic *GoodsLogic) Create() (err error) {
	return utils.CreateObj(logic.runtime.DB, &logic.Goods)
}

func (logic *GoodsLogic) Update() (err error) {
	tx := logic.runtime.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 1️⃣ 更新 Goods 自身
	if err = dao.GoodsDao.Update(tx, logic.Goods); err != nil {
		return
	}

	// 2️⃣ 同步更新在售批次里的 BatchGoods
	err = tx.Model(&model.BatchGoods{}).
		Where("goods_uuid = ?", logic.Goods.BaseModel.UID).
		Where(fmt.Sprintf(`
			batch_uuid IN (SELECT uid FROM %s WHERE status = ?)`, model.Batch{}.TableName()), 1).
		Updates(map[string]interface{}{
			"price":  logic.Goods.Price,
			"weight": logic.Goods.Weight,
		}).Error

	return
}

func (logic *GoodsLogic) LoadGoods(req request.GoodsListReq) (goodsList []model.Goods, err error) {
	goodsList, err = dao.GoodsDao.List(logic.runtime.DB, logic.OwnerUser, req)
	return
}
