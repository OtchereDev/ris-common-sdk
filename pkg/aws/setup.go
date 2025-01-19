package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Client     *s3.S3
	BucketName string
}

func Connect(ak, r, as, b string) (s *S3, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(r),
		Credentials: credentials.NewStaticCredentials(ak, as, ""),
	})

	if err != nil {
		return
	}

	svc := s3.New(sess)

	s = &S3{
		Client:     svc,
		BucketName: b,
	}

	return
}

func (s S3) UploadFileToS3(fileType string, base64File string, patientId string) (string, error) {
	date := time.Now().Format(time.RFC3339)
	key := fileType + "/" + patientId + "/" + date + ".pdf"

	b64data := base64File[strings.IndexByte(base64File, ',')+1:]

	decode, err := base64.StdEncoding.DecodeString(b64data)

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}
	_, err = s.Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(decode),
		ContentType: aws.String("application/pdf"),
	})

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}
	return key, nil
}

func (s S3) UploadApptFileToS3(fileType string, file multipart.File, patientId string, extn string, contentType string) (string, error) {
	date := time.Now().Format(time.RFC3339)
	key := fileType + "/" + patientId + "/" + date + extn

	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}
	return key, nil
}

func (s S3) DeleteFileFromS3(key string) error {

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}

	_, err := s.Client.DeleteObject(input)

	return err
}

func (s S3) GetFileLink(key string) (string, error) {

	params := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}

	req, _ := s.Client.GetObjectRequest(params)

	if req.Error != nil {
		fmt.Println(req.Error)
		return "", req.Error
	}

	url, err := req.Presign(30 * time.Minute)
	if err != nil {
		return "", err
	}

	return url, err
}
