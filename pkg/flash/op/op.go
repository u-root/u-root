// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package op contains available SPI opcods.
package op

// Op is sent at the beginning of a SPI transaction.
type Op byte

const (
	// Read reads from the flash chip.
	Read Op = 0x03
	// WriteDisable disables writing.
	WriteDisable Op = 0x04
	// ReadStatus reads the status register.
	ReadStatus Op = 0x05
	// WriteEnable enables writing.
	WriteEnable Op = 0x06
	// ReadSFDP reads from the SFDP.
	ReadSFDP Op = 0x5a
	// ReadID reads the JEDEC ID.
	ReadID Op = 0x9f
)
