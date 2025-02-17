package services

import (
	"image"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucketName = os.Getenv("AWS_BUCKET")

func GetClient() *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Endpoint:    aws.String(os.Getenv("AWS_ENDPOINT")),
	}))
	svc := s3.New(sess)

	return svc
}

func GetFile(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	svc := GetClient()
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}

func PutFile(file *multipart.FileHeader, mediaType string, userId uint) (string, error) {
	extension := filepath.Ext(file.Filename)
	now := time.Now()
	fileName := mediaType + "/" + strconv.Itoa(int(now.UnixNano())) + strconv.Itoa(int(userId)) + "-" + RandStringBytes(20, false) + extension
	svc := GetClient()

	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   fileContent,
	}

	_, err = svc.PutObject(params)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func PutImage(file *multipart.FileHeader, mediaType string, userId uint) (string, error, float32) {
	extension := filepath.Ext(file.Filename)
	now := time.Now()
	fileName := mediaType + "/" + strconv.Itoa(int(now.UnixNano())) + strconv.Itoa(int(userId)) + "-" + RandStringBytes(20, false) + extension
	svc := GetClient()
	var size float32
	size = 1.0

	fileContent, err := file.Open()
	if err != nil {
		return "", err, size
	}
	defer fileContent.Close()

	img, _, err := image.DecodeConfig(fileContent)

	if err != nil {
		size = float32(img.Width) / float32(img.Height)
		if math.IsNaN(float64(size)) {
			size = 1.0
		}
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   fileContent,
	}

	_, err = svc.PutObject(params)
	if err != nil {
		return "", err, size
	}

	return fileName, nil, size
}
