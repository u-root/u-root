//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package serial

func (port *unixPort) setSpecialBaudrate(speed uint32) error {
	// TODO: unimplemented
	return &PortError{code: InvalidSpeed}
}
