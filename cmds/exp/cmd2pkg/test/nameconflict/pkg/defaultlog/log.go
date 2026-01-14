// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deflog

import (
	"log"
	"os"
)

func Default() *log.Logger {
	return log.New(os.Stderr, "", 0)
}
