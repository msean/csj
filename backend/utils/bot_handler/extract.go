package bot_handler

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/msean/csj/backend/global/constant"
	"github.com/msean/csj/backend/model/bot"
	"github.com/msean/csj/backend/service/cache"
	botapi "github.com/msean/csj/backend/utils/bot_handler/bot_api"
	"golang.org/x/net/html"
)

type MediaItem struct {
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	FileID string `json:"file_id,omitempty"`
}

func ExtractImgsAndText(htmlStr string) (imgs []string, textWithoutImgs string) {
	// 找所有 <img ... src="..."> 并取 src
	imgRe := regexp.MustCompile(`(?i)<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches := imgRe.FindAllStringSubmatch(htmlStr, -1)
	for _, m := range matches {
		if len(m) >= 2 {
			imgs = append(imgs, m[1])
		}
	}
	// 去掉所有 <img> 标签，保留其余 HTML
	textWithoutImgs = imgRe.ReplaceAllString(htmlStr, "")
	// 去掉空的 <p></p> 等多余空白（可选）
	// 将 HTML 实体转回正常字符（比如 &gt; -> >）
	textWithoutImgs = html.UnescapeString(textWithoutImgs)
	textWithoutImgs = strings.TrimSpace(textWithoutImgs)
	return
}

// 辅助：创建 InlineKeyboardMarkup（如果没有按钮则返回 nil）
func BuildMarkupFromExtrend(raw json.RawMessage) *botapi.InlineKeyboardMarkup {
	if len(raw) == 0 {
		return nil
	}
	var btns []bot.ButtonItem
	if err := json.Unmarshal(raw, &btns); err != nil || len(btns) == 0 {
		return nil
	}
	var row []botapi.InlineKeyboardButton
	for _, b := range btns {
		// 只创建 URL 按钮（和你之前逻辑保持一致）
		row = append(row, botapi.NewInlineKeyboardButtonURL(b.Name, b.URL))
	}
	m := botapi.NewInlineKeyboardMarkup(row)
	return &m
}

func CleanHTMLForTelegram(htmlStr string) string {
	if htmlStr == "" {
		return ""
	}

	htmlStr = regexp.MustCompile(`(?i)<p[^>]*>`).ReplaceAllString(htmlStr, "")
	htmlStr = regexp.MustCompile(`(?i)</p>`).ReplaceAllString(htmlStr, "\n")
	htmlStr = regexp.MustCompile(`(?i)<div[^>]*>`).ReplaceAllString(htmlStr, "")
	htmlStr = regexp.MustCompile(`(?i)</div>`).ReplaceAllString(htmlStr, "\n")
	htmlStr = regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(htmlStr, "\n")

	allowed := map[string]bool{
		"b": true, "i": true, "strong": true, "em": true,
		"u": true, "s": true, "strike": true, "del": true,
		"a": true, "code": true, "pre": true,
	}

	tagRe := regexp.MustCompile(`(?i)</?([a-z0-9]+)[^>]*>`)
	htmlStr = tagRe.ReplaceAllStringFunc(htmlStr, func(tag string) string {
		m := tagRe.FindStringSubmatch(tag)
		if len(m) < 2 {
			return ""
		}
		name := strings.ToLower(m[1])
		if allowed[name] {
			return tag
		}
		return ""
	})

	htmlStr = html.UnescapeString(htmlStr)

	htmlStr = regexp.MustCompile(`\n{2,}`).ReplaceAllString(htmlStr, "\n\n")
	htmlStr = strings.TrimSpace(htmlStr)
	return htmlStr
}

func ExtractVideosFromHTML(content string) []string {
	var urls []string
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return urls
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "source" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					urls = append(urls, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls
}

func HandleTexWithMarup(chatID int64, token string, content string, markup any) (err error) {
	botAPI, err := botapi.NewBotAPI(token)
	if err != nil {
		return
	}

	// 处理内容
	imgs, text := ExtractImgsAndText(content)
	videos := ExtractVideosFromHTML(content)
	caption := CleanHTMLForTelegram(text)
	if len(caption) > 1024 {
		caption = caption[:1020] + "..."
	}

	// 处理 ReplyMarkup（支持 InlineKeyboard 和 ReplyKeyboard）
	var sendMarkup interface{}
	if markup != nil {
		switch kb := markup.(type) {
		case botapi.InlineKeyboardMarkup:
			if len(kb.InlineKeyboard) > 0 {
				sendMarkup = kb
			}
		case botapi.ReplyKeyboardMarkup:
			if len(kb.Keyboard) > 0 {
				sendMarkup = kb
			}
		}
	}

	firstMessage := true // 标记是否为第一条消息

	// 发送图片
	for _, img := range imgs {
		photo := botapi.NewPhoto(chatID, botapi.FileURL(img))
		if firstMessage {
			photo.Caption = caption
			photo.ParseMode = botapi.ModeHTML
			if sendMarkup != nil {
				photo.ReplyMarkup = sendMarkup
			}
			firstMessage = false
		}
		if _, err = botAPI.Send(photo); err != nil {
			return
		}
	}

	// 发送视频
	for _, vid := range videos {
		video := botapi.NewVideo(chatID, botapi.FileURL(vid))
		if firstMessage {
			video.Caption = caption
			video.ParseMode = botapi.ModeHTML
			if sendMarkup != nil {
				video.ReplyMarkup = sendMarkup
			}
			firstMessage = false
		}
		if _, err = botAPI.Send(video); err != nil {
			return
		}
	}

	// 发送纯文本
	if len(imgs) == 0 && len(videos) == 0 && caption != "" {
		msg := botapi.NewMessage(chatID, caption)
		msg.ParseMode = botapi.ModeHTML
		if sendMarkup != nil {
			msg.ReplyMarkup = sendMarkup
		}
		_, err = botAPI.Send(msg)
	}

	return err
}

func ParseContentFromCfg(cfg cache.BotCmdCache, buttonType int) (markup any) {
	switch buttonType {
	case constant.ButtonTypeKeyBoard: // 普通键盘（ReplyKeyboard）
		var keyboard [][]botapi.KeyboardButton

		if len(cfg.CmdButtons) > 0 {
			var buttons [][]struct {
				Name    string `json:"name"`
				BindCmd string `json:"bindCmd"`
			}
			_ = json.Unmarshal([]byte(cfg.CmdButtons), &buttons)

			for _, row := range buttons {
				kbRow := []botapi.KeyboardButton{}
				for _, btn := range row {
					kbRow = append(kbRow, botapi.NewKeyboardButton(btn.Name))
				}
				keyboard = append(keyboard, kbRow)
			}

			// 创建 ReplyKeyboardMarkup
			replyKeyboard := botapi.ReplyKeyboardMarkup{
				Keyboard:        keyboard,
				ResizeKeyboard:  true,
				OneTimeKeyboard: false,
			}
			markup = replyKeyboard
		}

	case constant.ButtonTypeInline: // 内联键盘（InlineKeyboard）
		var rows [][]struct {
			Name    string `json:"name"`
			BindCmd string `json:"bindCmd"`
		}
		_ = json.Unmarshal([]byte(cfg.CmdButtons), &rows)

		inlineRows := make([][]botapi.InlineKeyboardButton, 0)
		for _, row := range rows {
			btnRow := make([]botapi.InlineKeyboardButton, 0)
			for _, b := range row {
				var btn botapi.InlineKeyboardButton
				if strings.HasPrefix(b.BindCmd, "http://") || strings.HasPrefix(b.BindCmd, "https://") {
					btn = botapi.NewInlineKeyboardButtonURL(b.Name, b.BindCmd)
				} else {
					btn = botapi.NewInlineKeyboardButtonData(b.Name, b.BindCmd)
				}

				btnRow = append(btnRow, btn)
			}
			inlineRows = append(inlineRows, btnRow)
		}
		markup = botapi.NewInlineKeyboardMarkup(inlineRows...)
	}
	return
}

func GetChatUserID(update botapi.Update) (userId int64) {
	switch {
	case update.Message != nil:
		// 如果是用户发送的消息
		userId = int64(update.Message.From.ID)
	case update.CallbackQuery != nil:
		// 如果是用户点击了按钮
		userId = int64(update.CallbackQuery.From.ID)
	default:
		// 其他情况
		userId = 0
	}
	return
}

func GetUserName(update botapi.Update) (userName string) {
	var from *botapi.User

	switch {
	case update.Message != nil:
		from = update.Message.From

	case update.CallbackQuery != nil:
		from = update.CallbackQuery.From
	}

	if from == nil {
		return "Unknown"
	}

	// 优先使用 username
	if from.UserName != "" {
		userName = from.UserName
		return
	}

	// fallback 使用 姓名
	userName = strings.TrimSpace(from.FirstName + " " + from.LastName)

	if userName == "" {
		userName = "Unknown"
	}

	return
}

func GetChatID(update botapi.Update) (chatID int64) {
	switch {
	case update.Message != nil:
		// 如果是用户发送的消息
		chatID = update.Message.Chat.ID
	case update.CallbackQuery != nil:
		// 如果是用户点击了按钮
		chatID = update.CallbackQuery.Message.Chat.ID
	default:
		// 其他情况
		chatID = 0
	}
	return
}
