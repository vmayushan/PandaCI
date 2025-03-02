package storage

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TODO move env to env.go

type BucketClient struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func GetClient(ctx context.Context) (*BucketClient, error) {
	var accessKeyId = os.Getenv("S3_ACCESS_KEY_ID")
	var accessKeySecret = os.Getenv("S3_KEY_SECRET")

	config, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
	})

	presignClient := s3.NewPresignClient(s3Client)

	return &BucketClient{
		S3Client:      s3Client,
		PresignClient: presignClient,
	}, nil
}
