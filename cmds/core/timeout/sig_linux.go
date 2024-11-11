// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import "syscall"

var sigmap = map[string]syscall.Signal{
	"ABRT":   syscall.SIGABRT,
	"ALRM":   syscall.SIGALRM,
	"BUS":    syscall.SIGBUS,
	"CHLD":   syscall.SIGCHLD,
	"CLD":    syscall.SIGCHLD,
	"CONT":   syscall.SIGCONT,
	"FPE":    syscall.SIGFPE,
	"HUP":    syscall.SIGHUP,
	"ILL":    syscall.SIGILL,
	"INFO":   syscall.SIGPWR,
	"INT":    syscall.SIGINT,
	"IO":     syscall.SIGIO,
	"IOT":    syscall.SIGABRT,
	"KILL":   syscall.SIGKILL,
	"PIPE":   syscall.SIGPIPE,
	"POLL":   syscall.SIGIO,
	"PROF":   syscall.SIGPROF,
	"PWR":    syscall.SIGPWR,
	"QUIT":   syscall.SIGQUIT,
	"SEGV":   syscall.SIGSEGV,
	"STOP":   syscall.SIGSTOP,
	"TSTP":   syscall.SIGTSTP,
	"SYS":    syscall.SIGSYS,
	"TERM":   syscall.SIGTERM,
	"TRAP":   syscall.SIGTRAP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"URG":    syscall.SIGURG,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"VTALRM": syscall.SIGVTALRM,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
	"WINCH":  syscall.SIGWINCH,
}
