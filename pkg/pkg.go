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
)

type Runnable interface {
	Run(args ...string) int
	RunContext(ctx context.Context, args ...string) int
}
