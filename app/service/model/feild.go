package model

import (
	"app/pkg/utils"
	"app/service/common"
	"fmt"
)

type (
	BatchGoodsStatFeild struct {
		SurplusWeight   float64 `gorm:"-" json:"-"`       // 剩余重量 散装
		SurplusMount    int     `gorm:"-" json:"-"`       // 剩余件数 定装
		SellMount       int     `gorm:"-" json:"-"`       // 卖出件数 定装
		SellWeight      float64 `gorm:"-" json:"-"`       // 卖出总重量 散装
		SellTotal       float64 `gorm:"-" json:"-"`       // 销售总货款
		Surplus         string  `gorm:"-" json:"surplus"` // 剩余
		Type            int     `gorm:"-" json:"-"`
		SellAmount      string  `gorm:"-" json:"sellAmount"`      // 销量 件数/重量
		SellAvgPrice    string  `gorm:"-" json:"sellAvgPrice"`    // 销售均价
		SellTotalAmount string  `gorm:"-" json:"sellTotalAmount"` // 销售货款总金额
		SellProfit      float64 `gorm:"-" json:"sellProfit"`      // 销售利润
		InventoryPrice  float64 `gorm:"-" json:"-"`               // 入库价格
	}
	CustomerFeild struct {
		CustomerName  string `gorm:"-"  json:"customerName"`
		CustomerPhone string `gorm:"-"  json:"customerPhone"`
	}
	GoodsFeild struct {
		GoodsName   string  `gorm:"-" json:"name"`
		GoodsTyp    int     `gorm:"-" json:"type"`
		GoodsWeight float64 `gorm:"-" json:"goodsWeight"` // 定装是多少斤
	}
	PayFeild struct {
		PayFee  string `gorm:"-"  json:"payFee"`
		PayType int32  `gorm:"-"  json:"payType"`
		PaidFee string `gorm:"-"  json:"paidFee"`
	}
	BatchStatFeild struct {
		StatMount           string  `json:"statMount"`           // 入库总件数
		StatWeight          string  `json:"statWeight"`          // 入库总重量
		StatInventoryMount  string  `json:"statInventoryMount"`  // 库存 件数
		StatInventoryWeight string  `json:"statInventoryWeight"` // 库存 重量
		StatSalesAmount     string  `json:"statSalesAmount"`     // 卖货总金额
		StatSellProfit      float64 `json:"statSellProfit"`      // 盈利
		StatSellWeight      string  `json:"statSellWeight"`      // 卖出总重量
		StatSellMount       string  `json:"statSellMount"`       // 卖出总件数
	}
)

func (s *BatchGoodsStatFeild) Set() {
	s.SellTotalAmount = utils.FloatReserveStr(s.SellTotal, 1)

	if s.Type == common.GoodsTypeFix {
		s.Surplus = fmt.Sprintf("%d", s.SurplusMount)
		s.SellAmount = fmt.Sprintf("%d", s.SellMount)
		if s.SellMount != 0 {
			s.SellAvgPrice = utils.FloatReserveStr(s.SellTotal/float64(s.SellMount), 1)
		}
		return
	}
	s.Surplus = fmt.Sprintf("%.2f", s.SurplusWeight)
	s.SellAmount = fmt.Sprintf("%.2f", s.SellWeight)
	if s.SurplusWeight != 0 {
		s.SellAvgPrice = utils.FloatReserveStr(s.SellTotal/float64(s.SurplusWeight), 1)
	}
}
