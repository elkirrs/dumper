package utils

import (
	"fmt"
	"sync/atomic"
)

func Progress(done, total int64) {
	if total == 0 {
		fmt.Printf("\rDownloaded: %d bytes", done)
		return
	}
	percent := float64(done) / float64(total) * 100
	fmt.Printf("\rDownloading... %.1f%% [%d/%d bytes]\n", percent, done, total)
}

type GlobProgress struct {
	total     int64
	completed int64
}

func GlobalProgress(total int64) *GlobProgress {
	return &GlobProgress{total: total}
}

func (p *GlobProgress) Add(n int64) {
	atomic.AddInt64(&p.completed, n)
	p.Print()
}

func (p *GlobProgress) Print() {
	done := atomic.LoadInt64(&p.completed)
	percent := float64(done) / float64(p.total) * 100
	fmt.Printf("\rDownloading... %.2f%% [%d/%d bytes]", percent, done, p.total)
}
