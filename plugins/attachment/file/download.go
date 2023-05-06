package file

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fangbc5/gogo/attachment/download"
)

type Downloader struct {
	opts download.Options
}

func (d *Downloader) Init(options ...download.Option) error {
	d.opts = download.Options{
		Source: "file",
		Path:   "/gogo/upload",
		Bucket: "",
	}
	for _, opt := range options {
		opt(&d.opts)
	}
	return nil
}

func (d *Downloader) Options() download.Options {
	return d.opts
}

func (d *Downloader) Load(fileId string) []byte {
	if d.opts.Path == "" || d.opts.Bucket == "" {
		log.Panicln("path或bucket不能为空")
	}
	if fileId == "" {
		log.Panicln("fileId不能为空")
	}
	filePath := filepath.Join(d.opts.Path, d.opts.Bucket, fileId)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

func (d *Downloader) String() string {
	return "file"
}
