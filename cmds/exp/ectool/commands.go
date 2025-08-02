package main

/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Host communication command constants for Chrome EC */

const (

	/*
	 * Current version of this protocol
	 *
	 * TODO(crosbug.com/p/11223): This is effectively useless; protocol is
	 * determined in other ways.  Remove this once the kernel code no longer
	 * depends on it.
	 */
	ecProtoVersion = 0x00000002

	/* I/O addresses for ACPI commands */
	ecLpcAddrAcpiData = 0x62
	ecLpcAddrAcpiCmd  = 0x66

	/* I/O addresses for host command */
	ecLpcAddrHostData = 0x200
	ecLpcAddrHostCmd  = 0x204

	/* I/O addresses for host command args and params */
	/* Protocol version 2 */
	ecLpcAddrHostArgs  = 0x800 /* And 0x801, 0x802, 0x803 */
	ecLpcAddrHostParam = 0x804 /* For version 2 params; size is
	 * ecProto2MaxParamSize */
	/* Protocol version 3 */
	ecLpcAddrHostPacket = 0x800 /* Offset of version 3 packet */
	ecLpcHostPacketSize = 0x100 /* Max size of version 3 packet */

	/* The actual block is 0x800-0x8ff, but some BIOSes think it's 0x880-0x8ff
	 * and they tell the kernel that so we have to think of it as two parts. */
	ecHostCmdRegion0    = 0x800
	ecHostCmdRegion1    = 0x880
	ecHostCmdRegionSize = 0x80

	/* EC command register bit functions */
	ecLpcCmdrData     = (1 << 0) /* Data ready for host to read */
	ecLpcCmdrPending  = (1 << 1) /* Write pending to EC */
	ecLpcCmdrBusy     = (1 << 2) /* EC is busy processing a command */
	ecLpcCmdrCmd      = (1 << 3) /* Last host write was a command */
	ecLpcCmdrAcpiBrst = (1 << 4) /* Burst mode (not used) */
	ecLpcCmdrSci      = (1 << 5) /* SCI event is pending */
	ecLpcCmdrSmi      = (1 << 6) /* SMI event is pending */

	ecLpcAddrMemmap = 0x900
	ecMemmapSize    = 255 /* ACPI IO buffer max is 255 bytes */
	ecMemmapTextMax = 8   /* Size of a string in the memory map */

	/* The offset address of each type of data in mapped memory. */
	ecMemmapTempSensor      = 0x00 /* Temp sensors 0x00 - 0x0f */
	ecMemmapFan             = 0x10 /* Fan speeds 0x10 - 0x17 */
	ecMemmapTempSensorB     = 0x18 /* More temp sensors 0x18 - 0x1f */
	ecMemmapID              = 0x20 /* 0x20 == 'E', 0x21 == 'C' */
	ecMemmapIDVersion       = 0x22 /* Version of data in 0x20 - 0x2f */
	ecMemmapThermalVersion  = 0x23 /* Version of data in 0x00 - 0x1f */
	ecMemmapBatteryVersion  = 0x24 /* Version of data in 0x40 - 0x7f */
	ecMemmapSwitchesVersion = 0x25 /* Version of data in 0x30 - 0x33 */
	ecMemmapEventsVersion   = 0x26 /* Version of data in 0x34 - 0x3f */
	ecMemmapHostCmdFlags    = 0x27 /* Host cmd interface flags (8 bits) */
	/* Unused 0x28 - 0x2f */
	ecMemmapSwitches = 0x30 /* 8 bits */
	/* Unused 0x31 - 0x33 */
	ecMemmapHostEvents = 0x34 /* 32 bits */
	/* Reserve 0x38 - 0x3f for additional host event-related stuff */
	/* Battery values are all 32 bits */
	ecMemmapBattVolt = 0x40 /* Battery Present Voltage */
	ecMemmapBattRate = 0x44 /* Battery Present Rate */
	ecMemmapBattCap  = 0x48 /* Battery Remaining Capacity */
	ecMemmapBattFlag = 0x4c /* Battery State, defined below */
	ecMemmapBattDcap = 0x50 /* Battery Design Capacity */
	ecMemmapBattDvlt = 0x54 /* Battery Design Voltage */
	ecMemmapBattLfcc = 0x58 /* Battery Last Full Charge Capacity */
	ecMemmapBattCcnt = 0x5c /* Battery Cycle Count */
	/* Strings are all 8 bytes (ecMemmapTextMax) */
	ecMemmapBattMfgr   = 0x60 /* Battery Manufacturer String */
	ecMemmapBattModel  = 0x68 /* Battery Model Number String */
	ecMemmapBattSerial = 0x70 /* Battery Serial Number String */
	ecMemmapBattType   = 0x78 /* Battery Type String */
	ecMemmapAls        = 0x80 /* ALS readings in lux (2 X 16 bits) */
	/* Unused 0x84 - 0x8f */
	ecMemmapAccStatus = 0x90 /* Accelerometer status (8 bits )*/
	/* Unused 0x91 */
	ecMemmapAccData  = 0x92 /* Accelerometer data 0x92 - 0x9f */
	ecMemmapGyroData = 0xa0 /* Gyroscope data 0xa0 - 0xa5 */
	/* Unused 0xa6 - 0xdf */

	/*
	 * ACPI is unable to access memory mapped data at or above this offset due to
	 * limitations of the ACPI protocol. Do not place data in the range 0xe0 - 0xfe
	 * which might be needed by ACPI.
	 */
	ecMemmapNoAcpi = 0xe0

	/* Define the format of the accelerometer mapped memory status byte. */
	ecMemmapAccStatusSampleIDMask = 0x0f
	ecMemmapAccStatusBusyBit      = (1 << 4)
	ecMemmapAccStatusPresenceBit  = (1 << 7)

	/* Number of temp sensors at ecMemmapTempSensor */
	ecTempSensorEntries = 16
	/*
	 * Number of temp sensors at ecMemmapTempSensorB.
	 *
	 * Valid only if ecMemmapThermalVersion returns >= 2.
	 */
	ecTempSensorBEntries = 8

	/* Special values for mapped temperature sensors */
	ecTempSensorNotPresent    = 0xff
	ecTempSensorError         = 0xfe
	ecTempSensorNotPowered    = 0xfd
	ecTempSensorNotCalibrated = 0xfc
	/*
	 * The offset of temperature value stored in mapped memory.  This allows
	 * reporting a temperature range of 200K to 454K = -73C to 181C.
	 */
	ecTempSensorOffset = 200

	/*
	 * Number of ALS readings at ecMemmapAls
	 */
	ecAlsEntries = 2

	/*
	 * The default value a temperature sensor will return when it is present but
	 * has not been read this boot.  This is a reasonable number to avoid
	 * triggering alarms on the host.
	 */
	ecTempSensorDefault = (296 - ecTempSensorOffset)

	ecFanSpeedEntries    = 4      /* Number of fans at ecMemmapFan */
	ecFanSpeedNotPresent = 0xffff /* Entry not present */
	ecFanSpeedStalled    = 0xfffe /* Fan stalled */

	/* Battery bit flags at ecMemmapBattFlag. */
	ecBattFlagAcPresent     = 0x01
	ecBattFlagBattPresent   = 0x02
	ecBattFlagDischarging   = 0x04
	ecBattFlagCharging      = 0x08
	ecBattFlagLevelCritical = 0x10

	/* Switch flags at ecMemmapSwitches */
	ecSwitchLidOpen              = 0x01
	ecSwitchPowerButtonPressed   = 0x02
	ecSwitchWriteProtectDisabled = 0x04
	/* Was recovery requested via keyboard; now unused. */
	ecSwitchIgnore1 = 0x08
	/* Recovery requested via dedicated signal (from servo board) */
	ecSwitchDedicatedRecovery = 0x10
	/* Was fake developer mode switch; now unused.  Remove in next refactor. */
	ecSwitchIgnore0 = 0x20

	/* Host command interface flags */
	/* Host command interface supports LPC args (LPC interface only) */
	ecHostCmdFlagLpcArgsSupported = 0x01
	/* Host command interface supports version 3 protocol */
	ecHostCmdFlagVersion3 = 0x02

	/* Wireless switch flags */
	ecWirelessSwitchAll       = ^0x00 /* All flags */
	ecWirelessSwitchWlan      = 0x01  /* WLAN radio */
	ecWirelessSwitchBluetooth = 0x02  /* Bluetooth radio */
	ecWirelessSwitchWwan      = 0x04  /* WWAN power */
	ecWirelessSwitchWlanPower = 0x08  /* WLAN power */

	/*****************************************************************************/
	/*
	 * ACPI commands
	 *
	 * These are valid ONLY on the ACPI command/data port.
	 */

	/*
	 * ACPI Read Embedded Controller
	 *
	 * This reads from ACPI memory space on the EC (ecAcpiMem_*).
	 *
	 * Use the following sequence:
	 *
	 *    - Write ecCmdAcpiRead to ecLpcAddrAcpiCmd
	 *    - Wait for ecLpcCmdrPending bit to clear
	 *    - Write address to ecLpcAddrAcpiData
	 *    - Wait for ecLpcCmdrData bit to set
	 *    - Read value from ecLpcAddrAcpiData
	 */
	ecCmdAcpiRead = 0x80

	/*
	 * ACPI Write Embedded Controller
	 *
	 * This reads from ACPI memory space on the EC (ecAcpiMem_*).
	 *
	 * Use the following sequence:
	 *
	 *    - Write ecCmdAcpiWrite to ecLpcAddrAcpiCmd
	 *    - Wait for ecLpcCmdrPending bit to clear
	 *    - Write address to ecLpcAddrAcpiData
	 *    - Wait for ecLpcCmdrPending bit to clear
	 *    - Write value to ecLpcAddrAcpiData
	 */
	ecCmdAcpiWrite = 0x81

	/*
	 * ACPI Burst Enable Embedded Controller
	 *
	 * This enables burst mode on the EC to allow the host to issue several
	 * commands back-to-back. While in this mode, writes to mapped multi-byte
	 * data are locked out to ensure data consistency.
	 */
	ecCmdAcpiBurstEnable = 0x82

	/*
	 * ACPI Burst Disable Embedded Controller
	 *
	 * This disables burst mode on the EC and stops preventing EC writes to mapped
	 * multi-byte data.
	 */
	ecCmdAcpiBurstDisable = 0x83

	/*
	 * ACPI Query Embedded Controller
	 *
	 * This clears the lowest-order bit in the currently pending host events, and
	 * sets the result code to the 1-based index of the bit (event 0x00000001 = 1
	 * event 0x80000000 = 32), or 0 if no event was pending.
	 */
	ecCmdAcpiQueryEvent = 0x84

	/* Valid addresses in ACPI memory space, for read/write commands */

	/* Memory space version; set to ecAcpiMemVersionCurrent */
	ecAcpiMemVersion = 0x00
	/*
	 * Test location; writing value here updates test compliment byte to (0xff -
	 * value).
	 */
	ecAcpiMemTest = 0x01
	/* Test compliment; writes here are ignored. */
	ecAcpiMemTestCompliment = 0x02

	/* Keyboard backlight brightness percent (0 - 100) */
	ecAcpiMemKeyboardBacklight = 0x03
	/* DPTF Target Fan Duty (0-100, 0xff for auto/none) */
	ecAcpiMemFanDuty = 0x04

	/*
	 * DPTF temp thresholds. Any of the EC's temp sensors can have up to two
	 * independent thresholds attached to them. The current value of the ID
	 * register determines which sensor is affected by the THRESHOLD and COMMIT
	 * registers. The THRESHOLD register uses the same ecTempSensorOffset scheme
	 * as the memory-mapped sensors. The COMMIT register applies those settings.
	 *
	 * The spec does not mandate any way to read back the threshold settings
	 * themselves, but when a threshold is crossed the AP needs a way to determine
	 * which sensor(s) are responsible. Each reading of the ID register clears and
	 * returns one sensor ID that has crossed one of its threshold (in either
	 * direction) since the last read. A value of 0xFF means "no new thresholds
	 * have tripped". Setting or enabling the thresholds for a sensor will clear
	 * the unread event count for that sensor.
	 */
	ecAcpiMemTempID        = 0x05
	ecAcpiMemTempThreshold = 0x06
	ecAcpiMemTempCommit    = 0x07
	/*
	 * Here are the bits for the COMMIT register:
	 *   bit 0 selects the threshold index for the chosen sensor (0/1)
	 *   bit 1 enables/disables the selected threshold (0 = off, 1 = on)
	 * Each write to the commit register affects one threshold.
	 */
	ecAcpiMemTempCommitSelectMask = (1 << 0)
	ecAcpiMemTempCommitEnableMask = (1 << 1)
	/*
	 * Example:
	 *
	 * Set the thresholds for sensor 2 to 50 C and 60 C:
	 *   write 2 to [0x05]      --  select temp sensor 2
	 *   write 0x7b to [0x06]   --  CToK(50) - ecTempSensorOffset
	 *   write 0x2 to [0x07]    --  enable threshold 0 with this value
	 *   write 0x85 to [0x06]   --  CToK(60) - ecTempSensorOffset
	 *   write 0x3 to [0x07]    --  enable threshold 1 with this value
	 *
	 * Disable the 60 C threshold, leaving the 50 C threshold unchanged:
	 *   write 2 to [0x05]      --  select temp sensor 2
	 *   write 0x1 to [0x07]    --  disable threshold 1
	 */

	/* DPTF battery charging current limit */
	ecAcpiMemChargingLimit = 0x08

	/* Charging limit is specified in 64 mA steps */
	ecAcpiMemChargingLimitStepMa = 64
	/* Value to disable DPTF battery charging limit */
	ecAcpiMemChargingLimitDisabled = 0xff

	/*
	 * ACPI addresses 0x20 - 0xff map to ecMemmap offset 0x00 - 0xdf.  This data
	 * is read-only from the AP.  Added in ecAcpiMemVersion 2.
	 */
	ecAcpiMemMappedBegin = 0x20
	ecAcpiMemMappedSize  = 0xe0

	/* Current version of ACPI memory address space */
	ecAcpiMemVersionCurrent = 2

	/*
	 * This header file is used in coreboot both in C and ACPI code.  The ACPI code
	 * is pre-processed to handle constants but the ASL compiler is unable to
	 * handle actual C code so keep it separate.
	 */

	/* LPC command status byte masks */
	/* EC has written a byte in the data register and host hasn't read it yet */
	ecLpcStatusToHost = 0x01
	/* Host has written a command/data byte and the EC hasn't read it yet */
	ecLpcStatusFromHost = 0x02
	/* EC is processing a command */
	ecLpcStatusProcessing = 0x04
	/* Last write to EC was a command, not data */
	ecLpcStatusLastCmd = 0x08
	/* EC is in burst mode */
	ecLpcStatusBurstMode = 0x10
	/* SCI event is pending (requesting SCI query) */
	ecLpcStatusSciPending = 0x20
	/* SMI event is pending (requesting SMI query) */
	ecLpcStatusSmiPending = 0x40
	/* (reserved) */
	ecLpcStatusReserved = 0x80

	/*
	 * EC is busy.  This covers both the EC processing a command, and the host has
	 * written a new command but the EC hasn't picked it up yet.
	 */
	ecLpcStatusBusyMask = (ecLpcStatusFromHost | ecLpcStatusProcessing)
)

