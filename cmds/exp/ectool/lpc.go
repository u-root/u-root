package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	initialUdelay time.Duration = time.Microsecond
	maximumUdelay time.Duration = 20 * time.Microsecond
)

type lpc struct {
	ioport
	statusAddr   ioaddr
	status       ioaddr
	cmd          ioaddr
	initial, max time.Duration
	debugf
}

func init() {
	chips["lpc"] = newLPC
}

func newLPC(p ioport, sa ioaddr, i, m time.Duration, d debugf) ec {
	return &lpc{ioport: p, statusAddr: sa, initial: i, max: m, debugf: d}
}

func (l *lpc) Wait(timeout time.Duration) error {
	delay := l.initial

	for i := time.Duration(0); i < timeout; i += l.initial {
		// TODO: kill this static timeout and use chans.
		// But for now we clone what the C code did.
		if delay > l.max {
			delay = l.max
		}
		/*
		 * Delay first, in case we just sent out a command but the EC
		 * hasn't raised the busy flag.  However, I think this doesn't
		 * happen since the LPC commands are executed in order and the
		 * busy flag is set by hardware.  Minor issue in any case,
		 * since the initial delay is very short.
		 */
		time.Sleep(delay)
		v, err := l.Inb(l.statusAddr)
		if err != nil {
			return err
		}
		if v == 0 {
			return nil
		}
		if i > 20 && delay < l.max {
			delay *= 2
		}
	}
	return fmt.Errorf("LPC timed out")
}

func (l *lpc) Cleanup(timeout time.Duration) error { return nil }
func (l *lpc) Probe(timeout time.Duration) error   { return nil }

func (l *lpc) Command(c command, v version, idata []byte, outsize int, timeout time.Duration) ([]byte, error) {
	flags := uint8(ecHostArgsFlagFromHost)
	csum := flags + uint8(c) + uint8(v) + uint8(len(idata))

	for i, d := range idata {
		err := l.Outb(l.statusAddr+ioaddr(i), d)
		if err != nil {
			return nil, err
		}
		csum += d
	}

	cmd := []uint8{ecHostArgsFlagFromHost, uint8(v), uint8(len(idata)), uint8(csum)}
	_, err := l.Outs(l.statusAddr, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return nil, err
	}

	l.Outb(l.cmd, uint8(c))
	l.Outb(l.status, 0xa5)

	if err := l.Wait(10 * time.Second); err != nil {
		return nil, errors.New("timeout waiting for EC response")
	}

	/* Check result */
	i, err := l.Inb(ecLpcAddrHostData)
	if err != nil {
		return nil, err
	}

	if i != 0 {
		return nil, fmt.Errorf("bad EC value %d", i)
	}

	/* Read back args */
	dres, err := l.Ins(ecLpcAddrHostArgs, 4)
	if err != nil {
		return nil, err
	}

	/*
	 * If EC didn't modify args flags, then somehow we sent a new-style
	 * command to an old EC, which means it would have read its params
	 * from the wrong place.
	 */
	if flags&ecHostArgsFlagToHost == ecHostArgsFlagToHost {
		return nil, errors.New("EC appears to have reset (may be expected)")
	}

	if dres[2] > uint8(outsize) {
		return nil, errors.New("EC returned too much data")
	}

	/* Start calculating response checksum */
	csum = uint8(c) + dres[0] + uint8(v) + dres[2]

	/* Read response and update checksum */
	dout, err := l.Ins(ecLpcAddrHostParam, int(dres[2]))
	if err != nil {
		return nil, err
	}

	for _, i := range dout {
		csum += i
	}
	/* Verify checksum */
	if dout[3] != csum {
		return nil, errors.New("EC response has invalid checksum")
	}

	return dout, nil
}
