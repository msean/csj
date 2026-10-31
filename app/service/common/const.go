package common

import shared "github.com/msean/csj/shared/constants"

const (
	GoodsTypeFix           = shared.GoodsTypeFix
	GoodsTypeBulk          = shared.GoodsTypeBulk
	BatchOrderTemp         = shared.BatchOrderTemp
	BatchOrderedCredit     = shared.BatchOrderedCredit
	BatchOrderFinish       = shared.BatchOrderFinish
	BatchOrderCancel       = shared.BatchOrderCancel
	BatchOrderRefund       = shared.BatchOrderRefund
	BatchOrderReTurn       = shared.BatchOrderReTurn
	BatchOrderUnshare      = shared.BatchOrderUnshare
	BatchOrderShared       = shared.BatchOrderShared
	PayTypeWx              = shared.PayTypeWx
	PayTypeZFB             = shared.PayTypeZFB
	PayTypeBank            = shared.PayTypeBank
	PayTypeCash            = shared.PayTypeCash
	BatchStatusOnSellering = shared.BatchStatusOnSellering
	BatchStatusOffShelf    = shared.BatchStatusOffShelf
	BatchStatusSettled     = shared.BatchStatusSettled
	CashCustomerName       = shared.CashCustomerName
	MemberTypeMonth        = shared.MemberTypeMonth
	MemberTypeQuarter      = shared.MemberTypeQuarter
	MemberTypeYear         = shared.MemberTypeYear
	PayDetailTypeDebt      = shared.PayDetailTypeDebt
	PayDetailTypeOrder     = shared.PayDetailTypeOrder
)

var (
	ExCludeTempBatchOrder = shared.ExCludeTempBatchOrder
	ValidOrder            = shared.ValidOrder
	UnFinalStatusOrder    = shared.UnFinalStatusOrder
)
