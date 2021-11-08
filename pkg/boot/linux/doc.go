// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

//go:generate go run genpurg.go purgatories.go

// The kexec package contains kexec system call support as well as
// several purgatories. Callers may set the purgatory to use at runtime.
//
// There two major ways to use kexec: kexec_load and kexec_file_load.
// The latter is easy: give the kernel a bzimage (86) OR an Image (ARM*)
// and let it do its work.
//
// Sadly, in many cases, kexec_file_load will not work, because the 32-bit
// part of the kernel (maybe?) and/or the purgatory (??) do critical startup.
// Hence, the kexec_load classic system call has become increasingly important.
//
// kexec_load is conceptually simple to use: give the kernel some segments,
// tell it the entry point, let it go. If this sounds well matched to ELF,
// it is: the original kexec was designed to work with ELF, and was even hooked
// into the exec system call: starting a new kernel was as easy as typing
// ./vmlinux
// i.e., kexec was truly just a variant on exec!
//
// In the earliest kexec, as in all the early "kernel boots kernel" implementations (1),
// (and still in Plan 9 today),
// the kernel directly loaded, and started, the next kernel. For a number of reasons,
// kexec introduced the concept of a purgatory. The purgatory in principle is both simple
// and elegant: a small bit of code, supplied by user space, that manages the transition
// from one kernel to the next, and vanishes.
//
// The purgatory has a few main responsibilities:
// o (optionally) copy the new kernel over top of the old kernel
// o do any special device setup that neither kernel can manage (mainly console)
// o run a SHA256 over the kernel
// o communicate arguments to the new kernel (on x86, assembly linux params at 0x90000)
// o be able to return to the caller if things go wrong
// o run anywhere, because we may be booting a 16-bit kernel
//
// That last item is the one that causes a lot of trouble. In 2000, systems with 16 MiB
// were still common, Linux kernels had to load at 0x100000, memtest86 had to load in the
// low 640k, and finding a place to put the purgatory required that it
// be a position independent program. Rather than being written as such, it was instead
// compiled as a relocatable ELF, the relocation being done at kexec time. I.e.,
// kexec includes a link step.
//
// Because processors have changed a lot since 2000, when kexec was first written, these
// old assumptions are worth re-examining.
//
// First, systems that use kexec come with at least 1 GiB of memory nowadays. Further,
// newer kernels always avoid using the low 1MiB, since buggy BIOSes might corrupt memory.
// Finally, nobody cares about booting 16-bit kernels any more -- even memtest86 runs
// as a Linux binary.
// Hence: we can always get space in the low 1MiB for the purgatory, and in fact
// we can assume that memory is available at 0x3000. The low 640K must alway be there.
// That means we can link the purgatory to
// run at a fixed place -- since most kexec users load it at a fixed place anyway.
//
// Second, with relocatable kernels, the copy function is no longer needed.
//
// Third, we can dispense with ideas of returning. If things are so messed up that we
// can not kexec, it is likely time to reset the machine. Should we desire to implement return
// later, however, we need not use the messy mechanisms in the current purgatories to save
// registers. If we mark the function with a returnstwice attribute, gcc will use caller-save
// semantics for the call, not callee-save, removing any need to worry about saving registers.
//
// Hence, we can, should we care, arrange for the purgatory right up to the point that it
// drops to 32-bit unpaged mode. Because the number of operations from 32-bit to calling the next
// kernel are so few, we do not feel it is necessary to return past that point.
//
// Fourth, parameter passing is unnecessarily messy in the current purgatories. We can
// rewrite that contract: if we consider the first 8 uint64_t in the purgatory, the first can
// be used for a relative jump around the next 7, and those seven quadwords can be used for
// parameter passing.
//
// These changes should let us:
// o build the purgatory as a non-relative ELF, i.e. a statically linked program with one ELF program (segment)
// o and link it at 0x3000; the code was putting the current relative ELF in a fixed place anyway
// o use the ELF program header to tell us where to put the purgatory
// o communicate arguments in the seven quadwords mentioned above
// o rather than one does-it-all purgatory as we have today, we can provide several variants
//   so we get one suited to the job at hand.
//
// This should result in a dramatically simpler purgatory implementation. Also, being much simpler, it can be entirely
// Go assembly, obviating the need for a C compiler. This preserves a desired property of u-root: that it can
// always be built with only the Go toolchain.
//
// (1) "Give your bootstrap the boot: using the operating system to boot the operating system"
// Ron Minnich,  2004 IEEE International Conference on Cluster Computing (IEEE Cat. No.04EX935)
package linux
