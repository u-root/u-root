// Copyright 2025 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && tinygo

package cpuid

// In u-root when building cmdlets using tinygo we face the issue, that
// non-standard libraries such as golang.org/x/sys/cpu.cpuid rely on
// golang-style assembly to perform low-level operations. In tinygo we cannot do
// that so we have to use CGO inline assembly.
//
// Since this approach is not upstreamable, this makes use of the risky
// //go:linkname directive to overwrite the paths for xgetbv and cpuid in the
// aforementioned package we don't own.
// This behavior is only prevalent for tinygo builds, not regular go builds.
//
// TODO: @leongross
// Further, the cpuid has to be imported into the build context to take
// effect. In the future when we build the busybox we can ensure that in the
// build progress; for now to build individual cmdlets it should be imported
// into the respective cmds as follows:
//
// ```go
// import _ "github.com/u-root/cpuid"
// ```
//
// In the future, we still need to have this bridge in place, since tinygo
// will not have go-style assembly support. Anyway, it should be moved from
// this repo to the u-root repo, since it is not really part of the cpu package
// functionality but a workaround.

//go:linkname xgetbv vendor/golang.org/x/sys/cpu.xgetbv
func xgetbv(arg1 uint32) (eax, edx uint32) {
	return xgetbv_low(arg1)
}

//go:linkname cpuid vendor/golang.org/x/sys/cpu.cpuid
func cpuid(arg1, arg2 uint32) (eax, ebx, ecx, edx uint32) {
	return cpuid_low(arg1, arg2)
}
