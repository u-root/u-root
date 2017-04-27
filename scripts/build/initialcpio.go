// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
)

func init() {
	builders["initialcpio"] = initialCpioBuilder{}
}

type initialCpioBuilder struct {
}

// TODO: This builder is not yet implemented.
func (b initialCpioBuilder) generate(config Config) ([]file, error) {
	// TODO: read contents of an cpio and return the file array.
	return nil, errors.New("initialCpio builder not implemented yet")
}
