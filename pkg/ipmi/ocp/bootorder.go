// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ocp implements OCP/Facebook-specific IPMI client functions.
package ocp

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/systembooter"
	"github.com/u-root/u-root/pkg/ipmi"
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
	NETWORK_BOOT = 0x1
	LOCAL_BOOT   = 0x2
	INVALID_BOOT = 0xff

	// Default boot order systembooter configurations
	NETBOOTER_CONFIG   = "{\"type\":\"netboot\",\"method\":\"dhcpv6\"}"
	LOCALBOOTER_CONFIG = "{\"type\":\"localboot\",\"method\":\"grub\"}"
)

var (
	// BmcUpdatedBootorder is true when IPMI set boot order has been issued
	BmcUpdatedBootorder bool = false
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

// Currently we only support IPv6 Network boot (netboot) and SATA HDD GRUB (localboot),
// other types will be mapped to INVALID_BOOT and put to the end of the bootSeq.
func remapSortBootOrder(BootOrder *BootOrder) {
	var bootType byte
	sorted := [5]byte{INVALID_BOOT, INVALID_BOOT, INVALID_BOOT, INVALID_BOOT, INVALID_BOOT}
	var idx int
	for _, v := range BootOrder.bootSeq {
		bootType = v & 0x7
		if bootType == NETWORK_BOOT {
			if v&0x8 == 0 { // If IPv6 first bit is not set, set it
				v |= 0x8
			}
			sorted[idx] = v
			idx++
		} else if bootType == LOCAL_BOOT {
			sorted[idx] = v
			idx++
		}
	}
	copy(BootOrder.bootSeq[:], sorted[:])
}

func updateVPDBootOrder(i *ipmi.IPMI, BootOrder *BootOrder) error {
	var err error
	var key string
	var bootType byte
	var idx int
	for _, v := range BootOrder.bootSeq {
		key = fmt.Sprintf("Boot%04d", idx)
		bootType = v & 0x7
		if bootType == NETWORK_BOOT {
			log.Printf("VPD set %s to %s", key, NETBOOTER_CONFIG)
			BootEntries = append(BootEntries, systembooter.BootEntry{Name: key, Config: []byte(NETBOOTER_CONFIG)})
			if err = vpd.FlashromRWVpdSet(key, []byte(NETBOOTER_CONFIG), false); err != nil {
				return err
			}
			idx++
		} else if bootType == LOCAL_BOOT {
			log.Printf("VPD set %s to %s", key, LOCALBOOTER_CONFIG)
			BootEntries = append(BootEntries, systembooter.BootEntry{Name: key, Config: []byte(LOCALBOOTER_CONFIG)})
			if err = vpd.FlashromRWVpdSet(key, []byte(LOCALBOOTER_CONFIG), false); err != nil {
				return err
			}
			idx++
		} else if bootType == (INVALID_BOOT & 0x7) {
			// No need to write VPD
		} else {
			log.Printf("Ignoring unrecognized boot type: %x", bootType)
		}
	}

	if BmcUpdatedBootorder {
		// look for a Booter that supports the given configuration
		for idx, entry := range BootEntries {
			entry.Booter = systembooter.GetBooterFor(entry)
			BootEntries[idx] = entry
		}
	}
	// clear valid bit
	BootOrder.bootMode &^= 0x80
	return setBootOrder(i, BootOrder)
}

// CheckBMCBootOrder synchronize BIOS's boot order with BMC's boot order.
// When BMC IPMI sets boot order (valid bit 1), BIOS will update VPD boot
// order and create new BootEntries accordingly. If BMC didn't set boot order,
// BIOS would set its current boot order to BMC.
func CheckBMCBootOrder(i *ipmi.IPMI, bmcBootOverride bool) error {
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
		log.Printf("BMC set boot order valid bit is 1")
		BmcUpdatedBootorder = true
		return updateVPDBootOrder(i, &BMCBootOrder)
	}

	// Read VPD Boot entries and create BIOS BootOrder
	BIOSBootOrder.bootMode = BMCBootOrder.bootMode
	// Initialize with INVALID_BOOT
	idx := 0
	for idx = range BIOSBootOrder.bootSeq {
		BIOSBootOrder.bootSeq[idx] = INVALID_BOOT
	}
	bootEntries := systembooter.GetBootEntries()
	idx = 0
	var bootType string
	for _, entry := range bootEntries {
		if idx >= 5 {
			break
		}
		if bootType = entry.Booter.TypeName(); len(bootType) > 0 {
			if bootType == "netboot" {
				BIOSBootOrder.bootSeq[idx] = NETWORK_BOOT
				BIOSBootOrder.bootSeq[idx] |= 0x8
				idx++
			} else if bootType == "localboot" {
				BIOSBootOrder.bootSeq[idx] = LOCAL_BOOT
				idx++
			}
		}
	}
	// If there is no valid VPD boot order, write the default systembooter configurations
	if idx == 0 {
		log.Printf("No valid VPD boot order, set default boot orders to RW_VPD")
		BIOSBootOrder.bootSeq[0] = NETWORK_BOOT
		BIOSBootOrder.bootSeq[0] |= 0x8
		BIOSBootOrder.bootSeq[1] = LOCAL_BOOT
		return updateVPDBootOrder(i, &BIOSBootOrder)
	}
	// clear valid bit and set BIOS boot order to BMC
	BIOSBootOrder.bootMode &^= 0x80
	return setBootOrder(i, &BIOSBootOrder)
}
