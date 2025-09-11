// Package vm provides an exec-like interface to VMs running cpud. It is
// designed to make running commands in a vm as easy as using exec.
//
// For purposes of convenience, this package provides, via embed, a set of
// kernels and initramfs, for several architectures. Currently, it provides
// Linux kernels and initramfs for amd64, arm, arm64, and riscv64.
//
// The kernels are minimal, hence small; but enable enough options to run
// cpud and u-root programs. The initramfs only contain the cpud and
// dhclient commands, and are also small.
//
// The package asssumes that qemu for a target architecture is available on
// the system.
//
// Users can not yet provide their own kernel and initramfs, though
// that is planned.
package vm
