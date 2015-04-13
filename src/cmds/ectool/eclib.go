package main

import (
	"time"
)

type (
	Command int
	Version int
)

type ec interface {
	Probe(timeout time.Duration) error
	Cleanup(timeout time.Duration) error
	Wait(timeout time.Duration) error
	Command(command Command, version Version, idata []byte, outsize int, timeout time.Duration) ([]byte, error)
}
