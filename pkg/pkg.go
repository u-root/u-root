// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pkg defines an abstract os package interface and standardized RunMain func
// for supported commands to be executed in different modes: executing as the
// current process like normal, or to run inside the current process like a
// shell builtin with overridable env vars and working dir.

package pkg

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// RunMain is a standard interface for running any command that supports
// execution in the abstracted ExecContext environment.
type RunMain func(e ExecContext, args []string) int

// ExecContext represents the execution context of a command as an abstract
// subset of the os package.
type ExecContext interface {
	// Context returns a context.Context object for this command context.
	Context() context.Context

	// Getwd returns an absolute path name corresponding to the current
	// directory in the command context.
	Getwd() (string, error)

	// Chdir changes the current working directory to the named directory
	// in the command context.
	Chdir(string) error

	// Stdin returns an io.Reader representing the standard input in the
	// command context.
	Stdin() io.Reader

	// Stdout returns an io.Writer representing the standard output in the
	// command context.
	Stdout() io.Writer

	// Stderr returns an io.Writer representing the standard error in the
	// command context.
	Stderr() io.Writer

	// LookupEnv retrieves the value of the environment variable in the command
	// context the named by the key.
	LookupEnv(string) (string, bool)

	// Setenv sets the value of the environment variable in the command context
	// named by the key. It returns an error, if any.
	SetEnv(string, string) error

	// Unsetenv unsets a single environment variable in the command context.
	UnsetEnv(string) error
}

// OsContext returns an ExecContext which wraps the os funcs with the
// corresponding name. Used to make a command interact with the live
// process state normally.
func OsContext(ctx context.Context) ExecContext { return osctx{ctx} }

type osctx struct{ ctx context.Context }

func (e osctx) Context() context.Context          { return e.ctx }
func (osctx) Chdir(dir string) error              { return os.Chdir(dir) }
func (osctx) Getwd() (string, error)              { return os.Getwd() }
func (osctx) Stdin() io.Reader                    { return os.Stdin }
func (osctx) Stdout() io.Writer                   { return os.Stdout }
func (osctx) Stderr() io.Writer                   { return os.Stderr }
func (osctx) LookupEnv(key string) (string, bool) { return os.LookupEnv(key) }
func (osctx) SetEnv(key, value string) error      { return os.Setenv(key, value) }
func (osctx) UnsetEnv(key string) error           { return os.Unsetenv(key) }

// CustomContext returns an ExecContext with custom overrides for method
// return values. Used to run a command as a shell/interpretrer builtin, where
// its os environment does not interact with the host process' os environment.
func CustomContext(ctx context.Context, dir string, initialVars map[string]string, stdin io.Reader, stdout, stderr io.Writer) (ExecContext, error) {
	vars := new(sync.Map)
	for k, v := range initialVars {
		vars.Store(k, v)
	}
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	return &customctx{ctx, dir, vars, stdin, stdout, stderr}, nil
}

type customctx struct {
	ctx            context.Context
	dir            string
	vars           *sync.Map
	stdin          io.Reader
	stdout, stderr io.Writer
}

func (e *customctx) Context() context.Context { return e.ctx }
func (e *customctx) Getwd() (string, error)   { return e.dir, nil }
func (e *customctx) Chdir(dir string) error {
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(e.dir, dir)
	}
	f, err := os.Open(dir)
	if err != nil {
		return &os.PathError{Op: "chdir", Path: dir, Err: err}
	}
	return f.Close()
}
func (e *customctx) Stdin() io.Reader  { return e.stdin }
func (e *customctx) Stdout() io.Writer { return e.stdout }
func (e *customctx) Stderr() io.Writer { return e.stderr }
func (e *customctx) LookupEnv(key string) (string, bool) {
	if val, ok := e.vars.Load(key); ok {
		return val.(string), true
	}
	return "", false
}
func (e *customctx) SetEnv(key, value string) error {
	e.vars.Store(key, value)
	return nil
}
func (e *customctx) UnsetEnv(key string) error {
	e.vars.Delete(key)
	return nil
}
