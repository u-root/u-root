// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scuzz supports direct access to SCSI or SATA devices.
// SCSI and ATA used to be different, but for SATA, it's all the same look.
//
// In the long term we can use it to implement hdparm(1) and other
// Linux commands.
//
// This package only supports post-2003 48-bit lba addressing.
// Further, we only concern ourselves with ATA_16.
// For now it only works on Linux.
//
// Other info:
//
//	http://www.t13.org/ Technical Committee T13 AT Attachment (ATA/ATAPI) Interface.
//	http://www.serialata.org/ Serial ATA International Organization.
package scuzz
