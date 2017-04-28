// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
)

func init() {
	buildGenerators["bb"] = bbGenerator{}
}

type bbGenerator struct {
}

// TODO: This generator is not yet implemented.
func (g bbGenerator) generate(config Config) ([]file, error) {
	return nil, errors.New("bb generator not implemented yet")
}