/* Host command response codes */
type ecStatus uint8

const (
	ecResSuccess          ecStatus = 0
	ecResInvalidCommand   ecStatus = 1
	ecResError            ecStatus = 2
	ecResInvalidParam     ecStatus = 3
	ecResAccessDenied     ecStatus = 4
	ecResInvalidResponse  ecStatus = 5
	ecResInvalidVersion   ecStatus = 6
	ecResInvalidChecksum  ecStatus = 7
	ecResInProgress       ecStatus = 8  /* Accepted, command in progress */
	ecResUnavailable      ecStatus = 9  /* No response available */
	ecResTimeout          ecStatus = 10 /* We got a timeout */
	ecResOverflow         ecStatus = 11 /* Table / data overflow */
	ecResInvalidHeader    ecStatus = 12 /* Header contains invalid data */
	ecResRequestTruncated ecStatus = 13 /* Didn't get the entire request */
	ecResResponseTooBig   ecStatus = 14 /* Response was too big to handle */
	ecResBusError         ecStatus = 15 /* Communications bus error */
	ecResBusy             ecStatus = 16 /* Up but too busy.  Should retry */
)

/*
 * Host event codes.  Note these are 1-based, not 0-based, because ACPI query
 * EC command uses code 0 to mean "no event pending".  We explicitly specify
 * each value in the enum listing so they won't change if we delete/insert an
 * item or rearrange the list (it needs to be stable across platforms, not
 * just within a single compiled instance).
 */
type hostEventCode uint8

const (
	ecHostEventLidClosed        hostEventCode = 1
	ecHostEventLidOpen          hostEventCode = 2
	ecHostEventPowerButton      hostEventCode = 3
	ecHostEventAcConnected      hostEventCode = 4
	ecHostEventAcDisconnected   hostEventCode = 5
	ecHostEventBatteryLow       hostEventCode = 6
	ecHostEventBatteryCritical  hostEventCode = 7
	ecHostEventBattery          hostEventCode = 8
	ecHostEventThermalThreshold hostEventCode = 9
	ecHostEventThermalOverload  hostEventCode = 10
	ecHostEventThermal          hostEventCode = 11
	ecHostEventUsbCharger       hostEventCode = 12
	ecHostEventKeyPressed       hostEventCode = 13
	/*
	 * EC has finished initializing the host interface.  The host can check
	 * for this event following sending a ecCmdRebootEc command to
	 * determine when the EC is ready to accept subsequent commands.
	 */
	ecHostEventInterfaceReady = 14
	/* Keyboard recovery combo has been pressed */
	ecHostEventKeyboardRecovery = 15

	/* Shutdown due to thermal overload */
	ecHostEventThermalShutdown = 16
	/* Shutdown due to battery level too low */
	ecHostEventBatteryShutdown = 17

	/* Suggest that the AP throttle itself */
	ecHostEventThrottleStart = 18
	/* Suggest that the AP resume normal speed */
	ecHostEventThrottleStop = 19

	/* Hang detect logic detected a hang and host event timeout expired */
	ecHostEventHangDetect = 20
	/* Hang detect logic detected a hang and warm rebooted the AP */
	ecHostEventHangReboot = 21

	/* PD MCU triggering host event */
	ecHostEventPdMcu = 22

	/* Battery Status flags have changed */
	ecHostEventBatteryStatus = 23

	/* EC encountered a panic, triggering a reset */
	ecHostEventPanic = 24

	/*
	 * The high bit of the event mask is not used as a host event code.  If
	 * it reads back as set, then the entire event mask should be
	 * considered invalid by the host.  This can happen when reading the
	 * raw event status via ecMemmapHostEvents but the LPC interface is
	 * not initialized on the EC, or improperly configured on the host.
	 */
	ecHostEventInvalid = 32
)

/* Host event mask */
func ecHostEventMask(eventCode uint8) uint8 {
	return 1 << ((eventCode) - 1)
}

/* TYPE */
/* Arguments at ecLpcAddrHostArgs */
type ecLpcHostArgs struct {
	flags          uint8
	commandVersion uint8
	dataSize       uint8
	/*
	 * Checksum; sum of command + flags + commandVersion + dataSize +
	 * all params/response data bytes.
	 */
	checksum uint8
}

/* Flags for ecLpcHostArgs.flags */
/*
 * Args are from host.  Data area at ecLpcAddrHostParam contains command
 * params.
 *
 * If EC gets a command and this flag is not set, this is an old-style command.
 * Command version is 0 and params from host are at ecLpcAddrOldParam with
 * unknown length.  EC must respond with an old-style response (that is
 * withouth setting ecHostArgsFlagToHost).
 */
const ecHostArgsFlagFromHost = 0x01

/*
 * Args are from EC.  Data area at ecLpcAddrHostParam contains response.
 *
 * If EC responds to a command and this flag is not set, this is an old-style
 * response.  Command version is 0 and response data from EC is at
 * ecLpcAddrOldParam with unknown length.
 */
const ecHostArgsFlagToHost = 0x02

/*****************************************************************************/
/*
 * Byte codes returned by EC over SPI interface.
 *
 * These can be used by the AP to debug the EC interface, and to determine
 * when the EC is not in a state where it will ever get around to responding
 * to the AP.
 *
 * Example of sequence of bytes read from EC for a current good transfer:
 *   1. -                  - AP asserts chip select (CS#)
 *   2. ecSpiOldReady      - AP sends first byte(s) of request
 *   3. -                  - EC starts handling CS# interrupt
 *   4. ecSpiReceiving     - AP sends remaining byte(s) of request
 *   5. ecSpiProcessing    - EC starts processing request; AP is clocking in
 *                           bytes looking for ecSpiFrameStart
 *   6. -                  - EC finishes processing and sets up response
 *   7. ecSpiFrameStart    - AP reads frame byte
 *   8. (response packet)  - AP reads response packet
 *   9. ecSpiPastEnd       - Any additional bytes read by AP
 *   10 -                  - AP deasserts chip select
 *   11 -                  - EC processes CS# interrupt and sets up DMA for
 *                           next request
 *
 * If the AP is waiting for ecSpiFrameStart and sees any value other than
 * the following byte values:
 *   ecSpiOldReady
 *   ecSpiRxReady
 *   ecSpiReceiving
 *   ecSpiProcessing
 *
 * Then the EC found an error in the request, or was not ready for the request
 * and lost data.  The AP should give up waiting for ecSpiFrameStart
 * because the EC is unable to tell when the AP is done sending its request.
 */

/*
 * Framing byte which precedes a response packet from the EC.  After sending a
 * request, the AP will clock in bytes until it sees the framing byte, then
 * clock in the response packet.
 */
const (
	ecSpiFrameStart = 0xec

	/*
	 * Padding bytes which are clocked out after the end of a response packet.
	 */
	ecSpiPastEnd = 0xed

	/*
	 * EC is ready to receive, and has ignored the byte sent by the AP.  EC expects
	 * that the AP will send a valid packet header (starting with
	 * ecCommandProtocol3) in the next 32 bytes.
	 */
	ecSpiRxReady = 0xf8

	/*
	 * EC has started receiving the request from the AP, but hasn't started
	 * processing it yet.
	 */
	ecSpiReceiving = 0xf9

	/* EC has received the entire request from the AP and is processing it. */
	ecSpiProcessing = 0xfa

	/*
	 * EC received bad data from the AP, such as a packet header with an invalid
	 * length.  EC will ignore all data until chip select deasserts.
	 */
	ecSpiRxBadData = 0xfb

	/*
	 * EC received data from the AP before it was ready.  That is, the AP asserted
	 * chip select and started clocking data before the EC was ready to receive it.
	 * EC will ignore all data until chip select deasserts.
	 */
	ecSpiNotReady = 0xfc

	/*
	 * EC was ready to receive a request from the AP.  EC has treated the byte sent
	 * by the AP as part of a request packet, or (for old-style ECs) is processing
	 * a fully received packet but is not ready to respond yet.
	 */
	ecSpiOldReady = 0xfd

	/*****************************************************************************/

	/*
	 * Protocol version 2 for I2C and SPI send a request this way:
	 *
	 *	0	ecCmdVersion0 + (command version)
	 *	1	Command number
	 *	2	Length of params = N
	 *	3..N+2	Params, if any
	 *	N+3	8-bit checksum of bytes 0..N+2
	 *
	 * The corresponding response is:
	 *
	 *	0	Result code (ecRes_*)
	 *	1	Length of params = M
	 *	2..M+1	Params, if any
	 *	M+2	8-bit checksum of bytes 0..M+1
	 */
	ecProto2RequestHeaderBytes  = 3
	ecProto2RequestTrailerBytes = 1
	ecProto2RequestOverhead     = (ecProto2RequestHeaderBytes +
		ecProto2RequestTrailerBytes)

	ecProto2ResponseHeaderBytes  = 2
	ecProto2ResponseTrailerBytes = 1
	ecProto2ResponseOverhead     = (ecProto2ResponseHeaderBytes +
		ecProto2ResponseTrailerBytes)

	/* Parameter length was limited by the LPC interface */
	ecProto2MaxParamSize = 0xfc

	/* Maximum request and response packet sizes for protocol version 2 */
	ecProto2MaxRequestSize = (ecProto2RequestOverhead +
		ecProto2MaxParamSize)
	ecProto2MaxResponseSize = (ecProto2ResponseOverhead +
		ecProto2MaxParamSize)

	/*****************************************************************************/

	/*
	 * Value written to legacy command port / prefix byte to indicate protocol
	 * 3+ structs are being used.  Usage is bus-dependent.
	 */
	ecCommandProtocol3 = 0xda

	ecHostRequestVersion = 3
)

/* TYPE */
/* Version 3 request from host */
type ecHostRequest struct {
	/* Struct version (=3)
	 *
	 * EC will return ecResInvalidHeader if it receives a header with a
	 * version it doesn't know how to parse.
	 */
	structVersion uint8

	/*
	 * Checksum of request and data; sum of all bytes including checksum
	 * should total to 0.
	 */
	checksum uint8

	/* Command code */
	command uint16

	/* Command version */
	commandVersion uint8

	/* Unused byte in current protocol version; set to 0 */
	reserved uint8

	/* Length of data which follows this header */
	dataLen uint16
}

const ecHostResponseVersion = 3

