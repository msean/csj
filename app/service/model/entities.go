// Entity definitions live in the shared module. Aliases preserve app imports and methods.
package model

import shared "github.com/msean/csj/shared/model"

type (
	MemberConfig        = shared.MemberConfig
	ShareRecord         = shared.ShareRecord
	AppVersion          = shared.AppVersion
	User                = shared.User
	RechargeOrder       = shared.RechargeOrder
	BatchSerialCounter  = shared.BatchSerialCounter
	BatchOrderHistory   = shared.BatchOrderHistory
	Step                = shared.Step
	StepGoods           = shared.StepGoods
	BatchOrder          = shared.BatchOrder
	BatchOrderGoods     = shared.BatchOrderGoods
	BatchOrderOpr       = shared.BatchOrderOpr
	BatchOrderPay       = shared.BatchOrderPay
	PayDetail           = shared.PayDetail
	GoodsCategory       = shared.GoodsCategory
	Goods               = shared.Goods
	Batch               = shared.Batch
	BatchGoods          = shared.BatchGoods
	BaseModel           = shared.BaseModel
	Array               = shared.Array
	Map                 = shared.Map
	Ints                = shared.Ints
	Sms                 = shared.Sms
	BatchGoodsStatFeild = shared.BatchGoodsStatFeild
	CustomerFeild       = shared.CustomerFeild
	GoodsFeild          = shared.GoodsFeild
	PayFeild            = shared.PayFeild
	BatchStatFeild      = shared.BatchStatFeild
	MessageCenter       = shared.MessageCenter
	OfficialMessage     = shared.OfficialMessage
	OfficialMessageRead = shared.OfficialMessageRead
	Customer            = shared.Customer
)

const (
	HistoryStepOrder      = shared.HistoryStepOrder
	HistoryStepOrderFix   = shared.HistoryStepOrderFix
	HistoryStepOrderShare = shared.HistoryStepOrderShare
	HistoryStepCredit     = shared.HistoryStepCredit
	HistoryStepPay        = shared.HistoryStepPay
	HistoryStepCash       = shared.HistoryStepCash
	HistoryStepCrash      = shared.HistoryStepCrash
)
