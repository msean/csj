package main

import (
	"app/service"

	"github.com/spf13/pflag"
)

var (
	// 配置文件路径 --config /path/path  或者 -c /path/path
	cfg = pflag.StringP("config", "c", "", "haihe orchestration-service config file path.")
)

func main() {
	pflag.Parse()
	service.Run(*cfg)
}
