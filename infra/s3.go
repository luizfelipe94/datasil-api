package infra

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
)

type S3 struct {
	client *minio.Client
}

func NewS3(client *minio.Client) *S3 {
	return &S3{
		client: client,
	}
}

func (s3 *S3) CreateBucket(name string) error {
	return s3.client.MakeBucket(context.TODO(), name, minio.MakeBucketOptions{
		Region: "us-east-1",
	})
}

func (s3 *S3) UploadFile(bucketName string, objectName string, filePath string, size int64, contentType string, tags map[string]string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = s3.client.PutObject(context.TODO(), bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
		UserTags:    tags,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3 *S3) DeleteFile(bucketName string, objectName string) error {
	return s3.client.RemoveObject(context.TODO(), bucketName, objectName, minio.RemoveObjectOptions{})
}
