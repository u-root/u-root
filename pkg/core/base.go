// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package core

import (
	"io"
	"os"
	"path/filepath"
)

// Base is a struct that holds common fields for commands.
type Base struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	WorkingDir string
	LookupEnv  LookupEnvFunc
}

// Init initializes the Base command with default values.
func (b *Base) Init() {
	b.Stdin = os.Stdin
	b.Stdout = os.Stdout
	b.Stderr = os.Stderr
	b.WorkingDir = ""
	b.LookupEnv = os.LookupEnv
}

// SetIO sets the input/output streams for the command.
func (b *Base) SetIO(stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	b.Stdin = stdin
	b.Stdout = stdout
	b.Stderr = stderr
}

// SetWorkingDir sets the working directory for the command.
func (b *Base) SetWorkingDir(workingDir string) {
	b.WorkingDir = workingDir
}

// SetLookupEnv sets the function used to look up environment variables.
func (b *Base) SetLookupEnv(lookupEnv LookupEnvFunc) {
	b.LookupEnv = lookupEnv
}

// Getenv is a helper to retrieve an environment variable value without the
// extra bool return.
func (b *Base) Getenv(key string) string {
	v, _ := b.LookupEnv(key)
	return v
}

// ResolvePath resolves a path relative to the working directory.
func (b *Base) ResolvePath(path string) string {
	if filepath.IsAbs(path) || b.WorkingDir == "" {
		return path
	}
	return filepath.Join(b.WorkingDir, path)
}
