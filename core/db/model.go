package db

import "time"

var All = &Configuration{}

// Server 服务器配置
type Server struct {
	Namespace string
	Name      string
	Port      string
}

// Mysql MySQL配置结构体
type Mysql struct {
	Address  string
	Port     string
	Username string
	Password string
	Database string
}

// Redis Redis配置结构体
type Redis struct {
	Address     string
	Port     string
	Password string
	Database int
}

type Consul struct {
	Address string
	Port    int
	Config  string
}

type Auth struct {
	TokenLife int
	JwtSecret string
	JwtExpire time.Duration
}

type Configuration struct {
	Server
	Mysql
	Redis
	Consul
	Auth
}
