package common

import "errors"

// login && register
var (
	PhoneCannotBeBlankErr = errors.New("手机号码不能为空")
	VerifyCodeErr         = errors.New("验证码错误")
	PhoneObejectExistErr  = errors.New("该手机号码被注册")
	PhoneUnRegisterErr    = errors.New("该手机号码还没有被注册")
	UnRegisterErr         = errors.New("该用户没有被注册")
	TokenUnValidErr       = errors.New("非法用户")
	TokenUnGenerateErr    = errors.New("token生成失败")
)

// customer
var (
	CustomerDuplicateErr = errors.New("已经存在相同的客户名")
)

// Goods
var (
	GoodsAlreadyExistErr         = errors.New("该货品已经存在")
	GoodsCategoryAlreadyExistErr = errors.New("The goods is already exist")
)

// batch
var (
	BatchDuplicateErr = errors.New("当天的批次已存在，不可重复创建")
)

// batch_order
var (
	BatchOrderUUIDRequireErr  = errors.New("单次UUID不能为空")
	BatchUUIDRequireErr       = errors.New("批次UUID不能为空")
	BatchOrderGoodsRequireErr = errors.New("开单货品不能为空")
)

// common
var (
	ObjectNotExistErr = errors.New("object not found")
	RequestUIDMustErr = errors.New("UUID不能为空")
)
