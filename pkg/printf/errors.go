package printf

import (
	"errors"
)

var (
	ErrPrintf             = errors.New("printf")
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrUnimplemented      = errors.New("unimplemented")
	ErrInvalidDirective   = errors.New("invalid directive")
)
