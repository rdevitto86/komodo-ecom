package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	logger "komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	client  *s3.Client
	once    sync.Once
	mu      sync.RWMutex
	initErr error
)

type Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

// Initialize the S3 client
func Init(config Config) error {
	once.Do(func() {
		if config.Region == "" {
			logger.Error("s3 region is required", fmt.Errorf("s3 region is required"))
			initErr = fmt.Errorf("s3 region is required")
			return
		}

		ctx := context.Background()
		var cfg aws.Config

		cfgOpts := []func(*awsconfig.LoadOptions) error{
			awsconfig.WithRegion(config.Region),
		}

		if config.AccessKey != "" && config.SecretKey != "" {
			cfgOpts = append(cfgOpts, awsconfig.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					config.AccessKey,
					config.SecretKey,
					"",
				),
			))
		} else if config.Endpoint != "" {
			// For LocalStack, provide dummy credentials to avoid EC2 IMDS lookup
			cfgOpts = append(cfgOpts, awsconfig.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider("test", "test", ""),
			))
		}

		cfg, initErr = awsconfig.LoadDefaultConfig(ctx, cfgOpts...)
		if initErr != nil {
			logger.Error("s3 failed to load config", initErr)
			initErr = WrapError(initErr, "Init")
			return
		}

		opts := []func(*s3.Options){}
		if config.Endpoint != "" {
			opts = append(opts, func(s3Opts *s3.Options) {
				s3Opts.BaseEndpoint = aws.String(config.Endpoint)
				s3Opts.UsePathStyle = true
			})
		}

		mu.Lock()
		client = s3.NewFromConfig(cfg, opts...)
		mu.Unlock()
	})
	return initErr
}

// Check if the S3 client is initialized
func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return client != nil
}

// Retrieves an object from S3 as raw bytes
func GetObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	if client == nil {
		logger.Error("s3 client not initialized", fmt.Errorf("s3 client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "GetObject")
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Error("failed to get s3 object", err)
		return nil, WrapError(err, "GetObject")
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		logger.Error("failed to read s3 object body", err)
		return nil, WrapError(err, "GetObject read body")
	}
	return data, nil
}

// Retrieves an S3 object and unmarshals JSON into the provided output interface
func GetObjectAs(ctx context.Context, bucket string, key string, out interface{}) error {
	data, err := GetObject(ctx, bucket, key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		logger.Error("failed to unmarshal s3 object", err)
		return WrapError(err, "GetObjectAs unmarshal")
	}
	return nil
}

// Uploads an object to S3
func PutObject(ctx context.Context, bucket string, key string, data []byte, contentType string) error {
	if client == nil {
		logger.Error("s3 client not initialized", fmt.Errorf("s3 client not initialized"))
		return WrapError(ErrClientNotInitialized, "PutObject")
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	if _, err := client.PutObject(ctx, input); err != nil {
		logger.Error("failed to put s3 object", err)
		return WrapError(err, "PutObject")
	}
	return nil
}

// Deletes an object from S3
func DeleteObject(ctx context.Context, bucket string, key string) error {
	if client == nil {
		logger.Error("s3 client not initialized", fmt.Errorf("s3 client not initialized"))
		return WrapError(ErrClientNotInitialized, "DeleteObject")
	}

	if _, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		logger.Error("failed to delete s3 object", err)
		return WrapError(err, "DeleteObject")
	}
	return nil
}
