package download

import "context"

type Options struct {
	FileId   string
	Path  string
	Source string
	Bucket string
	// Alternative options
	Context context.Context
}

type Option func(*Options)

func WithSource(source string) Option {
	return func(args *Options) {
		args.Source = source
	}
}

func WithPath(path string) Option {
	return func(args *Options) {
		args.Path = path
	}
}

func WithBucket(bucket string) Option {
	return func(args *Options) {
		args.Bucket = bucket
	}
}

func WithFileId(fileId string) Option {
	return func(args *Options) {
		args.FileId = fileId
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