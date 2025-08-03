package drivers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	lampofsErrors "github.com/vanvanni/lampofs/errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Driver struct {
	client     *s3.Client
	bucketName string
}

type S3Options struct {
	Region          string
	BucketName      string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Timeout         time.Duration
}

func NewS3Driver(opts S3Options) (*S3Driver, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 10 * time.Second
	}

	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(opts.Region),
	}

	if opts.AccessKeyID != "" && opts.SecretAccessKey != "" {
		configOptions = append(configOptions,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				opts.AccessKeyID,
				opts.SecretAccessKey,
				"",
			)),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, configOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	if opts.Endpoint != "" {
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.EndpointResolver = s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               opts.Endpoint,
					HostnameImmutable: true,
					SigningRegion:     opts.Region,
				}, nil
			})
			o.UsePathStyle = true
		})
	}

	driver := &S3Driver{
		client:     client,
		bucketName: opts.BucketName,
	}

	if err := driver.testConnection(ctx); err != nil {
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	return driver, nil
}

func (d *S3Driver) testConnection(ctx context.Context) error {
	_, err := d.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(d.bucketName),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "NotFound", "NoSuchBucket":
				return fmt.Errorf("bucket %s does not exist: %w", d.bucketName, err)
			case "Forbidden":
				return fmt.Errorf("access denied to bucket %s: %w", d.bucketName, err)
			}
		}
		return fmt.Errorf("failed to connect to S3 bucket %s: %w", d.bucketName, err)
	}

	return nil
}

func (d *S3Driver) Read(path string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	result, err := d.client.GetObject(ctx, input)
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, lampofsErrors.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return result.Body, nil
}

func (d *S3Driver) Write(path string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := d.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	})

	if err == nil {
		return lampofsErrors.ErrFileExists
	}

	var noSuchKey *types.NoSuchKey
	if !errors.As(err, &noSuchKey) {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() != "NotFound" {
			return fmt.Errorf("unexpected error checking if file exists: %w", err)
		}
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	}

	_, err = d.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

func (d *S3Driver) Put(path string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	}

	_, err := d.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put file %s: %w", path, err)
	}

	return nil
}

func (d *S3Driver) Delete(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := d.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	})

	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return lampofsErrors.ErrFileNotFound
		}
		return fmt.Errorf("error checking if file exists: %w", err)
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	_, err = d.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}

	return nil
}

func (d *S3Driver) Update(path string, data []byte, prepend bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var existingData []byte

	input := &s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	result, err := d.client.GetObject(ctx, input)
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if !errors.As(err, &noSuchKey) {
			return fmt.Errorf("failed to get file for update: %w", err)
		}
	} else {
		defer result.Body.Close()
		existingData, err = io.ReadAll(result.Body)
		if err != nil {
			return fmt.Errorf("failed to read existing file content: %w", err)
		}
	}

	var newData []byte
	if prepend {
		newData = make([]byte, len(data)+len(existingData))
		copy(newData, data)
		copy(newData[len(data):], existingData)
	} else {
		newData = append(existingData, data...)
	}

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(newData),
	}

	_, err = d.client.PutObject(ctx, putInput)
	if err != nil {
		return fmt.Errorf("failed to update file %s: %w", path, err)
	}

	return nil
}
