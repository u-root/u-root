package core

import (
	"context"
	"io"
)

// LookupEnvFunc is a function type that looks up an environment variable.
type LookupEnvFunc func(string) (string, bool)

// Command is an interface that defines the methods for executing a command.
type Command interface {
	SetIO(stdin io.Reader, stdout io.Writer, stderr io.Writer)
	SetWorkingDir(workingDir string)
	SetLookupEnv(lookupEnv LookupEnvFunc)
	Run(ctx context.Context, args ...string) (int, error)
}
