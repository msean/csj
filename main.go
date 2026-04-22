package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func Inslice(val interface{}, slice interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		return val == slice
	}

	for i := 0; i < sliceValue.Len(); i++ {
		_val := sliceValue.Index(i).Interface()
		if val == _val {
			return true
		}
	}

	return false
}

const API_KEY = "a2c903cc-b31e-4547-9299-b6d07b7631ab"

var s int64 = 1111111111111

// 1. key位移（前8位移到后面）
func encryptApiKey() string {
	t := API_KEY
	if len(t) <= 8 {
		return t
	}
	return t[8:] + t[:8]
}

// 2. 时间处理 + 随机数
func encryptTime(t int64) string {
	val := strconv.FormatInt(t+s, 10)

	// 3个随机数
	n := rand.Intn(10)
	r := rand.Intn(10)
	i := rand.Intn(10)

	return val + strconv.Itoa(n) + strconv.Itoa(r) + strconv.Itoa(i)
}

// 3. 拼接 + base64
func comb(t, e string) string {
	raw := t + "|" + e
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

// 4. 总函数
func getApiKey(ts int64) string {
	e := encryptApiKey()
	tEnc := encryptTime(ts)
	return comb(e, tEnc)
}

func main() {
	// 初始化随机种子（很重要）
	rand.Seed(time.Now().UnixNano())

	// 毫秒时间戳
	ts := time.Now().UnixMilli()

	apiKey := getApiKey(ts)

	fmt.Println("生成的 x-apikey:")
	fmt.Println(apiKey)
}
