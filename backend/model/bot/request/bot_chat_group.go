package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BotChatGroupSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	BotID          int         `json:"botID" form:"botID"`
	Name           string      `json:"chatGroupName" form:"chatGroupName"`
	request.PageInfo
}

type BotChatGroupClassifySearch struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type UpdateClassifyParams struct {
	ClassifyID   int   `json:"classifyID" form:"classifyID"`
	ChatGroupIDs []int `json:"chatGroupIDs" form:"chatGroupIDs"`
}

// api/dto/bot_chat_message.go
type ChatMessageQuery struct {
	BotID       int64 `json:"botID" form:"botID" binding:"required"`
	ChatGroupID int64 `json:"chatGroupID" form:"chatGroupID" binding:"required"`

	UserID   int64  `json:"userId" form:"userId"`
	Username string `json:"username" form:"username"`
	Text     string `json:"text" form:"text"`

	SrcStartTime string `json:"startTime" form:"startTime"`
	SrcEndTime   string `json:"endTime" form:"endTime"`

	StartTime *time.Time
	EndTime   *time.Time

	// 🔥 游标字段
	BeforeID int64 `json:"beforeId" form:"beforeId"` // 加载更早消息
	AfterID  int64 `json:"afterId" form:"afterId"`   // 加载更新消息

	Limit int `json:"limit" form:"limit"` // 默认 50
}
