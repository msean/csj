package ticket

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// The container installs this font. A standalone deployment can supply print.font.
var defaultRasterFont = sync.OnceValues(func() (*opentype.Font, error) {
	return loadRasterFont("/usr/share/fonts/noto/NotoSansCJK-Regular.ttc")
})

func loadRasterFont(path string) (*opentype.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	collection, err := opentype.ParseCollection(data)
	if err != nil {
		return nil, err
	}
	for index := 0; index < collection.NumFonts(); index++ {
		f, err := collection.Font(index)
		if err != nil {
			return nil, err
		}
		name, _ := f.Name(nil, sfnt.NameIDFamily)
		if strings.HasSuffix(name, " SC") {
			return f, nil
		}
	}
	return collection.Font(0)
}

type raster struct {
	width  int
	height int
	pixels []byte
}

type textStyle struct {
	face font.Face
	bold bool
}

func (s textStyle) width(text string) int {
	if text == "" {
		return 0
	}
	width := font.MeasureString(s.face, text).Ceil()
	if s.bold {
		width++
	}
	return width
}

func (r *raster) add(height int) *image.Alpha {
	start := len(r.pixels)
	r.pixels = append(r.pixels, make([]byte, height*r.width)...)
	r.height += height
	return &image.Alpha{Pix: r.pixels[start:], Stride: r.width, Rect: image.Rect(0, 0, r.width, height)}
}

func (r *raster) rule(thickness int) {
	dst := r.add(8 + thickness)
	for y := 4; y < 4+thickness; y++ {
		for x := 0; x < r.width-4; x++ {
			dst.SetAlpha(x, y, color.Alpha{A: 255})
		}
	}
}

func (s textStyle) draw(dst *image.Alpha, x int, text string) {
	metrics := s.face.Metrics()
	baseline := (dst.Rect.Dy()-metrics.Ascent.Ceil()-metrics.Descent.Ceil())/2 + metrics.Ascent.Ceil()
	drawer := font.Drawer{Dst: dst, Src: image.Opaque, Face: s.face, Dot: fixed.P(x, baseline)}
	drawer.DrawString(text)
	if s.bold {
		drawer.Dot = fixed.P(x+1, baseline)
		drawer.DrawString(text)
	}
}

func (s textStyle) height() int {
	m := s.face.Metrics()
	return m.Ascent.Ceil() + m.Descent.Ceil() + 2
}

func (s textStyle) fit(text string, width int) string {
	if s.width(text) <= width {
		return text
	}
	end := 0
	for i := range text {
		if s.width(text[:i]+"…") > width {
			break
		}
		end = i
	}
	return text[:end] + "…"
}

func (s textStyle) wrap(text string, width int) []string {
	var rows []string
	start, last := 0, 0
	for i := range text {
		if last > start && s.width(text[start:i]) > width {
			rows = append(rows, text[start:last])
			start = last
		}
		last = i
	}
	if last > start && s.width(text[start:]) > width {
		rows = append(rows, text[start:last])
		start = last
	}
	return append(rows, text[start:])
}

func (r *raster) text(text string, style textStyle, alignment int) {
	if text = clean(text); text == "" {
		return
	}
	width := r.width - 4
	for _, part := range style.wrap(text, width) {
		x := (width - style.width(part)) * alignment / 2
		style.draw(r.add(style.height()), x, part)
	}
}

func (r *raster) pair(left string, leftStyle textStyle, right string, rightStyle textStyle) {
	if left == "" && right == "" {
		return
	}
	width := r.width - 4
	rightX := width - rightStyle.width(right)
	if left != "" && right != "" && leftStyle.width(left)+10 > rightX {
		r.text(left, leftStyle, 0)
		r.text(right, rightStyle, 2)
		return
	}
	dst := r.add(max(leftStyle.height(), rightStyle.height()))
	leftStyle.draw(dst, 0, leftStyle.fit(left, width))
	rightStyle.draw(dst, max(0, rightX), rightStyle.fit(right, width))
}

