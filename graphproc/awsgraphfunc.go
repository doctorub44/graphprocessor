package graphproc

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//AWSDownloadS3Bucket : graph stage that downloads a file from an S3 bucket and returns it as raw bytes
func AWSDownloadS3Bucket(s *State, payload *Payload) error {
	region, ok := s.Config("region")
	if !ok {
		return errors.New("unable to lookup configuration field: " + region)
	}
	bucket, ok := s.Config("bucket")
	if !ok {
		return errors.New("unable to lookup configuration field: " + bucket)
	}
	file, ok := s.Config("file")
	if !ok {
		return errors.New("unable to lookup configuration field: " + file)
	}
	data, err := downloadS3File(region, bucket, file)
	if err == nil {
		payload.Raw = append(payload.Raw[:0], data...)
	}
	return err
}

//downloadS3File : function to download the S3 file to the AWS Lambda /tmp directory and read and returns the file bytes
func downloadS3File(region, bucket, filename string) ([]byte, error) {
	file, err := os.Create("/tmp/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}
	return data, nil
}
