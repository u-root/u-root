// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"os"
)

type DevAPI interface {
	SendRequest(*request) error
	ReceiveResponse(int64, *response, []byte) ([]byte, error)
	GetFile() *os.File
}
