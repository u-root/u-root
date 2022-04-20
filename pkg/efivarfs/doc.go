// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// DISCLAIMER
// The code of pkg/efivarfs was originally based on this file
// github.com/canonical/go-efilib/blob/main/vars_linux.go so
// credits go their authors. The code has undergone several revisions
// and has been changed up several times, so that only most crucial parts
// still look familiar to the old code. At the date of writing the code
// in u-root has been simplified so much, that it doesn't resemble most of
// the old code anymore. If a more broader set of features and edge cases
// is required then users should look into importing canonicals library
// instead, as this code only covers the bare minimum!
// Canonical has been made aware this simplified version of their file
// exists and the reason a simplified version had to be created in the
// first place were license concerns.

// The code allows interaction with the efivarfs via u-root which means
// both read and write support. As the efivarfs is part of the Linux kernel,
// this code is only intended to be run under Linux environments.

package efivarfs
