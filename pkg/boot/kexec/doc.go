// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

//go:generate go run genpurg.go purgatories.go

// The kexec package contains kexec system call support as well as
// several purgatories. Callers may set the purgatory to use at runtime.
package kexec
