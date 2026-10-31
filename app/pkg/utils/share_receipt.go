package utils

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

//go:embed receipt.html.tmpl
var receiptTmpl string

// ReceiptGoodsItem is a view-layer item for the receipt template.
type ReceiptGoodsItem struct {
	GoodsName string
	MountStr  string // e.g. "60件" or "--"
	WeightStr string // e.g. "1.5斤" or "--"
	PriceStr  string
	TotalStr  string
	GoodType  int // 1-定装 2-散装；定装时重量列固定显示 "--"
}

// ReceiptOrder 单个订单的渲染数据（票号+时间+货品列表+三轮车费+合计）
type ReceiptOrder struct {
	SerialNo       string // 票号（订单号）
	Date           string // 开单时间
	GoodsList      []ReceiptGoodsItem
	TotalMount     string
	TotalWeight    string
	TotalAmountStr string
	TricycleFeeStr string
}

// ReceiptData is the data passed to the receipt HTML template.
type ReceiptData struct {
	Title                  string
	CustomerName           string
	SerialNo               string
	Date                   string
	Contact                string         // 联系方式（手机号）
	MarketName             string         // 市场名称（抬头第二行）
	UserName               string         // 用户名（抬头第三行）
	MarketAddress          string         // 市场地址（底部地址行）
	StallName              string         // 档口名称（抬头第四行）
	Orders                 []ReceiptOrder // 按订单分组，每单独立表格
	CreditAmountStr        string         // 今日赊欠（所有订单合计）
	TotalPreviousCreditStr string         // 前欠
	TotalDebtStr           string         // 总欠
}

// DartStringHashCode 复刻 Dart VM 中 String.hashCode 的算法：
// Jenkins one-at-a-time（按 UTF-16 code unit 累加），最后与 0x3fffffff 掩码，
// 保证结果落在 Dart Smi 范围内且恒为非负。
// 注意：Dart 的 String.hashCode 与 Java 的 String.hashCode 不是同一套算法。
func DartStringHashCode(s string) int32 {
	var h uint32
	for _, u := range utf16.Encode([]rune(s)) {
		h += uint32(u)
		h += h << 10
		h ^= h >> 6
	}
	h += h << 3
	h ^= h >> 11
	h += h << 15
	return int32(h & 0x3fffffff)
}

// UUIDToSerial 生成移动端展示的 6 位票号，与 Flutter 端保持一致：
// uuid.hashCode.abs().toString().padLeft(6, '0').substring(0, 6)
func UUIDToSerial(uuid string) string {
	n := int64(DartStringHashCode(uuid))
	// Dart hashCode 已做 30 位掩码，恒为非负，这里保留 .abs() 的语义
	if n < 0 {
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	if len(s) > 6 {
		s = s[:6]
	}
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

// GenerateHTML renders the receipt template and returns HTML string.
func GenerateHTML(data ReceiptData) (html string, err error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("receipt").Funcs(funcMap).Parse(receiptTmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}

	return buf.String(), nil
}

// GeneratePDFBytes renders the receipt template and converts it to PDF bytes via wkhtmltopdf.
// No files are persisted — a temp HTML file is written, wkhtmltopdf writes to stdout, and both are cleaned up.
func GeneratePDFBytes(data ReceiptData) (pdfBytes []byte, err error) {
	// 1. Render HTML into buffer
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("receipt").Funcs(funcMap).Parse(receiptTmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err = tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}

	// 2. Write temp HTML to OS temp dir (wkhtmltopdf needs a real file path for input)
	ts := strconv.FormatInt(time.Now().UnixNano(), 36)
	tmpDir := os.TempDir()
	htmlFile := filepath.Join(tmpDir, "receipt_"+ts+".html")
	if err = os.WriteFile(htmlFile, htmlBuf.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("write temp html: %w", err)
	}
	defer os.Remove(htmlFile)

	// 3. Run wkhtmltopdf, output to stdout (use "-" as pdf path)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("wkhtmltopdf",
		"--encoding", "utf-8",
		"--page-size", "A4",
		"--margin-top", "10mm",
		"--margin-bottom", "10mm",
		"--margin-left", "10mm",
		"--margin-right", "10mm",
		"--disable-smart-shrinking",
		"--no-background",
		"--quiet",
		htmlFile,
		"-", // output to stdout
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("wkhtmltopdf: %w — %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// FormatMount formats mount count into display string.
func FormatMount(mount int) string {
	if mount == 0 {
		return "--"
	}
	return fmt.Sprintf("%d件", mount)
}

// FormatWeight formats weight into display string.
func FormatWeight(weight float64) string {
	if weight == 0 {
		return "--"
	}
	return fmt.Sprintf("%.1f斤", weight)
}

// FormatTotalWeight formats the total weight for the receipt footer.
// Returns "--" when 0 (e.g. only 定装 goods, no bulk weight to sum).
func FormatTotalWeight(weight float64) string {
	if weight == 0 {
		return "--"
	}
	return FormatFloat(weight)
}

// FormatFloat formats a float64 to a clean display string (no trailing zeros).
func FormatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	// 去掉末尾多余的 0 和小数点，如 6.50 → 6.5，10.00 → 10
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// GenerateImageBytes renders the receipt template and converts it to PNG image bytes via wkhtmltoimage.
// Optimized for mobile sharing (WeChat, etc.) - returns PNG format.
func GenerateImageBytes(data ReceiptData) (imageBytes []byte, err error) {
	// 1. Render HTML into buffer
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("receipt").Funcs(funcMap).Parse(receiptTmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err = tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}

	// 2. Write temp HTML to OS temp dir
	ts := strconv.FormatInt(time.Now().UnixNano(), 36)
	tmpDir := os.TempDir()
	htmlFile := filepath.Join(tmpDir, "receipt_"+ts+".html")
	if err = os.WriteFile(htmlFile, htmlBuf.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("write temp html: %w", err)
	}
	defer os.Remove(htmlFile)

	// 3. Run wkhtmltoimage, output to stdout (use "-" as image path)
	// Mobile-optimized: 800px width, PNG format for WeChat sharing
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("wkhtmltoimage",
		"--encoding", "utf-8",
		"--width", "800", // Mobile-friendly width
		"--quality", "95", // High quality
		"--format", "png", // PNG for better quality
		"--no-background", // Clean background
		"--quiet",
		htmlFile,
		"-", // output to stdout
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("wkhtmltoimage: %w — %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