/* TYPE */
/* Version 3 response from EC */
type ecHostResponse struct {
	/* Struct version (=3) */
	structVersion uint8

	/*
	 * Checksum of response and data; sum of all bytes including checksum
	 * should total to 0.
	 */
	checksum uint8

	/* Result code (ecRes_*) */
	result uint16

	/* Length of data which follows this header */
	dataLen uint16

	/* Unused bytes in current protocol version; set to 0 */
	reserved uint16
}

/*****************************************************************************/
/*
 * Notes on commands:
 *
 * Each command is an 16-bit command value.  Commands which take params or
 * return response data specify structs for that data.  If no struct is
 * specified, the command does not input or output data, respectively.
 * Parameter/response length is implicit in the structs.  Some underlying
 * communication protocols (I2C, SPI) may add length or checksum headers, but
 * those are implementation-dependent and not defined here.
 */

/*****************************************************************************/
/* General / test commands */

/*
 * Get protocol version, used to deal with non-backward compatible protocol
 * changes.
 */
const ecCmdProtoVersion = 0x00

/* TYPE */
type ecResponseProtoVersion struct {
	version uint32
}

const (
	/*
	 * Hello.  This is a simple command to test the EC is responsive to
	 * commands.
	 */
	ecCmdHello = 0x01
)

/* TYPE */
type ecParamsHello struct {
	inData uint32 /* Pass anything here */
}

/* TYPE */
type ecResponseHello struct {
	outData uint32 /* Output will be inData + 0x01020304 */
}

const (
	/* Get version number */
	ecCmdGetVersion = 0x02
)

type ecCurrentImage uint8

const (
	ecImageUnknown ecCurrentImage = 0
	ecImageRo
	ecImageRw
)

/* TYPE */
type ecResponseGetVersion struct {
	/* Null-terminated version strings for RO, RW */
	versionStringRo [32]byte
	versionStringRw [32]byte
	reserved        [32]byte /* Was previously RW-B string */
	currentImage    uint32   /* One of ecCurrentImage */
}

const (
	/* Read test */
	ecCmdReadTest = 0x03
)

/* TYPE */
type ecParamsReadTest struct {
	offset uint32 /* Starting value for read buffer */
	size   uint32 /* Size to read in bytes */
}

/* TYPE */
type ecResponseReadTest struct {
	data [32]uint32
}

const (
	/*
	 * Get build information
	 *
	 * Response is null-terminated string.
	 */
	ecCmdGetBuildInfo = 0x04

	/* Get chip info */
	ecCmdGetChipInfo = 0x05
)

/* TYPE */
type ecResponseGetChipInfo struct {
	/* Null-terminated strings */
	vendor   [32]byte
	name     [32]byte
	revision [32]byte /* Mask version */
}

const (
	/* Get board HW version */
	ecCmdGetBoardVersion = 0x06
)

/* TYPE */
type ecResponseBoardVersion struct {
	boardVersion uint16 /* A monotonously incrementing number. */
}

/*
 * Read memory-mapped data.
 *
 * This is an alternate interface to memory-mapped data for bus protocols
 * which don't support direct-mapped memory - I2C, SPI, etc.
 *
 * Response is params.size bytes of data.
 */
const (
	ecCmdReadMemmap = 0x07
)

/* TYPE */
type ecParamsReadMemmap struct {
	offset uint8 /* Offset in memmap (ecMemmap_*) */
	size   uint8 /* Size to read in bytes */
}

/* Read versions supported for a command */
const ecCmdGetCmdVersions = 0x08

/* TYPE */
type ecParamsGetCmdVersions struct {
	cmd uint8 /* Command to check */
}

/* TYPE */
type ecParamsGetCmdVersionsV1 struct {
	cmd uint16 /* Command to check */
}

/* TYPE */
type ecResponseGetCmdVersions struct {
	/*
	 * Mask of supported versions; use ecVerMask() to compare with a
	 * desired version.
	 */
	versionMask uint32
}

/*
 * Check EC communcations status (busy). This is needed on i2c/spi but not
 * on lpc since it has its own out-of-band busy indicator.
 *
 * lpc must read the status from the command register. Attempting this on
 * lpc will overwrite the args/parameter space and corrupt its data.
 */
const ecCmdGetCommsStatus = 0x09

/* Avoid using ecStatus which is for return values */
type ecCommsStatus uint8

const (
	ecCommsStatusProcessing ecCommsStatus = 1 << 0 /* Processing cmd */
)

/* TYPE */
type ecResponseGetCommsStatus struct {
	flags uint32 /* Mask of enum ecCommsStatus */
}

const (
	/* Fake a variety of responses, purely for testing purposes. */
	ecCmdTestProtocol = 0x0a
)

/* TYPE */
/* Tell the EC what to send back to us. */
type ecParamsTestProtocol struct {
	ecResult uint32
	retLen   uint32
	buf      [32]uint8
}

/* TYPE */
/* Here it comes... */
type ecResponseTestProtocol struct {
	buf [32]uint8
}

/* Get prococol information */
const ecCmdGetProtocolInfo = 0x0b

/* Flags for ecResponseGetProtocolInfo.flags */
/* ecResInProgress may be returned if a command is slow */
const ecProtocolInfoInProgressSupported = (1 << 0)

/* TYPE */
type ecResponseGetProtocolInfo struct {
	/* Fields which exist if at least protocol version 3 supported */

	/* Bitmask of protocol versions supported (1 << n means version n)*/
	protocolVersions uint32

	/* Maximum request packet size, in bytes */
	maxRequestPacketSize uint16

	/* Maximum response packet size, in bytes */
	maxResponsePacketSize uint16

	/* Flags; see ecProtocolInfo_* */
	flags uint32
}

/*****************************************************************************/
/* Get/Set miscellaneous values */
const (
	/* The upper byte of .flags tells what to do (nothing means "get") */
	ecGsvSet = 0x80000000

	/* The lower three bytes of .flags identifies the parameter, if that has
	   meaning for an individual command. */
	ecGsvParamMask = 0x00ffffff
)

/* TYPE */
type ecParamsGetSetValue struct {
	flags uint32
	value uint32
}

/* TYPE */
type ecResponseGetSetValue struct {
	flags uint32
	value uint32
}

/* More than one command can use these structs to get/set parameters. */
const ecCmdGsvPauseInS5 = 0x0c

/*****************************************************************************/
/* Flash commands */

/* Get flash info */
const ecCmdFlashInfo = 0x10

/* TYPE */
/* Version 0 returns these fields */
type ecResponseFlashInfo struct {
	/* Usable flash size, in bytes */
	flashSize uint32
	/*
	 * Write block size.  Write offset and size must be a multiple
	 * of this.
	 */
	writeBlockSize uint32
	/*
	 * Erase block size.  Erase offset and size must be a multiple
	 * of this.
	 */
	eraseBlockSize uint32
	/*
	 * Protection block size.  Protection offset and size must be a
	 * multiple of this.
	 */
	protectBlockSize uint32
}

/* Flags for version 1+ flash info command */
/* EC flash erases bits to 0 instead of 1 */
const ecFlashInfoEraseTo0 = (1 << 0)

/* TYPE */
/*
 * Version 1 returns the same initial fields as version 0, with additional
 * fields following.
 *
 * gcc anonymous structs don't seem to get along with the  directive;
 * if they did we'd define the version 0 struct as a sub-struct of this one.
 */
type ecResponseFlashInfo1 struct {
	/* Version 0 fields; see above for description */
	flashSize        uint32
	writeBlockSize   uint32
	eraseBlockSize   uint32
	protectBlockSize uint32

	/* Version 1 adds these fields: */
	/*
	 * Ideal write size in bytes.  Writes will be fastest if size is
	 * exactly this and offset is a multiple of this.  For example, an EC
	 * may have a write buffer which can do half-page operations if data is
	 * aligned, and a slower word-at-a-time write mode.
	 */
	writeIdealSize uint32

	/* Flags; see ecFlashInfo_* */
	flags uint32
}

/*
 * Read flash
 *
 * Response is params.size bytes of data.
 */
const ecCmdFlashRead = 0x11

/* TYPE */
type ecParamsFlashRead struct {
	offset uint32 /* Byte offset to read */
	size   uint32 /* Size to read in bytes */
}

const (
	/* Write flash */
	ecCmdFlashWrite = 0x12
	ecVerFlashWrite = 1

	/* Version 0 of the flash command supported only 64 bytes of data */
	ecFlashWriteVer0Size = 64
)

/* TYPE */
type ecParamsFlashWrite struct {
	offset uint32 /* Byte offset to write */
	size   uint32 /* Size to write in bytes */
	/* Followed by data to write */
}

/* Erase flash */
const ecCmdFlashErase = 0x13

/* TYPE */
type ecParamsFlashErase struct {
	offset uint32 /* Byte offset to erase */
	size   uint32 /* Size to erase in bytes */
}

const (
	/*
	 * Get/set flash protection.
	 *
	 * If mask!=0, sets/clear the requested bits of flags.  Depending on the
	 * firmware write protect GPIO, not all flags will take effect immediately;
	 * some flags require a subsequent hard reset to take effect.  Check the
	 * returned flags bits to see what actually happened.
	 *
	 * If mask=0, simply returns the current flags state.
	 */
	ecCmdFlashProtect = 0x15
	ecVerFlashProtect = 1 /* Command version 1 */

	/* Flags for flash protection */
	/* RO flash code protected when the EC boots */
	ecFlashProtectRoAtBoot = (1 << 0)
	/*
	 * RO flash code protected now.  If this bit is set, at-boot status cannot
	 * be changed.
	 */
	ecFlashProtectRoNow = (1 << 1)
	/* Entire flash code protected now, until reboot. */
	ecFlashProtectAllNow = (1 << 2)
	/* Flash write protect GPIO is asserted now */
	ecFlashProtectGpioAsserted = (1 << 3)
	/* Error - at least one bank of flash is stuck locked, and cannot be unlocked */
	ecFlashProtectErrorStuck = (1 << 4)
	/*
	 * Error - flash protection is in inconsistent state.  At least one bank of
	 * flash which should be protected is not protected.  Usually fixed by
	 * re-requesting the desired flags, or by a hard reset if that fails.
	 */
	ecFlashProtectErrorInconsistent = (1 << 5)
	/* Entire flash code protected when the EC boots */
	ecFlashProtectAllAtBoot = (1 << 6)
)

/* TYPE */
type ecParamsFlashProtect struct {
	mask  uint32 /* Bits in flags to apply */
	flags uint32 /* New flags to apply */
}

/* TYPE */
type ecResponseFlashProtect struct {
	/* Current value of flash protect flags */
	flags uint32
	/*
	 * Flags which are valid on this platform.  This allows the caller
	 * to distinguish between flags which aren't set vs. flags which can't
	 * be set on this platform.
	 */
	validFlags uint32
	/* Flags which can be changed given the current protection state */
	writableFlags uint32
}

/*
 * Note: commands 0x14 - 0x19 version 0 were old commands to get/set flash
 * write protect.  These commands may be reused with version > 0.
 */

/* Get the region offset/size */
const (
	ecCmdFlashRegionInfo = 0x16
	ecVerFlashRegionInfo = 1
)

type ecFlashRegion uint8

const (
	/* Region which holds read-only EC image */
	ecFlashRegionRo ecFlashRegion = iota
	/* Region which holds rewritable EC image */
	ecFlashRegionRw
	/*
	 * Region which should be write-protected in the factory (a superset of
	 * ecFlashRegionRo)
	 */
	ecFlashRegionWpRo
	/* Number of regions */
	ecFlashRegionCount
)

/* TYPE */
type ecParamsFlashRegionInfo struct {
	region uint32 /* enum ecFlashRegion */
}

/* TYPE */
type ecResponseFlashRegionInfo struct {
	offset uint32
	size   uint32
}

const (
	/* Read/write VbNvContext */
	ecCmdVbnvContext = 0x17
	ecVerVbnvContext = 1
	ecVbnvBlockSize  = 16
)

type ecVbnvcontextOp uint8

const (
	ecVbnvContextOpRead ecVbnvcontextOp = iota
	ecVbnvContextOpWrite
)

/* TYPE */
type ecParamsVbnvcontext struct {
	op    uint32
	block [ecVbnvBlockSize]uint8
}

/* TYPE */
type ecResponseVbnvcontext struct {
	block [ecVbnvBlockSize]uint8
}

/*****************************************************************************/
/* PWM commands */

/* Get fan target RPM */
const ecCmdPwmGetFanTargetRpm = 0x20

/* TYPE */
type ecResponsePwmGetFanRpm struct {
	rpm uint32
}

/* Set target fan RPM */
const ecCmdPwmSetFanTargetRpm = 0x21

/* TYPE */
/* Version 0 of input params */
type ecParamsPwmSetFanTargetRpmV0 struct {
	rpm uint32
}

/* TYPE */
/* Version 1 of input params */
type ecParamsPwmSetFanTargetRpmV1 struct {
	rpm    uint32
	fanIdx uint8
}

