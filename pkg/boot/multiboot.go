// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/cpio"
)

// multibootImage is a multiboot-formated OSImage.
type multibootImage struct{}

var _ OSImage = &multibootImage{}

func newMultibootImage(a *cpio.Archive) (OSImage, error) {
	return nil, fmt.Errorf("multiboot images unimplemented")
}

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

// Pack implements OSImage.Pack.
func (multibootImage) Pack(sw cpio.RecordWriter) error {
	return fmt.Errorf("multiboot images unimplemented")
}
