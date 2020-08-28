package output

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// GetS3Uploader gets the uploader
func GetS3Uploader() (*s3manager.Uploader, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		fmt.Printf("Error, %v", err)
		return nil, fmt.Errorf("Error, %v", err)
	}

	uploader := s3manager.NewUploader(sess)

	return uploader, nil
}

// WriteToS3 uploads the file to s3
func WriteToS3(uploader *s3manager.Uploader, bucketName string, fileName string, jsonBytes []byte) error {
	// Upload the file to S3.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(jsonBytes),
	})

	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

// WriteLocationsToS3 writes to a bucket
func WriteLocationsToS3(uploader *s3manager.Uploader, bucketName string, fileName string, locations []TrackingLocation) error {
	jsonString, err := json.Marshal(TrackingLocationCollection{locations})
	if err != nil {
		panic(err)
	}

	WriteToS3(uploader, bucketName, fileName, jsonString)
	return nil
}

// WriteLatestTrackingStatusToS3 writes latest status to cloud storage
func WriteLatestTrackingStatusToS3(uploader *s3manager.Uploader, bucketName string, fileName string, locations []LatestTrackingStatus) error {
	jsonString, err := json.Marshal(LatestTrackingStatusCollection{locations})
	if err != nil {
		panic(err)
	}

	err = WriteToS3(uploader, bucketName, fileName, jsonString)
	if err != nil {
		panic(err)
	}

	return nil
}
