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

package abi

import (
	"encoding/binary"
	"syscall"

	"golang.org/x/sys/unix"
)

// From <linux/futex.h> and <sys/time.h>.
// Flags are used in syscall futex(2).
const (
	FUTEX_WAIT            = 0
	FUTEX_WAKE            = 1
	FUTEX_FD              = 2
	FUTEX_REQUEUE         = 3
	FUTEX_CMP_REQUEUE     = 4
	FUTEX_WAKE_OP         = 5
	FUTEX_LOCK_PI         = 6
	FUTEX_UNLOCK_PI       = 7
	FUTEX_TRYLOCK_PI      = 8
	FUTEX_WAIT_BITSET     = 9
	FUTEX_WAKE_BITSET     = 10
	FUTEX_WAIT_REQUEUE_PI = 11
	FUTEX_CMP_REQUEUE_PI  = 12

	FUTEX_PRIVATE_FLAG   = 128
	FUTEX_CLOCK_REALTIME = 256
)

// These are flags are from <linux/futex.h> and are used in FUTEX_WAKE_OP
// to define the operations.
const (
	FUTEX_OP_SET         = 0
	FUTEX_OP_ADD         = 1
	FUTEX_OP_OR          = 2
	FUTEX_OP_ANDN        = 3
	FUTEX_OP_XOR         = 4
	FUTEX_OP_OPARG_SHIFT = 8
	FUTEX_OP_CMP_EQ      = 0
	FUTEX_OP_CMP_NE      = 1
	FUTEX_OP_CMP_LT      = 2
	FUTEX_OP_CMP_LE      = 3
	FUTEX_OP_CMP_GT      = 4
	FUTEX_OP_CMP_GE      = 5
)

// FUTEX_TID_MASK is the TID portion of a PI futex word.
const FUTEX_TID_MASK = 0x3fffffff

// ptrace commands from include/uapi/linux/ptrace.h.
const (
	PTRACE_TRACEME              = 0
	PTRACE_PEEKTEXT             = 1
	PTRACE_PEEKDATA             = 2
	PTRACE_PEEKUSR              = 3
	PTRACE_POKETEXT             = 4
	PTRACE_POKEDATA             = 5
	PTRACE_POKEUSR              = 6
	PTRACE_CONT                 = 7
	PTRACE_KILL                 = 8
	PTRACE_SINGLESTEP           = 9
	PTRACE_ATTACH               = 16
	PTRACE_DETACH               = 17
	PTRACE_SYSCALL              = 24
	PTRACE_SETOPTIONS           = 0x4200
	PTRACE_GETEVENTMSG          = 0x4201
	PTRACE_GETSIGINFO           = 0x4202
	PTRACE_SETSIGINFO           = 0x4203
	PTRACE_GETREGSET            = 0x4204
	PTRACE_SETREGSET            = 0x4205
	PTRACE_SEIZE                = 0x4206
	PTRACE_INTERRUPT            = 0x4207
	PTRACE_LISTEN               = 0x4208
	PTRACE_PEEKSIGINFO          = 0x4209
	PTRACE_GETSIGMASK           = 0x420a
	PTRACE_SETSIGMASK           = 0x420b
	PTRACE_SECCOMP_GET_FILTER   = 0x420c
	PTRACE_SECCOMP_GET_METADATA = 0x420d
)

// ptrace commands from arch/x86/include/uapi/asm/ptrace-abi.h.
const (
	PTRACE_GETREGS           = 12
	PTRACE_SETREGS           = 13
	PTRACE_GETFPREGS         = 14
	PTRACE_SETFPREGS         = 15
	PTRACE_GETFPXREGS        = 18
	PTRACE_SETFPXREGS        = 19
	PTRACE_OLDSETOPTIONS     = 21
	PTRACE_GET_THREAD_AREA   = 25
	PTRACE_SET_THREAD_AREA   = 26
	PTRACE_ARCH_PRCTL        = 30
	PTRACE_SYSEMU            = 31
	PTRACE_SYSEMU_SINGLESTEP = 32
	PTRACE_SINGLEBLOCK       = 33
)

// ptrace event codes from include/uapi/linux/ptrace.h.
const (
	PTRACE_EVENT_FORK       = 1
	PTRACE_EVENT_VFORK      = 2
	PTRACE_EVENT_CLONE      = 3
	PTRACE_EVENT_EXEC       = 4
	PTRACE_EVENT_VFORK_DONE = 5
	PTRACE_EVENT_EXIT       = 6
	PTRACE_EVENT_SECCOMP    = 7
	PTRACE_EVENT_STOP       = 128
)

