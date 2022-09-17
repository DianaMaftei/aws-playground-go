package service

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/globalsign/mgo/bson"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func ConnectToAws() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create aws session: %w", err)
	}
	return sess, nil
}

func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader, dir string) (string, error) {
	size := fileHeader.Size
	buffer := make([]byte, size)
	_, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("unable to read file %w", err)
	}
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	tempFileName := dir + "/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	bucket := os.Getenv("IMAGE_BUCKET")
	_, err = s3manager.NewUploader(s).Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(tempFileName),
		Body:        fileBytes,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(fileType),
	})

	return tempFileName, fmt.Errorf("unable to upload file %w", err)
}
