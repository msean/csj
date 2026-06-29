package model

import (
	"app/pkg/utils"
	"app/service/common"
	"fmt"
)

type (
	BatchGoodsStatFeild struct {
		Weight          float64 `gorm:"-" json:"-"`
		Mount           int     `gorm:"-" json:"-"`
		SellMount       int     `gorm:"-" json:"-"`       // 总件数
		SellWeight      float64 `gorm:"-" json:"-"`       // 总重量
		SellTotal       float64 `gorm:"-" json:"-"`       // 销售总货款
		Surplus         string  `gorm:"-" json:"surplus"` // 剩余
		Type            int     `gorm:"-" json:"-"`
		SellAmount      string  `gorm:"-" json:"sellAmount"`      // 销量 件数/重量
		SellAvgPrice    string  `gorm:"-" json:"sellAvgPrice"`    // 销售均价
		SellTotalAmount string  `gorm:"-" json:"sellTotalAmount"` // 销售货款总金额
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
		StatMount           string  `json:"statMount"`           // 卖出总件数
		StatWeight          string  `json:"statWeight"`          // 卖出总重量
		StatInventoryMount  string  `json:"statInventoryMount"`  // 库存 件数
		StatInventoryWeight string  `json:"statInventoryWeight"` // 库存 重量
		StatSalesAmount     string  `json:"statSalesAmount"`     // 卖货总金额
		StatSellProfit      float64 `json:"statSellProfit"`      // 盈利
		// StatSellWeight      string `json:"statSellWeight"`
	}
)

func (s *BatchGoodsStatFeild) Set() {
	s.SellTotalAmount = utils.FloatReserveStr(s.SellTotal, 1)

	if s.Type == common.GoodsTypeFix {
		s.Surplus = fmt.Sprintf("%d", s.Mount)
		s.SellAmount = fmt.Sprintf("%d", s.SellMount)
		if s.SellMount != 0 {
			s.SellAvgPrice = utils.FloatReserveStr(s.SellTotal/float64(s.SellMount), 1)
		}
		return
	}
	s.Surplus = fmt.Sprintf("%.2f", s.Weight)
	s.SellAmount = fmt.Sprintf("%.2f", s.SellWeight)
	if s.Weight != 0 {
		s.SellAvgPrice = utils.FloatReserveStr(s.SellTotal/float64(s.Weight), 1)
	}
}
