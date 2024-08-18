package pkg

func TestAliSender() {
	ali := AliPlatfrom{
		Uid:      "msean",
		Password: "a5375302",
	}
	m := SmsMsg{
		TemplateCode: "100006",
		Mobile:       "15112534872",
		TemplateJson: map[string]any{
			"code": "123456",
		},
	}
	// ali.Send(m)
	_, _ = ali, m
	ali.Send(m)
}
