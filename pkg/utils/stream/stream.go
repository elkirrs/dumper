package stream

import (
	"context"
	"dumper/internal/connect"
	"dumper/pkg/utils/progress"
	"fmt"
	"io"
	"path/filepath"
)

func PipeReader(
	ctx context.Context,
	stdout io.Reader,
	fileSize int64,
) *io.PipeReader {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		buf := make([]byte, 32*1024)
		var uploaded int64

		for {
			if err := ctx.Err(); err != nil {
				_ = pw.CloseWithError(fmt.Errorf("upload cancelled: %w", err))
				return
			}

			n, err := stdout.Read(buf)
			if n > 0 {
				uploaded += int64(n)

				if gp, ok := ctx.Value("globalProgress").(*progress.GlobProgress); ok {
					gp.Add(int64(n))
				} else {
					progress.Progress(uploaded, fileSize)
				}

				if _, writeErr := pw.Write(buf[:n]); writeErr != nil {
					return
				}
			}

			if err == io.EOF {
				_ = pw.Close()
				return
			}
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()

	return pr
}

func SSHStreamer(
	ctx context.Context,
	conn *connect.Connect,
	dumpName string,
	fileSize int64,
) (*io.PipeReader, func() error, error) {
	session, err := conn.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := session.Start(fmt.Sprintf("cat %s", dumpName)); err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("failed to start remote command: %w", err)
	}

	pr := PipeReader(ctx, stdout, fileSize)

	closeFunc := func() error {
		_ = session.Wait()
		return session.Close()
	}

	return pr, closeFunc, nil
}

func TargetPath(dir, dumpName string) string {
	return filepath.Join(dir, filepath.Base(dumpName))
}
