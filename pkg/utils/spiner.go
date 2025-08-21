package utils

import (
	"fmt"
	"time"
)

func Spinner(stop chan struct{}) {
	chars := `-\|/`
	lenChars := len(chars)
	i := 0
	for {
		select {
		case <-stop:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\rCreating dump... %c", chars[i%lenChars])
			time.Sleep(200 * time.Millisecond)
			i++
		}
	}
}
