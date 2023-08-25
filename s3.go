package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"sync"
)

type ProgressReportingReader struct {
	fp      *os.File
	size    int64
	read    int64
	signMap map[int64]struct{}
	mux     sync.Mutex
}

func (r *ProgressReportingReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *ProgressReportingReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	r.mux.Lock()
	// if called first time then skip
	if _, ok := r.signMap[off]; ok {
		r.read += int64(n)
		log.Printf("\rtotal read: %d		progress: %d%%", r.read, int(float32(r.read*100)/float32(r.size)))
	} else {
		r.signMap[off] = struct{}{}
	}
	r.mux.Unlock()

	return n, err
}

func (r *ProgressReportingReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

func UploadFile(key, bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	file, err := os.Open(key)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	reader := &ProgressReportingReader{
		fp:      file,
		size:    fileInfo.Size(),
		signMap: map[int64]struct{}{},
	}

	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err == nil {
		log.Printf("Object \"%s\" uploaded successfully.", key)
	}

	return err
}

func DownloadFile(key, bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(client)

	file, err := os.Create(key)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = downloader.Download(context.Background(), file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		log.Printf("Object \"%s\" downloaded successfully.", key)
	}

	return err
}

func ObjectExists(key, bucket string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return false, err
	}

	client := s3.NewFromConfig(cfg)

	if _, err = client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return false, nil
	}

	return true, nil
}
