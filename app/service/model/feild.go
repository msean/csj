package model

import (
	"fmt"
)

type (
	SurplusField struct {
		Weight  float64 `gorm:"-" json:"-"`
		Mount   int32   `gorm:"-" json:"-"`
		Surplus string  `gorm:"-" json:"surplus"`
	}
	CustomerField struct {
		CustomerName  string `gorm:"-"  json:"customerName"`
		CustomerPhone string `gorm:"-"  json:"customerPhone"`
	}
	GoodsField struct {
		GoodsName string `gorm:"-" json:"name"`
		GoodsTyp  int32  `gorm:"-" json:"type"`
	}
	PayField struct {
		PayFee  float64 `gorm:"-"  json:"payFee"`
		PayType int32   `gorm:"-"  json:"payType"`
		PaidFee float64 `gorm:"-"  json:"paidFee"`
	}
)

func (s *SurplusField) Set() {
	if s.Mount != 0 {
		s.Surplus = fmt.Sprintf("%d", s.Mount)
		return
	}
	s.Surplus = fmt.Sprintf("%.2f", s.Weight)
	return
}
