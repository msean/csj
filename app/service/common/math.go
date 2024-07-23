package common

import (
	"fmt"
	"math"
	"strconv"
)

func FloatGreat(a, b float32) bool {
	if math.Abs(float64(a-b)) < 0.0001 {
		return true
	}
	return a > b
}

func FloatEqual(a, b float32) bool {
	if math.Abs(float64(a-b)) < 0.0001 {
		return true
	}
	return false
}

func Float32IsZero(value float32) bool {
	const epsilon = 1e-6 // 定义一个小的阈值
	return math.Abs(float64(value)) < epsilon
}

func Float32Preserve(f float32, places int) string {
	return strconv.FormatFloat(float64(f), 'f', places, 32)
}

func Float32ToString(f float32) string {
	return fmt.Sprintf("%.f", f)
}
