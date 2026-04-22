package config

type System struct {
	DbType             string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`    // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType            string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"` // Oss类型
	RouterPrefix       string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr               int    `mapstructure:"addr" json:"addr" yaml:"addr"` // 端口值
	LimitCountIP       int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP        int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseMultipoint      bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`                   // 多点登录拦截
	UseRedis           bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                                  // 使用redis
	UseStrictAuth      bool   `mapstructure:"use-strict-auth" json:"use-strict-auth" yaml:"use-strict-auth"`                // 使用树形角色分配模式
	DisableAutoMigrate bool   `mapstructure:"disable-auto-migrate" json:"disable-auto-migrate" yaml:"disable-auto-migrate"` // 自动迁移数据库表结构，生产环境建议设为false，手动迁移
	BotWeebhookPrefix  string `mapstructure:"bot-webhook-prefix" json:"bot-webhook-prefix" yaml:"bot-webhook-prefix"`       // 自动迁移数据库表结构，生产环境建议设为false，手动迁移
	Domain             string `mapstructure:"domain" json:"domain" yaml:"domain"`                                           // 自动迁移数据库表结构，生产环境建议设为false，手动迁移
	UnMatchPayment     bool   `mapstructure:"un-match-payment" json:"un-match-payment" yaml:"un-match-payment"`             // 自动迁移数据库表结构，生产环境建议设为false，手动迁移
	Env                string `mapstructure:"env" json:"env" yaml:"env"`
}
