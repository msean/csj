package main

import (
	"app/pkg"
	"encoding/json"
	"fmt"
)

func main() {
	ali := pkg.AliPlatfrom{
		Uid:      "msean",
		Password: "a5375302",
	}
	m := pkg.SmsMsg{
		TemplateCode: "100006",
		Mobile:       "15112534872",
		TemplateJson: map[string]any{
			"code": "123456",
		},
	}
	// ali.Send(m)
	_, _ = ali, m
	// 定义包含键值对的 map
	data := map[string]string{
		"code": "value",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("转换为 JSON 字符串时发生错误:", err)
		return
	}

	// 将 JSON 字符串输出
	fmt.Println(string(jsonData))
}
