package model

// ShareRecord 客户分享记录
// 分享以「客户 + 当天」为维度：同一客户当天分享的多个订单只记一条记录，
// 仅保存最近一次分享时间，用于判断客户当天是否已点击分享
type ShareRecord struct {
	BaseModel
	OwnerUser    string `gorm:"column:owner_user;index:idx_owner_customer,unique;comment:所属用户" json:"ownerUser"`
	CustomerUUID string `gorm:"column:customer_uuid;index:idx_owner_customer,unique;comment:客户UUID" json:"customerUUID"`
	ShareTime    int64  `gorm:"column:share_time;index;comment:最近一次分享时间(unix秒)" json:"shareTime"`
}
