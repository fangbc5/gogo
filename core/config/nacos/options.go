package nacos

type Options struct {
	Server   ServerConfig
	Auth     AuthConfig
	Nacos   NacosConfig
	Tracing  TracingConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

type Option func(*Options)

type ServerConfig struct {
	Name    string
	Port    string
	Version string
}

type AuthConfig struct {
	JwtSecret string
	TokenLife string
}

type NacosConfig struct {
	Addr   []string
	DataId string
	Group  string
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

func WithAddress(addr []string) Option {
	return func(o *Options) {
		o.Nacos.Addr = addr
	}
}

func WithDataId(dataId string) Option {
	return func(o *Options) {
		o.Nacos.DataId = dataId
	}
}

func WithGroup(group string) Option {
	return func(o *Options) {
		o.Nacos.Group = group
	}
}