// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"golang.org/x/sys/unix"
)

type Socket interface {
	SocketsString(bool, bool, *Output) (string, error)
	readData() error
}

func NewSocket(t Protocol) (Socket, error) {
	var ret Socket
	switch t {
	case PROT_TCP, PROT_TCP6, PROT_UDP, PROT_UDP6, PROT_UDPL, PROT_UDPL6, PROT_RAW, PROT_RAW6:
		ret = &NetSockets{
			Protocol: t,
		}
	case PROT_UNIX:
		ret = &UnixSockets{}
	}

	if err := ret.readData(); err != nil {
		return nil, err
	}
	return ret, nil
}

type NetSockets struct {
	Protocol
	Entries []netSocket
}

type netSocket struct {
	Protocol
	Index       uint64
	LocalAddr   IPAddress
	ForeignAddr IPAddress
	State       NetState
	RxQueue     uint64
	TxQueue     uint64
	TimerRun    uint8
	TimerLen    uint64
	Retr        uint64
	UID         uint32
	Timeout     uint64
	Inode       uint64
	Ref         uint64
}

func (n *NetSockets) readData() error {
	path := path.Join(ProcnetPath, n.Protocol.String())

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// We scan to get rid of the title line
	scanner.Scan()

	entries := make([]netSocket, 0)
	var localAddr, remAddr string

	for scanner.Scan() {
		line := scanner.Text()
		entry := netSocket{}
		_, err := fmt.Sscanf(line, "%d: %s %s %X %X:%X %d:%X %X %d %d %d %d",
			&entry.Index,
			&localAddr,
			&remAddr,
			&entry.State,
			&entry.RxQueue,
			&entry.TxQueue,
			&entry.TimerRun,
			&entry.TimerLen,
			&entry.Retr,
			&entry.UID,
			&entry.Timeout,
			&entry.Inode,
			&entry.Ref,
		)
		if err != nil {
			return err
		}

		entry.LocalAddr, err = newIPAddress(localAddr)
		if err != nil {
			return err
		}

		entry.ForeignAddr, err = newIPAddress(remAddr)
		if err != nil {
			return err
		}

		entry.Protocol = n.Protocol

		entries = append(entries, entry)
	}

	n.Entries = entries

	return nil
}

func (n *NetSockets) SocketsString(lst, all bool, outputfmt *Output) (string, error) {
	states := []NetState{TCP_ESTABLISHED, TCP_TIME_WAIT}

	if lst {
		states = []NetState{TCP_LISTEN}
	}

	outputfmt.InitIPSocketTitels()

	if all {
		states = []NetState{
			TCP_ESTABLISHED,
			TCP_SYN_SENT,
			TCP_SYN_RECV,
			TCP_FIN_WAIT1,
			TCP_FIN_WAIT2,
			TCP_TIME_WAIT,
			TCP_CLOSE,
			TCP_CLOSE_WAIT,
			TCP_LAST_ACK,
			TCP_LISTEN,
			TCP_CLOSING,
		}
	}

	// Table header
	for _, entry := range n.Entries {
		for _, state := range states {
			if entry.State.String() == state.String() {
				outputfmt.AddIPSocket(entry)
			}
		}
	}

	return outputfmt.String(), nil
}

type UnixSockets struct {
	Entries []unixSocket
}

type unixSocket struct {
	RefCnt uint32
	Proto  uint32
	Flags  uint32
	Type   SockType
	St     SockState
	Inode  uint32
	Path   string
}

type SockType int

func (s *SockType) String() string {
	var str strings.Builder

	switch *s {
	case unix.SOCK_STREAM:
		str.WriteString("STREAM")
	case unix.SOCK_DGRAM:
		str.WriteString("DGRAM")
	case unix.SOCK_RAW:
		str.WriteString("RAW")
	case unix.SOCK_RDM:
		str.WriteString("RDM")
	case unix.SOCK_SEQPACKET:
		str.WriteString("SEQPACKET")
	default:
		str.WriteString("UNKNOWN")
	}

	return str.String()
}

func (u *UnixSockets) SocketsString(lst, all bool, outputfmt *Output) (string, error) {
	states := []SockState{SSCONNECTED}

	if lst {
		states = []SockState{SSUNCONNECTED}
	}

	if all {
		states = []SockState{
			SSFREE,
			SSUNCONNECTED,
			SSCONNECTING,
			SSCONNECTED,
			SSDISCONNECTING,
		}
	}

	outputfmt.InitUnixSocketTitels()

	for _, entry := range u.Entries {
		for _, state := range states {
			// Parse State with 0 just for the state before printing the information.
			if entry.St.parseState(0) == state.parseState(0) {
				if lst && entry.Flags&SSACCEPTCON == 0 {
					continue
				}
				outputfmt.AddUnixSocket(entry)
			}
		}
	}

	return outputfmt.String(), nil
}

func (u *UnixSockets) readData() error {
	path := path.Join(ProcnetPath, "unix")

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	// We scan to get rid of the title line
	s.Scan()

	u.Entries = make([]unixSocket, 0)

	for s.Scan() {
		e := unixSocket{}
		line := s.Text()
		// Per the source: the format is %pK: %08X %08X %08X %04X %02X %5lu
		// e.g. '0000000000000000: 0000000B 00000000 00000000 0002 03  3238 /run/systemd/journal/socket'
		// The original code assumed some fields were decimal. Be careful!
		const fmtString = "%s %X %X %X %X %X %d %s"
		var dummy string
		if _, err := fmt.Sscanf(line, fmtString, &dummy, &e.RefCnt, &e.Proto, &e.Flags, &e.Type, &e.St, &e.Inode, &e.Path); err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("converting %q with %q:%w", line, fmtString, err)
		}

		u.Entries = append(u.Entries, e)
	}

	return nil
}

const (
	SSACCEPTCON = (1 << 16)
	SSWAITDATA  = (1 << 17)
	SSNOSPACE   = (1 << 18)
)

func parseUnixFlags(flags uint32) string {
	var s strings.Builder

	s.WriteString("[")
	if flags&SSACCEPTCON > 0 {
		s.WriteString("ACC")
	}

	if flags&SSWAITDATA > 0 {
		s.WriteString("W")
	}

	if flags&SSNOSPACE > 0 {
		s.WriteString("N")
	}
	s.WriteString("]")

	return s.String()
}
