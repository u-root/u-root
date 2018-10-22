// Copyright 2018 Google Inc.
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

package strace

import "fmt"

// FormatSpecifier values describe how an individual syscall argument should be
// formatted.
type FormatSpecifier int

// Valid FormatSpecifiers.
//
// Unless otherwise specified, values are formatted before syscall execution
// and not updated after syscall execution (the same value is output).
const (
	// Hex is just a hexadecimal number.
	Hex FormatSpecifier = iota

	// Oct is just an octal number.
	Oct

	// ReadBuffer is a buffer for a read-style call. The syscall return
	// value is used for the length.
	//
	// Formatted after syscall execution.
	ReadBuffer

	// WriteBuffer is a buffer for a write-style call. The following arg is
	// used for the length.
	//
	// Contents omitted after syscall execution.
	WriteBuffer

	// ReadIOVec is a pointer to a struct iovec for a writev-style call.
	// The following arg is used for the length. The return value is used
	// for the total length.
	//
	// Complete contents only formatted after syscall execution.
	ReadIOVec

	// WriteIOVec is a pointer to a struct iovec for a writev-style call.
	// The following arg is used for the length.
	//
	// Complete contents only formatted before syscall execution, omitted
	// after.
	WriteIOVec

	// IOVec is a generic pointer to a struct iovec. Contents are not dumped.
	IOVec

	// SendMsgHdr is a pointer to a struct msghdr for a sendmsg-style call.
	// Contents formatted only before syscall execution, omitted after.
	SendMsgHdr

	// RecvMsgHdr is a pointer to a struct msghdr for a recvmsg-style call.
	// Contents formatted only after syscall execution.
	RecvMsgHdr

	// Path is a pointer to a char* path.
	Path

	// PostPath is a pointer to a char* path, formatted after syscall
	// execution.
	PostPath

	// ExecveStringVector is a NULL-terminated array of strings. Enforces
	// the maximum execve array length.
	ExecveStringVector

	// PipeFDs is an array of two FDs, formatted after syscall execution.
	PipeFDs

	// Uname is a pointer to a struct uname, formatted after syscall execution.
	Uname

	// Stat is a pointer to a struct stat, formatted after syscall execution.
	Stat

	// SockAddr is a pointer to a struct sockaddr. The following arg is
	// used for length.
	SockAddr

	// PostSockAddr is a pointer to a struct sockaddr, formatted after
	// syscall execution. The following arg is a pointer to the socklen_t
	// length.
	PostSockAddr

	// SockLen is a pointer to a socklen_t, formatted before and after
	// syscall execution.
	SockLen

	// SockFamily is a socket protocol family value.
	SockFamily

	// SockType is a socket type and flags value.
	SockType

	// SockProtocol is a socket protocol value. Argument n-2 is the socket
	// protocol family.
	SockProtocol

	// SockFlags are socket flags.
	SockFlags

	// Timespec is a pointer to a struct timespec.
	Timespec

	// PostTimespec is a pointer to a struct timespec, formatted after
	// syscall execution.
	PostTimespec

	// UTimeTimespec is a pointer to a struct timespec. Formatting includes
	// UTIME_NOW and UTIME_OMIT.
	UTimeTimespec

	// ItimerVal is a pointer to a struct itimerval.
	ItimerVal

	// PostItimerVal is a pointer to a struct itimerval, formatted after
	// syscall execution.
	PostItimerVal

	// ItimerSpec is a pointer to a struct itimerspec.
	ItimerSpec

	// PostItimerSpec is a pointer to a struct itimerspec, formatted after
	// syscall execution.
	PostItimerSpec

	// Timeval is a pointer to a struct timeval, formatted before and after
	// syscall execution.
	Timeval

	// Utimbuf is a pointer to a struct utimbuf.
	Utimbuf

	// Rusage is a struct rusage, formatted after syscall execution.
	Rusage

	// CloneFlags are clone(2) flags.
	CloneFlags

	// OpenFlags are open(2) flags.
	OpenFlags

	// Mode is a mode_t.
	Mode

	// FutexOp is the futex(2) operation.
	FutexOp

	// PtraceRequest is the ptrace(2) request.
	PtraceRequest

	// ItimerType is an itimer type (ITIMER_REAL, etc).
	ItimerType
)

// defaultFormat is the syscall argument format to use if the actual format is
// not known. It formats all six arguments as hex.
var defaultFormat = []FormatSpecifier{Hex, Hex, Hex, Hex, Hex, Hex}

func defaultSyscallInfo(sysno int) *SyscallInfo {
	return &SyscallInfo{name: fmt.Sprintf("%d", sysno), format: defaultFormat}
}

// SyscallInfo captures the name and printing format of a syscall.
type SyscallInfo struct {
	// name is the name of the syscall.
	name string

	// format contains the format specifiers for each argument.
	//
	// Syscall calls can have up to six arguments. Arguments without a
	// corresponding entry in format will not be printed.
	format []FormatSpecifier
}

// makeSyscallInfo returns a SyscallInfo for a syscall.
func makeSyscallInfo(name string, f ...FormatSpecifier) SyscallInfo {
	return SyscallInfo{name: name, format: f}
}

// SyscallMap maps syscalls into names and printing formats.
type SyscallMap map[uintptr]SyscallInfo
