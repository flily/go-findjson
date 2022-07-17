package findjson

import "fmt"

type JsonError struct {
	Offset  int
	Message string
}

func (e *JsonError) Error() string {
	m := fmt.Sprintf("JSON error at %d: %s", e.Offset, e.Message)
	return m
}

func NewJsonError(Offset int, message string, args ...interface{}) *JsonError {
	e := &JsonError{
		Offset:  Offset,
		Message: fmt.Sprintf(message, args...),
	}

	return e
}
