package aws

import (
	"context"
	"dumper/internal/connect"
	"dumper/internal/domain/config/storage"
	"dumper/pkg/utils"
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/crypto/ssh"
)

type awsClient struct {
	ctx      context.Context
	Connect  *connect.Connect
	Storage  storage.Storage
	DumpName string
	FileSize int64
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
}

func NewClient(
	ctx context.Context,
	connect *connect.Connect,
	storage storage.Storage,
	dumpName string,
	fileSize int64,
) *awsClient {
	return &awsClient{
		ctx:      ctx,
		Connect:  connect,
		Storage:  storage,
		DumpName: dumpName,
		FileSize: fileSize,
	}
}

func (a *awsClient) Handler() error {
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

	if endpoint := a.endpoint(); endpoint != nil {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = endpoint
			o.UsePathStyle = true
		})
	}

	s3Client := s3.NewFromConfig(awsCfg, opts...)
	uploader := manager.NewUploader(s3Client)

	session, err := a.Connect.NewSession()

	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}

	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := session.Start(fmt.Sprintf("cat %s", a.DumpName)); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
	}

	targetPath := filepath.Join(
		a.Storage.Dir,
		filepath.Base(a.DumpName),
	)

	pr := utils.StreamToPipe(a.ctx, stdout, a.FileSize)

	_, err = uploader.Upload(a.ctx, &s3.PutObjectInput{
		Bucket: aws.String(a.Storage.Bucket),
		Key:    aws.String(targetPath),
		Body:   pr,
	})

	if err != nil {
		return fmt.Errorf("failed to upload to %s: %v", providerName, err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
	}

	utils.SafePrintln("[%s] Upload complete: %s", providerName, targetPath)
	return nil
}

func (a *awsClient) endpoint() *string {
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

	return aws.String(fmt.Sprintf(cfg.Endpoint, a.Storage.Bucket, a.Storage.Region))
}

func (a *awsClient) providerName() string {
	if cfg, ok := defaultAWS[a.Storage.Type]; ok {
		return cfg.Provider
	}
	return "AWS S3"
}
