package bot_handler

import (
	"strings"
	"time"

	BotAPI "github.com/msean/csj/backend/utils/bot_handler/bot_api"
)

type Bot struct {
	token string
	*BotAPI.BotAPI
}

func NewBot(token string) (bot *Bot, err error) {
	bot = &Bot{
		token: token,
	}
	if bot.BotAPI, err = BotAPI.NewBotAPI(token); err != nil {
		return
	}
	return
}

func (b *Bot) BanUser(chatID, userID int64, until int64) (err error) {
	// until := time.Now().Add(duration).Unix()
	// global.GVA_LOG.Info("util", zap.Int("util", int(until)))

	cfg := BotAPI.RestrictChatMemberConfig{
		ChatMemberConfig: BotAPI.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		Permissions: &BotAPI.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
		},
		UntilDate: until,
	}
	_, err = b.Request(cfg)
	return err
}

func (b *Bot) UnMuteUser(chatID int64, userID int64) error {
	cfg := BotAPI.RestrictChatMemberConfig{
		ChatMemberConfig: BotAPI.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: time.Now().Unix(),
		Permissions: &BotAPI.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
			CanChangeInfo:         true,
			CanInviteUsers:        true,
			CanPinMessages:        true,
		},
	}
	_, err := b.Request(cfg)
	return err
}

func RegisterWebhook(botToken, webhookURL string) error {
	bot, err := BotAPI.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	cfg := BotAPI.DeleteWebhookConfig{
		DropPendingUpdates: true, // true 表示丢弃所有未处理的消息
	}

	if _, err = bot.Request(cfg); err != nil {
		return err
	}

	wh, err := BotAPI.NewWebhook(webhookURL)
	if err != nil {
		return err
	}

	_, err = bot.Request(wh)
	if err != nil {
		return err
	}
	return nil
}

func UnRegisterWebhook(botToken string, dropPending bool) error {
	bot, err := BotAPI.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	cfg := BotAPI.DeleteWebhookConfig{
		DropPendingUpdates: dropPending, // true 表示丢弃所有未处理的消息
	}

	_, err = bot.Request(cfg)
	return err
}

func (b *Bot) DeleteMsg(chatID int64, msgID int) (err error) {
	cfg := BotAPI.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}
	_, err = b.Request(cfg)
	return
}

func (b *Bot) SendTextMessage(chatID int64, text string) (err error) {
	msg := BotAPI.NewMessage(chatID, text)
	_, err = b.Send(msg)
	return
}

func (b *Bot) SendMarkDownMessage(chatID int64, text string, button any) (err error) {
	msg := BotAPI.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = button
	_, err = b.Send(msg)
	return
}

func EscapeMarkdownV2CodeBlock(text string) string {
	specialChars := []string{"`", "\\"}
	for _, ch := range specialChars {
		text = strings.ReplaceAll(text, ch, "\\"+ch)
	}
	return text
}

// EscapeMarkdownV2 用于普通文本，转义 MarkdownV2 保留字符
func EscapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}",
		".", "!", // 普通文本里 . 和 ! 需要转义
	}
	for _, ch := range specialChars {
		text = strings.ReplaceAll(text, ch, "\\"+ch)
	}
	return text
}

func (b *Bot) Send(c BotAPI.Chattable) (msg BotAPI.Message, err error) {
	return b.BotAPI.Send(c)
}

func (b *Bot) DeleteOriginMessage(chatID int64, msgID int) error {
	deleteMsg := BotAPI.NewDeleteMessage(chatID, msgID)
	_, err := b.Send(deleteMsg)
	return err
}

// SendAdMessage 统一发送广告内容（文字 / 图片 + 文本 / 视频 + 文本）
// 如果 replyMarkup != nil，则用作按钮，否则不带按钮
func (b *Bot) TgSend(chatID int64, medias []MediaItem, replyMarkup interface{}) (tgMsg BotAPI.Message, err error) {
	var caption string
	var photoID string
	var videoID string

	// 整理用户发送的数据
	for _, m := range medias {
		switch m.Type {
		case "text":
			caption = m.Text
		case "photo":
			photoID = m.FileID
		case "video":
			videoID = m.FileID
		}
	}

	// PHOTO
	if photoID != "" {
		msg := BotAPI.NewPhoto(chatID, BotAPI.FileID(photoID))
		msg.Caption = caption

		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}

		return b.Send(msg)
	}

	// VIDEO
	if videoID != "" {
		msg := BotAPI.NewVideo(chatID, BotAPI.FileID(videoID))
		msg.Caption = caption

		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}

		return b.Send(msg)
	}

	// ONLY TEXT
	msg := BotAPI.NewMessage(chatID, caption)

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	return b.Send(msg)
}
