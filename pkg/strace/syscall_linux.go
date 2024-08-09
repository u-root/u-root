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

//go:build arm64 || amd64 || riscv64

package strace

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/strace/internal/abi"
	"golang.org/x/sys/unix"
)

// Task is a Linux process.
type Task interface {
	// Read reads from the process at Addr to the interface{}
	// and returns a byte count and error.
	Read(addr Addr, v interface{}) (int, error)

	// Name is a human-readable process identifier. E.g. PID or argv[0].
	Name() string
}

func path(t Task, addr Addr) string {
	path, err := ReadString(t, addr, unix.PathMax)
	if err != nil {
		return fmt.Sprintf("%#x (error decoding path: %s)", addr, err)
	}
	return fmt.Sprintf("%#x %s", addr, path)
}

func utimensTimespec(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	var tim unix.Timespec
	if _, err := t.Read(addr, &tim); err != nil {
		return fmt.Sprintf("%#x (error decoding timespec: %s)", addr, err)
	}

	var ns string
	switch tim.Nsec {
	case unix.UTIME_NOW:
		ns = "UTIME_NOW"
	case unix.UTIME_OMIT:
		ns = "UTIME_OMIT"
	default:
		ns = fmt.Sprintf("%v", tim.Nsec)
	}
	return fmt.Sprintf("%#x {sec=%v nsec=%s}", addr, tim.Sec, ns)
}

func timespec(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	var tim unix.Timespec
	if _, err := t.Read(addr, &tim); err != nil {
		return fmt.Sprintf("%#x (error decoding timespec: %s)", addr, err)
	}
	return fmt.Sprintf("%#x {sec=%v nsec=%v}", addr, tim.Sec, tim.Nsec)
}

func timeval(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	var tim unix.Timeval
	if _, err := t.Read(addr, &tim); err != nil {
		return fmt.Sprintf("%#x (error decoding timeval: %s)", addr, err)
	}

	return fmt.Sprintf("%#x {sec=%v usec=%v}", addr, tim.Sec, tim.Usec)
}

func utimbuf(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	var utim syscall.Utimbuf
	if _, err := t.Read(addr, &utim); err != nil {
		return fmt.Sprintf("%#x (error decoding utimbuf: %s)", addr, err)
	}

	return fmt.Sprintf("%#x {actime=%v, modtime=%v}", addr, utim.Actime, utim.Modtime)
}

func fileMode(mode uint32) string {
	return fmt.Sprintf("%#09o", mode&0x1ff)
}

func stat(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	var stat unix.Stat_t
	if _, err := t.Read(addr, &stat); err != nil {
		return fmt.Sprintf("%#x (error decoding stat: %s)", addr, err)
	}
	return fmt.Sprintf("%#x {dev=%d, ino=%d, mode=%s, nlink=%d, uid=%d, gid=%d, rdev=%d, size=%d, blksize=%d, blocks=%d, atime=%s, mtime=%s, ctime=%s}", addr, stat.Dev, stat.Ino, fileMode(stat.Mode), stat.Nlink, stat.Uid, stat.Gid, stat.Rdev, stat.Size, stat.Blksize, stat.Blocks, time.Unix(stat.Atim.Unix()), time.Unix(stat.Mtim.Unix()), time.Unix(stat.Ctim.Unix()))
}

func itimerval(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	interval := timeval(t, addr)
	value := timeval(t, addr+Addr(binary.Size(unix.Timeval{})))
	return fmt.Sprintf("%#x {interval=%s, value=%s}", addr, interval, value)
}

func itimerspec(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	interval := timespec(t, addr)
	value := timespec(t, addr+Addr(binary.Size(unix.Timespec{})))
	return fmt.Sprintf("%#x {interval=%s, value=%s}", addr, interval, value)
}

func stringVector(t Task, addr Addr) string {
	vs, err := ReadStringVector(t, addr, ExecMaxElemSize, ExecMaxTotalSize)
	if err != nil {
		return fmt.Sprintf("%#x {error copying vector: %v}", addr, err)
	}
	return fmt.Sprintf("%q", vs)
}

func rusage(t Task, addr Addr) string {
	if addr == 0 {
		return "null"
	}

	var ru unix.Rusage
	if _, err := t.Read(addr, &ru); err != nil {
		return fmt.Sprintf("%#x (error decoding rusage: %s)", addr, err)
	}
	return fmt.Sprintf("%#x %+v", addr, ru)
}

