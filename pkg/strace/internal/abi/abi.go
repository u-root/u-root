// Copyright 2018 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package abi describes the interface between a kernel and userspace.
// Most of this is from gvisor but I've reordered it a bit as some things
// are common to OSX. I've also followed Go practice for these mostly
// never-viewed files and put them into a small number of largish files.
package abi

import (
	"fmt"
)

// OS describes the target operating system for an ABI.
//
// Note that OS is architecture-independent. The details of the OS ABI will
// vary between architectures.
type OS int

const (
	// Linux is the Linux ABI.
	Linux OS = iota
)

// String implements fmt.Stringer.
func (o OS) String() string {
	switch o {
	case Linux:
		return "linux"
	default:
		return fmt.Sprintf("OS(%d)", o)
	}
}
