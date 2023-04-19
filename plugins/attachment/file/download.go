package file

import "github.com/fangbc5/gogo/attachment/download"

type Downloader struct {
	opts download.Options
}

func (d Downloader) Init(options ...download.Option) error {
	//TODO implement me
	panic("implement me")
}

func (d Downloader) Options() download.Options {
	return d.opts
}

func (d Downloader) Load() []byte {
	//TODO implement me
	panic("implement me")
}

func (d Downloader) String() string {
	return "file"
}
