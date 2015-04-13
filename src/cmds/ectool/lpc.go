package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	initial_udelay time.Duration = time.Microsecond
	maximum_udelay time.Duration = 20 * time.Microsecond
)

type lpc struct {
	ioport
	status_addr  ioaddr
	status       ioaddr
	cmd          ioaddr
	initial, max time.Duration
	debugf
}

func init() {
	chips["lpc"] = NewLPC
}

func NewLPC(p ioport, sa ioaddr, i, m time.Duration, d debugf) ec {
	return &lpc{ioport: p, status_addr: sa, initial: i, max: m, debugf: d}
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
		v, err := l.Inb(l.status_addr)
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

func (l *lpc) Command(c Command, v Version, idata []byte, outsize int, timeout time.Duration) ([]byte, error) {
	flags := uint8(EC_HOST_ARGS_FLAG_FROM_HOST)
	csum := flags + uint8(c) + uint8(v) + uint8(len(idata))

	for i := range idata {
		err := l.Outb(l.status_addr+ioaddr(i), idata[i])
		if err != nil {
			return nil, err
		}
		csum += idata[i]
	}

	cmd := []uint8{EC_HOST_ARGS_FLAG_FROM_HOST, uint8(v), uint8(len(idata)), uint8(csum)}
	_, err := l.Outs(l.status_addr, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return nil, err
	}

	l.Outb(l.cmd, uint8(c))
	l.Outb(l.status, 0xa5)

	if err := l.Wait(10 * time.Second); err != nil {
		return nil, errors.New("Timeout waiting for EC response")
	}

	/* Check result */
	i, err := l.Inb(EC_LPC_ADDR_HOST_DATA)
	if err != nil {
		return nil, err
	}

	if i != 0 {
		return nil, fmt.Errorf("Bad EC value %d", i)
	}

	/* Read back args */
	dres, err := l.Ins(EC_LPC_ADDR_HOST_ARGS, 4)
	if err != nil {
		return nil, err
	}

	/*
	 * If EC didn't modify args flags, then somehow we sent a new-style
	 * command to an old EC, which means it would have read its params
	 * from the wrong place.
	 */
	if flags&EC_HOST_ARGS_FLAG_TO_HOST == EC_HOST_ARGS_FLAG_TO_HOST {
		return nil, errors.New("EC appears to have reset (may be expected)")
	}

	if dres[2] > uint8(outsize) {
		return nil, errors.New("EC returned too much data")
	}

	/* Start calculating response checksum */
	csum = uint8(c) + dres[0] + uint8(v) + dres[2]

	/* Read response and update checksum */
	dout, err := l.Ins(EC_LPC_ADDR_HOST_PARAM, int(dres[2]))
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
