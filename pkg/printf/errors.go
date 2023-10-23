package printf

import (
	"errors"
	"fmt"
)

var (
	ErrPrintf             = errors.New("printf")
	ErrNotEnoughArguments = errors.New("not enough arguments")
)

type ErrInvalidDirective struct {
	directive string
}

func (e ErrInvalidDirective) Error() string {
	return fmt.Sprintf("%s: %s", "%"+e.directive, "invalid directive")
}

func NewErrInvalidDirective(directive string) error {
	return &ErrInvalidDirective{
		directive: directive,
	}
}
