package upload

import "context"

type Options struct {
	File   []byte
	Path    string
	Source string
	Bucket string
	// Alternative options
	Context context.Context
}

type Option func(*Options)

func WithFile(file []byte) Option {
	return func (o *Options) {
		o.File = file
	}
}

func WithPath(path string) Option {
	return func (o *Options) {
		o.Path = path
	}
}

func WithSource(source string) Option {
	return func (o *Options) {
		o.Source = source
	}
}

func WithBucket(bucket string) Option {
	return func (o *Options) {
		o.Bucket = bucket
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}