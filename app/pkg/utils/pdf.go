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
	"time"
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
}

// ReceiptData is the data passed to the receipt HTML template.
type ReceiptData struct {
	Title                  string
	CustomerName           string
	SerialNo               string
	Date                   string
	Contact                string
	GoodsList              []ReceiptGoodsItem
	TotalMount             string
	TotalWeight            string
	TotalAmountStr         string
	CreditAmountStr        string
	TotalPreviousCreditStr string
	TotalDebtStr           string
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

// FormatFloat formats a float64 to a clean display string (no trailing .00).
func FormatFloat(f float64) string {
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
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