// PTRACE_SETOPTIONS options from include/uapi/linux/ptrace.h.
const (
	PTRACE_O_TRACESYSGOOD    = 1
	PTRACE_O_TRACEFORK       = 1 << PTRACE_EVENT_FORK
	PTRACE_O_TRACEVFORK      = 1 << PTRACE_EVENT_VFORK
	PTRACE_O_TRACECLONE      = 1 << PTRACE_EVENT_CLONE
	PTRACE_O_TRACEEXEC       = 1 << PTRACE_EVENT_EXEC
	PTRACE_O_TRACEVFORKDONE  = 1 << PTRACE_EVENT_VFORK_DONE
	PTRACE_O_TRACEEXIT       = 1 << PTRACE_EVENT_EXIT
	PTRACE_O_TRACESECCOMP    = 1 << PTRACE_EVENT_SECCOMP
	PTRACE_O_EXITKILL        = 1 << 20
	PTRACE_O_SUSPEND_SECCOMP = 1 << 21
)

// from gvisor time.go
// Flags for clock_nanosleep(2).
const (
	TIMER_ABSTIME = 1
)

// Flags for timerfd syscalls (timerfd_create(2), timerfd_settime(2)).
const (
	// TFD_CLOEXEC is a timerfd_create flag.
	TFD_CLOEXEC = unix.O_CLOEXEC

	// TFD_NONBLOCK is a timerfd_create flag.
	TFD_NONBLOCK = unix.O_NONBLOCK

	// TFD_TIMER_ABSTIME is a timerfd_settime flag.
	TFD_TIMER_ABSTIME = 1
)

// TimeT represents time_t in <time.h>. It represents time in seconds.
type TimeT int64

// SizeOfTimeval is the size of a Timeval struct in bytes.
const SizeOfTimeval = 16

// ClockT represents type clock_t.
type ClockT int64

// Tms represents struct tms, used by times(2).
type Tms struct {
	UTime  ClockT
	STime  ClockT
	CUTime ClockT
	CSTime ClockT
}

// TimerID represents type timer_t, which identifies a POSIX per-process
// interval timer.
type TimerID int32

// ptrace

