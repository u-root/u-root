// Copyright 2020-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ocp implements OCP/Facebook-specific IPMI client functions.
package ocp

import (
	"fmt"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/systembooter"
	"github.com/u-root/u-root/pkg/ipmi"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/vpd"
)

const (
	// Boot order bit[2:0] defined as
	// 000b: USB device
	// 001b: Network
	// 010b: SATA HDD
	// 011b: SATA-CDROM
	// 100b: Other removable Device
	// If bit[2:0] is 001b (Network), bit3 determines IPv4/IPv6 order,
	// when bit3 is 0: IPv4 first, bit3 is 1: IPv6 first
	NETWORK_BOOT      = 0x1
	NETWORK_BOOT_IPV6 = 0x1 | 0x8
	LOCAL_BOOT        = 0x2
	INVALID_BOOT      = 0xff

	// Default boot order systembooter configurations
	NETBOOTER_CONFIG      = "{\"type\":\"pxeboot\"}" // IPV4 and IPV6 default to true
	NETBOOTER_IPV6_CONFIG = "{\"type\":\"pxeboot\",\"ipv6\":\"true\",\"ipv4\":\"false\"}"
	LOCALBOOTER_CONFIG    = "{\"type\":\"boot\"}"
)

var (
	// BmcUpdatedBootorder is true when IPMI set boot order has been issued
	BmcUpdatedBootorder = false
	// BootEntries is created with the new boot order and will be used to boot this time
	BootEntries []systembooter.BootEntry
)

type BootOrder struct {
	bootMode byte
	bootSeq  [5]byte // Boot sequence, Boot0000, Boot0001...Boot0004
}

func getBootOrder(i *ipmi.IPMI, BootOrder *BootOrder) error {
	recv, err := i.SendRecv(_IPMI_FB_OEM_NET_FUNCTION1, _FB_OEM_GET_BIOS_BOOT_ORDER, nil)
	if err != nil {
		return err
	}
	BootOrder.bootMode = recv[1]
	copy(BootOrder.bootSeq[:], recv[2:])
	return nil
}

func setBootOrder(i *ipmi.IPMI, BootOrder *BootOrder) error {
	msg := ipmi.Msg{
		Netfn:   _IPMI_FB_OEM_NET_FUNCTION1,
		Cmd:     _FB_OEM_SET_BIOS_BOOT_ORDER,
		Data:    unsafe.Pointer(BootOrder),
		DataLen: 6,
	}

	if _, err := i.RawSendRecv(msg); err != nil {
		return err
	}
	return nil
}

// Currently we only support Network boot (pxeboot) and SATA HDD GRUB (boot),
// other types will be mapped to INVALID_BOOT and put to the end of the bootSeq.
func remapSortBootOrder(BootOrder *BootOrder) {
	sorted := [5]byte{INVALID_BOOT, INVALID_BOOT, INVALID_BOOT, INVALID_BOOT, INVALID_BOOT}
	var idx int
	for _, bootType := range BootOrder.bootSeq {
		if bootType == NETWORK_BOOT || bootType == NETWORK_BOOT_IPV6 || bootType == LOCAL_BOOT {
			sorted[idx] = bootType
			idx++
		}
	}
	copy(BootOrder.bootSeq[:], sorted[:])
}