/* Get keyboard backlight */
const ecCmdPwmGetKeyboardBacklight = 0x22

/* TYPE */
type ecResponsePwmGetKeyboardBacklight struct {
	percent uint8
	enabled uint8
}

/* Set keyboard backlight */
const ecCmdPwmSetKeyboardBacklight = 0x23

/* TYPE */
type ecParamsPwmSetKeyboardBacklight struct {
	percent uint8
}

/* Set target fan PWM duty cycle */
const ecCmdPwmSetFanDuty = 0x24

/* TYPE */
/* Version 0 of input params */
type ecParamsPwmSetFanDutyV0 struct {
	percent uint32
}

/* TYPE */
/* Version 1 of input params */
type ecParamsPwmSetFanDutyV1 struct {
	percent uint32
	fanIdx  uint8
}

/*****************************************************************************/
/*
 * Lightbar commands. This looks worse than it is. Since we only use one HOST
 * command to say "talk to the lightbar", we put the "and tell it to do X" part
 * into a subcommand. We'll make separate structs for subcommands with
 * different input args, so that we know how much to expect.
 */
const ecCmdLightbarCmd = 0x28

/* TYPE */
type rgbS struct {
	r, g, b uint8
}

const lbBatteryLevels = 4

/* TYPE */
/* List of tweakable parameters. NOTE: It's  so it can be sent in a
 * host command, but the alignment is the same regardless. Keep it that way.
 */
type lightbarParamsV0 struct {
	/* Timing */
	googleRampUp   int32
	googleRampDown int32
	s3s0RampUp     int32
	s0TickDelay    [2]int32 /* AC=0/1 */
	s0aTickDelay   [2]int32 /* AC=0/1 */
	s0s3RampDown   int32
	s3SleepFor     int32
	s3RampUp       int32
	s3RampDown     int32

	/* Oscillation */
	newS0  uint8
	oscMin [2]uint8 /* AC=0/1 */
	oscMax [2]uint8 /* AC=0/1 */
	wOfs   [2]uint8 /* AC=0/1 */

	/* Brightness limits based on the backlight and AC. */
	brightBlOffFixed [2]uint8 /* AC=0/1 */
	brightBlOnMin    [2]uint8 /* AC=0/1 */
	brightBlOnMax    [2]uint8 /* AC=0/1 */

	/* Battery level thresholds */
	batteryhreshold [lbBatteryLevels - 1]uint8

	/* Map [AC][batteryLevel] to color index */
	s0Idx [2][lbBatteryLevels]uint8 /* AP is running */
	s3Idx [2][lbBatteryLevels]uint8 /* AP is sleeping */

	/* Color palette */
	color [8]rgbS /* 0-3 are Google colors */
}

/* TYPE */
type lightbarParamsV1 struct {
	/* Timing */
	googleRampUp   int32
	googleRampDown int32
	s3s0RampUp     int32
	s0TickDelay    [2]int32 /* AC=0/1 */
	s0aTickDelay   [2]int32 /* AC=0/1 */
	s0s3RampDown   int32
	s3SleepFor     int32
	s3RampUp       int32
	s3RampDown     int32
	s5RampUp       int32
	s5RampDown     int32
	tapTickDelay   int32
	tapGateDelay   int32
	tapDisplayTime int32

	/* Tap-for-battery params */
	tapPctRed   uint8
	tapPctGreen uint8
	tapSegMinOn uint8
	tapSegMaxOn uint8
	tapSegOsc   uint8
	tapIdx      [3]uint8

	/* Oscillation */
	oscMin [2]uint8 /* AC=0/1 */
	oscMax [2]uint8 /* AC=0/1 */
	wOfs   [2]uint8 /* AC=0/1 */

	/* Brightness limits based on the backlight and AC. */
	brightBlOffFixed [2]uint8 /* AC=0/1 */
	brightBlOnMin    [2]uint8 /* AC=0/1 */
	brightBlOnMax    [2]uint8 /* AC=0/1 */

	/* Battery level thresholds */
	batteryhreshold [lbBatteryLevels - 1]uint8

	/* Map [AC][batteryLevel] to color index */
	s0Idx [2][lbBatteryLevels]uint8 /* AP is running */
	s3Idx [2][lbBatteryLevels]uint8 /* AP is sleeping */

	/* s5: single color pulse on inhibited power-up */
	s5Idx uint8

	/* Color palette */
	color [8]rgbS /* 0-3 are Google colors */
}

/* TYPE */
/* Lightbar command params v2
 * crbug.com/467716
 *
 * lightbarParmsV1 was too big for i2c, therefore in v2, we split them up by
 * logical groups to make it more manageable ( < 120 bytes).
 *
 * NOTE: Each of these groups must be less than 120 bytes.
 */

type lightbarParamsV2Timing struct {
	/* Timing */
	googleRampUp   int32
	googleRampDown int32
	s3s0RampUp     int32
	s0TickDelay    [2]int32 /* AC=0/1 */
	s0aTickDelay   [2]int32 /* AC=0/1 */
	s0s3RampDown   int32
	s3SleepFor     int32
	s3RampUp       int32
	s3RampDown     int32
	s5RampUp       int32
	s5RampDown     int32
	tapTickDelay   int32
	tapGateDelay   int32
	tapDisplayTime int32
}

/* TYPE */
type lightbarParamsV2Tap struct {
	/* Tap-for-battery params */
	tapPctRed   uint8
	tapPctGreen uint8
	tapSegMinOn uint8
	tapSegMaxOn uint8
	tapSegOsc   uint8
	tapIdx      [3]uint8
}

/* TYPE */
type lightbarParamsV2Oscillation struct {
	/* Oscillation */
	oscMin [2]uint8 /* AC=0/1 */
	oscMax [2]uint8 /* AC=0/1 */
	wOfs   [2]uint8 /* AC=0/1 */
}

/* TYPE */
type lightbarParamsV2Brightness struct {
	/* Brightness limits based on the backlight and AC. */
	brightBlOffFixed [2]uint8 /* AC=0/1 */
	brightBlOnMin    [2]uint8 /* AC=0/1 */
	brightBlOnMax    [2]uint8 /* AC=0/1 */
}

/* TYPE */
type lightbarParamsV2Thresholds struct {
	/* Battery level thresholds */
	batteryhreshold [lbBatteryLevels - 1]uint8
}

/* TYPE */
type lightbarParamsV2Colors struct {
	/* Map [AC][batteryLevel] to color index */
	s0Idx [2][lbBatteryLevels]uint8 /* AP is running */
	s3Idx [2][lbBatteryLevels]uint8 /* AP is sleeping */

	/* s5: single color pulse on inhibited power-up */
	s5Idx uint8

	/* Color palette */
	color [8]rgbS /* 0-3 are Google colors */
}

/* Lightbyte program. */
const ecLbProgLen = 192

/* TYPE */
type lightbarProgram struct {
	size uint8
	data [ecLbProgLen]uint8
}

/* TYPE */
/* this one is messy
type ecParamsLightbar struct {
cmd uint8		      /* Command (see enum lightbarCommand)
	union {
		struct {
			/* no args
		} dump, off, on, init, getSeq, getParamsV0, getParamsV1
			version, getBrightness, getDemo, suspend, resume
			getParamsV2Timing, getParamsV2Tap
			getParamsV2Osc, getParamsV2Bright
			getParamsV2Thlds, getParamsV2Colors;

		struct {
		uint8 num;
		} setBrightness, seq, demo;

		struct {
		uint8 ctrl, reg, value;
		} reg;

		struct {
		uint8 led, red, green, blue;
		} setRgb;

		struct {
		uint8 led;
		} getRgb;

		struct {
		uint8 enable;
		} manualSuspendCtrl;

		struct lightbarParamsV0 setParamsV0;
		struct lightbarParamsV1 setParamsV1;

		struct lightbarParamsV2Timing setV2parTiming;
		struct lightbarParamsV2Tap setV2parTap;
		struct lightbarParamsV2Oscillation setV2parOsc;
		struct lightbarParamsV2Brightness setV2parBright;
		struct lightbarParamsV2Thresholds setV2parThlds;
		struct lightbarParamsV2Colors setV2parColors;

		struct lightbarProgram setProgram;
	};
} ;
*/
/* TYPE */
/*
type ecResponseLightbar struct {
	union {
		struct {
			struct {
			uint8 reg;
			uint8 ic0;
			uint8 ic1;
			} vals[23];
		} dump;

		struct  {
		uint8 num;
		} getSeq, getBrightness, getDemo;

		struct lightbarParamsV0 getParamsV0;
		struct lightbarParamsV1 getParamsV1;


		struct lightbarParamsV2Timing getParamsV2Timing;
		struct lightbarParamsV2Tap getParamsV2Tap;
		struct lightbarParamsV2Oscillation getParamsV2Osc;
		struct lightbarParamsV2Brightness getParamsV2Bright;
		struct lightbarParamsV2Thresholds getParamsV2Thlds;
		struct lightbarParamsV2Colors getParamsV2Colors;

		struct {
		uint32 num;
		uint32 flags;
		} version;

		struct {
		uint8 red, green, blue;
		} getRgb;

		struct {
			/* no return params *
		} off, on, init, setBrightness, seq, reg, setRgb
			demo, setParamsV0, setParamsV1
			setProgram, manualSuspendCtrl, suspend, resume
			setV2parTiming, setV2parTap
			setV2parOsc, setV2parBright, setV2parThlds
			setV2parColors;
	};
} ;
*/
/* Lightbar commands */
type lightbarCommand uint8

const (
	lightbarCmdDump                   lightbarCommand = 0
	lightbarCmdOff                                    = 1
	lightbarCmdOn                                     = 2
	lightbarCmdInit                                   = 3
	lightbarCmdSetBrightness                          = 4
	lightbarCmdSeq                                    = 5
	lightbarCmdReg                                    = 6
	lightbarCmdSetRgb                                 = 7
	lightbarCmdGetSeq                                 = 8
	lightbarCmdDemo                                   = 9
	lightbarCmdGetParamsV0                            = 10
	lightbarCmdSetParamsV0                            = 11
	lightbarCmdVersion                                = 12
	lightbarCmdGetBrightness                          = 13
	lightbarCmdGetRgb                                 = 14
	lightbarCmdGetDemo                                = 15
	lightbarCmdGetParamsV1                            = 16
	lightbarCmdSetParamsV1                            = 17
	lightbarCmdSetProgram                             = 18
	lightbarCmdManualSuspendCtrl                      = 19
	lightbarCmdSuspend                                = 20
	lightbarCmdResume                                 = 21
	lightbarCmdGetParamsV2Timing                      = 22
	lightbarCmdSetParamsV2Timing                      = 23
	lightbarCmdGetParamsV2Tap                         = 24
	lightbarCmdSetParamsV2Tap                         = 25
	lightbarCmdGetParamsV2Oscillation                 = 26
	lightbarCmdSetParamsV2Oscillation                 = 27
	lightbarCmdGetParamsV2Brightness                  = 28
	lightbarCmdSetParamsV2Brightness                  = 29
	lightbarCmdGetParamsV2Thresholds                  = 30
	lightbarCmdSetParamsV2Thresholds                  = 31
	lightbarCmdGetParamsV2Colors                      = 32
	lightbarCmdSetParamsV2Colors                      = 33
	lightbarNumCmds
)

/*****************************************************************************/
/* LED control commands */

const ecCmdLedControl = 0x29

type ecLedID uint8

const (
	/* LED to indicate battery state of charge */
	ecLedIDBatteryLed ecLedID = 0
	/*
	 * LED to indicate system power state (on or in suspend).
	 * May be on power button or on C-panel.
	 */
	ecLedIDPowerLed
	/* LED on power adapter or its plug */
	ecLedIDAdapterLed

	ecLedIDCount
)

const (
	/* LED control flags */
	ecLedFlagsQuery = (1 << iota) /* Query LED capability only */
	ecLedFlagsAuto                /* Switch LED back to automatic control */
)

type ecLedColors uint8

const (
	ecLedColorRed ecLedColors = 0
	ecLedColorGreen
	ecLedColorBlue
	ecLedColorYellow
	ecLedColorWhite

	ecLedColorCount
)

/* TYPE */
type ecParamsLedControl struct {
	ledID uint8 /* Which LED to control */
	flags uint8 /* Control flags */

	brightness [ecLedColorCount]uint8
}

/* TYPE */
type ecResponseLedControl struct {
	/*
	 * Available brightness value range.
	 *
	 * Range 0 means color channel not present.
	 * Range 1 means on/off control.
	 * Other values means the LED is control by PWM.
	 */
	brightnessRange [ecLedColorCount]uint8
}

/*****************************************************************************/
/* Verified boot commands */

/*
 * Note: command code 0x29 version 0 was VBOOTCmd in Link EVT; it may be
 * reused for other purposes with version > 0.
 */