// pre fills in the pre-execution arguments for a system call. If an argument
// cannot be interpreted before the system call is executed, then a hex value
// will be used. Note that a full output slice will always be provided, that is
// len(return) == len(args).
func (i *SyscallInfo) pre(t Task, args SyscallArguments, maximumBlobSize uint) []string {
	var output []string
	for arg := range args {
		if arg >= len(i.format) {
			break
		}
		switch i.format[arg] {
		case WriteBuffer:
			output = append(output, dump(t, args[arg].Pointer(), args[arg+1].SizeT(), maximumBlobSize))
		case WriteIOVec:
			output = append(output, iovecs(t, args[arg].Pointer(), int(args[arg+1].Int()), true /* content */, uint64(maximumBlobSize)))
		case IOVec:
			output = append(output, iovecs(t, args[arg].Pointer(), int(args[arg+1].Int()), false /* content */, uint64(maximumBlobSize)))
		case SendMsgHdr:
			output = append(output, msghdr(t, args[arg].Pointer(), true /* content */, uint64(maximumBlobSize)))
		case RecvMsgHdr:
			output = append(output, msghdr(t, args[arg].Pointer(), false /* content */, uint64(maximumBlobSize)))
		case Path:
			output = append(output, path(t, args[arg].Pointer()))
		case ExecveStringVector:
			output = append(output, stringVector(t, args[arg].Pointer()))
		case SockAddr:
			output = append(output, sockAddr(t, args[arg].Pointer(), uint32(args[arg+1].Uint64())))
		case SockLen:
			output = append(output, sockLenPointer(t, args[arg].Pointer()))
		case SockFamily:
			output = append(output, abi.SocketFamily.Parse(uint64(args[arg].Int())))
		case SockType:
			output = append(output, abi.SockType(args[arg].Int()))
		case SockProtocol:
			output = append(output, abi.SockProtocol(args[arg-2].Int(), args[arg].Int()))
		case SockFlags:
			output = append(output, abi.SockFlags(args[arg].Int()))
		case Timespec:
			output = append(output, timespec(t, args[arg].Pointer()))
		case UTimeTimespec:
			output = append(output, utimensTimespec(t, args[arg].Pointer()))
		case ItimerVal:
			output = append(output, itimerval(t, args[arg].Pointer()))
		case ItimerSpec:
			output = append(output, itimerspec(t, args[arg].Pointer()))
		case Timeval:
			output = append(output, timeval(t, args[arg].Pointer()))
		case Utimbuf:
			output = append(output, utimbuf(t, args[arg].Pointer()))
		case CloneFlags:
			output = append(output, abi.CloneFlagSet.Parse(uint64(args[arg].Uint())))
		case OpenFlags:
			output = append(output, abi.Open(uint64(args[arg].Uint())))
		case Mode:
			output = append(output, os.FileMode(args[arg].Uint()).String())
		case FutexOp:
			output = append(output, abi.Futex(uint64(args[arg].Uint())))
		case PtraceRequest:
			output = append(output, abi.PtraceRequestSet.Parse(args[arg].Uint64()))
		case ItimerType:
			output = append(output, abi.ItimerTypes.Parse(uint64(args[arg].Int())))
		case Oct:
			output = append(output, "0o"+strconv.FormatUint(args[arg].Uint64(), 8))
		case Hex:
			fallthrough
		default:
			output = append(output, "0x"+strconv.FormatUint(args[arg].Uint64(), 16))
		}
	}

	return output
}

// post fills in the post-execution arguments for a system call. This modifies
// the given output slice in place with arguments that may only be interpreted
// after the system call has been executed.
func (i *SyscallInfo) post(t Task, args SyscallArguments, rval SyscallArgument, output []string, maximumBlobSize uint) {
	for arg := range output {
		if arg >= len(i.format) {
			break
		}
		switch i.format[arg] {
		case ReadBuffer:
			output[arg] = dump(t, args[arg].Pointer(), uint(rval.Uint64()), maximumBlobSize)
		case ReadIOVec:
			printLength := rval.Uint()
			if printLength > uint32(maximumBlobSize) {
				printLength = uint32(maximumBlobSize)
			}
			output[arg] = iovecs(t, args[arg].Pointer(), int(args[arg+1].Int()), true /* content */, uint64(printLength))
		case WriteIOVec, IOVec, WriteBuffer:
			// We already have a big blast from write.
			output[arg] = "..."
		case SendMsgHdr:
			output[arg] = msghdr(t, args[arg].Pointer(), false /* content */, uint64(maximumBlobSize))
		case RecvMsgHdr:
			output[arg] = msghdr(t, args[arg].Pointer(), true /* content */, uint64(maximumBlobSize))
		case PostPath:
			output[arg] = path(t, args[arg].Pointer())
		case PipeFDs:
			output[arg] = fdpair(t, args[arg].Pointer())
		case Uname:
			output[arg] = uname(t, args[arg].Pointer())
		case Stat:
			output[arg] = stat(t, args[arg].Pointer())
		case PostSockAddr:
			output[arg] = postSockAddr(t, args[arg].Pointer(), args[arg+1].Pointer())
		case SockLen:
			output[arg] = sockLenPointer(t, args[arg].Pointer())
		case PostTimespec:
			output[arg] = timespec(t, args[arg].Pointer())
		case PostItimerVal:
			output[arg] = itimerval(t, args[arg].Pointer())
		case PostItimerSpec:
			output[arg] = itimerspec(t, args[arg].Pointer())
		case Timeval:
			output[arg] = timeval(t, args[arg].Pointer())
		case Rusage:
			output[arg] = rusage(t, args[arg].Pointer())
		}
	}
}