// PtraceRequestSet are the possible ptrace(2) requests.
var PtraceRequestSet = FlagSet{
	&Value{
		Value: PTRACE_TRACEME,
		Name:  "PTRACE_TRACEME",
	},
	&Value{
		Value: PTRACE_PEEKTEXT,
		Name:  "PTRACE_PEEKTEXT",
	},
	&Value{
		Value: PTRACE_PEEKDATA,
		Name:  "PTRACE_PEEKDATA",
	},
	&Value{
		Value: PTRACE_PEEKUSR,
		Name:  "PTRACE_PEEKUSR",
	},
	&Value{
		Value: PTRACE_POKETEXT,
		Name:  "PTRACE_POKETEXT",
	},
	&Value{
		Value: PTRACE_POKEDATA,
		Name:  "PTRACE_POKEDATA",
	},
	&Value{
		Value: PTRACE_POKEUSR,
		Name:  "PTRACE_POKEUSR",
	},
	&Value{
		Value: PTRACE_CONT,
		Name:  "PTRACE_CONT",
	},
	&Value{
		Value: PTRACE_KILL,
		Name:  "PTRACE_KILL",
	},
	&Value{
		Value: PTRACE_SINGLESTEP,
		Name:  "PTRACE_SINGLESTEP",
	},
	&Value{
		Value: PTRACE_ATTACH,
		Name:  "PTRACE_ATTACH",
	},
	&Value{
		Value: PTRACE_DETACH,
		Name:  "PTRACE_DETACH",
	},
	&Value{
		Value: PTRACE_SYSCALL,
		Name:  "PTRACE_SYSCALL",
	},
	&Value{
		Value: PTRACE_SETOPTIONS,
		Name:  "PTRACE_SETOPTIONS",
	},
	&Value{
		Value: PTRACE_GETEVENTMSG,
		Name:  "PTRACE_GETEVENTMSG",
	},
	&Value{
		Value: PTRACE_GETSIGINFO,
		Name:  "PTRACE_GETSIGINFO",
	},
	&Value{
		Value: PTRACE_SETSIGINFO,
		Name:  "PTRACE_SETSIGINFO",
	},
	&Value{
		Value: PTRACE_GETREGSET,
		Name:  "PTRACE_GETREGSET",
	},
	&Value{
		Value: PTRACE_SETREGSET,
		Name:  "PTRACE_SETREGSET",
	},
	&Value{
		Value: PTRACE_SEIZE,
		Name:  "PTRACE_SEIZE",
	},
	&Value{
		Value: PTRACE_INTERRUPT,
		Name:  "PTRACE_INTERRUPT",
	},
	&Value{
		Value: PTRACE_LISTEN,
		Name:  "PTRACE_LISTEN",
	},
	&Value{
		Value: PTRACE_PEEKSIGINFO,
		Name:  "PTRACE_PEEKSIGINFO",
	},
	&Value{
		Value: PTRACE_GETSIGMASK,
		Name:  "PTRACE_GETSIGMASK",
	},
	&Value{
		Value: PTRACE_SETSIGMASK,
		Name:  "PTRACE_SETSIGMASK",
	},
	&Value{
		Value: PTRACE_GETREGS,
		Name:  "PTRACE_GETREGS",
	},
	&Value{
		Value: PTRACE_SETREGS,
		Name:  "PTRACE_SETREGS",
	},
	&Value{
		Value: PTRACE_GETFPREGS,
		Name:  "PTRACE_GETFPREGS",
	},
	&Value{
		Value: PTRACE_SETFPREGS,
		Name:  "PTRACE_SETFPREGS",
	},
	&Value{
		Value: PTRACE_GETFPXREGS,
		Name:  "PTRACE_GETFPXREGS",
	},
	&Value{
		Value: PTRACE_SETFPXREGS,
		Name:  "PTRACE_SETFPXREGS",
	},
	&Value{
		Value: PTRACE_OLDSETOPTIONS,
		Name:  "PTRACE_OLDSETOPTIONS",
	},
	&Value{
		Value: PTRACE_GET_THREAD_AREA,
		Name:  "PTRACE_GET_THREAD_AREA",
	},
	&Value{
		Value: PTRACE_SET_THREAD_AREA,
		Name:  "PTRACE_SET_THREAD_AREA",
	},
	&Value{
		Value: PTRACE_ARCH_PRCTL,
		Name:  "PTRACE_ARCH_PRCTL",
	},
	&Value{
		Value: PTRACE_SYSEMU,
		Name:  "PTRACE_SYSEMU",
	},
	&Value{
		Value: PTRACE_SYSEMU_SINGLESTEP,
		Name:  "PTRACE_SYSEMU_SINGLESTEP",
	},
	&Value{
		Value: PTRACE_SINGLEBLOCK,
		Name:  "PTRACE_SINGLEBLOCK",
	},
}

// clone

// CloneFlagSet is the set of clone(2) flags.
var CloneFlagSet = FlagSet{
	&BitFlag{
		Value: syscall.CLONE_VM,
		Name:  "CLONE_VM",
	},
	&BitFlag{
		Value: syscall.CLONE_FS,
		Name:  "CLONE_FS",
	},
	&BitFlag{
		Value: syscall.CLONE_FILES,
		Name:  "CLONE_FILES",
	},
	&BitFlag{
		Value: syscall.CLONE_SIGHAND,
		Name:  "CLONE_SIGHAND",
	},
	&BitFlag{
		Value: syscall.CLONE_PTRACE,
		Name:  "CLONE_PTRACE",
	},
	&BitFlag{
		Value: syscall.CLONE_VFORK,
		Name:  "CLONE_VFORK",
	},
	&BitFlag{
		Value: syscall.CLONE_PARENT,
		Name:  "CLONE_PARENT",
	},
	&BitFlag{
		Value: syscall.CLONE_THREAD,
		Name:  "CLONE_THREAD",
	},
	&BitFlag{
		Value: syscall.CLONE_NEWNS,
		Name:  "CLONE_NEWNS",
	},
	&BitFlag{
		Value: syscall.CLONE_SYSVSEM,
		Name:  "CLONE_SYSVSEM",
	},
	&BitFlag{
		Value: syscall.CLONE_SETTLS,
		Name:  "CLONE_SETTLS",
	},
	&BitFlag{
		Value: syscall.CLONE_PARENT_SETTID,
		Name:  "CLONE_PARENT_SETTID",
	},
	&BitFlag{
		Value: syscall.CLONE_CHILD_CLEARTID,
		Name:  "CLONE_CHILD_CLEARTID",
	},
	&BitFlag{
		Value: syscall.CLONE_DETACHED,
		Name:  "CLONE_DETACHED",
	},
	&BitFlag{
		Value: syscall.CLONE_UNTRACED,
		Name:  "CLONE_UNTRACED",
	},
	&BitFlag{
		Value: syscall.CLONE_CHILD_SETTID,
		Name:  "CLONE_CHILD_SETTID",
	},
	&BitFlag{
		Value: syscall.CLONE_NEWUTS,
		Name:  "CLONE_NEWUTS",
	},
	&BitFlag{
		Value: syscall.CLONE_NEWIPC,
		Name:  "CLONE_NEWIPC",
	},
	&BitFlag{
		Value: syscall.CLONE_NEWUSER,
		Name:  "CLONE_NEWUSER",
	},
	&BitFlag{
		Value: syscall.CLONE_NEWPID,
		Name:  "CLONE_NEWPID",
	},
	&BitFlag{
		Value: syscall.CLONE_NEWNET,
		Name:  "CLONE_NEWNET",
	},
	&BitFlag{
		Value: syscall.CLONE_IO,
		Name:  "CLONE_IO",
	},
}