func renderRaster(receipt Receipt, width int, fontPath string) (*raster, error) {
	var rasterFont *opentype.Font
	var err error
	if fontPath == "" {
		rasterFont, err = defaultRasterFont()
	} else {
		rasterFont, err = loadRasterFont(fontPath)
	}
	if err != nil {
		return nil, fmt.Errorf("load receipt font: %w", err)
	}
	regular, err := opentype.NewFace(rasterFont, &opentype.FaceOptions{Size: 22, DPI: 72})
	if err != nil {
		return nil, err
	}
	defer regular.Close()
	large, err := opentype.NewFace(rasterFont, &opentype.FaceOptions{Size: 28, DPI: 72})
	if err != nil {
		return nil, err
	}
	defer large.Close()
	tableSize := 20.0
	if width == 384 {
		tableSize = 16
	}
	table, err := opentype.NewFace(rasterFont, &opentype.FaceOptions{Size: tableSize, DPI: 72})
	if err != nil {
		return nil, err
	}
	defer table.Close()

	r := &raster{width: width}
	normal, bold, big := textStyle{face: regular}, textStyle{regular, true}, textStyle{large, true}
	if clean(receipt.Title) == "" {
		receipt.Title = "广东食用农产品批发市场销售票据"
	}
	r.text(receipt.Title, normal, 1)
	r.text(receipt.Market, normal, 1)
	r.text(receipt.Stall, big, 1)
	if receipt.Repayment > 0 {
		r.text("还 款 单", bold, 1)
	}
	if customer := clean(receipt.Customer); customer != "" {
		r.text("客户:"+customer, bold, 0)
	}
	r.rule(1)
	serial := clean(receipt.Serial)
	if serial != "" {
		serial = "票号:" + serial
	}
	r.pair(serial, bold, clean(receipt.Time), normal)
	r.table(receipt.tableRows(), table)

	amount := func(label string, value float64) {
		if value > 0 {
			r.text(label+money(value), normal, 2)
		}
	}
	amount("三轮车费: ", receipt.TricycleFee)
	amount("优惠: -", receipt.Discount)
	r.rule(2)
	if receipt.Repayment > 0 {
		amount("还款金额: ", receipt.Repayment)
	} else {
		amount("本单合计: ", receipt.GoodsTotal)
	}
	amount("前欠: ", receipt.PreviousCredit)
	amount("总欠: ", receipt.TotalCredit)
	amount("实收: ", receipt.Received)
	if method := clean(receipt.PaymentMethod); method != "" {
		r.text("支付方式: "+method, normal, 2)
	}

	r.rule(1)
	if receipt.Contact != "" {
		normal.draw(r.add(normal.height()), 0, "联系方式: "+receipt.Contact)
	}
	r.add(8)
	return r, nil
}

func (r *raster) table(rows [][5]string, face font.Face) {
	const gap = 10
	regular, bold := textStyle{face: face}, textStyle{face, true}
	var widths, minimum [5]int
	minimum[0] = bold.width("货品名称")
	for col := 1; col < len(widths); col++ {
		minimum[col] = regular.width(rows[0][col]) + 1
		widths[col] = minimum[col]
		for i := 1; i < len(rows); i++ {
			style := regular
			if i == len(rows)-1 {
				style = bold
			}
			widths[col] = max(widths[col], style.width(rows[i][col])+1)
		}
	}
	widths[0] = r.width - 4 - 4*gap
	for _, width := range widths[1:] {
		widths[0] -= width
	}
	// When data is too wide, wrap those cells while keeping the headers intact.
	for col := len(widths) - 1; col > 0 && widths[0] < minimum[0]; col-- {
		reclaim := min(minimum[0]-widths[0], widths[col]-minimum[col])
		widths[col] -= reclaim
		widths[0] += reclaim
	}
	for i, row := range rows {
		style := regular
		if i > 0 && i == len(rows)-1 {
			r.rule(1)
			style = bold
		}
		var cells [5][]string
		cells[0] = []string{style.fit(row[0], widths[0])}
		lineCount := 1
		for col := 1; col < len(row); col++ {
			cells[col] = style.wrap(row[col], widths[col])
			lineCount = max(lineCount, len(cells[col]))
		}
		for line := 0; line < lineCount; line++ {
			dst, x := r.add(style.height()), 0
			for col, parts := range cells {
				if line < len(parts) {
					offset := 0
					if col > 0 {
						offset = widths[col] - style.width(parts[line])
					}
					style.draw(dst, x+offset, parts[line])
				}
				x += widths[col] + gap
			}
		}
	}
}

func writeRasterBands(output *bytes.Buffer, bitmap *raster) {
	const bandHeight = 24
	width := bitmap.width
	for band := 0; band < bitmap.height; band += bandHeight {
		output.Write([]byte{0x1b, 0x2a, 33, byte(width), byte(width >> 8)})
		for x := 0; x < width; x++ {
			for group := 0; group < 3; group++ {
				var packed byte
				for bit := 0; bit < 8; bit++ {
					y := band + group*8 + bit
					// Black drawn on white has gray = 255 - alpha.
					if y < bitmap.height && bitmap.pixels[y*width+x] >= 88 {
						packed |= 0x80 >> bit
					}
				}
				output.WriteByte(packed)
			}
		}
	}
}
