// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import "log"

// Booter is an interface that defines custom boot types. Implementations can be
// like network boot, local boot, etc.
type Booter interface {
	Boot() error
	TypeName() string
}

// NullBooter is a dummy booter that does nothing. It is used when no other
// booter has been found
type NullBooter struct {
}

// TypeName returns the name of the booter type
func (nb *NullBooter) TypeName() string {
	return "null"
}

// Boot will run the boot procedure. In the case of this NullBooter it will do
// nothing
func (nb *NullBooter) Boot() error {
	log.Printf("Null booter does nothing")
	return nil
}
