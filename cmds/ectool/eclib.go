package main

import (
	"time"
)

type (
	command int
	version int
)

type ec interface {
	Probe(timeout time.Duration) error
	Cleanup(timeout time.Duration) error
	Wait(timeout time.Duration) error
	Command(command command, version version, idata []byte, outsize int, timeout time.Duration) ([]byte, error)
}
