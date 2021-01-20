// The MIT License (MIT)
//
// Copyright (c) 2018 The Genuinetools Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package container

import (
	"fmt"

	"github.com/opencontainers/runc/libcontainer/configs"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	rLimitCPU        = iota // CPU time in sec
	rLimitFsize             // Maximum filesize
	rLimitData              // max data size
	rLimitStack             // max stack size
	rLimitCore              // max core file size
	rLimitRss               // max resident set size
	rLimitNproc             // max number of processes
	rLimitNofile            // max number of open files
	rLimitMemlock           // max locked-in-memory address space
	rLimitAs                // address space limit
	rLimitLocks             // maximum file locks held
	rLimitSigpending        // max number of pending signals
	rLimitMsgqueue          // maximum bytes in POSIX mqueues
	rLimitNice              // max nice prio allowed to raise to
	rLimitRtprio            // maximum realtime priority
	rLimitRttime            // timeout for RT tasks in us
)

var rlimitMap = map[string]int{
	"RLIMIT_CPU":        rLimitCPU,
	"RLIMIT_FSIZE":      rLimitFsize,
	"RLIMIT_DATA":       rLimitData,
	"RLIMIT_STACK":      rLimitStack,
	"RLIMIT_CORE":       rLimitCore,
	"RLIMIT_RSS":        rLimitRss,
	"RLIMIT_NPROC":      rLimitNproc,
	"RLIMIT_NOFILE":     rLimitNofile,
	"RLIMIT_MEMLOCK":    rLimitMemlock,
	"RLIMIT_AS":         rLimitAs,
	"RLIMIT_LOCKS":      rLimitLocks,
	"RLIMIT_SIGPENDING": rLimitSigpending,
	"RLIMIT_MSGQUEUE":   rLimitMsgqueue,
	"RLIMIT_NICE":       rLimitNice,
	"RLIMIT_RTPRIO":     rLimitRtprio,
	"RLIMIT_RTTIME":     rLimitRttime,
}

func strToRlimit(key string) (int, error) {
	rl, ok := rlimitMap[key]
	if !ok {
		return 0, fmt.Errorf("wrong rlimit value: %s", key)
	}
	return rl, nil
}

func createLibContainerRlimit(rlimit specs.POSIXRlimit) (configs.Rlimit, error) {
	rl, err := strToRlimit(rlimit.Type)
	if err != nil {
		return configs.Rlimit{}, err
	}
	return configs.Rlimit{
		Type: rl,
		Hard: rlimit.Hard,
		Soft: rlimit.Soft,
	}, nil
}
