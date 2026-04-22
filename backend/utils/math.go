package utils

import (
	"fmt"
	"math"
	"strconv"
)

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

func Abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
