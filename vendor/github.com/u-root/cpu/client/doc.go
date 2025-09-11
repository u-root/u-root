// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package client provides an exec.Command and ssh like interface for cpu sessions.
// It attempts to cleave as much as possible to the original.
// The choice between options and environment variables mirrors this effort.
// For example, the nonce for the mount protocol back is an environment variable.
// command name and arguments are passed in os.Args
// The only required parameter for Command() is a host name; if os.Args is empty,
// the remote server reads SHELL and starts a shell.
// Similarly, because the root for the client namespace is known only to the client.
// it is settable in the Cmd struct.
package client
