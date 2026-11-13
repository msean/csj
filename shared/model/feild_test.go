package model

import (
	"testing"

	common "github.com/msean/csj/shared/constants"
)

func TestBatchGoodsStatAveragePriceUsesSoldQuantity(t *testing.T) {
	for _, tc := range []struct {
		name string
		stat BatchGoodsStatFeild
		want string
	}{
		{"bulk with negative stock", BatchGoodsStatFeild{
			Type: common.GoodsTypeBulk, SurplusWeight: -25, SellWeight: 100, SellTotal: 160,
		}, "1.6"},
		{"bulk with zero stock", BatchGoodsStatFeild{
			Type: common.GoodsTypeBulk, SellWeight: 100, SellTotal: 160,
		}, "1.6"},
		{"bulk with no sales", BatchGoodsStatFeild{
			Type: common.GoodsTypeBulk, SurplusWeight: 100,
		}, ""},
		{"fixed by piece", BatchGoodsStatFeild{
			Type: common.GoodsTypeFix, SellMount: 5, SellWeight: 100, SellTotal: 160,
		}, "32"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.stat.Set()
			if tc.stat.SellAvgPrice != tc.want {
				t.Errorf("SellAvgPrice = %q, want %q", tc.stat.SellAvgPrice, tc.want)
			}
		})
	}
}
