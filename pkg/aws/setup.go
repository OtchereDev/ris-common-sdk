package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Client     *s3.Client
	BucketName string
}

func Connect(bucket string, cfg aws.Config, optFns ...func(*s3.Options)) (s *S3, err error) {

	s3Client := s3.NewFromConfig(cfg, optFns...)

	s = &S3{
		Client:     s3Client,
		BucketName: bucket,
	}

	return
}

func (s *S3) UploadFileToS3(ctx context.Context, fileName, base64File, contentType string) (string, error) {

	b64data := base64File[strings.IndexByte(base64File, ',')+1:]

	decode, err := base64.StdEncoding.DecodeString(b64data)

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}
	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(decode),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to S3: %w", err)
	}

	return fileName, nil
}

func (s *S3) UploadApptFileToS3(ctx context.Context, filename string, file io.Reader, contentType string) (string, error) {
	uploader := manager.NewUploader(s.Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}
	return filename, nil
}

func (s *S3) DeleteFileFromS3(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error deleting file from S3: %w", err)
	}
	return nil
}

func (s *S3) GetFileLink(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(s.Client)
	presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(7*24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("error generating presigned URL: %w", err)
	}
	return presignedURL.URL, nil
}
