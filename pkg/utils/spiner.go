package utils

import (
	"fmt"
	"time"
)

func Spinner(stop chan struct{}) {
	chars := []string{"-", "\\", "|", "/"}
	lenChars := len(chars)
	i := 0

	start := time.Now()

	spinTicker := time.NewTicker(200 * time.Millisecond)
	defer spinTicker.Stop()

	timeTicker := time.NewTicker(50 * time.Millisecond)
	defer timeTicker.Stop()

	var elapsed float64

	for {
		select {
		case <-stop:
			fmt.Print("\r")
			return
		case <-spinTicker.C:
			fmt.Printf("\rCreating dump... %s  |  Elapsed: %.2f sec ", chars[i%lenChars], elapsed)
			i++
		case <-timeTicker.C:
			elapsed = time.Since(start).Seconds()
		}
	}
}
