package s3

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func Uploader() *manager.Uploader {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return uploader
}

func UploadToBucket(uploader *manager.Uploader, bucket string, uploadFile string) string {
	key := filepath.Base(uploadFile)
	file, err := os.Open(uploadFile)
	//fmt.Println("upload attempt: " + uploadFile)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer file.Close()
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return result.Location
}
