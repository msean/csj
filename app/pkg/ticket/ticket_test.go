package ticket

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestReceiptTableRows(t *testing.T) {
	receipt := Receipt{Items: []Item{
		{Name: "土豆", Quantity: 2, QuantityUnit: "件", Price: 45.50, Amount: 91},
		{Name: "香菜", Weight: 8.5, Price: 2.30, Amount: 19.55},
		{Name: "大蒜", Quantity: 3, QuantityUnit: "包", Price: 12, Amount: 36},
	}}
	want := [][5]string{
		{"货品", "数量(件)", "重量(斤)", "单价(元)", "小记(元)"},
		{"土豆", "2", "--", "45.5", "91"},
		{"香菜", "--", "8.5", "2.3", "19.55"},
		{"大蒜", "3包", "--", "12", "36"},
		{"合计", "2 3包", "8.5", "", "146.55"},
	}
	got := receipt.tableRows()
	if len(got) != len(want) {
		t.Fatalf("rows: got %v, want %v", got, want)
	}
	for i, row := range want {
		if got[i] != row {
			t.Errorf("row %d: got %v, want %v", i, got[i], row)
		}
	}
	for _, i := range []int{0, 1} {
		single := Receipt{Items: receipt.Items[i : i+1]}
		rows := single.tableRows()
		total := rows[len(rows)-1]
		if total[1] != want[i+1][1] || total[2] != want[i+1][2] || total[3] != "" {
			t.Errorf("single item total: got %v", total)
		}
	}
}

func TestRasterBandPacking(t *testing.T) {
	bitmap := &raster{width: 2}
	bitmap.add(25)
	for _, point := range []struct {
		x, y  int
		alpha byte
	}{
		{0, 0, 88}, {0, 7, 255}, {0, 8, 255}, {1, 0, 87}, {1, 23, 255}, {1, 24, 255},
	} {
		bitmap.pixels[point.y*bitmap.width+point.x] = point.alpha
	}
	var output bytes.Buffer
	writeRasterBands(&output, bitmap)
	want := []byte{
		0x1b, 0x2a, 33, 2, 0, 0x81, 0x80, 0, 0, 0, 1,
		0x1b, 0x2a, 33, 2, 0, 0, 0, 0, 0x80, 0, 0,
	}
	if !bytes.Equal(output.Bytes(), want) {
		t.Fatalf("got %x, want %x", output.Bytes(), want)
	}
}

func TestReceiptImageBandsAndCut(t *testing.T) {
	var fontPath string
	for _, path := range []string{"/usr/share/fonts/noto/NotoSansCJK-Regular.ttc", "/System/Library/Fonts/PingFang.ttc"} {
		if _, err := os.Stat(path); err == nil {
			fontPath = path
			break
		}
	}
	if fontPath == "" {
		t.Skip("需要 Noto Sans CJK 或 PingFang 进行中文位图测试")
	}
	receipt := Receipt{
		Market: "农贸市场", Stall: "新鲜蔬菜", Customer: "王老板",
		Items:      []Item{{Name: "非常长的有机蔬菜品名", Quantity: 3, Price: 12, Amount: 36}},
		GoodsTotal: 36, PreviousCredit: 12, TotalCredit: 32, Received: 16,
		Contact: "13800138000",
	}
	for _, size := range []struct{ paper, width int }{{58, 384}, {80, 576}} {
		for _, ending := range []struct {
			feed int
			tail []byte
		}{
			{0, []byte{0x1b, 0x32, 0x1d, 0x56, 0x30}},
			{5, []byte{0x1b, 0x32, 0x1b, 0x64, 5}},
		} {
			options := Options{PaperWidth: size.paper, PrintMode: "image", FeedLines: ending.feed}
			raw, err := Render(receipt, options, fontPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(raw, []byte{0x1b, 0x40, 0x1d, 0x4c, 0, 0, 0x1b, 0x33, 0}) ||
				!bytes.HasSuffix(raw, ending.tail) {
				t.Fatal("invalid initialization or ending")
			}
			body := raw[len(initialize) : len(raw)-len(ending.tail)]
			bands, ink := 0, false
			bandBytes := 5 + size.width*3
			for len(body) > 0 {
				if len(body) < bandBytes || !bytes.Equal(body[:3], []byte{0x1b, 0x2a, 33}) {
					t.Fatal("invalid raster band")
				}
				if width := int(binary.LittleEndian.Uint16(body[3:5])); width != size.width {
					t.Fatalf("invalid raster width: %d", width)
				}
				ink = ink || len(bytes.Trim(body[5:bandBytes], "\x00")) > 0
				bands++
				body = body[bandBytes:]
			}
			if !ink || bands < 2 {
				t.Fatalf("missing visible raster bands: %d", bands)
			}
			reprint, err := Render(receipt, options, fontPath)
			if err != nil || !bytes.Equal(raw, reprint) {
				t.Fatal("same data and options must render identical bytes")
			}
		}
	}
}
