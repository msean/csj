// Package ticket renders sales receipts as ESC/POS printer data.
package ticket

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Item struct {
	Name, QuantityUnit string
	Quantity           int
	Weight, Price      float64
	Amount             float64
}

type Receipt struct {
	Title, Market, Stall, Customer, Serial, Time, Contact string
	Items                                                 []Item
	TricycleFee, Discount, GoodsTotal                     float64
	Repayment, PreviousCredit                             float64
	TotalCredit, Received                                 float64
	PaymentMethod                                         string
}

type Options struct {
	PaperWidth int
	PrintMode  string
	FeedLines  int
}

var initialize = []byte{0x1b, 0x40, 0x1d, 0x4c, 0, 0, 0x1b, 0x33, 0}

func Render(receipt Receipt, options Options, fontPath string) ([]byte, error) {
	var dots int
	switch options.PaperWidth {
	case 58:
		dots = 384
	case 80:
		dots = 576
	default:
		return nil, fmt.Errorf("unsupported paper width: %d", options.PaperWidth)
	}
	if options.PrintMode != "image" {
		return nil, fmt.Errorf("unsupported print mode: %q", options.PrintMode)
	}
	if options.FeedLines != 0 && options.FeedLines != 5 {
		return nil, fmt.Errorf("unsupported feed lines: %d", options.FeedLines)
	}

	bitmap, err := renderRaster(receipt, dots, fontPath)
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	output.Write(initialize)
	writeRasterBands(&output, bitmap)
	output.Write([]byte{0x1b, 0x32})
	if options.FeedLines == 0 {
		output.Write([]byte{0x1d, 0x56, 0x30})
	} else {
		output.Write([]byte{0x1b, 0x64, byte(options.FeedLines)})
	}
	return output.Bytes(), nil
}

func (r Receipt) tableRows() [][5]string {
	rows := [][5]string{{"货品", "数量(件)", "重量(斤)", "单价(元)", "小记(元)"}}
	quantities := make(map[string]int)
	var weight, amount float64
	var hasWeight bool
	for _, item := range r.Items {
		quantity, weightText := "--", "--"
		if item.Quantity != 0 {
			unit := clean(item.QuantityUnit)
			if unit == "件" {
				unit = ""
			}
			quantity = strconv.Itoa(item.Quantity) + unit
			quantities[unit] += item.Quantity
		}
		if item.Weight != 0 {
			weightText = money(item.Weight)
			weight += item.Weight
			hasWeight = true
		}
		rows = append(rows, [5]string{
			clean(item.Name), quantity, weightText, money(item.Price), money(item.Amount),
		})
		amount += item.Amount
	}
	if len(r.Items) == 0 {
		return rows
	}
	units := make([]string, 0, len(quantities))
	for unit := range quantities {
		units = append(units, unit)
	}
	sort.Strings(units)
	totals := make([]string, 0, len(units))
	for _, unit := range units {
		totals = append(totals, strconv.Itoa(quantities[unit])+unit)
	}
	totalQuantity, totalWeight := "--", "--"
	if len(totals) > 0 {
		totalQuantity = strings.Join(totals, " ")
	}
	if hasWeight {
		totalWeight = money(weight)
	}
	rows = append(rows, [5]string{"合计", totalQuantity, totalWeight, "", money(amount)})
	return rows
}

func money(value float64) string {
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 2, 64), "0"), ".")
}

func clean(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, strings.TrimSpace(value))
}
