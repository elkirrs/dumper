package console

import (
	"fmt"
	"sync"
)

var consoleMu sync.Mutex

func SafePrintln(format string, a ...interface{}) {
	consoleMu.Lock()
	defer consoleMu.Unlock()

	fmt.Printf("\r\033[K")
	fmt.Printf(format+"\n", a...)
}
