package consul

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
	Version      = "v1.0.0"
	ConsulAddr   = "127.0.0.1:8500"
	ConsulPrefix = "/micro/config"
	ENV          = "env"
)

var cfg *Options

func WithName(name string) Option {
	return func(c *Options) {
		c.Server.Name = name
	}
}
func WithConsul(addr string, prefix string) Option {
	return func(c *Options) {
		c.Consul.Addr = addr
		c.Consul.Prefix = prefix
	}
}

func WithEnv(env string) Option {
	return func(c *Options) {
		c.Env = env
	}
}

func Get() *Options {
	return cfg
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

func Init(opts ...Option) error {
	cfg := &Options{
		Env:     ENV,
		Profile: Profile,
		Server: ServerConfig{
			Name:    Name,
			Port:    Port,
			Version: Version,
		},
		Consul: ConsulConfig{
			Addr:   ConsulAddr,
			Prefix: ConsulPrefix,
		},
	}
	//设置参数
	for _, opt := range opts {
		opt(cfg)
	}
	//加载profile
	configor, err := config.NewConfig(config.WithSource(source.NewSource(
		source.WithAddress(cfg.Consul.Addr),
		source.WithPrefix(cfg.Consul.Prefix),
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
	configor.Get(cfg.Profile, cfg.Server.Name).Scan(cfg)
	//监听配置变化
	w, err := configor.Watch(cfg.Profile, cfg.Server.Name)
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
