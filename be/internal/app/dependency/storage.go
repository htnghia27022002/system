package dependency

import (
	"fmt"

	"be/pkg/storage"
)

// NewStorage builds the STORAGE_DRIVER object store.
func NewStorage(infra *Infra) (storage.Storage, error) {
	cfg := infra.Config.Storage
	switch cfg.Driver {
	case "", "local":
		return storage.NewLocal(infra.Config.UploadDir, cfg.PublicURL), nil
	case "s3":
		return storage.NewS3(storage.S3Config{
			Endpoint:        cfg.S3Endpoint,
			Region:          cfg.S3Region,
			Bucket:          cfg.S3Bucket,
			AccessKeyID:     cfg.S3AccessKeyID,
			SecretAccessKey: cfg.S3SecretKey,
			UseSSL:          cfg.S3UseSSL,
			ForcePathStyle:  cfg.S3ForcePathStyle,
			PublicURL:       cfg.S3PublicURL,
		})
	default:
		return nil, fmt.Errorf("unknown STORAGE_DRIVER %q (want local or s3)", cfg.Driver)
	}
}
