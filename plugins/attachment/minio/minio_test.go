package minio

import (
	"fmt"
	"testing"

	"github.com/fangbc5/gogo/attachment/download"
	"github.com/fangbc5/gogo/attachment/upload"
)

func TestSave(t *testing.T) {
	u := new(Uploader)
	u.Init(
		upload.WithBucket("file"),
		upload.SetOption("endpoint", "43.143.136.7:9000"),
		upload.SetOption("accessKeyID", "fangbc"),
		upload.SetOption("secretAccessKey", "fangbc928"),
	)

	s := u.Save([]byte("发到你赛道"))
	fmt.Println(s)
}

func TestLoad(t *testing.T) {
	// 2023/05/06/1683339463493646800
	d := new(Downloader)
	d.Init(
		download.WithBucket("file"),
		download.SetOption("endpoint", "43.143.136.7:9000"),
		download.SetOption("accessKeyID", "fangbc"),
		download.SetOption("secretAccessKey", "fangbc928"),
	)

	s := d.Load("2023/05/06/1683342279352791400")
	fmt.Println(string(s))
}
