package minio

import (
	"fmt"
	"io"
	"log"

	"github.com/fangbc5/gogo/attachment/download"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Downloader struct {
	opts        download.Options
	minioClient *minio.Client
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
	if d.opts.Context == nil {
		log.Panicln("context不能为空，请调用Init方法时传入upload.SetOption并设置endpoint，accessKeyID，secretAccessKey")
	}
	if d.opts.Context.Value("endpoint") == nil {
		log.Panicln("请调用Init方法时传入upload.SetOption并设置endpoint")
	}
	if d.opts.Context.Value("accessKeyID") == nil {
		log.Panicln("请调用Init方法时传入upload.SetOption并设置accessKeyID")
	}
	if d.opts.Context.Value("secretAccessKey") == nil {
		log.Panicln("请调用Init方法时传入upload.SetOption并设置secretAccessKey")
	}
	endpoint := fmt.Sprint(d.opts.Context.Value("endpoint"))
	accessKeyID := fmt.Sprint(d.opts.Context.Value("accessKeyID"))
	secretAccessKey := fmt.Sprint(d.opts.Context.Value("secretAccessKey"))
	useSSL := false
	if d.opts.Context.Value("useSSL") != nil {
		v, ok := d.opts.Context.Value("useSSL").(bool)
		if ok {
			useSSL = v
		} else {
			log.Panicln("useSSL请使用bool值")
		}
	}
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	d.minioClient = minioClient

	return nil
}

func (d *Downloader) Options() download.Options {
	return d.opts
}

func (d *Downloader) Load(fileId string) []byte {
	if d.opts.Bucket == "" {
		log.Panicln("path或bucket不能为空")
	}
	if fileId == "" {
		log.Panicln("fileId不能为空")
	}

	obj, err := d.minioClient.GetObject(d.opts.Context, d.opts.Bucket, fileId, minio.GetObjectOptions{})
	if err != nil {
		log.Panicln("文件读取失败！！！", err)
	}
	content, err := io.ReadAll(obj)
	if err != nil {
		log.Panicln("文件读取失败！！！", err)
	}

	return content
}

func (d *Downloader) String() string {
	return "minio"
}
