package file

import "github.com/fangbc5/gogo/attachment/upload"

type Uploader struct {
	opts upload.Options
}

func (u Uploader) Init(options ...upload.Option) error {
	u.opts = upload.Options{
		Source: "file",
		Path: "/gogo/upload",
		Bucket: "",
	}
	for _, opt := range options {
		opt(&u.opts)
	}
	return nil
}

func (u Uploader) Options() upload.Options {
	return u.opts
}

func (u Uploader) Save(bytes []byte) string {
	return "fileId"
}

func (u Uploader) String() string {
	return "file"
}
