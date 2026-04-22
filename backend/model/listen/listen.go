package listen

import "time"

type BotChatMap struct {
	GroupID   int64     `gorm:"primaryKey;column:group_id" json:"group_id"`
	GroupType int16     `gorm:"column:group_type" json:"group_type"`
	GroupName string    `gorm:"column:group_name" json:"group_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (BotChatMap) TableName() string {
	return "bot_chat_map"
}

type BotMessageVO struct {
	ID        int64  `json:"id"`
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`

	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`

	MessageType string `json:"message_type"`
	Text        string `json:"text"`
	Caption     string `json:"caption"`

	Timestamp time.Time `json:"timestamp"`
}

type ListenQueryReq struct {
	// ===== 群 / 频道 =====
	GroupID   int64 `json:"groupId" form:"groupId" binding:"required"`
	GroupType int16 `json:"groupType" form:"groupType"` // 可选，1=群 2=频道

	// ===== 内容搜索 =====
	Keyword string `json:"keyword" form:"keyword"`

	// ===== 时间范围 =====
	StartTime *time.Time `json:"startTime" form:"startTime"`
	EndTime   *time.Time `json:"endTime" form:"endTime"`

	// ===== 分页 =====
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`

	// ===== 导出专用 =====
	IsExport bool `json:"isExport" form:"isExport"`
}
