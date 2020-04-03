// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"unsafe"
)

const (
	_IPMI_FB_OEM_NET_FUNCTION1 = 0x30

	_FB_OEM_SET_BIOS_BOOT_ORDER = 0x52
	_FB_OEM_GET_BIOS_BOOT_ORDER = 0x53
)

// Get BIOS boot order data and check if CMOS clear bit and valid bit are both set
func (i *IPMI) IsCMOSClearSet() (bool, []byte, error) {
	req := &req{}
	req.msg.cmd = _FB_OEM_GET_BIOS_BOOT_ORDER
	req.msg.netfn = _IPMI_FB_OEM_NET_FUNCTION1

	recv, err := i.sendrecv(req)
	if err != nil {
		return false, nil, err
	}
	// recv[1] bit 1: CMOS clear, bit 7: valid bit, check if both are set
	if len(recv) > 6 && (recv[1]&0x82) == 0x82 {
		return true, recv[1:], nil
	}
	return false, nil, nil
}

// Set BIOS boot order with both CMOS clear and valid bits cleared
func (i *IPMI) ClearCMOSClearValidBits(data []byte) error {
	req := &req{}
	req.msg.cmd = _FB_OEM_SET_BIOS_BOOT_ORDER
	req.msg.netfn = _IPMI_FB_OEM_NET_FUNCTION1
	// Clear bit 1 and bit 7
	data[0] &= 0x7d
	req.msg.data = unsafe.Pointer(&data[0])
	req.msg.dataLen = 6

	if _, err := i.sendrecv(req); err != nil {
		return err
	}
	return nil
}
