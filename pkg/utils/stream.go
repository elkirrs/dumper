package utils

import (
	"context"
	"fmt"
	"io"
)

func StreamToPipe(
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
			select {
			case <-ctx.Done():
				_ = pw.CloseWithError(fmt.Errorf("upload cancelled by context"))
				return
			default:
			}

			n, err := stdout.Read(buf)
			if n > 0 {
				uploaded += int64(n)
				if gp, ok := ctx.Value("globalProgress").(*GlobProgress); ok {
					gp.Add(int64(n))
				} else {
					Progress(uploaded, fileSize)
				}

				if _, writeErr := pw.Write(buf[:n]); writeErr != nil {
					return
				}
			}

			if err == io.EOF {
				break
			}
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()

	return pr
}
