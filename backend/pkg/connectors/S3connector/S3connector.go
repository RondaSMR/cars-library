package S3connector

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client *s3.Client
	bucket   string
	endpoint string
}

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewS3Client(cfg S3Config) (Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return Client{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Указываем кастомный Endpoint при создании клиента S3
	s3Client := s3.New(s3.Options{
		Credentials:  awsCfg.Credentials,
		Region:       cfg.Region,
		BaseEndpoint: aws.String(cfg.Endpoint), // Новый метод для задания кастомного Endpoint
	})

	return Client{s3Client: s3Client, bucket: cfg.Bucket, endpoint: cfg.Endpoint}, nil
}

func (c Client) UploadFile(ctx context.Context, fileName string, fileData []byte) (string, error) {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", c.endpoint, c.bucket, fileName)
	return fileURL, nil
}

func (c Client) DeleteFile(ctx context.Context, fileName string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}
