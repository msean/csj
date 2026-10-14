package model

type MessageCenter struct {
	BaseModel
	OwnerUser    string `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
	CustomerUUID string `gorm:"column:customer_uuid;comment:客户UUID" json:"customerUUID"`
	Type         int    `gorm:"column:type;comment:消息类型 1-还款 2-撤销还款" json:"type"`
	Event        string `gorm:"column:event;comment:事件名称" json:"event"`
	Content      string `gorm:"column:content;type:text;comment:消息内容" json:"content"`
	IsRead       int    `gorm:"column:is_read;default:0;comment:是否已读 0-未读 1-已读" json:"isRead"`
	RelatedUUID  string `gorm:"column:related_uuid;comment:关联UUID（还款记录UUID）" json:"relatedUUID"`
}
