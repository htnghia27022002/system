package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Config configures an S3-compatible bucket.
type S3Config struct {
	// Endpoint is host[:port] without scheme, e.g. "s3.amazonaws.com",
	// "<account>.r2.cloudflarestorage.com", "minio:9000".
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	// ForcePathStyle addresses the bucket as endpoint/bucket (MinIO); otherwise
	// the SDK picks virtual-host style where supported.
	ForcePathStyle bool
	// PublicURL is the base clients fetch objects from (bucket website, CDN or
	// public bucket URL), e.g. "https://cdn.example.com". Required.
	PublicURL string
}

// S3 stores objects in an S3-compatible bucket. The bucket (or the CDN in front
// of it) must allow public reads of the objects it serves.
type S3 struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

var _ Storage = (*S3)(nil)

func NewS3(cfg S3Config) (*S3, error) {
	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.PublicURL == "" {
		return nil, fmt.Errorf("storage: S3 needs S3_ENDPOINT, S3_BUCKET and S3_PUBLIC_URL")
	}
	lookup := minio.BucketLookupAuto
	if cfg.ForcePathStyle {
		lookup = minio.BucketLookupPath
	}
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure:       cfg.UseSSL,
		Region:       cfg.Region,
		BucketLookup: lookup,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: S3 client: %w", err)
	}
	return &S3{client: client, bucket: cfg.Bucket, publicURL: strings.TrimRight(cfg.PublicURL, "/")}, nil
}

func (s *S3) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	key, err := CleanKey(key)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *S3) Open(ctx context.Context, key string) (io.ReadCloser, Object, error) {
	key, err := CleanKey(key)
	if err != nil {
		return nil, Object{}, ErrNotFound
	}
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, Object{}, mapS3Err(err)
	}
	info, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, Object{}, mapS3Err(err)
	}
	return obj, Object{Key: key, Size: info.Size, ContentType: info.ContentType}, nil
}

func (s *S3) Delete(ctx context.Context, key string) error {
	key, err := CleanKey(key)
	if err != nil {
		return err
	}
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *S3) URL(key string) string {
	parts := strings.Split(key, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return s.publicURL + "/" + strings.Join(parts, "/")
}

func (s *S3) KeyFromURL(u string) (string, bool) { return keyFromPrefixedURL(s.publicURL, u) }

// Ping checks the bucket is reachable (readiness).
func (s *S3) Ping(ctx context.Context) error {
	ok, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("storage: bucket %q does not exist", s.bucket)
	}
	return nil
}

func mapS3Err(err error) error {
	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return ErrNotFound
	}
	return err
}
