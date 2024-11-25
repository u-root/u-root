package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	lifecycleConfig = `{
		"rule": [
			{
				"action": {
					"type": "Delete"
				},
				"condition": {
					"age": 1
				}
			}
		]
	}`
)

func main() {
	key := flag.String("key", "key.json", "google service account key")
	proj := flag.String("project", "hpc-benchmarking-sandbox", "project")
	bucketName := flag.String("bucket", "", "bucket name")
	flag.Parse()

	if len(*bucketName) == 0 {
		*bucketName = "snapcheck" + strings.ReplaceAll(strings.ToLower(time.Now().Format(time.RFC3339)), ":", "-")
	}
	ctx := context.Background()

	// Initialize Google Cloud Storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(*key))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create or get the bucket
	fmt.Printf("gs://%s/%s", *proj, *bucketName)
	bucket := client.Bucket(*proj)

	if err := bucket.Create(ctx, *proj, nil); err != nil { // && err != storage.ErrBucketExists {
		gerr, ok := err.(*googleapi.Error)
		if ok && gerr.Code != 409 {
			log.Fatalf("create bucket: %v", err)
		}
		if !ok {
			log.Fatalf("create bucket: unknown error type %T:%v", err, err)
		}
	}

	// Upload a JSON file to the bucket
	if err := uploadJSON(ctx, bucket, os.Stdin, *bucketName); err != nil {
		log.Fatalf("uploading to %q:%v", *bucketName, err)
	}

	// Apply lifecycle policy
	if err := setBucketLifecyclePolicy(ctx, bucket); err != nil {
		log.Fatalf("setting life cycle:%v", err)
	}
}

func uploadJSON(ctx context.Context, bucket *storage.BucketHandle, r io.Reader, destFileName string) error {
	// Create a writer to the destination file in the bucket
	wc := bucket.Object(destFileName).NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Failed to close writer: %v", err)
	}

	return nil
}

func setBucketLifecyclePolicy(ctx context.Context, bucket *storage.BucketHandle) error {
	// Define the lifecycle policy
	lc := &storage.Lifecycle{
		Rules: []storage.LifecycleRule{
			{
				Action:    storage.LifecycleAction{Type: "Delete"},
				Condition: storage.LifecycleCondition{AgeInDays: 1},
			},
		},
	}

	// Update the bucket with the lifecycle configuration
	attrsToUpdate := storage.BucketAttrsToUpdate{
		Lifecycle: lc,
	}

	if _, err := bucket.Update(ctx, attrsToUpdate); err != nil {
		return fmt.Errorf("failed to update bucket: %v", err)
	}

	return nil
}
