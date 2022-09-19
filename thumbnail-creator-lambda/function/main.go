package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
)

type s3Event struct {
	Bucket struct {
		Name string `json:"name"`
	} `json:"bucket"`
	Object events.S3Object `json:"object"`
}

type deps struct {
	downloader s3manageriface.DownloaderAPI
	uploader   s3manageriface.UploaderAPI
}

var decode = image.Decode
var resizeImg = resize.Resize
var encode = jpeg.Encode

func main() {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)

	deps := deps{
		downloader,
		uploader,
	}
	lambda.Start(deps.LambdaHandler)
}

func (d *deps) LambdaHandler(cwEvent events.CloudWatchEvent) error {
	var s3Event s3Event
	json.Unmarshal(cwEvent.Detail, &s3Event)

	srcBucket := s3Event.Bucket.Name
	dstBucket := srcBucket + "-resized"
	srcKey := s3Event.Object.Key
	dstKey := "resized-" + srcKey

	img, err := d.getObject(srcBucket, srcKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get image from bucket: %v\n", err)
		return err
	}

	resized, err := resizeImage(img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to resize image: %v\n", err)
		return err
	}

	err = d.uploadObject(resized, dstBucket, dstKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to upload image: %v\n", err)
		return err
	}

	fmt.Printf("Successfully resized %s/%s and uploaded to %s/%s \n", srcBucket, srcKey, dstBucket, dstKey)
	return nil
}

func (d *deps) getObject(bucket, objKey string) ([]byte, error) {
	buff := &aws.WriteAtBuffer{}
	_, err := d.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file, %w", err)
	}
	return buff.Bytes(), nil
}

func (d *deps) uploadObject(data []byte, bucket, fileKey string) error {
	reader := bytes.NewReader(data)

	_, err := d.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %w", err)
	}
	return nil
}

func resizeImage(data []byte) ([]byte, error) {
	img, _, err := decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode img, %w", err)
	}

	newImage := resizeImg(160, 0, img, resize.Lanczos3)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err = encode(w, newImage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode img, %w", err)
	}

	return b.Bytes(), nil
}