/* Verified boot hash command */
const ecCmdVbootHash = 0x2A

/* TYPE */
type ecParamsVbootHash struct {
	cmd       uint8     /* enum ecVbootHashCmd */
	hashype   uint8     /* enum ecVbootHashType */
	nonceSize uint8     /* Nonce size; may be 0 */
	reserved0 uint8     /* Reserved; set 0 */
	offset    uint32    /* Offset in flash to hash */
	size      uint32    /* Number of bytes to hash */
	nonceData [64]uint8 /* Nonce data; ignored if nonceSize=0 */
}

/* TYPE */
type ecResponseVbootHash struct {
	status     uint8     /* enum ecVbootHashStatus */
	hashype    uint8     /* enum ecVbootHashType */
	digestSize uint8     /* Size of hash digest in bytes */
	reserved0  uint8     /* Ignore; will be 0 */
	offset     uint32    /* Offset in flash which was hashed */
	size       uint32    /* Number of bytes hashed */
	hashDigest [64]uint8 /* Hash digest data */
}

type ecVbootHashCmd uint8

const (
	ecVbootHashGet    ecVbootHashCmd = 0 /* Get current hash status */
	ecVbootHashAbort  ecVbootHashCmd = 1 /* Abort calculating current hash */
	ecVbootHashStart  ecVbootHashCmd = 2 /* Start computing a new hash */
	ecVbootHashRecalc ecVbootHashCmd = 3 /* Synchronously compute a new hash */
)

type ecVbootHashType uint8

const (
	ecVbootHashTypeSha256 ecVbootHashType = 0 /* SHA-256 */
)

type ecVbootHashStatus uint8

const (
	ecVbootHashStatusNone ecVbootHashStatus = 0 /* No hash (not started, or aborted) */
	ecVbootHashStatusDone ecVbootHashStatus = 1 /* Finished computing a hash */
	ecVbootHashStatusBusy ecVbootHashStatus = 2 /* Busy computing a hash */
)

/*
 * Special values for offset for ecVbootHashStart and ecVbootHashRecalc.
 * If one of these is specified, the EC will automatically update offset and
 * size to the correct values for the specified image (RO or RW).
 */
const (
	ecVbootHashOffsetRo = 0xfffffffe
	ecVbootHashOffsetRw = 0xfffffffd
)

/*****************************************************************************/
/*
 * Motion sense commands. We'll make separate structs for sub-commands with
 * different input args, so that we know how much to expect.
 */
const ecCmdMotionSenseCmd = 0x2B

/* Motion sense commands */
type motionsenseCommand uint8

const (
	/* Dump command returns all motion sensor data including motion sense
	 * module flags and individual sensor flags.
	 */
	motionsenseCmdDump motionsenseCommand = iota

	/*
	 * Info command returns data describing the details of a given sensor
	 * including enum motionsensorType, enum motionsensorLocation, and
	 * enum motionsensorChip.
	 */
	motionsenseCmdInfo

	/*
	 * EC Rate command is a setter/getter command for the EC sampling rate
	 * of all motion sensors in milliseconds.
	 */
	motionsenseCmdEcRate

	/*
	 * Sensor ODR command is a setter/getter command for the output data
	 * rate of a specific motion sensor in millihertz.
	 */
	motionsenseCmdSensorOdr

	/*
	 * Sensor range command is a setter/getter command for the range of
	 * a specified motion sensor in +/-G's or +/- deg/s.
	 */
	motionsenseCmdSensorRange

	/*
	 * Setter/getter command for the keyboard wake angle. When the lid
	 * angle is greater than this value, keyboard wake is disabled in S3
	 * and when the lid angle goes less than this value, keyboard wake is
	 * enabled. Note, the lid angle measurement is an approximate
	 * un-calibrated value, hence the wake angle isn't exact.
	 */
	motionsenseCmdKbWakeAngle

	/* Number of motionsense sub-commands. */
	motionsenseNumCmds
)

/* List of motion sensor types. */
type motionsensorType uint8

const (
	motionsenseTypeAccel motionsensorType = 0
	motionsenseTypeGyro  motionsensorType = 1
)

/* List of motion sensor locations. */
type motionsensorLocation uint8

const (
	motionsenseLocBase motionsensorLocation = 0
	motionsenseLocLid  motionsensorLocation = 1
)

/* List of motion sensor chips. */
type motionsensorChip uint8

const (
	motionsenseChipKxcj9   motionsensorChip = 0
	motionsenseChipLsm6ds0 motionsensorChip = 1
)

/* Module flag masks used for the dump sub-command. */
const motionsenseModuleFlagActive = (1 << 0)

/* Sensor flag masks used for the dump sub-command. */
const motionsenseSensorFlagPresent = (1 << 0)

/*
 * Send this value for the data element to only perform a read. If you
 * send any other value, the EC will interpret it as data to set and will
 * return the actual value set.
 */
const ecMotionSenseNoValue = -1

/* some other time
type ecParamsMotionSense struct {
cmd uint8
	union {
		/* Used for MOTIONSENSECmdDump * /
		struct {
			/*
			 * Maximal number of sensor the host is expecting.
			 * 0 means the host is only interested in the number
			 * of sensors controlled by the EC.
			 * /
		uint8 maxSensorCount;
		} dump;

		/*
		 * Used for MOTIONSENSECmdEcRate and
		 * MOTIONSENSECmdKbWakeAngle.
		 * /
		struct {
			/* Data to set or ecMotionSenseNoValue to read. * /
			data int16
		} ecRate, kbWakeAngle;

		/* Used for MOTIONSENSECmdInfo. * /
		struct {
		uint8 sensorNum;
		} info;

		/*
		 * Used for MOTIONSENSECmdSensorOdr and
		 * MOTIONSENSECmdSensorRange.
		 * /
		struct {
		uint8 sensorNum;

			/* Rounding flag, true for round-up, false for down. * /
		uint8 roundup;

		uint16 reserved;

			/* Data to set or ecMotionSenseNoValue to read. * /
			data int32
		} sensorOdr, sensorRange;
	};
} ;
*/
/* TYPE */
type ecResponseMotionSensorData struct {
	/* Flags for each sensor. */
	flags   uint8
	padding uint8

	/* Each sensor is up to 3-axis. */
	data [3]int16
}

/* TYPE */
/* some other time

type ecResponseMotionSense struct {
	union {
		/* Used for MOTIONSENSECmdDump * /
		struct {
			/* Flags representing the motion sensor module. * /
		uint8 moduleFlags;

			/* Number of sensors managed directly by the EC * /
		uint8 sensorCount;

			/*
			 * sensor data is truncated if responseMax is too small
			 * for holding all the data.
			 * /
			struct ecResponseMotionSensorData sensor[0];
		} dump;

		/* Used for MOTIONSENSECmdInfo. * /
		struct {
			/* Should be element of enum motionsensorType. * /
		uint8 type;

			/* Should be element of enum motionsensorLocation. * /
		uint8 location;

			/* Should be element of enum motionsensorChip. * /
		uint8 chip;
		} info;

		/*
		 * Used for MOTIONSENSECmdEcRate, MOTIONSENSECmdSensorOdr
		 * MOTIONSENSECmdSensorRange, and
		 * MOTIONSENSECmdKbWakeAngle.
		 * /
		struct {
			/* Current value of the parameter queried. * /
			ret int32
		} ecRate, sensorOdr, sensorRange, kbWakeAngle;
	};
} ;
*/
/*****************************************************************************/
/* Force lid open command */

/* Make lid event always open */
const ecCmdForceLidOpen = 0x2c

/* TYPE */
type ecParamsForceLidOpen struct {
	enabled uint8
}

/*****************************************************************************/
/* USB charging control commands */

/* Set USB port charging mode */
const ecCmdUsbChargeSetMode = 0x30

/* TYPE */
type ecParamsUsbChargeSetMode struct {
	usbPortID uint8
	mode      uint8
}

/*****************************************************************************/
/* Persistent storage for host */

/* Maximum bytes that can be read/written in a single command */
const ecPstoreSizeMax = 64

/* Get persistent storage info */
const ecCmdPstoreInfo = 0x40

/* TYPE */
type ecResponsePstoreInfo struct {
	/* Persistent storage size, in bytes */
	pstoreSize uint32
	/* Access size; read/write offset and size must be a multiple of this */
	accessSize uint32
}

/*
 * Read persistent storage
 *
 * Response is params.size bytes of data.
 */
const ecCmdPstoreRead = 0x41

/* TYPE */
type ecParamsPstoreRead struct {
	offset uint32 /* Byte offset to read */
	size   uint32 /* Size to read in bytes */
}

/* Write persistent storage */
const ecCmdPstoreWrite = 0x42

/* TYPE */
type ecParamsPstoreWrite struct {
	offset uint32 /* Byte offset to write */
	size   uint32 /* Size to write in bytes */
	data   [ecPstoreSizeMax]uint8
}

/* TYPE */
/*****************************************************************************/
/* Real-time clock */

/* RTC params and response structures */
type ecParamsRtc struct {
	time uint32
}

/* TYPE */
type ecResponseRtc struct {
	time uint32
}

/* These use ecResponseRtc */
const (
	ecCmdRtcGetValue = 0x44
	ecCmdRtcGetAlarm = 0x45
)

/* These all use ecParamsRtc */
const (
	ecCmdRtcSetValue = 0x46
	ecCmdRtcSetAlarm = 0x47
)

/*****************************************************************************/
/* Port80 log access */

/* Maximum entries that can be read/written in a single command */
const ecPort80SizeMax = 32

/* Get last port80 code from previous boot */
const (
	ecCmdPort80LastBoot = 0x48
	ecCmdPort80Read     = 0x48
)

type ecPort80Subcmd uint8

const (
	ecPort80GetInfo ecPort80Subcmd = 0
	ecPort80ReadBuffer
)

/* TYPE */
type ecParamsPort80Read struct {
	subcmd     uint16
	offset     uint32
	numEntries uint32
}

/* TYPE */
type ecResponsePort80Read struct {
	/*
		struct {
		uint32 writes;
		uint32 historySize;
		uint32 lastBoot;
		} getInfo;*/

	codes [ecPort80SizeMax]uint16
}

/* TYPE */
type ecResponsePort80LastBoot struct {
	code uint16
}

/*****************************************************************************/
/* Thermal engine commands. Note that there are two implementations. We'll
 * reuse the command number, but the data and behavior is incompatible.
 * Version 0 is what originally shipped on Link.
 * Version 1 separates the CPU thermal limits from the fan control.
 */

const (
	ecCmdThermalSetThreshold = 0x50
	ecCmdThermalGetThreshold = 0x51
)

/* TYPE */
/* The version 0 structs are opaque. You have to know what they are for
 * the get/set commands to make any sense.
 */

/* Version 0 - set */
type ecParamsThermalSetThreshold struct {
	sensorype   uint8
	thresholdID uint8
	value       uint16
}

/* TYPE */
/* Version 0 - get */
type ecParamsThermalGetThreshold struct {
	sensorype   uint8
	thresholdID uint8
}

/* TYPE */
type ecResponseThermalGetThreshold struct {
	value uint16
}

/* The version 1 structs are visible. */
type ecTempThresholds uint8

const (
	ecTempThreshWarn ecTempThresholds = 0
	ecTempThreshHigh
	ecTempThreshHalt

	ecTempThreshCount
)

/* TYPE */
/* Thermal configuration for one temperature sensor. Temps are in degrees K.
 * Zero values will be silently ignored by the thermal task.
 */
type ecThermalConfig struct {
	tempHost   [ecTempThreshCount]uint32 /* levels of hotness */
	tempFanOff uint32                    /* no active cooling needed */
	tempFanMax uint32                    /* max active cooling needed */
}

/* TYPE */
/* Version 1 - get config for one sensor. */
type ecParamsThermalGetThresholdV1 struct {
	sensorNum uint32
}

/* TYPE */
/* This returns a struct ecThermalConfig */

/* Version 1 - set config for one sensor.
 * Use read-modify-write for best results! */
type ecParamsThermalSetThresholdV1 struct {
	sensorNum uint32
	cfg       ecThermalConfig
}

/* This returns no data */

/****************************************************************************/

/* Toggle automatic fan control */
const ecCmdThermalAutoFanCtrl = 0x52

/* TYPE */
/* Version 1 of input params */
type ecParamsAutoFanCtrlV1 struct {
	fanIdx uint8
}

/* Get/Set TMP006 calibration data */
const (
	ecCmdTmp006GetCalibration = 0x53
	ecCmdTmp006SetCalibration = 0x54
)

/* TYPE */
/*
 * The original TMP006 calibration only needed four params, but now we need
 * more. Since the algorithm is nothing but magic numbers anyway, we'll leave
 * the params opaque. The v1 "get" response will include the algorithm number
 * and how many params it requires. That way we can change the EC code without
 * needing to update this file. We can also use a different algorithm on each
 * sensor.
 */

