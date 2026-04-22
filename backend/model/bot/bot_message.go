package bot

import (
	"time"

	"gorm.io/datatypes"
)

type (
	TGMessageRecord struct {
		ID               int64          `json:"id" form:"id" gorm:"primaryKey;comment:消息ID(message_id);column:id"`
		ChatID           int64          `json:"chatId" form:"chatId" gorm:"comment:群聊ID;column:chat_id"`
		UserID           int64          `json:"userId" form:"userId" gorm:"comment:用户ID;column:user_id"`
		Username         string         `json:"username" form:"username" gorm:"type:varchar(64);comment:用户名;column:username"`
		FirstName        string         `json:"firstName" form:"firstName" gorm:"type:varchar(128);comment:名;column:first_name"`
		LastName         string         `json:"lastName" form:"lastName" gorm:"type:varchar(128);comment:姓;column:last_name"`
		NickName         string         `json:"nickName" form:"nickName" gorm:"type:varchar(64);comment:昵称;column:nick_name"`
		IsBot            bool           `json:"isBot" form:"isBot" gorm:"comment:是否机器人;column:is_bot"`
		ReplyToMessageID int64          `json:"replyToMessageId" form:"replyToMessageId" gorm:"comment:回复消息ID;column:reply_to_message_id"`
		MessageType      string         `json:"messageType" form:"messageType" gorm:"type:varchar(50);comment:消息类型;text/photo/video;column:message_type"`
		Text             string         `json:"text" form:"text" gorm:"type:text;comment:文本内容;column:text"`
		Caption          string         `json:"caption" form:"caption" gorm:"type:text;comment:媒体说明;column:caption"`
		FileID           string         `json:"fileId" form:"fileId" gorm:"type:text;comment:文件ID;column:file_id"`
		FileUniqueID     string         `json:"fileUniqueId" form:"fileUniqueId" gorm:"type:text;comment:文件唯一ID;column:file_unique_id"`
		FileType         string         `json:"fileType" form:"fileType" gorm:"type:varchar(32);comment:文件类型;column:file_type"`
		Timestamp        time.Time      `json:"timestamp" form:"timestamp" gorm:"comment:Telegram消息时间;column:timestamp"`
		Raw              datatypes.JSON `json:"raw" form:"raw" gorm:"type:jsonb;comment:Telegram原始JSON;column:raw"`
		CreatedAt        time.Time      `json:"createdAt" form:"createdAt" gorm:"comment:创建时间;column:created_at"`
	}

	TgChatMessageV1 struct {
		TGMessageRecord
		ReplyID          *int64  `json:"replyId"`
		ReplyUserID      *int64  `json:"replyUserId"`
		ReplyUsername    *string `json:"replyUsername"`
		ReplyText        *string `json:"replyText"`
		ReplyMessageType *string `json:"replyMessageType"`
		FileUrl          *string `json:"fileUrl" form:"fileUrl"`
	}
)
