package storage

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/rs/zerolog/log"
)

func (basics BucketClient) KeyExists(ctx context.Context, bucketName string, objectKey string) (bool, error) {
	_, err := basics.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				return false, nil
			}
		}

		log.Error().Err(err).Msgf("Couldn't check if key exists %v:%v", bucketName, objectKey)
		return false, err
	}

	return true, nil
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (basics BucketClient) UploadFile(ctx context.Context, bucketName string, objectKey string, file []byte, contentType string) error {

	_, err := basics.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(file),
		ContentType: &contentType,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't upload file to %v:%v",
			bucketName, objectKey)
	}

	return err
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (basics BucketClient) DownloadFile(ctx context.Context, bucketName string, objectKey string) ([]byte, error) {
	result, err := basics.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't get object %v:%v", bucketName, objectKey)
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't read object body from %v", objectKey)
	}
	return body, err
}
