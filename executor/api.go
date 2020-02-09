package executor

import "fmt"

type Task interface {
	Pre() error
	PerformAction() error
	Post() error
	Name() string
}

type Plan interface {
	Create() ([]Task, error)
}

type FatalError struct {
	msg string
	err error
}

func (e *FatalError) Error() string {
	return fmt.Sprintf("FatalError: %s", e.msg)
}

func (e *FatalError) Unwrap() error {
	return e.err
}
