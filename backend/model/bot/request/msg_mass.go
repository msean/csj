package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BotMsgMassSearch struct {
	BotID          *int64      `json:"botID" form:"botID"`
	ChatGroupID    *int64      `json:"chatGroupID" form:"chatGroupID"`
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}

type BotMsgMassSend struct {
	Msg     string `json:"msg" binding:"required"`
	IDs     []uint `json:"ids" binding:"required"`
	AtUsers bool   `json:"atUsers"`
}
