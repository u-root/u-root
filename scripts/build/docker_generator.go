// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
)

func init() {
	archiveGenerators["docker"] = dockerGenerator{}
}

type dockerGenerator struct {
}

// TODO: Generate a docker image.
func (g dockerGenerator) generate(config Config, files []file) error {
	return errors.New("docker generator not implemented yet")
}

// TODO: Run the docker image.
func (g dockerGenerator) run(config Config) error {
	return errors.New("docker generator not implemented yet")
}
