package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func FloatReserveStr(f float64, reserve int) string {
	formatted := strconv.FormatFloat(f, 'f', reserve, 64)
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return formatted
}

func Violent2String(in any) string {
	if f, ok := in.(float64); ok {
		return FloatReserveStr(f, 4)
	}
	if f, ok := in.(float32); ok {
		return FloatReserveStr(float64(f), 4)
	}

	return fmt.Sprintf("%v", in)
}