func updateVPDBootOrder(i *ipmi.IPMI, BootOrder *BootOrder, l ulog.Logger) error {
	var err error
	var key string
	var idx int
	for _, bootType := range BootOrder.bootSeq {
		key = fmt.Sprintf("Boot%04d", idx)
		if bootType == NETWORK_BOOT {
			l.Printf("VPD set %s to %s", key, NETBOOTER_CONFIG)
			BootEntries = append(BootEntries, systembooter.BootEntry{Name: key, Config: []byte(NETBOOTER_CONFIG)})
			if err = vpd.FlashromRWVpdSet(key, []byte(NETBOOTER_CONFIG), false); err != nil {
				return err
			}
			idx++
		} else if bootType == NETWORK_BOOT_IPV6 {
			l.Printf("VPD set %s to %s", key, NETBOOTER_IPV6_CONFIG)
			BootEntries = append(BootEntries, systembooter.BootEntry{Name: key, Config: []byte(NETBOOTER_IPV6_CONFIG)})
			if err = vpd.FlashromRWVpdSet(key, []byte(NETBOOTER_IPV6_CONFIG), false); err != nil {
				return err
			}
			idx++
		} else if bootType == LOCAL_BOOT {
			l.Printf("VPD set %s to %s", key, LOCALBOOTER_CONFIG)
			BootEntries = append(BootEntries, systembooter.BootEntry{Name: key, Config: []byte(LOCALBOOTER_CONFIG)})
			if err = vpd.FlashromRWVpdSet(key, []byte(LOCALBOOTER_CONFIG), false); err != nil {
				return err
			}
			idx++
		} else if bootType == (INVALID_BOOT) {
			// No need to write VPD
		} else {
			l.Printf("Ignoring unrecognized boot type: %x", bootType)
		}
	}

	// Update the BootEntries with booters to match the new VPD
	for idx, entry := range BootEntries {
		entry.Booter, _ = systembooter.GetBooterFor(entry, l)
		BootEntries[idx] = entry
	}

	BmcUpdatedBootorder = true

	// clear valid bit
	BootOrder.bootMode &^= 0x80
	return setBootOrder(i, BootOrder)
}

// CheckBMCBootOrder synchronize BIOS's boot order with BMC's boot order.
// When BMC IPMI sets boot order (valid bit 1), BIOS will update VPD boot
// order and create new BootEntries accordingly. If BMC didn't set boot order,
// BIOS would set its current boot order to BMC.
func CheckBMCBootOrder(i *ipmi.IPMI, bmcBootOverride bool, l ulog.Logger) error {
	var BMCBootOrder, BIOSBootOrder BootOrder
	// Read boot order from BMC
	if err := getBootOrder(i, &BMCBootOrder); err != nil {
		return err
	}
	remapSortBootOrder(&BMCBootOrder)
	// If valid bit is set, IPMI set boot order has been issued so update
	// VPD boot order accordingly. For now the booter configurations will be
	// set to the default ones.
	if bmcBootOverride && BMCBootOrder.bootMode&0x80 != 0 {
		l.Printf("BMC set boot order valid bit is 1. Update VPD boot order.")
		return updateVPDBootOrder(i, &BMCBootOrder, l)
	}

	// Read VPD Boot entries and create BIOS BootOrder
	BIOSBootOrder.bootMode = BMCBootOrder.bootMode
	// Initialize with INVALID_BOOT
	idx := 0
	for idx = range BIOSBootOrder.bootSeq {
		BIOSBootOrder.bootSeq[idx] = INVALID_BOOT
	}
	bootEntries := systembooter.GetBootEntries(l)
	idx = 0
	var bootType string
	for _, entry := range bootEntries {
		if idx >= 5 {
			break
		}
		if bootType = entry.Booter.TypeName(); len(bootType) > 0 {
			if bootType == "pxeboot" {
				// Note: Does IPv6 only when BIOS sets its boot order to BMC, we could extend
				// Booter interface with a MethodName() to differentiate IPv4 and IPv6 in the future.
				BIOSBootOrder.bootSeq[idx] = NETWORK_BOOT_IPV6
				idx++
			} else if bootType == "boot" {
				BIOSBootOrder.bootSeq[idx] = LOCAL_BOOT
				idx++
			}
		}
	}
	// If there is no valid VPD boot order, write the default systembooter configurations
	if idx == 0 {
		l.Printf("No valid VPD boot order, set default boot orders to RW_VPD")
		BIOSBootOrder.bootSeq[0] = NETWORK_BOOT_IPV6
		BIOSBootOrder.bootSeq[1] = LOCAL_BOOT
		return updateVPDBootOrder(i, &BIOSBootOrder, l)
	}
	// clear valid bit and set BIOS boot order to BMC
	BIOSBootOrder.bootMode &^= 0x80
	return setBootOrder(i, &BIOSBootOrder)
}
