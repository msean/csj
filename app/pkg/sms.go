package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	SmsMsg struct {
		TemplateCode string
		Mobile       string
		TemplateJson map[string]any
	}
	AliPlatfrom struct {
		Password string
		Uid      string
		Secret   string
	}
	SmsSender interface {
		Send(msg SmsMsg) error
	}
)

// http://api.sms.cn/sms/?ac=send
// &uid=msean&pwd=接口密码[获取密码]
// &template=100006&mobile=填写要发送的手机号
// &content={"code":"value"}
func (ap *AliPlatfrom) Send(msg SmsMsg) (err error) {
	params := url.Values{}
	params.Set("ac", "send")
	params.Set("uid", ap.Uid)
	pwd := md5.Sum([]byte(ap.Password + ap.Uid))
	params.Set("pwd", hex.EncodeToString(pwd[:]))
	params.Set("template", msg.TemplateCode)
	params.Set("mobile", msg.Mobile)

	var data []byte
	data, err = json.Marshal(msg.TemplateJson)
	if err != nil {
		fmt.Println("转换为 JSON 字符串时发生错误:", err)
		return
	}
	params.Set("content", string(data))
	apiURL := "http://api.sms.cn/sms/?" + params.Encode()
	fmt.Println(apiURL)

	_, err = http.Get(apiURL)
	return
}