/* This is the same struct for both v0 and v1. */
type ecParamsTmp006GetCalibration struct {
	index uint8
}

/* TYPE */
/* Version 0 */
type ecResponseTmp006GetCalibrationV0 struct {
	s0, b0, b1, bw float32
}

/* TYPE */
type ecParamsTmp006SetCalibrationV0 struct {
	index          uint8
	reserved       [3]uint8
	s0, b0, b1, b2 float32
}

/* TYPE */
/* Version 1 */
type ecResponseTmp006GetCalibrationV1 struct {
	algorithm uint8
	numParams uint8
	reserved  [2]uint8
	val       []float32
}

/* TYPE */
type ecParamsTmp006SetCalibrationV1 struct {
	index     uint8
	algorithm uint8
	numParams uint8
	reserved  uint8
	val       []float32
}

/* Read raw TMP006 data */
const ecCmdTmp006GetRaw = 0x55

/* TYPE */
type ecParamsTmp006GetRaw struct {
	index uint8
}

/* TYPE */
type ecResponseTmp006GetRaw struct {
	t int32 /* In 1/100 K */
	v int32 /* In nV */
}

/*****************************************************************************/
/* MKBP - Matrix KeyBoard Protocol */

/*
 * Read key state
 *
 * Returns raw data for keyboard cols; see ecResponseMkbpInfo.cols for
 * expected response size.
 */
const ecCmdMkbpState = 0x60

/* Provide information about the matrix : number of rows and columns */
const ecCmdMkbpInfo = 0x61

/* TYPE */
type ecResponseMkbpInfo struct {
	rows     uint32
	cols     uint32
	switches uint8
}

/* Simulate key press */
const ecCmdMkbpSimulateKey = 0x62

/* TYPE */
type ecParamsMkbpSimulateKey struct {
	col     uint8
	row     uint8
	pressed uint8
}

/* Configure keyboard scanning */
const (
	ecCmdMkbpSetConfig = 0x64
	ecCmdMkbpGetConfig = 0x65
)

/* flags */
type mkbpConfigFlags uint8

const (
	ecMkbpFlagsEnable mkbpConfigFlags = 1 /* Enable keyboard scanning */
)

type mkbpConfigValid uint8

const (
	ecMkbpValidScanPeriod       mkbpConfigValid = 1 << 0
	ecMkbpValidPollTimeout                      = 1 << 1
	ecMkbpValidMinPostScanDelay                 = 1 << 3
	ecMkbpValidOutputSettle                     = 1 << 4
	ecMkbpValidDebounceDown                     = 1 << 5
	ecMkbpValidDebounceUp                       = 1 << 6
	ecMkbpValidFifoMaxDepth                     = 1 << 7
)

/* TYPE */
/* Configuration for our key scanning algorithm */
type ecMkbpConfig struct {
	validMask    uint32 /* valid fields */
	flags        uint8  /* some flags (enum mkbpConfigFlags) */
	validFlags   uint8  /* which flags are valid */
	scanPeriodUs uint16 /* period between start of scans */
	/* revert to interrupt mode after no activity for this long */
	pollimeoutUs uint32
	/*
	 * minimum post-scan relax time. Once we finish a scan we check
	 * the time until we are due to start the next one. If this time is
	 * shorter this field, we use this instead.
	 */
	minPostScanDelayUs uint16
	/* delay between setting up output and waiting for it to settle */
	outputSettleUs uint16
	debounceDownUs uint16 /* time for debounce on key down */
	debounceUpUs   uint16 /* time for debounce on key up */
	/* maximum depth to allow for fifo (0 = no keyscan output) */
	fifoMaxDepth uint8
}

/* TYPE */
type ecParamsMkbpSetConfig struct {
	config ecMkbpConfig
}

/* TYPE */
type ecResponseMkbpGetConfig struct {
	config ecMkbpConfig
}

/* Run the key scan emulation */
const ecCmdKeyscanSeqCtrl = 0x66

type ecKeyscanSeqCmd uint8

const (
	ecKeyscanSeqStatus  ecKeyscanSeqCmd = 0 /* Get status information */
	ecKeyscanSeqClear   ecKeyscanSeqCmd = 1 /* Clear sequence */
	ecKeyscanSeqAdd     ecKeyscanSeqCmd = 2 /* Add item to sequence */
	ecKeyscanSeqStart   ecKeyscanSeqCmd = 3 /* Start running sequence */
	ecKeyscanSeqCollect ecKeyscanSeqCmd = 4 /* Collect sequence summary data */
)

type ecCollectFlags uint8

const (
	/* Indicates this scan was processed by the EC. Due to timing, some
	 * scans may be skipped.
	 */
	ecKeyscanSeqFlagDone ecCollectFlags = 1 << iota
)

/* TYPE */
type ecCollectItem struct {
	flags uint8 /* some flags (enum ecCollectFlags) */
}

/* TYPE */
/* later
type ecParamsKeyscanSeqCtrl struct {
cmd uint8	/* Command to send (enum ecKeyscanSeqCmd) * /
	union {
		struct {
		uint8 active;		/* still active * /
		uint8 numItems;	/* number of items * /
			/* Current item being presented * /
		uint8 curItem;
		} status;
		struct {
			/*
			 * Absolute time for this scan, measured from the
			 * start of the sequence.
			 * /
		uint32 timeUs;
		uint8 scan[0];	/* keyscan data * /
		} add;
		struct {
		uint8 startItem;	/* First item to return * /
		uint8 numItems;	/* Number of items to return * /
		} collect;
	};
} ;
*/
/* TYPE */
/* lter
type ecResultKeyscanSeqCtrl struct {
	union {
		struct {
		uint8 numItems;	/* Number of items *
			/* Data for each item *
			struct ecCollectItem item[0]
		} collect;
	};
} ;
*/
/*
 * Get the next pending MKBP event.
 *
 * Returns ecResUnavailable if there is no event pending.
 */
const ecCmdGetNextEvent = 0x67

type ecMkbpEvent uint8

const (
	/* Keyboard matrix changed. The event data is the new matrix state. */
	ecMkbpEventKeyMatrix = iota

	/* New host event. The event data is 4 bytes of host event flags. */
	ecMkbpEventHostEvent

	/* Number of MKBP events */
	ecMkbpEventCount
)

/* TYPE */
type ecResponseGetNextEvent struct {
	eventype uint8
	/* Followed by event data if any */
}

/*****************************************************************************/
/* Temperature sensor commands */

/* Read temperature sensor info */
const ecCmdTempSensorGetInfo = 0x70

/* TYPE */
type ecParamsTempSensorGetInfo struct {
	id uint8
}

/* TYPE */
type ecResponseTempSensorGetInfo struct {
	sensorName [32]byte
	sensorype  uint8
}

/* TYPE */
/*****************************************************************************/

/*
 * Note: host commands 0x80 - 0x87 are reserved to avoid conflict with ACPI
 * commands accidentally sent to the wrong interface.  See the ACPI section
 * below.
 */

/*****************************************************************************/
/* Host event commands */

/*
 * Host event mask params and response structures, shared by all of the host
 * event commands below.
 */
type ecParamsHostEventMask struct {
	mask uint32
}

/* TYPE */
type ecResponseHostEventMask struct {
	mask uint32
}

/* These all use ecResponseHostEventMask */
const (
	ecCmdHostEventGetB        = 0x87
	ecCmdHostEventGetSmiMask  = 0x88
	ecCmdHostEventGetSciMask  = 0x89
	ecCmdHostEventGetWakeMask = 0x8d
)

/* These all use ecParamsHostEventMask */
const (
	ecCmdHostEventSetSmiMask  = 0x8a
	ecCmdHostEventSetSciMask  = 0x8b
	ecCmdHostEventClear       = 0x8c
	ecCmdHostEventSetWakeMask = 0x8e
	ecCmdHostEventClearB      = 0x8f
)

/*****************************************************************************/
/* Switch commands */

/* Enable/disable LCD backlight */
const ecCmdSwitchEnableBklight = 0x90

/* TYPE */
type ecParamsSwitchEnableBacklight struct {
	enabled uint8
}

/* Enable/disable WLAN/Bluetooth */
const (
	ecCmdSwitchEnableWireless = 0x91
	ecVerSwitchEnableWireless = 1
)

/* TYPE */
/* Version 0 params; no response */
type ecParamsSwitchEnableWirelessV0 struct {
	enabled uint8
}

/* TYPE */
/* Version 1 params */
type ecParamsSwitchEnableWirelessV1 struct {
	/* Flags to enable now */
	nowFlags uint8

	/* Which flags to copy from nowFlags */
	nowMask uint8

	/*
	 * Flags to leave enabled in S3, if they're on at the S0->S3
	 * transition.  (Other flags will be disabled by the S0->S3
	 * transition.)
	 */
	suspendFlags uint8

	/* Which flags to copy from suspendFlags */
	suspendMask uint8
}

/* TYPE */
/* Version 1 response */
type ecResponseSwitchEnableWirelessV1 struct {
	/* Flags to enable now */
	nowFlags uint8

	/* Flags to leave enabled in S3 */
	suspendFlags uint8
}

/*****************************************************************************/
/* GPIO commands. Only available on EC if write protect has been disabled. */

/* Set GPIO output value */
const ecCmdGpioSet = 0x92

/* TYPE */
type ecParamsGpioSet struct {
	name [32]byte
	val  uint8
}

/* Get GPIO value */
const ecCmdGpioGet = 0x93

/* TYPE */
/* Version 0 of input params and response */
type ecParamsGpioGet struct {
	name [32]byte
}

/* TYPE */
type ecResponseGpioGet struct {
	val uint8
}

/* TYPE */
/* Version 1 of input params and response */
type ecParamsGpioGetV1 struct {
	subcmd uint8
	data   [32]byte
}

/* TYPE */
/* later
type ecResponseGpioGetV1 struct {
	union {
		struct {
		uint8 val;
		} getValueByName, getCount;
		struct {
		uint8 val;
			char name[32];
		uint32 flags;
		} getInfo;
	};
} ;
*/
type gpioGetSubcmd uint8

const (
	ecGpioGetByName gpioGetSubcmd = 0
	ecGpioGetCount  gpioGetSubcmd = 1
	ecGpioGetInfo   gpioGetSubcmd = 2
)

/*****************************************************************************/
/* I2C commands. Only available when flash write protect is unlocked. */

/*
 * TODO(crosbug.com/p/23570): These commands are deprecated, and will be
 * removed soon.  Use ecCmdI2cXfer instead.
 */

/* Read I2C bus */
const ecCmdI2cRead = 0x94

/* TYPE */
type ecParamsI2cRead struct {
	addr     uint16 /* 8-bit address (7-bit shifted << 1) */
	readSize uint8  /* Either 8 or 16. */
	port     uint8
	offset   uint8
}

/* TYPE */
type ecResponseI2cRead struct {
	data uint16
}

/* Write I2C bus */
const ecCmdI2cWrite = 0x95

/* TYPE */
type ecParamsI2cWrite struct {
	data      uint16
	addr      uint16 /* 8-bit address (7-bit shifted << 1) */
	writeSize uint8  /* Either 8 or 16. */
	port      uint8
	offset    uint8
}

/*****************************************************************************/
/* Charge state commands. Only available when flash write protect unlocked. */

/* Force charge state machine to stop charging the battery or force it to
 * discharge the battery.
 */
const (
	ecCmdChargeControl = 0x96
	ecVerChargeControl = 1
)

type ecChargeControlMode uint8

const (
	chargeControlNormal ecChargeControlMode = 0
	chargeControlIdle
	chargeControlDischarge
)

/* TYPE */
type ecParamsChargeControl struct {
	mode uint32 /* enum chargeControlMode */
}

/*****************************************************************************/
/* Console commands. Only available when flash write protect is unlocked. */

/* Snapshot console output buffer for use by ecCmdConsoleRead. */
const ecCmdConsoleSnapshot = 0x97

/*
 * Read next chunk of data from saved snapshot.
 *
 * Response is null-terminated string.  Empty string, if there is no more
 * remaining output.
 */
const ecCmdConsoleRead = 0x98

/*****************************************************************************/

/*
 * Cut off battery power immediately or after the host has shut down.
 *
 * return ecResInvalidCommand if unsupported by a board/battery.
 *	  ecResSuccess if the command was successful.
 *	  ecResError if the cut off command failed.
 */
const ecCmdBatteryCutOff = 0x99

const ecBatteryCutoffFlagAtShutdown = (1 << 0)

/* TYPE */
type ecParamsBatteryCutoff struct {
	flags uint8
}

/*****************************************************************************/
/* USB port mux control. */

/*
 * Switch USB mux or return to automatic switching.
 */
const ecCmdUsbMux = 0x9a

/* TYPE */
type ecParamsUsbMux struct {
	mux uint8
}

/*****************************************************************************/
/* LDOs / FETs control. */

