package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ValidPhone(phone string) bool {
	if phone == "" {
		return false
	}
	return regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(phone)
}

func IsBlankString(in string) bool {
	return in == ""
}

func StringInSlice(s string, list []string) bool {
	for _, _s := range list {
		if s == _s {
			return true
		}
	}
	return false
}

func Split2Int64(s string, sep string) (int64s []int64) {
	splits := strings.Split(s, sep)
	for _, s := range splits {
		if num, err := strconv.ParseInt(s, 10, 64); err != nil {
			int64s = append(int64s, num)
		}
	}
	return
}

func FormatTime(time time.Time) string {
	if time.IsZero() {
		return ""
	}

	return time.Format("2006-01-02 15:04:05")
}

func FormatStruct(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	// remove space
	reg := regexp.MustCompile(`( )+|(\n)+`)
	return reg.ReplaceAllString(string(bytes), "$1$2")
}

func RemoveSpace(data string) string {
	reg := regexp.MustCompile(`( )+|(\n)+`)
	return reg.ReplaceAllString(data, "$1$2")
}

func IsValidURL(_url string) bool {
	parsedURL, err := url.Parse(_url)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
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

func ViolentJsonString(in any) (out string) {
	var err error
	var _bytes []byte
	if _bytes, err = json.Marshal(in); err != nil {
		return
	}
	return string(_bytes)
}

func ProductCode(isp, area int, value int) string {
	if area == 0 {
		return fmt.Sprintf("%d00%04d", isp, value)
	}
	return fmt.Sprintf("%d%d%04d", isp, area, value)
}

func StrToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func InCombineString(dst string, src string, symbol string) bool {
	stringList := strings.Split(src, symbol)
	for _, str := range stringList {
		if str == dst {
			return true
		}
	}
	return false
}
