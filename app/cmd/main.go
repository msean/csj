package main

import (
	"app/service"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/pflag"
)

var (
	// 配置文件路径 --config /path/path  或者 -c /path/path
	cfg = pflag.StringP("config", "c", "", "haihe orchestration-service config file path.")
)

func main() {
	pflag.Parse()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	service.Run(*cfg)

}
