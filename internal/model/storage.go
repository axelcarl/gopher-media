package model

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func UploadFile(object string, file io.Reader) error {
	bucket := os.Getenv("BUCKET_NAME")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)

	if _, err = io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	return nil
}
