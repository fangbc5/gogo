package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fangbc5/gogo/attachment/upload"
)

type Uploader struct {
	opts upload.Options
}

func (u *Uploader) Init(options ...upload.Option) error {
	u.opts = upload.Options{
		Source: "file",
		Path:   "/gogo/upload",
		Bucket: "",
	}
	for _, opt := range options {
		opt(&u.opts)
	}
	return nil
}

func (u *Uploader) Options() upload.Options {
	return u.opts
}

func (u *Uploader) Save(bytes []byte) string {
	if u.opts.Path == "" || u.opts.Bucket == "" {
		log.Panicln("path或bucket不能为空")
	}
	// 生成文件名
	millisecond := time.Now().UnixNano()
	filename := strconv.Itoa(int(millisecond))

	// 打开文件
	dir := u.DirSplit()
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 字节数组写入到文件
	file.Write(bytes)
	//返回路径
	return strings.ReplaceAll(filePath, filepath.Join(u.opts.Path, u.opts.Bucket), "")
}

func (u *Uploader) String() string {
	return "file"
}

func (u *Uploader) DirSplit() string {
	now := time.Now()
	year, month, day := now.Date()

	// 以当前年份作为根目录
	yearDir := filepath.Join(u.opts.Path, u.opts.Bucket, strconv.Itoa(year))
	if _, err := os.Stat(yearDir); os.IsNotExist(err) {
		os.Mkdir(yearDir, os.ModePerm)
	}

	// 以当前月份作为子目录
	monthDir := filepath.Join(yearDir, fmt.Sprintf("%02d", int(month)))
	if _, err := os.Stat(monthDir); os.IsNotExist(err) {
		os.Mkdir(monthDir, os.ModePerm)
	}

	// 以当前日期作为子目录
	dayDir := filepath.Join(monthDir, fmt.Sprintf("%02d", day))
	if _, err := os.Stat(dayDir); os.IsNotExist(err) {
		os.Mkdir(dayDir, os.ModePerm)
	}
	return dayDir
}
