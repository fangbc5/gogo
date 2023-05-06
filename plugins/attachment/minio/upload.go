package minio

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fangbc5/gogo/attachment/upload"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Uploader struct {
	opts        upload.Options
	minioClient *minio.Client
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
	if u.opts.Context == nil {
		log.Panicln("context不能为空，请调用Init方法时传入upload.SetOption并设置endpoint，accessKeyID，secretAccessKey")
	}
	if u.opts.Context.Value("endpoint") == nil {
		log.Panicln("请调用Init方法时传入upload.SetOption并设置endpoint")
	}
	if u.opts.Context.Value("accessKeyID") == nil {
		log.Panicln("请调用Init方法时传入upload.SetOption并设置accessKeyID")
	}
	if u.opts.Context.Value("secretAccessKey") == nil {
		log.Panicln("请调用Init方法时传入upload.SetOption并设置secretAccessKey")
	}
	endpoint := fmt.Sprint(u.opts.Context.Value("endpoint"))
	accessKeyID := fmt.Sprint(u.opts.Context.Value("accessKeyID"))
	secretAccessKey := fmt.Sprint(u.opts.Context.Value("secretAccessKey"))
	useSSL := false
	if u.opts.Context.Value("useSSL") != nil {
		v, ok := u.opts.Context.Value("useSSL").(bool)
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
	u.minioClient = minioClient

	return nil
}

func (u *Uploader) Options() upload.Options {
	return u.opts
}

func (u *Uploader) Save(obj []byte) string {
	if u.opts.Bucket == "" {
		log.Panicln("path或bucket不能为空")
	}
	// 生成文件名
	millisecond := time.Now().UnixNano()
	filename := strconv.Itoa(int(millisecond))

	// 打开文件
	dir := u.DirSplit()
	objectPath := filepath.Join(dir, filename)
	objectPath = strings.ReplaceAll(objectPath, "\\", "/")
	if ok, _ := u.minioClient.BucketExists(u.opts.Context, u.opts.Bucket); !ok {
		err := u.minioClient.MakeBucket(u.opts.Context, u.opts.Bucket, minio.MakeBucketOptions{Region: "us-east-8"})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := u.minioClient.BucketExists(u.opts.Context, u.opts.Bucket)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", u.opts.Bucket)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Successfully created %s\n", u.opts.Bucket)
		}
	}
	contentType := "application/octet-stream"
	contentLength := int64(len(obj))
	reader := bytes.NewReader(obj)
	info, err := u.minioClient.PutObject(u.opts.Context, u.opts.Bucket, objectPath, reader, contentLength, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", filename, info.Size)
	//返回路径
	return objectPath
}

func (u *Uploader) String() string {
	return "minio"
}

func (u *Uploader) DirSplit() string {
	now := time.Now()
	year, month, day := now.Date()

	// 以当前年份作为根目录
	yearDir := filepath.Join(strconv.Itoa(year))
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
