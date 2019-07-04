// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

// Recoverer interface offers recovering
// from critical errors in different ways.
// Currently permissiverecoverer with log
// output and securerecovery with shutdown
// capabilities are supported.
type Recoverer interface {
	Recover(message string) error
}
