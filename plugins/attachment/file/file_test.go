package file

import (
	"fmt"
	"testing"

	"github.com/fangbc5/gogo/attachment/download"
	"github.com/fangbc5/gogo/attachment/upload"
)

func TestSave(t *testing.T) {
	u := new(Uploader)
	u.Init(upload.WithPath("D:/programs"), upload.WithBucket("file"))
	fileId := u.Save([]byte("abcd"))
	fmt.Println(fileId)
}

func TestLoad(t *testing.T) {
	d := new(Downloader)
	d.Init(download.WithPath("D:/programs"), download.WithBucket("file"))
	b := d.Load("2023\\05\\05\\1683276899484262500")
	fmt.Println(string(b))
}

func TestDirSplit(t *testing.T) {
	u := new(Uploader)
	u.Init(upload.WithPath("D:/programs"), upload.WithBucket("file"))
	u.DirSplit()
}