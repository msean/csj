package utils

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed payment.html.tmpl
var paymentTmpl string

// PaymentData is the data passed to the payment receipt HTML template.
type PaymentData struct {
	Title           string
	CustomerName    string
	Date            string
	Contact         string
	PreCreditAmount string
	PaymentAmount   string
	Discount        string // 空字符串时不显示优惠行
	CurrentAmount   string
}

// GeneratePaymentHTML renders the payment receipt template and returns HTML string.
func GeneratePaymentHTML(data PaymentData) (html string, err error) {
	tmpl, err := template.New("payment").Parse(paymentTmpl)
	if err != nil {
		return "", fmt.Errorf("parse payment template: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render payment template: %w", err)
	}

	return buf.String(), nil
}
