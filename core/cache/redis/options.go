package redis

// Options Redis配置结构体
type Options struct {
	Address  string
	Port     string
	Password string
	Database int
}

type Option func(options *Options)

func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func WithPort(port string) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithDatabase(database int) Option {
	return func(o *Options) {
		o.Database = database
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Address:  "127.0.0.1",
		Port:     "6379",
		Database: 0,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
