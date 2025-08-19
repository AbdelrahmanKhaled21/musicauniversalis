package config

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinIO(config *Config) error {
	var err error
	MinioClient, err = minio.New(config.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIOAccessKey, config.MinIOSecretKey, ""),
		Secure: false, // Set to true for HTTPS
	})

	if err != nil {
		return err
	}

	log.Println("Connected to MinIO successfully")

	// Ensure the bucket exists
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, config.MinIOBucket)
	if err != nil {
		return err
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, config.MinIOBucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Created bucket: %s", config.MinIOBucket)
	} else {
		log.Printf("Bucket %s already exists", config.MinIOBucket)
	}

	return nil
}
