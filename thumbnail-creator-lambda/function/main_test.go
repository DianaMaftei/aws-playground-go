package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"image"
	"image/jpeg"
	"io"
	"os"
	"testing"
)

type DownloaderMock struct {
	mock.Mock
}

type UploaderMock struct {
	mock.Mock
}

func decodeMock(io.Reader) (image.Image, string, error) {
	return nil, "", nil
}
func resizeMock(uint, uint, image.Image, resize.InterpolationFunction) image.Image {
	return nil
}

func encodeMock(io.Writer, image.Image, *jpeg.Options) error {
	return nil
}

func (mock *UploaderMock) Upload(input *s3manager.UploadInput, f ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := mock.Called(input, f)
	return args.Get(0).(*s3manager.UploadOutput), args.Error(1)
}

func (mock *DownloaderMock) Download(writer io.WriterAt, input *s3.GetObjectInput, f ...func(*s3manager.Downloader)) (int64, error) {
	args := mock.Called(writer, input, f)
	return args.Get(0).(int64), args.Error(1)
}

func TestLambdaHandler(t *testing.T) {
	// mock AWS S3 service
	mD := DownloaderMock{}
	mU := UploaderMock{}
	d := deps{downloader: &mD, uploader: &mU}

	mD.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil)
	mU.On("Upload", mock.Anything, mock.Anything).Return(&s3manager.UploadOutput{}, nil)

	// get sample cloud watch event for S3 PutObject
	testEventFileName := "testEvent.json"
	var testEvent events.CloudWatchEvent

	file, err := readTestFile(testEventFileName)
	require.NoError(t, err)
	err = json.Unmarshal(file, &testEvent)
	require.NoError(t, err)

	// mock image resizing functions
	decode = decodeMock
	resizeImg = resizeMock
	encode = encodeMock

	// call Lambda function handler
	d.LambdaHandler(testEvent)

	// assert that the image is downloaded and then uploaded
	mD.AssertCalled(t, "Download", mock.Anything, mock.Anything, mock.Anything)
	mU.AssertCalled(t, "Upload", mock.Anything, mock.Anything)
}

func readTestFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file %s\": %v\n", fileName, err)
		return nil, err
	}

	defer f.Close()

	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read file %s\": %v\n", fileName, err)
		return nil, err
	}
	return data, nil
}

func (mock *DownloaderMock) DownloadWithContext(aws.Context, io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error) {
	panic("dummy")
}

func (mock *UploaderMock) UploadWithContext(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	panic("dummy")
}
