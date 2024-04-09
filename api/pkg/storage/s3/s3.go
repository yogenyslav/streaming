package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yogenyslav/logger"
)

type S3 struct {
	conn *minio.Client
}

func MustNew(cfg *Config) *S3 {
	minioClient, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Ssl,
	})
	if err != nil {
		logger.Panic(err)
	}

	return &S3{
		conn: minioClient,
	}
}
