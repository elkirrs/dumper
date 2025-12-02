package s3

import (
	"context"
	"dumper/internal/domain/storage"
	"dumper/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/crypto/ssh"
)

type S3 struct {
	ctx    context.Context
	config *storage.Config
}

func NewApp(
	ctx context.Context,
	config *storage.Config,
) *S3 {
	return &S3{
		ctx:    ctx,
		config: config,
	}
}

func (s *S3) Save() error {
	cred := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			s.config.Config.AccessKey,
			s.config.Config.SecretKey,
			"",
		),
	)

	awsCfg, err := config.LoadDefaultConfig(
		s.ctx,
		config.WithRegion(s.config.Config.Region),
		config.WithCredentialsProvider(cred),
		config.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	uploader := manager.NewUploader(s3Client)

	session, err := s.config.Conn.NewSession()
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

	if err := session.Start(fmt.Sprintf("cat %s", s.config.DumpName)); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
	}

	key := filepath.Join(s.config.Config.Dir, filepath.Base(s.config.DumpName))

	pr, pw := io.Pipe()

	go func() {
		defer func(pw *io.PipeWriter) {
			_ = pw.Close()
		}(pw)
		buf := make([]byte, 32*1024)
		var uploaded int64

		for {
			select {
			case <-s.ctx.Done():
				_ = pw.CloseWithError(fmt.Errorf("s3 upload cancelled by context"))
				return
			default:
			}

			n, readErr := stdout.Read(buf)
			if n > 0 {
				uploaded += int64(n)
				if gp, ok := s.ctx.Value("globalProgress").(*utils.GlobProgress); ok {
					gp.Add(int64(n))
				} else {
					utils.Progress(uploaded, s.config.FileSize)
				}
				if _, err := pw.Write(buf[:n]); err != nil {
					return
				}
			}
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				_ = pw.CloseWithError(readErr)
				return
			}
		}
	}()

	_, err = uploader.Upload(s.ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Config.Bucket),
		Key:    aws.String(key),
		Body:   pr,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
	}

	fmt.Println("\n[S3] Upload complete:", key)
	return nil
}
