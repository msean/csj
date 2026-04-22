package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func Violent2String(in any) string {
	if f, ok := in.(float64); ok {
		return FloatReserveStr(f, 4)
	}
	if f, ok := in.(float32); ok {
		return FloatReserveStr(float64(f), 4)
	}

	return fmt.Sprintf("%v", in)
}

func FloatReserveStr(f float64, reserve int) string {
	formatted := strconv.FormatFloat(f, 'f', reserve, 64)
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return formatted
}

func StringsToIntsIgnoreError(strs []string) []int {
	ints := make([]int, 0, len(strs))
	for _, s := range strs {
		if n, err := strconv.Atoi(s); err == nil {
			ints = append(ints, n)
		}
	}
	return ints
}

func IntSliceToAnySlice(ints []int) []any {
	res := make([]any, len(ints))
	for i, v := range ints {
		res[i] = v
	}
	return res
}
