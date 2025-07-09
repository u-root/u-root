// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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
	Run(args ...string) error
	RunContext(ctx context.Context, args ...string) error
}
