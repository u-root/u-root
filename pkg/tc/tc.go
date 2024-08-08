// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"errors"

	"github.com/florianl/go-tc"
)

var (
	ErrNotEnoughArgs  = errors.New("not enough argument")
	ErrInvalidArg     = errors.New("invalid argument in list")
	ErrNotImplemented = errors.New("not implemented")
	ErrOutOfBounds    = errors.New("integer argument out of bounds")
)

type Trafficctl struct {
	*tc.Tc
}