type ecLdoState uint8

const (
	ecLdoStateOff ecLdoState = 0 /* the LDO / FET is shut down */
	ecLdoStateOn  ecLdoState = 1 /* the LDO / FET is ON / providing power */
)

/*
 * Switch on/off a LDO.
 */
const ecCmdLdoSet = 0x9b

/* TYPE */
type ecParamsLdoSet struct {
	index uint8
	state uint8
}

/*
 * Get LDO state.
 */
const ecCmdLdoGet = 0x9c

/* TYPE */
type ecParamsLdoGet struct {
	index uint8
}

/* TYPE */
type ecResponseLdoGet struct {
	state uint8
}

/*****************************************************************************/
/* Power info. */

/*
 * Get power info.
 */
const ecCmdPowerInfo = 0x9d

/* TYPE */
type ecResponsePowerInfo struct {
	usbDevype       uint32
	voltageAc       uint16
	voltageSystem   uint16
	currentSystem   uint16
	usbCurrentLimit uint16
}

/*****************************************************************************/
/* I2C passthru command */

const ecCmdI2cPassthru = 0x9e

/* Read data; if not present, message is a write */
const ecI2cFlagRead = (1 << 15)

/* Mask for address */
const ecI2cAddrMask = 0x3ff

const (
	ecI2cStatusNak     = (1 << 0) /* Transfer was not acknowledged */
	ecI2cStatusTimeout = (1 << 1) /* Timeout during transfer */
)

/* Any error */
const ecI2cStatusError = (ecI2cStatusNak | ecI2cStatusTimeout)

/* TYPE */
type ecParamsI2cPassthruMsg struct {
	addrFlags uint16 /* I2C slave address (7 or 10 bits) and flags */
	len       uint16 /* Number of bytes to read or write */
}

/* TYPE */
type ecParamsI2cPassthru struct {
	port    uint8 /* I2C port number */
	numMsgs uint8 /* Number of messages */
	msg     []ecParamsI2cPassthruMsg
	/* Data to write for all messages is concatenated here */
}

/* TYPE */
type ecResponseI2cPassthru struct {
	i2cStatus uint8   /* Status flags (ecI2cStatus_...) */
	numMsgs   uint8   /* Number of messages processed */
	data      []uint8 /* Data read by messages concatenated here */
}

/*****************************************************************************/
/* Power button hang detect */

const ecCmdHangDetect = 0x9f

/* Reasons to start hang detection timer */
/* Power button pressed */
const ecHangStartOnPowerPress = (1 << 0)

/* Lid closed */
const ecHangStartOnLidClose = (1 << 1)

/* Lid opened */
const ecHangStartOnLidOpen = (1 << 2)

/* Start of AP S3->S0 transition (booting or resuming from suspend) */
const ecHangStartOnResume = (1 << 3)

/* Reasons to cancel hang detection */

/* Power button released */
const ecHangStopOnPowerRelease = (1 << 8)

/* Any host command from AP received */
const ecHangStopOnHostCommand = (1 << 9)

/* Stop on end of AP S0->S3 transition (suspending or shutting down) */
const ecHangStopOnSuspend = (1 << 10)

/*
 * If this flag is set, all the other fields are ignored, and the hang detect
 * timer is started.  This provides the AP a way to start the hang timer
 * without reconfiguring any of the other hang detect settings.  Note that
 * you must previously have configured the timeouts.
 */
const ecHangStartNow = (1 << 30)

/*
 * If this flag is set, all the other fields are ignored (including
 * ecHangStartNow).  This provides the AP a way to stop the hang timer
 * without reconfiguring any of the other hang detect settings.
 */
const ecHangStopNow = (1 << 31)

/* TYPE */
type ecParamsHangDetect struct {
	/* Flags; see ecHang_* */
	flags uint32

	/* Timeout in msec before generating host event, if enabled */
	hostEventimeoutMsec uint16

	/* Timeout in msec before generating warm reboot, if enabled */
	warmRebootimeoutMsec uint16
}

/*****************************************************************************/
/* Commands for battery charging */

/*
 * This is the single catch-all host command to exchange data regarding the
 * charge state machine (v2 and up).
 */
const ecCmdChargeState = 0xa0

/* Subcommands for this host command */
type chargeStateCommand uint8

const (
	chargeStateCmdGetState chargeStateCommand = iota
	chargeStateCmdGetParam
	chargeStateCmdSetParam
	chargeStateNumCmds
)

/*
 * Known param numbers are defined here. Ranges are reserved for board-specific
 * params, which are handled by the particular implementations.
 */
type chargeStateParams uint8

const (
	csParamChgVoltage      chargeStateParams = iota /* charger voltage limit */
	csParamChgCurrent                               /* charger current limit */
	csParamChgInputCurrent                          /* charger input current limit */
	csParamChgStatus                                /* charger-specific status */
	csParamChgOption                                /* charger-specific options */
	/* How many so far? */
	csNumBaseParams

	/* Range for CONFIGChargerProfileOverride params */
	csParamCustomProfileMin = 0x10000
	csParamCustomProfileMax = 0x1ffff

	/* Other custom param ranges go here... */
)

/* TYPE */
/* ler
type ecParamsChargeState struct {
cmd uint8				/* enum chargeStateCommand * /
	union {
		struct {
			/* no args * /
		} getState;

		struct {
		uint32 param;		/* enum chargeStateParam * /
		} getParam;

		struct {
		uint32 param;		/* param to set * /
		uint32 value;		/* value to set * /
		} setParam;
	};
} ;

/* TYPE */
/* later
type ecResponseChargeState struct {
	union {
		struct {
			int ac;
			int chgVoltage;
			int chgCurrent;
			int chgInputCurrent;
			int battStateOfCharge;
		} getState;

		struct {
		uint32 value;
		} getParam;
		struct {
			/* no return values *
		} setParam;
	};
} ;
*/

/*
 * Set maximum battery charging current.
 */
const ecCmdChargeCurrentLimit = 0xa1

/* TYPE */
type ecParamsCurrentLimit struct {
	limit uint32 /* in mA */
}

/*
 * Set maximum external power current.
 */
const ecCmdExtPowerCurrentLimit = 0xa2

/* TYPE */
type ecParamsExtPowerCurrentLimit struct {
	limit uint32 /* in mA */
}

/*****************************************************************************/
/* Smart battery pass-through */

/* Get / Set 16-bit smart battery registers */
const (
	ecCmdSbReadWord  = 0xb0
	ecCmdSbWriteWord = 0xb1
)

/* Get / Set string smart battery parameters
 * formatted as SMBUS "block".
 */
const (
	ecCmdSbReadBlock  = 0xb2
	ecCmdSbWriteBlock = 0xb3
)

/* TYPE */
type ecParamsSbRd struct {
	reg uint8
}

/* TYPE */
type ecResponseSbRdWord struct {
	value uint16
}

/* TYPE */
type ecParamsSbWrWord struct {
	reg   uint8
	value uint16
}

/* TYPE */
type ecResponseSbRdBlock struct {
	data [32]uint8
}

/* TYPE */
type ecParamsSbWrBlock struct {
	reg  uint8
	data [32]uint16
}

/*****************************************************************************/
/* Battery vendor parameters
 *
 * Get or set vendor-specific parameters in the battery. Implementations may
 * differ between boards or batteries. On a set operation, the response
 * contains the actual value set, which may be rounded or clipped from the
 * requested value.
 */

const ecCmdBatteryVendorParam = 0xb4

type ecBatteryVendorParamMode uint8

const (
	batteryVendorParamModeGet ecBatteryVendorParamMode = 0
	batteryVendorParamModeSet
)

/* TYPE */
type ecParamsBatteryVendorParam struct {
	param uint32
	value uint32
	mode  uint8
}

/* TYPE */
type ecResponseBatteryVendorParam struct {
	value uint32
}

/*****************************************************************************/
/*
 * Smart Battery Firmware Update Commands
 */
const ecCmdSbFwUpdate = 0xb5

type ecSbFwUpdateSubcmd uint8

const (
	ecSbFwUpdatePrepare ecSbFwUpdateSubcmd = 0x0
	ecSbFwUpdateInfo    ecSbFwUpdateSubcmd = 0x1 /*query sb info */
	ecSbFwUpdateBegin   ecSbFwUpdateSubcmd = 0x2 /*check if protected */
	ecSbFwUpdateWrite   ecSbFwUpdateSubcmd = 0x3 /*check if protected */
	ecSbFwUpdateEnd     ecSbFwUpdateSubcmd = 0x4
	ecSbFwUpdateStatus  ecSbFwUpdateSubcmd = 0x5
	ecSbFwUpdateProtect ecSbFwUpdateSubcmd = 0x6
	ecSbFwUpdateMax     ecSbFwUpdateSubcmd = 0x7
)

const (
	sbFwUpdateCmdWriteBlockSize = 32
	sbFwUpdateCmdStatusSize     = 2
	sbFwUpdateCmdInfoSize       = 8
)

/* TYPE */
type ecSbFwUpdateHeader struct {
	subcmd uint16 /* enum ecSbFwUpdateSubcmd */
	fwID   uint16 /* firmware id */
}

/* TYPE */
type ecParamsSbFwUpdate struct {
	hdr ecSbFwUpdateHeader
	/* no args. */
	/* ecSbFwUpdatePrepare  = 0x0 */
	/* ecSbFwUpdateInfo     = 0x1 */
	/* ecSbFwUpdateBegin    = 0x2 */
	/* ecSbFwUpdateEnd      = 0x4 */
	/* ecSbFwUpdateStatus   = 0x5 */
	/* ecSbFwUpdateProtect  = 0x6 */
	/* or ... */
	/* ecSbFwUpdateWrite    = 0x3 */
	data [sbFwUpdateCmdWriteBlockSize]uint8
}

/* TYPE */
type ecResponseSbFwUpdate struct {
	data []uint8
	/* ecSbFwUpdateInfo     = 0x1 */
	//uint8 data[SBFwUpdateCmdInfoSize];
	/* ecSbFwUpdateStatus   = 0x5 */
	//uint8 data[SBFwUpdateCmdStatusSize];
}

/*
 * Entering Verified Boot Mode Command
 * Default mode is VBOOTModeNormal if EC did not receive this command.
 * Valid Modes are: normal, developer, and recovery.
 */
const ecCmdEnteringMode = 0xb6

/* TYPE */
type ecParamsEnteringMode struct {
	vbootMode int
}

const (
	vbootModeNormal    = 0
	vbootModeDeveloper = 1
	vbootModeRecovery  = 2
)

/*****************************************************************************/
/* System commands */

/*
 * TODO(crosbug.com/p/23747): This is a confusing name, since it doesn't
 * necessarily reboot the EC.  Rename to "image" or something similar?
 */
const ecCmdRebootEc = 0xd2

/* Command */
type ecRebootCmd uint8

const (
	ecRebootCancel ecRebootCmd = 0 /* Cancel a pending reboot */
	ecRebootJumpRo ecRebootCmd = 1 /* Jump to RO without rebooting */
	ecRebootJumpRw ecRebootCmd = 2 /* Jump to RW without rebooting */
	/* (command 3 was jump to RW-B) */
	ecRebootCold        ecRebootCmd = 4 /* Cold-reboot */
	ecRebootDisableJump ecRebootCmd = 5 /* Disable jump until next reboot */
	ecRebootHibernate   ecRebootCmd = 6 /* Hibernate EC */
)

/* Flags for ecParamsRebootEc.rebootFlags */
const (
	ecRebootFlagReserved0    = (1 << 0) /* Was recovery request */
	ecRebootFlagOnApShutdown = (1 << 1) /* Reboot after AP shutdown */)

/* TYPE */
type ecParamsRebootEc struct {
	cmd   uint8 /* enum ecRebootCmd */
	flags uint8 /* See ecRebootFlag_* */
}

/*
 * Get information on last EC panic.
 *
 * Returns variable-length platform-dependent panic information.  See panic.h
 * for details.
 */
const ecCmdGetPanicInfo = 0xd3

/*****************************************************************************/
/*
 * Special commands
 *
 * These do not follow the normal rules for commands.  See each command for
 * details.
 */

/*
 * Reboot NOW
 *
 * This command will work even when the EC LPC interface is busy, because the
 * reboot command is processed at interrupt level.  Note that when the EC
 * reboots, the host will reboot too, so there is no response to this command.
 *
 * Use ecCmdRebootEc to reboot the EC more politely.
 */
const ecCmdReboot = 0xd1 /* Think "die" */

/*
 * Resend last response (not supported on LPC).
 *
 * Returns ecResUnavailable if there is no response available - for example
 * there was no previous command, or the previous command's response was too
 * big to save.
 */
const ecCmdResendResponse = 0xdb

