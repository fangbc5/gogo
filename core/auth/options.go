package auth

import (
	"context"
	"github.com/fangbc5/gogo/core/logger"
	"time"
)

type Options struct {
	// Namespace the service belongs to
	Namespace string
	// ID is the services auth ID
	ID string
	// Secret is used to authenticate the service
	Secret string
	// Token is the services token used to authenticate itself
	Token *Token
	// PublicKey for decoding JWTs
	PublicKey string
	// PrivateKey for encoding JWTs
	PrivateKey string
	// Addrs sets the addresses of auth
	Addrs []string
	// Logger is the underline logger
	Logger logger.Logger
	// Store to persist the tokens
	//Store store.Store
}

type Option func(o *Options)

// WithAddrs is the auth addresses to use.
func WithAddrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// WithNamespace the service belongs to.
func WithNamespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}

// WithStore sets the token providers store.
//func WithStore(s store.Store) Option {
//	return func(o *Options) {
//o.Store = s
//}
//}

// WithPublicKey sets the JWT public key.
func WithPublicKey(key string) Option {
	return func(o *Options) {
		o.PublicKey = key
	}
}

// WithPrivateKey sets the JWT private key.
func WithPrivateKey(key string) Option {
	return func(o *Options) {
		o.PrivateKey = key
	}
}

// WithLogger sets the underline logger.
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithCredentials sets the auth credentials.
func WithCredentials(id, secret string) Option {
	return func(o *Options) {
		o.ID = id
		o.Secret = secret
	}
}

// WithToken sets the auth token to use when making requests.
func WithToken(token *Token) Option {
	return func(o *Options) {
		o.Token = token
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Logger: logger.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type GenerateOptions struct {
	// Metadata associated with the account
	Metadata map[string]string
	// Scopes the account has access too
	Scopes []string
	// Provider of the account, e.g. oauth
	Provider string
	// Type of the account, e.g. user
	Type string
	// Secret used to authenticate the account
	Secret string
	// RefreshToken is used to refresh a token
	RefreshToken string
	// Expiry for the token
	Expiry time.Duration
}

type GenerateOption func(o *GenerateOptions)

// WithSecret for the generated account.
func WithSecret(s string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Secret = s
	}
}

// WithType for the generated account.
func WithType(t string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Type = t
	}
}

// WithMetadata for the generated account.
func WithMetadata(md map[string]string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Metadata = md
	}
}

// WithProvider for the generated account.
func WithProvider(p string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Provider = p
	}
}

// WithScopes for the generated account.
func WithScopes(s ...string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Scopes = s
	}
}

func WithRefreshToken(rt string) GenerateOption {
	return func(o *GenerateOptions) {
		o.RefreshToken = rt
	}
}

// WithExpiry for the generated account's token expires.
func WithExpiry(d time.Duration) GenerateOption {
	return func(o *GenerateOptions) {
		o.Expiry = d
	}
}

// NewGenerateOptions from a slice of options.
func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	var options GenerateOptions
	for _, o := range opts {
		o(&options)
	}

	// set default expiry of token
	if options.Expiry == 0 {
		options.Expiry = time.Minute * 15
	}

	return options
}

type VerifyOptions struct {
	Context context.Context
}

type VerifyOption func(o *VerifyOptions)

func VerifyContext(ctx context.Context) VerifyOption {
	return func(o *VerifyOptions) {
		o.Context = ctx
	}
}

type ListOptions struct {
	Context context.Context
}

type ListOption func(o *ListOptions)

func RulesContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}
