package nacos

import (
	"fmt"

	source "github.com/go-micro/plugins/v4/config/source/nacos"
	"github.com/go-micro/plugins/v4/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
)

const (
	Name        = "myservice"
	Port        = "38080"
	Version     = "v1.0.0"
	NacosAddr   = "127.0.0.1:8848"
	NacosDataId = "myservice-dev.yml"
	NacosGroup  = "gogo"
	Profile     = "DEV"
)

var cfg *Options

func WithName(name string) Option {
	return func(c *Options) {
		c.Server.Name = name
	}
}
func WithNacos(nacosConfig NacosConfig) Option {
	return func(c *Options) {
		c.Nacos = nacosConfig
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
		Server: ServerConfig{
			Name:    Name,
			Port:    Port,
			Version: Version,
		},
		Nacos: NacosConfig{
			Addr:   []string{NacosAddr},
			DataId: NacosDataId,
			Group:  NacosGroup,
		},
	}
	//设置参数
	for _, opt := range opts {
		opt(cfg)
	}
	//加载profile
	configor, err := config.NewConfig(config.WithSource(source.NewSource(
		source.WithAddress(cfg.Nacos.Addr),
		source.WithClientConfig(constant.ClientConfig{
			Username:    "nacos",
			Password:    "nacos",
		}),
		source.WithDataId(cfg.Nacos.DataId),
		source.WithGroup(cfg.Nacos.Group),
	)))
	if err != nil {
		return errors.Wrap(err, "configor.New")
	}
	if err := configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load")
	}
	return nil
}

func NacosRegistry() registry.Registry {
	return nacos.NewRegistry(registry.Addrs(cfg.Nacos.Addr...))
}
