package utils

import (
	"context"
	"fmt"
	"log"
	// "net/url"
	// "time"
	"os"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7/pkg/credentials"
)

var publicEndpoint = os.Getenv("MINIO_PUBLIC_ENDPOINT")

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader, bucketName, objectName string, minioClient *minio.Client) (string, error) {

	ctx := context.Background()

	// Check if the bucket exists
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Upload the file
	uploadInfo, err := minioClient.PutObject(ctx, bucketName, objectName, file, fileHeader.Size,
		minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")})
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", uploadInfo.Key, uploadInfo.Size)

	// Make the bucket public
	err = makeBucketPublic(minioClient, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to make bucket public: %w", err)
	}

	// Construct the public URL
	url := fmt.Sprintf("http://%s/%s/%s", publicEndpoint, bucketName, objectName)

	return url, nil
}

func makeBucketPublic(minioClient *minio.Client, bucketName string) error {

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


// func GetPresignedURL(minioClient *minio.Client, bucketName, objectName string) (string, error) {

// 	const expiresIn = time.Minute * 5

// 	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expiresIn, nil)
// 	if err != nil {
// 		return "", fmt.Errorf("error generating presigned URL: %v", err)
// 	}

// 	parsedURL, err := url.Parse(presignedURL.String())
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing presigned URL: %v", err)
// 	}

// 	parsedURL.Host = publicEndpoint
// 	// parsedURL.Scheme = "http"

// 	return parsedURL.String(), nil
// }