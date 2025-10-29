package connect_error

import "fmt"

type ConnectError struct {
	Addr string
	Err  error
}

func (e *ConnectError) Error() string {
	return fmt.Sprintf("connection error to %s: %v", e.Addr, e.Err)
}

func (e *ConnectError) Unwrap() error {
	return e.Err
}