// Socket defines. Some of these might move to abi_unix.go
// Address families, from linux/socket.h.
const (
	AF_UNSPEC     = 0
	AF_UNIX       = 1
	AF_INET       = 2
	AF_AX25       = 3
	AF_IPX        = 4
	AF_APPLETALK  = 5
	AF_NETROM     = 6
	AF_BRIDGE     = 7
	AF_ATMPVC     = 8
	AF_X25        = 9
	AF_INET6      = 10
	AF_ROSE       = 11
	AF_DECnet     = 12
	AF_NETBEUI    = 13
	AF_SECURITY   = 14
	AF_KEY        = 15
	AF_NETLINK    = 16
	AF_PACKET     = 17
	AF_ASH        = 18
	AF_ECONET     = 19
	AF_ATMSVC     = 20
	AF_RDS        = 21
	AF_SNA        = 22
	AF_IRDA       = 23
	AF_PPPOX      = 24
	AF_WANPIPE    = 25
	AF_LLC        = 26
	AF_IB         = 27
	AF_MPLS       = 28
	AF_CAN        = 29
	AF_TIPC       = 30
	AF_BLUETOOTH  = 31
	AF_IUCV       = 32
	AF_RXRPC      = 33
	AF_ISDN       = 34
	AF_PHONET     = 35
	AF_IEEE802154 = 36
	AF_CAIF       = 37
	AF_ALG        = 38
	AF_NFC        = 39
	AF_VSOCK      = 40
)

// sendmsg(2)/recvmsg(2) flags, from linux/socket.h.
const (
	MSG_OOB              = 0x1
	MSG_PEEK             = 0x2
	MSG_DONTROUTE        = 0x4
	MSG_TRYHARD          = 0x4
	MSG_CTRUNC           = 0x8
	MSG_PROBE            = 0x10
	MSG_TRUNC            = 0x20
	MSG_DONTWAIT         = 0x40
	MSG_EOR              = 0x80
	MSG_WAITALL          = 0x100
	MSG_FIN              = 0x200
	MSG_EOF              = MSG_FIN
	MSG_SYN              = 0x400
	MSG_CONFIRM          = 0x800
	MSG_RST              = 0x1000
	MSG_ERRQUEUE         = 0x2000
	MSG_NOSIGNAL         = 0x4000
	MSG_MORE             = 0x8000
	MSG_WAITFORONE       = 0x10000
	MSG_SENDPAGE_NOTLAST = 0x20000
	MSG_REINJECT         = 0x8000000
	MSG_ZEROCOPY         = 0x4000000
	MSG_FASTOPEN         = 0x20000000
	MSG_CMSG_CLOEXEC     = 0x40000000
)

// SOL_SOCKET is from socket.h
const SOL_SOCKET = 1

// Socket types, from linux/net.h.
const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_RDM       = 4
	SOCK_SEQPACKET = 5
	SOCK_DCCP      = 6
	SOCK_PACKET    = 10
)

// SOCK_TYPE_MASK covers all of the above socket types. The remaining bits are
// flags. From linux/net.h.
const SOCK_TYPE_MASK = 0xf

