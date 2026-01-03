package minio

import (
	"context"
	"dumper/internal/domain/storage"
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

type Minio struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *Minio {
	return &Minio{
		ctx:    ctx,
		config: config,
	}
}

func (m *Minio) Save() error {
	cred := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			m.config.Config.AccessKey,
			m.config.Config.SecretKey,
			"",
		),
	)

	awsCfg, err := config.LoadDefaultConfig(
		m.ctx,
		config.WithRegion(m.config.Config.Region),
		config.WithCredentialsProvider(cred),
	)

	if err != nil {
		return fmt.Errorf("failed to load MinIO config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(m.config.Config.Endpoint)
	})

	uploader := manager.NewUploader(s3Client)

	session, err := m.config.Conn.NewSession()

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

	if err := session.Start(fmt.Sprintf("cat %s", m.config.DumpName)); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
	}

	targetPath := filepath.Join(
		m.config.Config.Dir,
		filepath.Base(m.config.DumpName),
	)

	pr := utils.StreamToPipe(m.ctx, stdout, m.config.FileSize)

	_, err = uploader.Upload(m.ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.config.Config.Bucket),
		Key:    aws.String(targetPath),
		Body:   pr,
	})

	if err != nil {
		return fmt.Errorf("failed to upload to MinIO: %v", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
	}

	utils.SafePrintln("[MinIO] Upload complete: %s", targetPath)
	return nil
}
