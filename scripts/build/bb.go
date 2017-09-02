// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
)

func init() {
	builders["bb"] = bbBuilder{}
}

type bbBuilder struct {
}

// TODO: This builder is not yet implemented.
func (b bbBuilder) generate(config Config) ([]file, error) {
	return nil, errors.New("bb builder not implemented yet")
}
