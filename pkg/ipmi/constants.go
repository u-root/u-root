// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

const (
	_IPMI_BMC_CHANNEL                = 0xf
	_IPMI_BUF_SIZE                   = 1024
	_IPMI_IOC_MAGIC                  = 'i'
	_IPMI_OPENIPMI_READ_TIMEOUT      = 15
	_IPMI_SYSTEM_INTERFACE_ADDR_TYPE = 0x0c

	// Net functions
	_IPMI_NETFN_CHASSIS   NetFn = 0x0
	_IPMI_NETFN_APP       NetFn = 0x6
	_IPMI_NETFN_STORAGE   NetFn = 0xA
	_IPMI_NETFN_TRANSPORT NetFn = 0xC

	// IPM Device "Global" Commands
	BMC_GET_DEVICE_ID Command = 0x01

	// BMC Device and Messaging Commands
	BMC_SET_WATCHDOG_TIMER     Command = 0x24
	BMC_GET_WATCHDOG_TIMER     Command = 0x25
	BMC_SET_GLOBAL_ENABLES     Command = 0x2E
	BMC_GET_GLOBAL_ENABLES     Command = 0x2F
	SET_SYSTEM_INFO_PARAMETERS Command = 0x58
	BMC_ADD_SEL                Command = 0x44

	// Chassis Device Commands
	BMC_GET_CHASSIS_STATUS Command = 0x01

	// SEL device Commands
	BMC_GET_SEL_INFO Command = 0x40

	// LAN Device Commands
	BMC_GET_LAN_CONFIG Command = 0x02

	IPM_WATCHDOG_NO_ACTION    = 0x00
	IPM_WATCHDOG_SMS_OS       = 0x04
	IPM_WATCHDOG_CLEAR_SMS_OS = 0x10

	ADTL_SEL_DEVICE         = 0x04
	EN_SYSTEM_EVENT_LOGGING = 0x08

	// SEL
	// STD_TYPE  = 0x02
	OEM_NTS_TYPE = 0xFB

	_SYSTEM_INFO_BLK_SZ = 16

	_SYSTEM_FW_VERSION = 1

	_ASCII = 0

	// Set 62 Bytes (4 sets) as the maximal string length
	strlenMax = 62
)
