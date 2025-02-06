// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64

package universalpayload

func getPhysicalAddressSizes() (uint8, error) {
	return 0, nil
}

func constructTrampoline(buf []uint8, hobAddr uint64, entry uint64) []uint8 {
	return nil
}

func archGetAcpiRsdpData() (uint64, []byte, error) {
	return 0xDEADBEEF, nil, nil
}

func appendAddonMemMap(_ *EFIMemoryMapHOB) uint64 {
	return 0
}
