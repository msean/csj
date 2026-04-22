package config

import (
	"strings"
)

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	config := m.Config
	if strings.TrimSpace(config) == "" {
		config = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + config
}
