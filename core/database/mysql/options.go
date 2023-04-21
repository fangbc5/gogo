package mysql

import "time"

// Options Mysql配置结构体
type Options struct {
	Address  string
	Port     string
	Username string
	Password string
	Database string

	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
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

func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithDatabase(database string) Option {
	return func(o *Options) {
		o.Database = database
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(o *Options) {
		o.MaxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(o *Options) {
		o.MaxOpenConns = maxOpenConns
	}
}

func WithConnMaxIdleTime(connMaxIdleTime time.Duration) Option {
	return func(o *Options) {
		o.ConnMaxIdleTime = connMaxIdleTime
	}
}

func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(o *Options) {
		o.ConnMaxLifetime = connMaxLifetime
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Address:  "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "123456",
		Database: "gogo",

		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxIdleTime: time.Second * 10,
		ConnMaxLifetime: time.Hour,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
