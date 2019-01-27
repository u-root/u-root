// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"errors"
	"io"
)

// PrintDTS prints the FDT in the .dts format.
// TODO: not yet implemented
func (fdt *FDT) PrintDTS(f io.Writer) error {
	return errors.New("not yet implemented")
}
