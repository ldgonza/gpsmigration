package output

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/storage"
)

var ctx = context.Background()

// WriteLocationsToCloudStorage writes to a bucket
func WriteLocationsToCloudStorage(bucketName string, fileName string, locations []TrackingLocation) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		// log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
		panic(err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	wc := bucket.Object(fileName).NewWriter(ctx)
	wc.ContentType = "application/json"

	encoder := json.NewEncoder(wc)
	if err := encoder.Encode(TrackingLocationCollection{locations}); err != nil {
		panic(err)
	}

	if err := wc.Close(); err != nil {
		panic(err)
		//		log.Errorf(ctx, "createFile: unable to close bucket %q, file %q: %v", bucket, fileName, err)
	}

	return nil
}

// WriteLatestTrackingStatusToCloudStorage writes latest status to cloud storage
func WriteLatestTrackingStatusToCloudStorage(bucketName string, fileName string, locations []LatestTrackingStatus) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	wc := bucket.Object(fileName).NewWriter(ctx)
	wc.ContentType = "application/json"

	encoder := json.NewEncoder(wc)
	if err := encoder.Encode(LatestTrackingStatusCollection{locations}); err != nil {
		panic(err)
	}

	if err := wc.Close(); err != nil {
		panic(err)
	}

	return nil
}
