package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Driver stores files in an Amazon S3 (or compatible) bucket.
type S3Driver struct {
	Client   *s3.Client
	Bucket   string
	Region   string
	Endpoint string // optional custom endpoint for S3-compatible services
}

// NewS3Driver creates an S3Driver from environment variables.
//
// Required env vars: S3_BUCKET, S3_REGION, S3_KEY, S3_SECRET.
// Optional: S3_ENDPOINT (for MinIO, DigitalOcean Spaces, etc.).
func NewS3Driver() (*S3Driver, error) {
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("S3_REGION")
	key := os.Getenv("S3_KEY")
	secret := os.Getenv("S3_SECRET")
	endpoint := os.Getenv("S3_ENDPOINT")

	if bucket == "" || region == "" || key == "" || secret == "" {
		return nil, errors.New("S3_BUCKET, S3_REGION, S3_KEY, and S3_SECRET are required")
	}

	opts := func(o *s3.Options) {
		o.Region = region
		o.Credentials = credentials.NewStaticCredentialsProvider(key, secret, "")
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	}

	client := s3.New(s3.Options{}, opts)

	return &S3Driver{
		Client:   client,
		Bucket:   bucket,
		Region:   region,
		Endpoint: endpoint,
	}, nil
}

// safePath validates the path doesn't contain traversal attempts.
func (d *S3Driver) safePath(path string) (string, error) {
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "." || clean == "" {
		return "", errors.New("empty path")
	}
	if strings.HasPrefix(clean, "..") || strings.Contains(clean, "/../") {
		return "", errors.New("path traversal detected")
	}
	clean = strings.TrimPrefix(clean, "/")
	if clean == "" {
		return "", errors.New("empty path")
	}
	return clean, nil
}

// Put uploads content to the S3 bucket at the given key.
func (d *S3Driver) Put(path string, content io.Reader) (string, error) {
	key, err := d.safePath(path)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	_, err = d.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(key),
		Body:   content,
	})
	if err != nil {
		return "", fmt.Errorf("s3 put %q: %w", key, err)
	}
	return key, nil
}

// Get retrieves an object from S3 and returns its body as io.ReadCloser.
func (d *S3Driver) Get(path string) (io.ReadCloser, error) {
	key, err := d.safePath(path)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	out, err := d.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get %q: %w", key, err)
	}
	return out.Body, nil
}

// Delete removes an object from S3.
func (d *S3Driver) Delete(path string) error {
	key, err := d.safePath(path)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = d.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete %q: %w", key, err)
	}
	return nil
}

// URL returns the public URL for the given S3 object.
func (d *S3Driver) URL(path string) string {
	key := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "/")
	if d.Endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", strings.TrimRight(d.Endpoint, "/"), d.Bucket, key)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", d.Bucket, d.Region, key)
}
