package conf

import (
	"gogo/core/config"
)

var Config = &Configuration{}

type Configuration struct {
	UserTableName string
	VerifyCode    bool
	config.Auth
	config.Server
	config.Mysql
	config.Redis
	config.Consul
}
