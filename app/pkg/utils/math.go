package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func FloatGreat(a, b float64) bool {
	if math.Abs(float64(a-b)) < 0.0001 {
		return true
	}
	return a > b
}

func FloatEqual(a, b float64) bool {
	if math.Abs(float64(a-b)) < 0.0001 {
		return true
	}
	return false
}

func IntSecurityDiv(numerator int, denominator int, reserve uint) (_f float64) {
	if denominator == 0 {
		return 0
	}
	return FloatReserve(float64(numerator)/float64(denominator), reserve)
}

func Float64SecurityDiv(numerator float64, denominator float64, reserve uint) (_f float64) {
	if denominator == 0 {
		return 0
	}
	return FloatReserve(float64(numerator)/float64(denominator), reserve)
}

func FloatReserve(f float64, reserve uint) float64 {
	factor := math.Pow(10, float64(reserve))
	rounded := math.Round(f*factor) / factor

	roundedStr := fmt.Sprintf("%.*f", reserve, rounded)

	_f, _ := strconv.ParseFloat(roundedStr, 64)
	return _f
}

func FloatReserveString(f float64, reserve int) string {
	formatted := strconv.FormatFloat(f*100, 'f', reserve, 64)
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return formatted + "%"
}

func FloatReserveStr(f float64, reserve int) string {
	formatted := strconv.FormatFloat(f, 'f', reserve, 64)
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return formatted
}