// printEntry prints the given system call entry.
func (i *SyscallInfo) printEnter(t Task, args SyscallArguments) string {
	o := i.pre(t, args, LogMaximumSize)
	switch len(o) {
	case 0:
		return fmt.Sprintf("%s E %s()", t.Name(), i.name)
	case 1:
		return fmt.Sprintf("%s E %s(%s)", t.Name(), i.name, o[0])
	case 2:
		return fmt.Sprintf("%s E %s(%s, %s)", t.Name(), i.name, o[0], o[1])
	case 3:
		return fmt.Sprintf("%s E %s(%s, %s, %s)", t.Name(), i.name, o[0], o[1], o[2])
	case 4:
		return fmt.Sprintf("%s E %s(%s, %s, %s, %s)", t.Name(), i.name, o[0], o[1], o[2], o[3])
	case 5:
		return fmt.Sprintf("%s E %s(%s, %s, %s, %s, %s)", t.Name(), i.name, o[0], o[1], o[2], o[3], o[4])
	case 6:
		return fmt.Sprintf("%s E %s(%s, %s, %s, %s, %s, %s)", t.Name(), i.name, o[0], o[1], o[2], o[3], o[4], o[5])
	default:
		return fmt.Sprintf("%s E %s(%s, %s, %s, %s, %s, %s, ...)", t.Name(), i.name, o[0], o[1], o[2], o[3], o[4], o[5])
	}
}

// SysCallEnter is called each time a system call enter event happens.
func SysCallEnter(t Task, s *SyscallEvent) string {
	i := defaultSyscallInfo(s.Sysno)
	if v, ok := syscalls[uintptr(s.Sysno)]; ok {
		*i = v
	}
	return i.printEnter(t, s.Args)
}

// SysCallExit is called each time a system call exit event happens.
func SysCallExit(t Task, s *SyscallEvent) string {
	i := defaultSyscallInfo(s.Sysno)
	if v, ok := syscalls[uintptr(s.Sysno)]; ok {
		*i = v
	}
	return i.printExit(t, s.Duration, s.Args, s.Ret[0], s.Errno)
}

// printExit prints the given system call exit.
func (i *SyscallInfo) printExit(t Task, elapsed time.Duration, args SyscallArguments, retval SyscallArgument, errno unix.Errno) string {
	// Eventually, we'll be able to cache o and look at the entry record's output.
	o := i.pre(t, args, LogMaximumSize)
	var rval string
	if errno == 0 {
		// Fill in the output after successful execution.
		i.post(t, args, retval, o, LogMaximumSize)
		rval = fmt.Sprintf("%#x (%v)", retval.Uint64(), elapsed)
	} else {
		rval = fmt.Sprintf("%s (%#x) (%v)", errno, errno, elapsed)
	}

	switch len(o) {
	case 0:
		return fmt.Sprintf("%s X %s() = %s", t.Name(), i.name, rval)
	case 1:
		return fmt.Sprintf("%s X %s(%s) = %s", t.Name(), i.name, o[0], rval)
	case 2:
		return fmt.Sprintf("%s X %s(%s, %s) = %s", t.Name(), i.name, o[0], o[1], rval)
	case 3:
		return fmt.Sprintf("%s X %s(%s, %s, %s) = %s", t.Name(), i.name, o[0], o[1], o[2], rval)
	case 4:
		return fmt.Sprintf("%s X %s(%s, %s, %s, %s) = %s", t.Name(), i.name, o[0], o[1], o[2], o[3], rval)
	case 5:
		return fmt.Sprintf("%s X %s(%s, %s, %s, %s, %s) = %s", t.Name(), i.name, o[0], o[1], o[2], o[3], o[4], rval)
	case 6:
		return fmt.Sprintf("%s X %s(%s, %s, %s, %s, %s, %s) = %s", t.Name(), i.name, o[0], o[1], o[2], o[3], o[4], o[5], rval)
	default:
		return fmt.Sprintf("%s X %s(%s, %s, %s, %s, %s, %s, ...) = %s", t.Name(), i.name, o[0], o[1], o[2], o[3], o[4], o[5], rval)
	}
}
