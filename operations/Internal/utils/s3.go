package utils

import (
	"context"
	"fmt"
	"os"
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var s3Client *s3.Client
var s3Bucket string

// InitS3 initializes the S3 client and bucket name from environment variables
func InitS3() error {
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return fmt.Errorf("AWS_S3_BUCKET environment variable not set")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	s3Bucket = bucket
	return nil
}

// UploadChunk uploads a file chunk to S3
func UploadChunk(ctx context.Context, key string, data []byte, contentType string) error {
	if s3Client == nil {
		return fmt.Errorf("S3 client not initialized")
	}
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		// ACL:         types.ObjectCannedACLPublicRead, // or private
	})
	return err
} 