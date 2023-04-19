package consul

type Options struct {
	Env      string
	Profile  string
	Server   ServerConfig
	Auth     AuthConfig
	Consul   ConsulConfig
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