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

	// Completion codes.
	// See Intelligent Platform Management Interface Specification v2.0 rev. 1.1, section 5.2.
	IPMI_CC_OK                          CompletionCode = 0x00
	IPMI_CC_NODE_BUSY                   CompletionCode = 0xc0
	IPMI_CC_INV_CMD                     CompletionCode = 0xc1
	IPMI_CC_INV_CMD_FOR_LUN             CompletionCode = 0xc2
	IPMI_CC_TIMEOUT                     CompletionCode = 0xc3
	IPMI_CC_OUT_OF_SPACE                CompletionCode = 0xc4
	IPMI_CC_RES_CANCELED                CompletionCode = 0xc5
	IPMI_CC_REQ_DATA_TRUNC              CompletionCode = 0xc6
	IPMI_CC_REQ_DATA_INV_LENGTH         CompletionCode = 0xc7
	IPMI_CC_REQ_DATA_FIELD_EXCEED       CompletionCode = 0xc8
	IPMI_CC_PARAM_OUT_OF_RANGE          CompletionCode = 0xc9
	IPMI_CC_CANT_RET_NUM_REQ_BYTES      CompletionCode = 0xca
	IPMI_CC_REQ_DATA_NOT_PRESENT        CompletionCode = 0xcb
	IPMI_CC_INV_DATA_FIELD_IN_REQ       CompletionCode = 0xcc
	IPMI_CC_ILL_SENSOR_OR_RECORD        CompletionCode = 0xcd
	IPMI_CC_RESP_COULD_NOT_BE_PRV       CompletionCode = 0xce
	IPMI_CC_CANT_RESP_DUPLI_REQ         CompletionCode = 0xcf
	IPMI_CC_CANT_RESP_SDRR_UPDATE       CompletionCode = 0xd0
	IPMI_CC_CANT_RESP_FIRM_UPDATE       CompletionCode = 0xd1
	IPMI_CC_CANT_RESP_BMC_INIT          CompletionCode = 0xd2
	IPMI_CC_DESTINATION_UNAVAILABLE     CompletionCode = 0xd3
	IPMI_CC_INSUFFICIENT_PRIVILEGES     CompletionCode = 0xd4
	IPMI_CC_NOT_SUPPORTED_PRESENT_STATE CompletionCode = 0xd5
	IPMI_CC_ILLEGAL_COMMAND_DISABLED    CompletionCode = 0xd6
	IPMI_CC_UNSPECIFIED_ERROR           CompletionCode = 0xff

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
