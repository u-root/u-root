package iodev

import (
	"fmt"
)

type FWDebug struct{}

func (f *FWDebug) Read(port uint64, data []byte) error {
	if len(data) == 1 {
		// This magic value is read from the Port to indicate the availability of the debug port.
		// See: https://github.com/cloud-hypervisor/edk2/blob/ch/OvmfPkg/Library/PlatformDebugLibIoPort/DebugIoPortQemu.c
		data[0] = 0xE9
	} else {
		return errDataLenInvalid
	}

	return nil
}

func (f *FWDebug) Write(port uint64, data []byte) error {
	if len(data) == 1 {
		if data[0] == '\000' {
			fmt.Printf("\r\n")
		} else {
			fmt.Printf("%c", data[0])
		}
	} else {
		return errDataLenInvalid
	}

	return nil
}

func (f *FWDebug) IOPort() uint64 {
	// https://github.com/tianocore/edk2/commit/bf23b44d926982dfc9ecc7785cea17e0889a9297
	return 0x402
}

func (f *FWDebug) Size() uint64 {
	return 0x1
}
