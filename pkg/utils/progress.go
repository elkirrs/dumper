package utils

import "fmt"

func Progress(done, total int64) {
	if total == 0 {
		fmt.Printf("\rDownloaded: %d bytes", done)
		return
	}
	percent := float64(done) / float64(total) * 100
	fmt.Printf("\rDownloading... %.1f%% [%d/%d bytes]", percent, done, total)
}
