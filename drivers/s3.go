package drivers

import (
	"bytes"
	"context"
	"github.com/vanvanni/lampofs/errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Driver struct {
	client     *s3.Client
	bucketName string
}

type S3Options struct {
	Region     string
	BucketName string
	Endpoint   string
}

func NewS3Driver(opts S3Options) (*S3Driver, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(opts.Region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	if opts.Endpoint != "" {
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.EndpointResolver = s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: opts.Endpoint,
				}, nil
			})
		})
	}

	return &S3Driver{
		client:     client,
		bucketName: opts.BucketName,
	}, nil
}

func (d *S3Driver) Read(path string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	result, err := d.client.GetObject(context.TODO(), input)
	if err != nil {
		// Check if it's a "not found" error
		// In a real implementation, you'd check the specific error type
		return nil, errors.ErrFileNotFound
	}

	return result.Body, nil
}

func (d *S3Driver) Write(path string, data []byte) error {
	_, err := d.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	})

	if err == nil {
		return errors.ErrFileExists
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	}

	_, err = d.client.PutObject(context.TODO(), input)
	return err
}

func (d *S3Driver) Put(path string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	}

	_, err := d.client.PutObject(context.TODO(), input)
	return err
}

func (d *S3Driver) Delete(path string) error {
	_, err := d.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	})

	if err != nil {
		return errors.ErrFileNotFound
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	_, err = d.client.DeleteObject(context.TODO(), input)
	return err
}

func (d *S3Driver) Update(path string, data []byte, prepend bool) error {
	var existingData []byte

	input := &s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	result, err := d.client.GetObject(context.TODO(), input)
	if err == nil {
		defer result.Body.Close()
		existingData, err = io.ReadAll(result.Body)
		if err != nil {
			return err
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

	_, err = d.client.PutObject(context.TODO(), putInput)
	return err
}