/*
 * This header byte on a command indicate version 0. Any header byte less
 * than this means that we are talking to an old EC which doesn't support
 * versioning. In that case, we assume version 0.
 *
 * Header bytes greater than this indicate a later version. For example
 * ecCmdVersion0 + 1 means we are using version 1.
 *
 * The old EC interface must not use commands 0xdc or higher.
 */
const ecCmdVersion0 = 0xdc

/*****************************************************************************/
/*
 * PD commands
 *
 * These commands are for PD MCU communication.
 */

/* EC to PD MCU exchange status command */
const ecCmdPdExchangeStatus = 0x100

type pdChargeState uint8

const (
	pdChargeNoChange pdChargeState = 0 /* Don't change charge state */
	pdChargeNone                       /* No charging allowed */
	pdCharge5v                         /* 5V charging only */
	pdChargeMax                        /* Charge at max voltage */
)

/* TYPE */
/* Status of EC being sent to PD */
type ecParamsPdStatus struct {
	battSoc     int8  /* battery state of charge */
	chargeState uint8 /* charging state (from enum pdChargeState) */
}

/* Status of PD being sent back to EC */
const (
	pdStatusHostEvent     = (1 << 0) /* Forward host event to AP */
	pdStatusInRw          = (1 << 1) /* Running RW image */
	pdStatusJumpedToImage = (1 << 2) /* Current image was jumped to */
)

/* TYPE */
type ecResponsePdStatus struct {
	status           uint32 /* PD MCU status */
	currLimMa        uint32 /* input current limit */
	activeChargePort int32  /* active charging port */
}

/* AP to PD MCU host event status command, cleared on read */
const ecCmdPdHostEventStatus = 0x104

/* PD MCU host event status bits */
const (
	pdEventUpdateDevice     = (1 << 0)
	pdEventPowerChange      = (1 << 1)
	pdEventIdentityReceived = (1 << 2)
)

/* TYPE */
type ecResponseHostEventStatus struct {
	status uint32 /* PD MCU host event status */
}

/* Set USB type-C port role and muxes */
const ecCmdUsbPdControl = 0x101

type usbPdControlRole uint8

const (
	usbPdCtrlRoleNoChange    usbPdControlRole = 0
	usbPdCtrlRoleToggleOn                     = 1 /* == AUTO */
	usbPdCtrlRoleToggleOff                    = 2
	usbPdCtrlRoleForceSink                    = 3
	usbPdCtrlRoleForceSource                  = 4
	usbPdCtrlRoleCount
)

type usbPdControlMux uint8

const (
	usbPdCtrlMuxNoChange usbPdControlMux = 0
	usbPdCtrlMuxNone                     = 1
	usbPdCtrlMuxUsb                      = 2
	usbPdCtrlMuxDp                       = 3
	usbPdCtrlMuxDock                     = 4
	usbPdCtrlMuxAuto                     = 5
	usbPdCtrlMuxCount
)

/* TYPE */
type ecParamsUsbPdControl struct {
	port uint8
	role uint8
	mux  uint8
}

/* TYPE */
type ecResponseUsbPdControl struct {
	enabled  uint8
	role     uint8
	polarity uint8
	state    uint8
}

/* TYPE */
type ecResponseUsbPdControlV1 struct {
	enabled  uint8
	role     uint8 /* [0] power: 0=SNK/1=SRC [1] data: 0=UFP/1=DFP */
	polarity uint8
	state    [32]byte
}

const ecCmdUsbPdPorts = 0x102

/* TYPE */
type ecResponseUsbPdPorts struct {
	numPorts uint8
}

const ecCmdUsbPdPowerInfo = 0x103

const pdPowerChargingPort = 0xff

/* TYPE */
type ecParamsUsbPdPowerInfo struct {
	port uint8
}

type usbChgType uint8

const (
	usbChgTypeNone usbChgType = iota
	usbChgTypePd
	usbChgTypeC
	usbChgTypeProprietary
	usbChgTypeBc12Dcp
	usbChgTypeBc12Cdp
	usbChgTypeBc12Sdp
	usbChgTypeOther
	usbChgTypeVbus
	usbChgTypeUnknown
)

type usbPowerRoles uint8

const (
	usbPdPortPowerDisconnected usbPowerRoles = iota
	usbPdPortPowerSource
	usbPdPortPowerSink
	usbPdPortPowerSinkNotCharging
)

/* TYPE */
type usbChgMeasures struct {
	voltageMax uint16
	voltageNow uint16
	currentMax uint16
	currentLim uint16
}

/* TYPE */
type ecResponseUsbPdPowerInfo struct {
	role      uint8
	etype     uint8
	dualrole  uint8
	reserved1 uint8
	meas      usbChgMeasures
	maxPower  uint32
}

/* Write USB-PD device FW */
const ecCmdUsbPdFwUpdate = 0x110

type usbPdFwUpdateCmds uint8

const (
	usbPdFwReboot usbPdFwUpdateCmds = iota
	usbPdFwFlashErase
	usbPdFwFlashWrite
	usbPdFwEraseSig
)

/* TYPE */
type ecParamsUsbPdFwUpdate struct {
	devID uint16
	cmd   uint8
	port  uint8
	size  uint32 /* Size to write in bytes */
	/* Followed by data to write */
}

/* Write USB-PD Accessory RWHash table entry */
const ecCmdUsbPdRwHashEntry = 0x111

/* RW hash is first 20 bytes of SHA-256 of RW section */
const pdRwHashSize = 20

/* TYPE */
type ecParamsUsbPdRwHashEntry struct {
	devID        uint16
	devRwHash    [pdRwHashSize]uint8
	reserved     uint8  /* For alignment of currentImage */
	currentImage uint32 /* One of ecCurrentImage */
}

/* Read USB-PD Accessory info */
const ecCmdUsbPdDevInfo = 0x112

/* TYPE */
type ecParamsUsbPdInfoRequest struct {
	port uint8
}

/* Read USB-PD Device discovery info */
const ecCmdUsbPdDiscovery = 0x113

/* TYPE */
type ecParamsUsbPdDiscoveryEntry struct {
	vid   uint16 /* USB-IF VID */
	pid   uint16 /* USB-IF PID */
	ptype uint8  /* product type (hub,periph,cable,ama) */
}

/* Override default charge behavior */
const ecCmdPdChargePortOverride = 0x114

/* Negative port parameters have special meaning */
type usbPdOverridePorts int8

const (
	overrideDontCharge usbPdOverridePorts = -2
	overrideOff        usbPdOverridePorts = -1
	/* [0, pdPortCount): Port# */
)

/* TYPE */
type ecParamsChargePortOverride struct {
	overridePort int16 /* Override port# */
}

/* Read (and delete) one entry of PD event log */
const ecCmdPdGetLogEntry = 0x115

/* TYPE */
type ecResponsePdLog struct {
	timestamp uint32  /* relative timestamp in milliseconds */
	etype     uint8   /* event type : see pdEventXx below */
	sizePort  uint8   /* [7:5] port number [4:0] payload size in bytes */
	data      uint16  /* type-defined data payload */
	payload   []uint8 /* optional additional data payload: 0..16 bytes */
}

/* The timestamp is the microsecond counter shifted to get about a ms. */
const (
	pdLogTimestampShift = 10 /* 1 LSB = 1024us */
	pdLogSizeMask       = 0x1F
	pdLogPortMask       = 0xE0
	pdLogPortShift      = 5
)

func pdLogPortSize(port, size uint8) uint8 {
	return (port << pdLogPortShift) | (size & pdLogSizeMask)
}

func pdLogPort(sizePort uint8) uint8 {
	return sizePort >> pdLogPortShift
}

func pdLogSize(sizePort uint8) uint8 {
	return sizePort & pdLogSizeMask
}

/* PD event log : entry types */
/* PD MCU events */
const (
	pdEventMcuBase    = 0x00
	pdEventMcuCharge  = (pdEventMcuBase + 0)
	pdEventMcuConnect = (pdEventMcuBase + 1)
)

/* Reserved for custom board event */
const pdEventMcuBoardCustom = (pdEventMcuBase + 2)

/* PD generic accessory events */
const (
	pdEventAccBase    = 0x20
	pdEventAccRwFail  = (pdEventAccBase + 0)
	pdEventAccRwErase = (pdEventAccBase + 1)
)

/* PD power supply events */
const (
	pdEventPsBase  = 0x40
	pdEventPsFault = (pdEventPsBase + 0)
)

/* PD video dongles events */
const (
	pdEventVideoBase   = 0x60
	pdEventVideoDpMode = (pdEventVideoBase + 0)
	pdEventVideoCodec  = (pdEventVideoBase + 1)
)

/* Returned in the "type" field, when there is no entry available */
const pdEventNoEntry = 0xFF

/*
 * pdEventMcuCharge event definition :
 * the payload is "struct usbChgMeasures"
 * the data field contains the port state flags as defined below :
 */
/* Port partner is a dual role device */
const chargeFlagsDualRole = (1 << 15)

/* Port is the pending override port */
const chargeFlagsDelayedOverride = (1 << 14)

/* Port is the override port */
const chargeFlagsOverride = (1 << 13)

/* Charger type */
const (
	chargeFlagsTypeShift = 3
	chargeFlagsTypeMask  = (0xF << chargeFlagsTypeShift)
)

/* Power delivery role */
const chargeFlagsRoleMask = (7 << 0)

/*
 * pdEventPsFault data field flags definition :
 */
const (
	psFaultOcp     = 1
	psFaultFastOcp = 2
	psFaultOvp     = 3
	psFaultDisch   = 4
)

/* TYPE */
/*
 * pdEventVideoCodec payload is "struct mcdpInfo".
 */
type mcdpVersion struct {
	major uint8
	minor uint8
	build uint16
}

/* TYPE */
type mcdpInfo struct {
	family [2]uint8
	chipid [2]uint8
	irom   mcdpVersion
	fw     mcdpVersion
}

/* struct mcdpInfo field decoding */
func mcdpChipid(chipid []uint8) uint16 {
	return (uint16(chipid[0]) << 8) | uint16(chipid[1])
}

func mcdpFamily(family []uint8) uint16 {
	return (uint16(family[0]) << 8) | uint16(family[1])
}

/* Get/Set USB-PD Alternate mode info */
const ecCmdUsbPdGetAmode = 0x116

/* TYPE */
type ecParamsUsbPdGetModeRequest struct {
	svidIdx uint16 /* SVID index to get */
	port    uint8  /* port */
}

/* TYPE */
type ecParamsUsbPdGetModeResponse struct {
	svid uint16    /* SVID */
	opos uint16    /* Object Position */
	vdo  [6]uint32 /* Mode VDOs */
}

const ecCmdUsbPdSetAmode = 0x117

type pdModeCmd uint8

const (
	pdExitMode  pdModeCmd = 0
	pdEnterMode pdModeCmd = 1
	/* Not a command.  Do NOT remove. */
	pdModeCmdCount
)

/* TYPE */
type ecParamsUsbPdSetModeRequest struct {
	cmd  uint32 /* enum pdModeCmd */
	svid uint16 /* SVID to set */
	opos uint8  /* Object Position */
	port uint8  /* port */
}

/* Ask the PD MCU to record a log of a requested type */
const ecCmdPdWriteLogEntry = 0x118

/* TYPE */
type ecParamsPdWriteLogEntry struct {
	etype uint8 /* event type : see pdEventXx above */
	port  uint8 /* port#, or 0 for events unrelated to a given port */
}

/*****************************************************************************/
/*
 * Passthru commands
 *
 * Some platforms have sub-processors chained to each other.  For example.
 *
 *     AP <--> EC <--> PD MCU
 *
 * The top 2 bits of the command number are used to indicate which device the
 * command is intended for.  Device 0 is always the device receiving the
 * command; other device mapping is board-specific.
 *
 * When a device receives a command to be passed to a sub-processor, it passes
 * it on with the device number set back to 0.  This allows the sub-processor
 * to remain blissfully unaware of whether the command originated on the next
 * device up the chain, or was passed through from the AP.
 *
 * In the above example, if the AP wants to send command 0x0002 to the PD MCU
 *     AP sends command 0x4002 to the EC
 *     EC sends command 0x0002 to the PD MCU
 *     EC forwards PD MCU response back to the AP
 */

/* Offset and max command number for sub-device n */
func ecCmdPassthruOffset(n uint) uint {
	return 0x4000 * n
}

func ecCmdPassthruMax(n uint) uint {
	return ecCmdPassthruOffset(n) + 0x3fff
}

/*****************************************************************************/
/*
 * Deprecated constants. These constants have been renamed for clarity. The
 * meaning and size has not changed. Programs that use the old names should
 * switch to the new names soon, as the old names may not be carried forward
 * forever.
 */
const (
	ecHostParamSize   = ecProto2MaxParamSize
	ecLpcAddrOldParam = ecHostCmdRegion1
	ecOldParamSize    = ecHostCmdRegionSize
)

func ecVerMask(version uint8) uint8 {
	/* Command version mask */
	return 1 << version
}
