// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build !linux

package rtc

import (
	"errors"
)

func OpenRTC() (*RTC, error) {
	return nil, errors.New("not implemented on this platform")
}
