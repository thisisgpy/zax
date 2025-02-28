package util

import "fmt"

type ZaxError struct {
	Message string `json:"message"`
}

func (e *ZaxError) Error() string {
	return e.Message
}

func NewZaxError(message string) *ZaxError {
	return &ZaxError{Message: message}
}

func NewZaxErrorf(format string, args ...interface{}) *ZaxError {
	return &ZaxError{Message: fmt.Sprintf(format, args...)}
}
