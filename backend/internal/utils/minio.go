package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

var publicEndpoint = os.Getenv("MINIO_PUBLIC_ENDPOINT")

func IsBucketExists(ctx context.Context, minioClient *minio.Client, bucketName string) error {

	// Check if the bucket exists
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}

func MakeBucketPublic(minioClient *minio.Client, bucketName string) error {

    ctx := context.Background()

    // Check if the bucket exists
    exists, err := minioClient.BucketExists(ctx, bucketName)
    if err != nil {
        return fmt.Errorf("failed to check bucket: %w", err)
    }

    if !exists {
        err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return fmt.Errorf("failed to create bucket: %w", err)
        }
    }

    // Create a policy to make the bucket public
    policy := fmt.Sprintf(`
	{
	"Version":"2012-10-17",
	"Statement":[
		{
		"Effect":"Allow",
		"Principal": {"AWS":["*"]},
		"Action":["s3:GetObject"],
		"Resource":["arn:aws:s3:::%s/*"]
		}
	]
	}`, 
	bucketName)

    err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
    if err != nil {
        return fmt.Errorf("failed to set bucket policy: %w", err)
    }

    return nil
}

func GetBucketAndSubBuckets(role string) ([]string, []string) {
	if role == "restaurant" {
		return []string{"restaurant-pictures"}, []string{"restaurants", "menu-items"}
	}

	if role == "customer" {
		return []string{"customer-pictures"}, []string{"qr-codes"}
	}

	return []string{}, []string{}
}

func CreateBucketAndSubBuckets(minioClient *minio.Client, bucketName string, subBuckets []string) error {
	ctx := context.Background()

	// Create main bucket if not exists
	err := IsBucketExists(ctx, minioClient, bucketName)
	if err != nil {
		return fmt.Errorf("error creating bucket %s: %v", bucketName, err)
	}
	log.Printf("Bucket %s is ready", bucketName)

	// Make the bucket public
	err = MakeBucketPublic(minioClient, bucketName)
	if err != nil {
		return fmt.Errorf("error making bucket %s public: %v", bucketName, err)
	}
	log.Printf("Bucket %s is public", bucketName)

	// Create sub-buckets (folders)
	for _, subFolder := range subBuckets {
		objectName := fmt.Sprintf("%s/", subFolder)
		_, err := minioClient.PutObject(ctx, bucketName, objectName, nil, 0, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("error creating folder %s in bucket %s: %v", subFolder, bucketName, err)
		}
		log.Printf("Folder %s/ created inside bucket %s", subFolder, bucketName)
	}
	return nil
}