// socket(2)/socketpair(2)/accept4(2) flags, from linux/net.h.
const (
	SOCK_CLOEXEC  = unix.O_CLOEXEC
	SOCK_NONBLOCK = unix.O_NONBLOCK
)

// shutdown(2) how commands, from <linux/net.h>.
const (
	SHUT_RD   = 0
	SHUT_WR   = 1
	SHUT_RDWR = 2
)

// Socket options from socket.h.
const (
	SO_ERROR       = 4
	SO_KEEPALIVE   = 9
	SO_LINGER      = 13
	SO_MARK        = 36
	SO_PASSCRED    = 16
	SO_PEERCRED    = 17
	SO_PEERNAME    = 28
	SO_PROTOCOL    = 38
	SO_RCVBUF      = 8
	SO_RCVTIMEO    = 20
	SO_REUSEADDR   = 2
	SO_SNDBUF      = 7
	SO_SNDTIMEO    = 21
	SO_TIMESTAMP   = 29
	SO_TIMESTAMPNS = 35
	SO_TYPE        = 3
)

// SockAddrMax is the maximum size of a struct sockaddr, from
// uapi/linux/socket.h.
const SockAddrMax = 128

// SockAddrInt is struct sockaddr_in, from uapi/linux/in.h.
type SockAddrInet struct {
	Family uint16
	Port   uint16
	Addr   [4]byte
	Zero   [8]uint8 // pad to sizeof(struct sockaddr).
}

// SockAddrInt6 is struct sockaddr_in6, from uapi/linux/in6.h.
type SockAddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte
	Scope_id uint32
}

// UnixPathMax is the maximum length of the path in an AF_UNIX socket.
//
// From uapi/linux/un.h.
const UnixPathMax = 108

// SockAddrUnix is struct sockaddr_un, from uapi/linux/un.h.
type SockAddrUnix struct {
	Family uint16
	Path   [UnixPathMax]int8
}

// TCPInfo is a collection of TCP statistics.
//
// From uapi/linux/tcp.h.
type TCPInfo struct {
	State       uint8
	CaState     uint8
	Retransmits uint8
	Probes      uint8
	Backoff     uint8
	Options     uint8
	// WindowScale is the combination of snd_wscale (first 4 bits) and rcv_wscale (second 4 bits)
	WindowScale uint8
	// DeliveryRateAppLimited is a boolean and only the first bit is meaningful.
	DeliveryRateAppLimited uint8

	RTO    uint32
	ATO    uint32
	SndMss uint32
	RcvMss uint32

	Unacked uint32
	Sacked  uint32
	Lost    uint32
	Retrans uint32
	Fackets uint32

	// Times.
	LastDataSent uint32
	LastAckSent  uint32
	LastDataRecv uint32
	LastAckRecv  uint32

	// Metrics.
	PMTU        uint32
	RcvSsthresh uint32
	RTT         uint32
	RTTVar      uint32
	SndSsthresh uint32
	SndCwnd     uint32
	Advmss      uint32
	Reordering  uint32

	RcvRTT   uint32
	RcvSpace uint32

	TotalRetrans uint32

	PacingRate    uint64
	MaxPacingRate uint64
	// BytesAcked is RFC4898 tcpEStatsAppHCThruOctetsAcked.
	BytesAcked uint64
	// BytesReceived is RFC4898 tcpEStatsAppHCThruOctetsReceived.
	BytesReceived uint64
	// SegsOut is RFC4898 tcpEStatsPerfSegsOut.
	SegsOut uint32
	// SegsIn is RFC4898 tcpEStatsPerfSegsIn.
	SegsIn uint32

	NotSentBytes uint32
	MinRTT       uint32
	// DataSegsIn is RFC4898 tcpEStatsDataSegsIn.
	DataSegsIn uint32
	// DataSegsOut is RFC4898 tcpEStatsDataSegsOut.
	DataSegsOut uint32

	DeliveryRate uint64

	// BusyTime is the time in microseconds busy sending data.
	BusyTime uint64
	// RwndLimited is the time in microseconds limited by receive window.
	RwndLimited uint64
	// SndBufLimited is the time in microseconds limited by send buffer.
	SndBufLimited uint64
}

// SizeOfTCPInfo is the binary size of a TCPInfo struct (104 bytes).
var SizeOfTCPInfo = binary.Size(TCPInfo{})

