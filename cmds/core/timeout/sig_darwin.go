// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import "syscall"

var sigmap = map[string]syscall.Signal{
	"HUP":    syscall.SIGHUP,
	"INT":    syscall.SIGINT,
	"QUIT":   syscall.SIGQUIT,
	"ILL":    syscall.SIGILL,
	"TRAP":   syscall.SIGTRAP,
	"ABRT":   syscall.SIGABRT,
	"EMT":    syscall.SIGEMT,
	"FPE":    syscall.SIGFPE,
	"KILL":   syscall.SIGKILL,
	"BUS":    syscall.SIGBUS,
	"SEGV":   syscall.SIGSEGV,
	"SYS":    syscall.SIGSYS,
	"PIPE":   syscall.SIGPIPE,
	"ALRM":   syscall.SIGALRM,
	"TERM":   syscall.SIGTERM,
	"URG":    syscall.SIGURG,
	"STOP":   syscall.SIGSTOP,
	"TSTP":   syscall.SIGTSTP,
	"CONT":   syscall.SIGCONT,
	"CHLD":   syscall.SIGCHLD,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"IO":     syscall.SIGIO,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
	"VTALRM": syscall.SIGVTALRM,
	"PROF":   syscall.SIGPROF,
	"WINCH":  syscall.SIGWINCH,
	"INFO":   syscall.SIGINFO,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
}
