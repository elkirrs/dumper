package aws

import (
	"context"
	"dumper/internal/connect"
	domainConfigStorage "dumper/internal/domain/config/storage"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils/console"
	"dumper/pkg/utils/stream"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	ctx      context.Context
	Connect  *connect.Connect
	Storage  domainConfigStorage.Storage
	DumpName string
	FileSize int64
	Backend  string
}

type DefaultAWS struct {
	Endpoint string
	Region   string
	Provider string
}

var defaultAWS = map[string]DefaultAWS{
	"s3": {
		Region:   "us-east-1",
		Provider: "AWS S3",
	},
	"r2": {
		Endpoint: "https://%s.r2.cloudflarestorage.com",
		Region:   "us-east-1",
		Provider: "Cloudflare R2",
	},
	"minio": {
		Region:   "us-east-1",
		Provider: "MinIO",
	},
	"b2": {
		Endpoint: "https://%s.s3.%s.backblazeb2.com",
		Region:   "us-east-001",
		Provider: "Backblaze B2",
	},
	"spaces": {
		Endpoint: "https://%s.%s.digitaloceanspaces.com",
		Region:   "fra1",
		Provider: "DigitalOcean Spaces",
	},
	"yandex": {
		Endpoint: "https://storage.yandexcloud.net",
		Region:   "ru-central1",
		Provider: "Yandex Cloud",
	},
}

func NewClient(
	ctx context.Context,
	connect *connect.Connect,
	storage domainConfigStorage.Storage,
	dumpName string,
	fileSize int64,
	backend string,
) *Client {
	return &Client{
		ctx:      ctx,
		Connect:  connect,
		Storage:  storage,
		DumpName: dumpName,
		FileSize: fileSize,
		Backend:  backend,
	}
}

func (a *Client) Handler() error {
	cred := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			a.Storage.AccessKey,
			a.Storage.SecretKey,
			"",
		),
	)

	awsCfg, err := config.LoadDefaultConfig(
		a.ctx,
		config.WithRegion(a.Storage.Region),
		config.WithCredentialsProvider(cred),
	)

	providerName := a.providerName()
	if err != nil {
		return fmt.Errorf("failed to load %s config: %v", providerName, err)
	}

	var opts []func(*s3.Options)

	if endpoint := a.endpoint(); endpoint != nil {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = endpoint
			pathStyle := true
			if a.Storage.Type == "s3" {
				pathStyle = false
			}
			o.UsePathStyle = pathStyle
		})
	}

	s3Client := s3.NewFromConfig(awsCfg, opts...)

	_, err = s3Client.HeadBucket(a.ctx, &s3.HeadBucketInput{
		Bucket: aws.String(a.Storage.Bucket),
	})
	if err != nil {
		return &storage.UploadError{
			Backend: a.Backend,
			Err:     fmt.Errorf("[%s] bucket %s is not accessible: %w", providerName, a.Storage.Bucket, err),
		}
	}

	uploader := manager.NewUploader(s3Client)

	pr, closeSSH, err := stream.SSHStreamer(a.ctx, a.Connect, a.DumpName, a.FileSize)

	if err != nil {
		return &storage.UploadError{
			Backend: a.Backend,
			Err:     fmt.Errorf("failed to create SSH session: %v", err),
		}
	}

	defer func() {
		_ = closeSSH()
	}()

	targetPath := stream.TargetPath(a.Storage.Dir, a.DumpName)

	_, err = uploader.Upload(a.ctx, &s3.PutObjectInput{
		Bucket: aws.String(a.Storage.Bucket),
		Key:    aws.String(targetPath),
		Body:   pr,
	})

	if err != nil {
		return &storage.UploadError{Backend: a.Backend, Err: err}
	}

	console.SafePrintln("[%s] Upload complete: %s", providerName, targetPath)
	return nil
}

func (a *Client) endpoint() *string {
	if a.Storage.Endpoint != "" {
		return aws.String(a.Storage.Endpoint)
	}

	cfg, ok := defaultAWS[a.Storage.Type]
	if !ok || cfg.Endpoint == "" {
		return nil
	}

	if a.Storage.Type == "r2" {
		return aws.String(fmt.Sprintf(cfg.Endpoint, a.Storage.AccountID))
	}

	if a.Storage.Type == "yandex" {
		return aws.String(cfg.Endpoint)
	}

	return aws.String(fmt.Sprintf(cfg.Endpoint, a.Storage.Bucket, a.Storage.Region))
}

func (a *Client) providerName() string {
	if cfg, ok := defaultAWS[a.Storage.Type]; ok {
		return cfg.Provider
	}
	return "AWS S3"
}
