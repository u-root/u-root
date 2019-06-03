// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"log"
)

// multibootImage is a multiboot-formated OSImage.
type multibootImage struct{}

var _ OSImage = &multibootImage{}

// ExecutionInfo implements OSImage.ExecutionInfo.
func (multibootImage) ExecutionInfo(log *log.Logger) {
	log.Printf("Multiboot images are unsupported")
}

// Execute implements OSImage.Execute.
func (multibootImage) Execute() error {
	return fmt.Errorf("multiboot images unimplemented")
}

// String implements fmt.Stringer.
func (multibootImage) String() string {
	return fmt.Sprintf("multiboot images unimplemented")
}
