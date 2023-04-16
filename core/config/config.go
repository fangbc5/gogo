package config

import (
	"fmt"

	source "github.com/go-micro/plugins/v4/config/source/consul"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
)

const (
	Profile      = "dev"
	Name         = "myservice"
	Port         = "38080"
	Version      = "1.0.0"
	ConsulAddr   = "127.0.0.1:8500"
	ConsulPrefix = "/micro/config"
	ENV          = "env"
)

type Config struct {
	Env      string
	Profile  string
	Server   ServerConfig
	Auth     AuthConfig
	Consul   ConsulConfig
	Tracing  TracingConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

type ServerConfig struct {
	Name    string
	Port    string
	Version string
}

type AuthConfig struct {
	JwtSecret string
	TokenLife string
}

type ConsulConfig struct {
	Addr   string
	Prefix string
}

type TracingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

type DatabaseConfig struct {
	Address  string
	Port     string
	Username string
	Password string
	Database string
}

type CacheConfig struct {
	Address  string
	Port     string
	Database string
	Password string
}

var cfg *Config = &Config{
	Env:     ENV,
	Profile: Profile,
	Server: ServerConfig{
		Name:    Name,
		Port:    Port,
		Version: Version,
	},
	Consul: ConsulConfig{
		Addr:   ConsulAddr,
		Prefix: ConsulPrefix + "/" + Profile + "/" + Name,
	},
}

func Get() Config {
	return *cfg
}

func WithConsul(addr string, prefix string) {
	cfg.Consul.Addr = addr
	cfg.Consul.Prefix = prefix
}

func SetEnv(env string) {
	cfg.Env = env
}

func GetName() string {
	return cfg.Server.Name
}

func GetVersion() string {
	return cfg.Server.Version
}

func GetAddress() string {
	return fmt.Sprintf(":%v", cfg.Server.Port)
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

func Load() error {
	//加载profile
	configor, err := config.NewConfig(config.WithSource(source.NewSource(
		source.WithAddress(cfg.Consul.Addr),
		source.WithPrefix(ConsulPrefix),
		source.StripPrefix(true),
	)))
	if err != nil {
		return errors.Wrap(err, "configor.New")
	}
	if err := configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load")
	}
	if err := configor.Get(cfg.Env).Scan(cfg); err != nil {
		return errors.Wrap(err, "configor.Scan")
	}
	configor.Get(cfg.Profile, Name).Scan(cfg)
	//监听配置变化
	w, err := configor.Watch(cfg.Profile, Name)
	if err != nil {
		return errors.Wrap(err, "configor.Watch")
	}
	go func() {
		for {
			v, err := w.Next()
			if err != nil {
				fmt.Println(errors.Wrap(err, "configor.WatchNext"))
			}
			v.Scan(cfg)
			fmt.Println(cfg)
		}
	}()
	// db.MysqlConn(db.Mysql{Address: cfg.Database.Address, Port: cfg.Database.Port, Username: cfg.Database.Username, Password: cfg.Database.Password, Database: cfg.Database.Database})
	return nil
}

func ConsulRegistry() registry.Registry {
	return consul.NewRegistry(registry.Addrs(cfg.Consul.Addr))
}