// Control message types, from linux/socket.h.
const (
	SCM_CREDENTIALS = 0x2
	SCM_RIGHTS      = 0x1
)

// A ControlMessageHeader is the header for a socket control message.
//
// ControlMessageHeader represents struct cmsghdr from linux/socket.h.
type ControlMessageHeader struct {
	Length uint64
	Level  int32
	Type   int32
}

// SizeOfControlMessageHeader is the binary size of a ControlMessageHeader
// struct.
var SizeOfControlMessageHeader = int(binary.Size(ControlMessageHeader{}))

// A ControlMessageCredentials is an SCM_CREDENTIALS socket control message.
//
// ControlMessageCredentials represents struct ucred from linux/socket.h.
type ControlMessageCredentials struct {
	PID int32
	UID uint32
	GID uint32
}

// SizeOfControlMessageCredentials is the binary size of a
// ControlMessageCredentials struct.
var SizeOfControlMessageCredentials = int(binary.Size(ControlMessageCredentials{}))

// A ControlMessageRights is an SCM_RIGHTS socket control message.
type ControlMessageRights []int32

// SizeOfControlMessageRight is the size of a single element in
// ControlMessageRights.
const SizeOfControlMessageRight = 4

// SCM_MAX_FD is the maximum number of FDs accepted in a single sendmsg call.
// From net/scm.h.
const SCM_MAX_FD = 253

// itimer

// itimer types for getitimer(2) and setitimer(2), from
// include/uapi/linux/time.h.
const (
	ITIMER_REAL    = 0
	ITIMER_VIRTUAL = 1
	ITIMER_PROF    = 2
)

// ItimerTypes are the possible itimer types.
var ItimerTypes = FlagSet{
	&Value{
		Value: ITIMER_REAL,
		Name:  "ITIMER_REAL",
	},
	&Value{
		Value: ITIMER_VIRTUAL,
		Name:  "ITIMER_VIRTUAL",
	},
	&Value{
		Value: ITIMER_PROF,
		Name:  "ITIMER_PROF",
	},
}

// FutexCmd are the possible futex(2) commands.
//
// from gvisor futex.go
var FutexCmd = FlagSet{
	&Value{
		Value: FUTEX_WAIT,
		Name:  "FUTEX_WAIT",
	},
	&Value{
		Value: FUTEX_WAKE,
		Name:  "FUTEX_WAKE",
	},
	&Value{
		Value: FUTEX_FD,
		Name:  "FUTEX_FD",
	},
	&Value{
		Value: FUTEX_REQUEUE,
		Name:  "FUTEX_REQUEUE",
	},
	&Value{
		Value: FUTEX_CMP_REQUEUE,
		Name:  "FUTEX_CMP_REQUEUE",
	},
	&Value{
		Value: FUTEX_WAKE_OP,
		Name:  "FUTEX_WAKE_OP",
	},
	&Value{
		Value: FUTEX_LOCK_PI,
		Name:  "FUTEX_LOCK_PI",
	},
	&Value{
		Value: FUTEX_UNLOCK_PI,
		Name:  "FUTEX_UNLOCK_PI",
	},
	&Value{
		Value: FUTEX_TRYLOCK_PI,
		Name:  "FUTEX_TRYLOCK_PI",
	},
	&Value{
		Value: FUTEX_WAIT_BITSET,
		Name:  "FUTEX_WAIT_BITSET",
	},
	&Value{
		Value: FUTEX_WAKE_BITSET,
		Name:  "FUTEX_WAKE_BITSET",
	},
	&Value{
		Value: FUTEX_WAIT_REQUEUE_PI,
		Name:  "FUTEX_WAIT_REQUEUE_PI",
	},
	&Value{
		Value: FUTEX_CMP_REQUEUE_PI,
		Name:  "FUTEX_CMP_REQUEUE_PI",
	},
}

func Futex(op uint64) string {
	cmd := op &^ (FUTEX_PRIVATE_FLAG | FUTEX_CLOCK_REALTIME)
	clockRealtime := (op & FUTEX_CLOCK_REALTIME) == FUTEX_CLOCK_REALTIME
	private := (op & FUTEX_PRIVATE_FLAG) == FUTEX_PRIVATE_FLAG

	s := FutexCmd.Parse(cmd)
	if clockRealtime {
		s += "|FUTEX_CLOCK_REALTIME"
	}
	if private {
		s += "|FUTEX_PRIVATE_FLAG"
	}
	return s
}
