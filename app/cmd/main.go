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
	// 是否执行数据库迁移 --migrate 或 -m
	migrate = pflag.BoolP("migrate", "m", false, "enable database migration (default: false)")
	// 是否执行分表迁移 --shard-migrate 或 -s
	shardMigrate = pflag.BoolP("shard-migrate", "s", false, "enable sharding table migration (default: false)")
)

func main() {
	pflag.Parse()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	service.Run(*cfg, *migrate, *shardMigrate)

}
