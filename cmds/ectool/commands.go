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
	EC_PROTO_VERSION = 0x00000002

	/* I/O addresses for ACPI commands */
	EC_LPC_ADDR_ACPI_DATA = 0x62
	EC_LPC_ADDR_ACPI_CMD  = 0x66

	/* I/O addresses for host command */
	EC_LPC_ADDR_HOST_DATA = 0x200
	EC_LPC_ADDR_HOST_CMD  = 0x204

	/* I/O addresses for host command args and params */
	/* Protocol version 2 */
	EC_LPC_ADDR_HOST_ARGS  = 0x800 /* And 0x801, 0x802, 0x803 */
	EC_LPC_ADDR_HOST_PARAM = 0x804 /* For version 2 params; size is
	 * EC_PROTO2_MAX_PARAM_SIZE */
	/* Protocol version 3 */
	EC_LPC_ADDR_HOST_PACKET = 0x800 /* Offset of version 3 packet */
	EC_LPC_HOST_PACKET_SIZE = 0x100 /* Max size of version 3 packet */

	/* The actual block is 0x800-0x8ff, but some BIOSes think it's 0x880-0x8ff
	 * and they tell the kernel that so we have to think of it as two parts. */
	EC_HOST_CMD_REGION0     = 0x800
	EC_HOST_CMD_REGION1     = 0x880
	EC_HOST_CMD_REGION_SIZE = 0x80

	/* EC command register bit functions */
	EC_LPC_CMDR_DATA      = (1 << 0) /* Data ready for host to read */
	EC_LPC_CMDR_PENDING   = (1 << 1) /* Write pending to EC */
	EC_LPC_CMDR_BUSY      = (1 << 2) /* EC is busy processing a command */
	EC_LPC_CMDR_CMD       = (1 << 3) /* Last host write was a command */
	EC_LPC_CMDR_ACPI_BRST = (1 << 4) /* Burst mode (not used) */
	EC_LPC_CMDR_SCI       = (1 << 5) /* SCI event is pending */
	EC_LPC_CMDR_SMI       = (1 << 6) /* SMI event is pending */

	EC_LPC_ADDR_MEMMAP = 0x900
	EC_MEMMAP_SIZE     = 255 /* ACPI IO buffer max is 255 bytes */
	EC_MEMMAP_TEXT_MAX = 8   /* Size of a string in the memory map */

	/* The offset address of each type of data in mapped memory. */
	EC_MEMMAP_TEMP_SENSOR      = 0x00 /* Temp sensors 0x00 - 0x0f */
	EC_MEMMAP_FAN              = 0x10 /* Fan speeds 0x10 - 0x17 */
	EC_MEMMAP_TEMP_SENSOR_B    = 0x18 /* More temp sensors 0x18 - 0x1f */
	EC_MEMMAP_ID               = 0x20 /* 0x20 == 'E', 0x21 == 'C' */
	EC_MEMMAP_ID_VERSION       = 0x22 /* Version of data in 0x20 - 0x2f */
	EC_MEMMAP_THERMAL_VERSION  = 0x23 /* Version of data in 0x00 - 0x1f */
	EC_MEMMAP_BATTERY_VERSION  = 0x24 /* Version of data in 0x40 - 0x7f */
	EC_MEMMAP_SWITCHES_VERSION = 0x25 /* Version of data in 0x30 - 0x33 */
	EC_MEMMAP_EVENTS_VERSION   = 0x26 /* Version of data in 0x34 - 0x3f */
	EC_MEMMAP_HOST_CMD_FLAGS   = 0x27 /* Host cmd interface flags (8 bits) */
	/* Unused 0x28 - 0x2f */
	EC_MEMMAP_SWITCHES = 0x30 /* 8 bits */
	/* Unused 0x31 - 0x33 */
	EC_MEMMAP_HOST_EVENTS = 0x34 /* 32 bits */
	/* Reserve 0x38 - 0x3f for additional host event-related stuff */
	/* Battery values are all 32 bits */
	EC_MEMMAP_BATT_VOLT = 0x40 /* Battery Present Voltage */
	EC_MEMMAP_BATT_RATE = 0x44 /* Battery Present Rate */
	EC_MEMMAP_BATT_CAP  = 0x48 /* Battery Remaining Capacity */
	EC_MEMMAP_BATT_FLAG = 0x4c /* Battery State, defined below */
	EC_MEMMAP_BATT_DCAP = 0x50 /* Battery Design Capacity */
	EC_MEMMAP_BATT_DVLT = 0x54 /* Battery Design Voltage */
	EC_MEMMAP_BATT_LFCC = 0x58 /* Battery Last Full Charge Capacity */
	EC_MEMMAP_BATT_CCNT = 0x5c /* Battery Cycle Count */
	/* Strings are all 8 bytes (EC_MEMMAP_TEXT_MAX) */
	EC_MEMMAP_BATT_MFGR   = 0x60 /* Battery Manufacturer String */
	EC_MEMMAP_BATT_MODEL  = 0x68 /* Battery Model Number String */
	EC_MEMMAP_BATT_SERIAL = 0x70 /* Battery Serial Number String */
	EC_MEMMAP_BATT_TYPE   = 0x78 /* Battery Type String */
	EC_MEMMAP_ALS         = 0x80 /* ALS readings in lux (2 X 16 bits) */
	/* Unused 0x84 - 0x8f */
	EC_MEMMAP_ACC_STATUS = 0x90 /* Accelerometer status (8 bits )*/
	/* Unused 0x91 */
	EC_MEMMAP_ACC_DATA  = 0x92 /* Accelerometer data 0x92 - 0x9f */
	EC_MEMMAP_GYRO_DATA = 0xa0 /* Gyroscope data 0xa0 - 0xa5 */
	/* Unused 0xa6 - 0xdf */

	/*
	 * ACPI is unable to access memory mapped data at or above this offset due to
	 * limitations of the ACPI protocol. Do not place data in the range 0xe0 - 0xfe
	 * which might be needed by ACPI.
	 */
	EC_MEMMAP_NO_ACPI = 0xe0

	/* Define the format of the accelerometer mapped memory status byte. */
	EC_MEMMAP_ACC_STATUS_SAMPLE_ID_MASK = 0x0f
	EC_MEMMAP_ACC_STATUS_BUSY_BIT       = (1 << 4)
	EC_MEMMAP_ACC_STATUS_PRESENCE_BIT   = (1 << 7)

	/* Number of temp sensors at EC_MEMMAP_TEMP_SENSOR */
	EC_TEMP_SENSOR_ENTRIES = 16
	/*
	 * Number of temp sensors at EC_MEMMAP_TEMP_SENSOR_B.
	 *
	 * Valid only if EC_MEMMAP_THERMAL_VERSION returns >= 2.
	 */
	EC_TEMP_SENSOR_B_ENTRIES = 8

	/* Special values for mapped temperature sensors */
	EC_TEMP_SENSOR_NOT_PRESENT    = 0xff
	EC_TEMP_SENSOR_ERROR          = 0xfe
	EC_TEMP_SENSOR_NOT_POWERED    = 0xfd
	EC_TEMP_SENSOR_NOT_CALIBRATED = 0xfc
	/*
	 * The offset of temperature value stored in mapped memory.  This allows
	 * reporting a temperature range of 200K to 454K = -73C to 181C.
	 */
	EC_TEMP_SENSOR_OFFSET = 200

	/*
	 * Number of ALS readings at EC_MEMMAP_ALS
	 */
	EC_ALS_ENTRIES = 2

	/*
	 * The default value a temperature sensor will return when it is present but
	 * has not been read this boot.  This is a reasonable number to avoid
	 * triggering alarms on the host.
	 */
	EC_TEMP_SENSOR_DEFAULT = (296 - EC_TEMP_SENSOR_OFFSET)

	EC_FAN_SPEED_ENTRIES     = 4      /* Number of fans at EC_MEMMAP_FAN */
	EC_FAN_SPEED_NOT_PRESENT = 0xffff /* Entry not present */
	EC_FAN_SPEED_STALLED     = 0xfffe /* Fan stalled */

	/* Battery bit flags at EC_MEMMAP_BATT_FLAG. */
	EC_BATT_FLAG_AC_PRESENT     = 0x01
	EC_BATT_FLAG_BATT_PRESENT   = 0x02
	EC_BATT_FLAG_DISCHARGING    = 0x04
	EC_BATT_FLAG_CHARGING       = 0x08
	EC_BATT_FLAG_LEVEL_CRITICAL = 0x10

	/* Switch flags at EC_MEMMAP_SWITCHES */
	EC_SWITCH_LID_OPEN               = 0x01
	EC_SWITCH_POWER_BUTTON_PRESSED   = 0x02
	EC_SWITCH_WRITE_PROTECT_DISABLED = 0x04
	/* Was recovery requested via keyboard; now unused. */
	EC_SWITCH_IGNORE1 = 0x08
	/* Recovery requested via dedicated signal (from servo board) */
	EC_SWITCH_DEDICATED_RECOVERY = 0x10
	/* Was fake developer mode switch; now unused.  Remove in next refactor. */
	EC_SWITCH_IGNORE0 = 0x20

	/* Host command interface flags */
	/* Host command interface supports LPC args (LPC interface only) */
	EC_HOST_CMD_FLAG_LPC_ARGS_SUPPORTED = 0x01
	/* Host command interface supports version 3 protocol */
	EC_HOST_CMD_FLAG_VERSION_3 = 0x02

	/* Wireless switch flags */
	EC_WIRELESS_SWITCH_ALL        = ^0x00 /* All flags */
	EC_WIRELESS_SWITCH_WLAN       = 0x01  /* WLAN radio */
	EC_WIRELESS_SWITCH_BLUETOOTH  = 0x02  /* Bluetooth radio */
	EC_WIRELESS_SWITCH_WWAN       = 0x04  /* WWAN power */
	EC_WIRELESS_SWITCH_WLAN_POWER = 0x08  /* WLAN power */

	/*****************************************************************************/
	/*
	 * ACPI commands
	 *
	 * These are valid ONLY on the ACPI command/data port.
	 */

	/*
	 * ACPI Read Embedded Controller
	 *
	 * This reads from ACPI memory space on the EC (EC_ACPI_MEM_*).
	 *
	 * Use the following sequence:
	 *
	 *    - Write EC_CMD_ACPI_READ to EC_LPC_ADDR_ACPI_CMD
	 *    - Wait for EC_LPC_CMDR_PENDING bit to clear
	 *    - Write address to EC_LPC_ADDR_ACPI_DATA
	 *    - Wait for EC_LPC_CMDR_DATA bit to set
	 *    - Read value from EC_LPC_ADDR_ACPI_DATA
	 */
	EC_CMD_ACPI_READ = 0x80

	/*
	 * ACPI Write Embedded Controller
	 *
	 * This reads from ACPI memory space on the EC (EC_ACPI_MEM_*).
	 *
	 * Use the following sequence:
	 *
	 *    - Write EC_CMD_ACPI_WRITE to EC_LPC_ADDR_ACPI_CMD
	 *    - Wait for EC_LPC_CMDR_PENDING bit to clear
	 *    - Write address to EC_LPC_ADDR_ACPI_DATA
	 *    - Wait for EC_LPC_CMDR_PENDING bit to clear
	 *    - Write value to EC_LPC_ADDR_ACPI_DATA
	 */
	EC_CMD_ACPI_WRITE = 0x81

	/*
	 * ACPI Burst Enable Embedded Controller
	 *
	 * This enables burst mode on the EC to allow the host to issue several
	 * commands back-to-back. While in this mode, writes to mapped multi-byte
	 * data are locked out to ensure data consistency.
	 */
	EC_CMD_ACPI_BURST_ENABLE = 0x82

	/*
	 * ACPI Burst Disable Embedded Controller
	 *
	 * This disables burst mode on the EC and stops preventing EC writes to mapped
	 * multi-byte data.
	 */
	EC_CMD_ACPI_BURST_DISABLE = 0x83

	/*
	 * ACPI Query Embedded Controller
	 *
	 * This clears the lowest-order bit in the currently pending host events, and
	 * sets the result code to the 1-based index of the bit (event 0x00000001 = 1
	 * event 0x80000000 = 32), or 0 if no event was pending.
	 */
	EC_CMD_ACPI_QUERY_EVENT = 0x84

	/* Valid addresses in ACPI memory space, for read/write commands */

	/* Memory space version; set to EC_ACPI_MEM_VERSION_CURRENT */
	EC_ACPI_MEM_VERSION = 0x00
	/*
	 * Test location; writing value here updates test compliment byte to (0xff -
	 * value).
	 */
	EC_ACPI_MEM_TEST = 0x01
	/* Test compliment; writes here are ignored. */
	EC_ACPI_MEM_TEST_COMPLIMENT = 0x02

	/* Keyboard backlight brightness percent (0 - 100) */
	EC_ACPI_MEM_KEYBOARD_BACKLIGHT = 0x03
	/* DPTF Target Fan Duty (0-100, 0xff for auto/none) */
	EC_ACPI_MEM_FAN_DUTY = 0x04

	/*
	 * DPTF temp thresholds. Any of the EC's temp sensors can have up to two
	 * independent thresholds attached to them. The current value of the ID
	 * register determines which sensor is affected by the THRESHOLD and COMMIT
	 * registers. The THRESHOLD register uses the same EC_TEMP_SENSOR_OFFSET scheme
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
	EC_ACPI_MEM_TEMP_ID        = 0x05
	EC_ACPI_MEM_TEMP_THRESHOLD = 0x06
	EC_ACPI_MEM_TEMP_COMMIT    = 0x07
	/*
	 * Here are the bits for the COMMIT register:
	 *   bit 0 selects the threshold index for the chosen sensor (0/1)
	 *   bit 1 enables/disables the selected threshold (0 = off, 1 = on)
	 * Each write to the commit register affects one threshold.
	 */
	EC_ACPI_MEM_TEMP_COMMIT_SELECT_MASK = (1 << 0)
	EC_ACPI_MEM_TEMP_COMMIT_ENABLE_MASK = (1 << 1)
	/*
	 * Example:
	 *
	 * Set the thresholds for sensor 2 to 50 C and 60 C:
	 *   write 2 to [0x05]      --  select temp sensor 2
	 *   write 0x7b to [0x06]   --  C_TO_K(50) - EC_TEMP_SENSOR_OFFSET
	 *   write 0x2 to [0x07]    --  enable threshold 0 with this value
	 *   write 0x85 to [0x06]   --  C_TO_K(60) - EC_TEMP_SENSOR_OFFSET
	 *   write 0x3 to [0x07]    --  enable threshold 1 with this value
	 *
	 * Disable the 60 C threshold, leaving the 50 C threshold unchanged:
	 *   write 2 to [0x05]      --  select temp sensor 2
	 *   write 0x1 to [0x07]    --  disable threshold 1
	 */

	/* DPTF battery charging current limit */
	EC_ACPI_MEM_CHARGING_LIMIT = 0x08

	/* Charging limit is specified in 64 mA steps */
	EC_ACPI_MEM_CHARGING_LIMIT_STEP_MA = 64
	/* Value to disable DPTF battery charging limit */
	EC_ACPI_MEM_CHARGING_LIMIT_DISABLED = 0xff

	/*
	 * ACPI addresses 0x20 - 0xff map to EC_MEMMAP offset 0x00 - 0xdf.  This data
	 * is read-only from the AP.  Added in EC_ACPI_MEM_VERSION 2.
	 */
	EC_ACPI_MEM_MAPPED_BEGIN = 0x20
	EC_ACPI_MEM_MAPPED_SIZE  = 0xe0

	/* Current version of ACPI memory address space */
	EC_ACPI_MEM_VERSION_CURRENT = 2

	/*
	 * This header file is used in coreboot both in C and ACPI code.  The ACPI code
	 * is pre-processed to handle constants but the ASL compiler is unable to
	 * handle actual C code so keep it separate.
	 */

	/* LPC command status byte masks */
	/* EC has written a byte in the data register and host hasn't read it yet */
	EC_LPC_STATUS_TO_HOST = 0x01
	/* Host has written a command/data byte and the EC hasn't read it yet */
	EC_LPC_STATUS_FROM_HOST = 0x02
	/* EC is processing a command */
	EC_LPC_STATUS_PROCESSING = 0x04
	/* Last write to EC was a command, not data */
	EC_LPC_STATUS_LAST_CMD = 0x08
	/* EC is in burst mode */
	EC_LPC_STATUS_BURST_MODE = 0x10
	/* SCI event is pending (requesting SCI query) */
	EC_LPC_STATUS_SCI_PENDING = 0x20
	/* SMI event is pending (requesting SMI query) */
	EC_LPC_STATUS_SMI_PENDING = 0x40
	/* (reserved) */
	EC_LPC_STATUS_RESERVED = 0x80

	/*
	 * EC is busy.  This covers both the EC processing a command, and the host has
	 * written a new command but the EC hasn't picked it up yet.
	 */
	EC_LPC_STATUS_BUSY_MASK = (EC_LPC_STATUS_FROM_HOST | EC_LPC_STATUS_PROCESSING)
)

/* Host command response codes */
type ec_status uint8

const (
	EC_RES_SUCCESS           ec_status = 0
	EC_RES_INVALID_COMMAND             = 1
	EC_RES_ERROR                       = 2
	EC_RES_INVALID_PARAM               = 3
	EC_RES_ACCESS_DENIED               = 4
	EC_RES_INVALID_RESPONSE            = 5
	EC_RES_INVALID_VERSION             = 6
	EC_RES_INVALID_CHECKSUM            = 7
	EC_RES_IN_PROGRESS                 = 8  /* Accepted, command in progress */
	EC_RES_UNAVAILABLE                 = 9  /* No response available */
	EC_RES_TIMEOUT                     = 10 /* We got a timeout */
	EC_RES_OVERFLOW                    = 11 /* Table / data overflow */
	EC_RES_INVALID_HEADER              = 12 /* Header contains invalid data */
	EC_RES_REQUEST_TRUNCATED           = 13 /* Didn't get the entire request */
	EC_RES_RESPONSE_TOO_BIG            = 14 /* Response was too big to handle */
	EC_RES_BUS_ERROR                   = 15 /* Communications bus error */
	EC_RES_BUSY                        = 16 /* Up but too busy.  Should retry */
)

/*
 * Host event codes.  Note these are 1-based, not 0-based, because ACPI query
 * EC command uses code 0 to mean "no event pending".  We explicitly specify
 * each value in the enum listing so they won't change if we delete/insert an
 * item or rearrange the list (it needs to be stable across platforms, not
 * just within a single compiled instance).
 */
type host_event_code uint8

const (
	EC_HOST_EVENT_LID_CLOSED        host_event_code = 1
	EC_HOST_EVENT_LID_OPEN                          = 2
	EC_HOST_EVENT_POWER_BUTTON                      = 3
	EC_HOST_EVENT_AC_CONNECTED                      = 4
	EC_HOST_EVENT_AC_DISCONNECTED                   = 5
	EC_HOST_EVENT_BATTERY_LOW                       = 6
	EC_HOST_EVENT_BATTERY_CRITICAL                  = 7
	EC_HOST_EVENT_BATTERY                           = 8
	EC_HOST_EVENT_THERMAL_THRESHOLD                 = 9
	EC_HOST_EVENT_THERMAL_OVERLOAD                  = 10
	EC_HOST_EVENT_THERMAL                           = 11
	EC_HOST_EVENT_USB_CHARGER                       = 12
	EC_HOST_EVENT_KEY_PRESSED                       = 13
	/*
	 * EC has finished initializing the host interface.  The host can check
	 * for this event following sending a EC_CMD_REBOOT_EC command to
	 * determine when the EC is ready to accept subsequent commands.
	 */
	EC_HOST_EVENT_INTERFACE_READY = 14
	/* Keyboard recovery combo has been pressed */
	EC_HOST_EVENT_KEYBOARD_RECOVERY = 15

	/* Shutdown due to thermal overload */
	EC_HOST_EVENT_THERMAL_SHUTDOWN = 16
	/* Shutdown due to battery level too low */
	EC_HOST_EVENT_BATTERY_SHUTDOWN = 17

	/* Suggest that the AP throttle itself */
	EC_HOST_EVENT_THROTTLE_START = 18
	/* Suggest that the AP resume normal speed */
	EC_HOST_EVENT_THROTTLE_STOP = 19

	/* Hang detect logic detected a hang and host event timeout expired */
	EC_HOST_EVENT_HANG_DETECT = 20
	/* Hang detect logic detected a hang and warm rebooted the AP */
	EC_HOST_EVENT_HANG_REBOOT = 21

	/* PD MCU triggering host event */
	EC_HOST_EVENT_PD_MCU = 22

	/* Battery Status flags have changed */
	EC_HOST_EVENT_BATTERY_STATUS = 23

	/* EC encountered a panic, triggering a reset */
	EC_HOST_EVENT_PANIC = 24

	/*
	 * The high bit of the event mask is not used as a host event code.  If
	 * it reads back as set, then the entire event mask should be
	 * considered invalid by the host.  This can happen when reading the
	 * raw event status via EC_MEMMAP_HOST_EVENTS but the LPC interface is
	 * not initialized on the EC, or improperly configured on the host.
	 */
	EC_HOST_EVENT_INVALID = 32
)

/* Host event mask */
func ec_host_event_mask(event_code uint8) uint8 {
	return 1 << ((event_code) - 1)
}

/* TYPE */
/* Arguments at EC_LPC_ADDR_HOST_ARGS */
type ec_lpc_host_args struct {
	flags           uint8
	command_version uint8
	data_size       uint8
	/*
	 * Checksum; sum of command + flags + command_version + data_size +
	 * all params/response data bytes.
	 */
	checksum uint8
}

/* Flags for ec_lpc_host_args.flags */
/*
 * Args are from host.  Data area at EC_LPC_ADDR_HOST_PARAM contains command
 * params.
 *
 * If EC gets a command and this flag is not set, this is an old-style command.
 * Command version is 0 and params from host are at EC_LPC_ADDR_OLD_PARAM with
 * unknown length.  EC must respond with an old-style response (that is
 * withouth setting EC_HOST_ARGS_FLAG_TO_HOST).
 */
const EC_HOST_ARGS_FLAG_FROM_HOST = 0x01

/*
 * Args are from EC.  Data area at EC_LPC_ADDR_HOST_PARAM contains response.
 *
 * If EC responds to a command and this flag is not set, this is an old-style
 * response.  Command version is 0 and response data from EC is at
 * EC_LPC_ADDR_OLD_PARAM with unknown length.
 */
const EC_HOST_ARGS_FLAG_TO_HOST = 0x02

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
 *   2. EC_SPI_OLD_READY   - AP sends first byte(s) of request
 *   3. -                  - EC starts handling CS# interrupt
 *   4. EC_SPI_RECEIVING   - AP sends remaining byte(s) of request
 *   5. EC_SPI_PROCESSING  - EC starts processing request; AP is clocking in
 *                           bytes looking for EC_SPI_FRAME_START
 *   6. -                  - EC finishes processing and sets up response
 *   7. EC_SPI_FRAME_START - AP reads frame byte
 *   8. (response packet)  - AP reads response packet
 *   9. EC_SPI_PAST_END    - Any additional bytes read by AP
 *   10 -                  - AP deasserts chip select
 *   11 -                  - EC processes CS# interrupt and sets up DMA for
 *                           next request
 *
 * If the AP is waiting for EC_SPI_FRAME_START and sees any value other than
 * the following byte values:
 *   EC_SPI_OLD_READY
 *   EC_SPI_RX_READY
 *   EC_SPI_RECEIVING
 *   EC_SPI_PROCESSING
 *
 * Then the EC found an error in the request, or was not ready for the request
 * and lost data.  The AP should give up waiting for EC_SPI_FRAME_START
 * because the EC is unable to tell when the AP is done sending its request.
 */

/*
 * Framing byte which precedes a response packet from the EC.  After sending a
 * request, the AP will clock in bytes until it sees the framing byte, then
 * clock in the response packet.
 */
const (
	EC_SPI_FRAME_START = 0xec

	/*
	 * Padding bytes which are clocked out after the end of a response packet.
	 */
	EC_SPI_PAST_END = 0xed

	/*
	 * EC is ready to receive, and has ignored the byte sent by the AP.  EC expects
	 * that the AP will send a valid packet header (starting with
	 * EC_COMMAND_PROTOCOL_3) in the next 32 bytes.
	 */
	EC_SPI_RX_READY = 0xf8

	/*
	 * EC has started receiving the request from the AP, but hasn't started
	 * processing it yet.
	 */
	EC_SPI_RECEIVING = 0xf9

	/* EC has received the entire request from the AP and is processing it. */
	EC_SPI_PROCESSING = 0xfa

	/*
	 * EC received bad data from the AP, such as a packet header with an invalid
	 * length.  EC will ignore all data until chip select deasserts.
	 */
	EC_SPI_RX_BAD_DATA = 0xfb

	/*
	 * EC received data from the AP before it was ready.  That is, the AP asserted
	 * chip select and started clocking data before the EC was ready to receive it.
	 * EC will ignore all data until chip select deasserts.
	 */
	EC_SPI_NOT_READY = 0xfc

	/*
	 * EC was ready to receive a request from the AP.  EC has treated the byte sent
	 * by the AP as part of a request packet, or (for old-style ECs) is processing
	 * a fully received packet but is not ready to respond yet.
	 */
	EC_SPI_OLD_READY = 0xfd

	/*****************************************************************************/

	/*
	 * Protocol version 2 for I2C and SPI send a request this way:
	 *
	 *	0	EC_CMD_VERSION0 + (command version)
	 *	1	Command number
	 *	2	Length of params = N
	 *	3..N+2	Params, if any
	 *	N+3	8-bit checksum of bytes 0..N+2
	 *
	 * The corresponding response is:
	 *
	 *	0	Result code (EC_RES_*)
	 *	1	Length of params = M
	 *	2..M+1	Params, if any
	 *	M+2	8-bit checksum of bytes 0..M+1
	 */
	EC_PROTO2_REQUEST_HEADER_BYTES  = 3
	EC_PROTO2_REQUEST_TRAILER_BYTES = 1
	EC_PROTO2_REQUEST_OVERHEAD      = (EC_PROTO2_REQUEST_HEADER_BYTES +
		EC_PROTO2_REQUEST_TRAILER_BYTES)

	EC_PROTO2_RESPONSE_HEADER_BYTES  = 2
	EC_PROTO2_RESPONSE_TRAILER_BYTES = 1
	EC_PROTO2_RESPONSE_OVERHEAD      = (EC_PROTO2_RESPONSE_HEADER_BYTES +
		EC_PROTO2_RESPONSE_TRAILER_BYTES)

	/* Parameter length was limited by the LPC interface */
	EC_PROTO2_MAX_PARAM_SIZE = 0xfc

	/* Maximum request and response packet sizes for protocol version 2 */
	EC_PROTO2_MAX_REQUEST_SIZE = (EC_PROTO2_REQUEST_OVERHEAD +
		EC_PROTO2_MAX_PARAM_SIZE)
	EC_PROTO2_MAX_RESPONSE_SIZE = (EC_PROTO2_RESPONSE_OVERHEAD +
		EC_PROTO2_MAX_PARAM_SIZE)

	/*****************************************************************************/

	/*
	 * Value written to legacy command port / prefix byte to indicate protocol
	 * 3+ structs are being used.  Usage is bus-dependent.
	 */
	EC_COMMAND_PROTOCOL_3 = 0xda

	EC_HOST_REQUEST_VERSION = 3
)

/* TYPE */
/* Version 3 request from host */
type ec_host_request struct {
	/* Struct version (=3)
	 *
	 * EC will return EC_RES_INVALID_HEADER if it receives a header with a
	 * version it doesn't know how to parse.
	 */
	struct_version uint8

	/*
	 * Checksum of request and data; sum of all bytes including checksum
	 * should total to 0.
	 */
	checksum uint8

	/* Command code */
	command uint16

	/* Command version */
	command_version uint8

	/* Unused byte in current protocol version; set to 0 */
	reserved uint8

	/* Length of data which follows this header */
	data_len uint16
}

const EC_HOST_RESPONSE_VERSION = 3

/* TYPE */
/* Version 3 response from EC */
type ec_host_response struct {
	/* Struct version (=3) */
	struct_version uint8

	/*
	 * Checksum of response and data; sum of all bytes including checksum
	 * should total to 0.
	 */
	checksum uint8

	/* Result code (EC_RES_*) */
	result uint16

	/* Length of data which follows this header */
	data_len uint16

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
const EC_CMD_PROTO_VERSION = 0x00

/* TYPE */
type ec_response_proto_version struct {
	version uint32
}

const (
	/*
	 * Hello.  This is a simple command to test the EC is responsive to
	 * commands.
	 */
	EC_CMD_HELLO = 0x01
)

/* TYPE */
type ec_params_hello struct {
	in_data uint32 /* Pass anything here */
}

/* TYPE */
type ec_response_hello struct {
	out_data uint32 /* Output will be in_data + 0x01020304 */
}

const (
	/* Get version number */
	EC_CMD_GET_VERSION = 0x02
)

type ec_current_image uint8

const (
	EC_IMAGE_UNKNOWN ec_current_image = 0
	EC_IMAGE_RO
	EC_IMAGE_RW
)

/* TYPE */
type ec_response_get_version struct {
	/* Null-terminated version strings for RO, RW */
	version_string_ro [32]byte
	version_string_rw [32]byte
	reserved          [32]byte /* Was previously RW-B string */
	current_image     uint32   /* One of ec_current_image */
}

const (
	/* Read test */
	EC_CMD_READ_TEST = 0x03
)

/* TYPE */
type ec_params_read_test struct {
	offset uint32 /* Starting value for read buffer */
	size   uint32 /* Size to read in bytes */
}

/* TYPE */
type ec_response_read_test struct {
	data [32]uint32
}

const (
	/*
	 * Get build information
	 *
	 * Response is null-terminated string.
	 */
	EC_CMD_GET_BUILD_INFO = 0x04

	/* Get chip info */
	EC_CMD_GET_CHIP_INFO = 0x05
)

/* TYPE */
type ec_response_get_chip_info struct {
	/* Null-terminated strings */
	vendor   [32]byte
	name     [32]byte
	revision [32]byte /* Mask version */
}

const (
	/* Get board HW version */
	EC_CMD_GET_BOARD_VERSION = 0x06
)

/* TYPE */
type ec_response_board_version struct {
	board_version uint16 /* A monotonously incrementing number. */
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
	EC_CMD_READ_MEMMAP = 0x07
)

/* TYPE */
type ec_params_read_memmap struct {
	offset uint8 /* Offset in memmap (EC_MEMMAP_*) */
	size   uint8 /* Size to read in bytes */
}

/* Read versions supported for a command */
const EC_CMD_GET_CMD_VERSIONS = 0x08

/* TYPE */
type ec_params_get_cmd_versions struct {
	cmd uint8 /* Command to check */
}

/* TYPE */
type ec_params_get_cmd_versions_v1 struct {
	cmd uint16 /* Command to check */
}

/* TYPE */
type ec_response_get_cmd_versions struct {
	/*
	 * Mask of supported versions; use EC_VER_MASK() to compare with a
	 * desired version.
	 */
	version_mask uint32
}

/*
 * Check EC communcations status (busy). This is needed on i2c/spi but not
 * on lpc since it has its own out-of-band busy indicator.
 *
 * lpc must read the status from the command register. Attempting this on
 * lpc will overwrite the args/parameter space and corrupt its data.
 */
const EC_CMD_GET_COMMS_STATUS = 0x09

/* Avoid using ec_status which is for return values */
type ec_comms_status uint8

const (
	EC_COMMS_STATUS_PROCESSING ec_comms_status = 1 << 0 /* Processing cmd */
)

/* TYPE */
type ec_response_get_comms_status struct {
	flags uint32 /* Mask of enum ec_comms_status */
}

const (
	/* Fake a variety of responses, purely for testing purposes. */
	EC_CMD_TEST_PROTOCOL = 0x0a
)

/* TYPE */
/* Tell the EC what to send back to us. */
type ec_params_test_protocol struct {
	ec_result uint32
	ret_len   uint32
	buf       [32]uint8
}

/* TYPE */
/* Here it comes... */
type ec_response_test_protocol struct {
	buf [32]uint8
}

/* Get prococol information */
const EC_CMD_GET_PROTOCOL_INFO = 0x0b

/* Flags for ec_response_get_protocol_info.flags */
/* EC_RES_IN_PROGRESS may be returned if a command is slow */
const EC_PROTOCOL_INFO_IN_PROGRESS_SUPPORTED = (1 << 0)

/* TYPE */
type ec_response_get_protocol_info struct {
	/* Fields which exist if at least protocol version 3 supported */

	/* Bitmask of protocol versions supported (1 << n means version n)*/
	protocol_versions uint32

	/* Maximum request packet size, in bytes */
	max_request_packet_size uint16

	/* Maximum response packet size, in bytes */
	max_response_packet_size uint16

	/* Flags; see EC_PROTOCOL_INFO_* */
	flags uint32
}

/*****************************************************************************/
/* Get/Set miscellaneous values */
const (
	/* The upper byte of .flags tells what to do (nothing means "get") */
	EC_GSV_SET = 0x80000000

	/* The lower three bytes of .flags identifies the parameter, if that has
	   meaning for an individual command. */
	EC_GSV_PARAM_MASK = 0x00ffffff
)

/* TYPE */
type ec_params_get_set_value struct {
	flags uint32
	value uint32
}

/* TYPE */
type ec_response_get_set_value struct {
	flags uint32
	value uint32
}

/* More than one command can use these structs to get/set paramters. */
const EC_CMD_GSV_PAUSE_IN_S5 = 0x0c

/*****************************************************************************/
/* Flash commands */

/* Get flash info */
const EC_CMD_FLASH_INFO = 0x10

/* TYPE */
/* Version 0 returns these fields */
type ec_response_flash_info struct {
	/* Usable flash size, in bytes */
	flash_size uint32
	/*
	 * Write block size.  Write offset and size must be a multiple
	 * of this.
	 */
	write_block_size uint32
	/*
	 * Erase block size.  Erase offset and size must be a multiple
	 * of this.
	 */
	erase_block_size uint32
	/*
	 * Protection block size.  Protection offset and size must be a
	 * multiple of this.
	 */
	protect_block_size uint32
}

/* Flags for version 1+ flash info command */
/* EC flash erases bits to 0 instead of 1 */
const EC_FLASH_INFO_ERASE_TO_0 = (1 << 0)

/* TYPE */
/*
 * Version 1 returns the same initial fields as version 0, with additional
 * fields following.
 *
 * gcc anonymous structs don't seem to get along with the  directive;
 * if they did we'd define the version 0 struct as a sub-struct of this one.
 */
type ec_response_flash_info_1 struct {
	/* Version 0 fields; see above for description */
	flash_size         uint32
	write_block_size   uint32
	erase_block_size   uint32
	protect_block_size uint32

	/* Version 1 adds these fields: */
	/*
	 * Ideal write size in bytes.  Writes will be fastest if size is
	 * exactly this and offset is a multiple of this.  For example, an EC
	 * may have a write buffer which can do half-page operations if data is
	 * aligned, and a slower word-at-a-time write mode.
	 */
	write_ideal_size uint32

	/* Flags; see EC_FLASH_INFO_* */
	flags uint32
}

/*
 * Read flash
 *
 * Response is params.size bytes of data.
 */
const EC_CMD_FLASH_READ = 0x11

/* TYPE */
type ec_params_flash_read struct {
	offset uint32 /* Byte offset to read */
	size   uint32 /* Size to read in bytes */
}

const (
	/* Write flash */
	EC_CMD_FLASH_WRITE = 0x12
	EC_VER_FLASH_WRITE = 1

	/* Version 0 of the flash command supported only 64 bytes of data */
	EC_FLASH_WRITE_VER0_SIZE = 64
)

/* TYPE */
type ec_params_flash_write struct {
	offset uint32 /* Byte offset to write */
	size   uint32 /* Size to write in bytes */
	/* Followed by data to write */
}

/* Erase flash */
const EC_CMD_FLASH_ERASE = 0x13

/* TYPE */
type ec_params_flash_erase struct {
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
	EC_CMD_FLASH_PROTECT = 0x15
	EC_VER_FLASH_PROTECT = 1 /* Command version 1 */

	/* Flags for flash protection */
	/* RO flash code protected when the EC boots */
	EC_FLASH_PROTECT_RO_AT_BOOT = (1 << 0)
	/*
	 * RO flash code protected now.  If this bit is set, at-boot status cannot
	 * be changed.
	 */
	EC_FLASH_PROTECT_RO_NOW = (1 << 1)
	/* Entire flash code protected now, until reboot. */
	EC_FLASH_PROTECT_ALL_NOW = (1 << 2)
	/* Flash write protect GPIO is asserted now */
	EC_FLASH_PROTECT_GPIO_ASSERTED = (1 << 3)
	/* Error - at least one bank of flash is stuck locked, and cannot be unlocked */
	EC_FLASH_PROTECT_ERROR_STUCK = (1 << 4)
	/*
	 * Error - flash protection is in inconsistent state.  At least one bank of
	 * flash which should be protected is not protected.  Usually fixed by
	 * re-requesting the desired flags, or by a hard reset if that fails.
	 */
	EC_FLASH_PROTECT_ERROR_INCONSISTENT = (1 << 5)
	/* Entire flash code protected when the EC boots */
	EC_FLASH_PROTECT_ALL_AT_BOOT = (1 << 6)
)

/* TYPE */
type ec_params_flash_protect struct {
	mask  uint32 /* Bits in flags to apply */
	flags uint32 /* New flags to apply */
}

/* TYPE */
type ec_response_flash_protect struct {
	/* Current value of flash protect flags */
	flags uint32
	/*
	 * Flags which are valid on this platform.  This allows the caller
	 * to distinguish between flags which aren't set vs. flags which can't
	 * be set on this platform.
	 */
	valid_flags uint32
	/* Flags which can be changed given the current protection state */
	writable_flags uint32
}

/*
 * Note: commands 0x14 - 0x19 version 0 were old commands to get/set flash
 * write protect.  These commands may be reused with version > 0.
 */

/* Get the region offset/size */
const EC_CMD_FLASH_REGION_INFO = 0x16
const EC_VER_FLASH_REGION_INFO = 1

type ec_flash_region uint8

const (
	/* Region which holds read-only EC image */
	EC_FLASH_REGION_RO ec_flash_region = iota
	/* Region which holds rewritable EC image */
	EC_FLASH_REGION_RW
	/*
	 * Region which should be write-protected in the factory (a superset of
	 * EC_FLASH_REGION_RO)
	 */
	EC_FLASH_REGION_WP_RO
	/* Number of regions */
	EC_FLASH_REGION_COUNT
)

/* TYPE */
type ec_params_flash_region_info struct {
	region uint32 /* enum ec_flash_region */
}

/* TYPE */
type ec_response_flash_region_info struct {
	offset uint32
	size   uint32
}

const (
	/* Read/write VbNvContext */
	EC_CMD_VBNV_CONTEXT = 0x17
	EC_VER_VBNV_CONTEXT = 1
	EC_VBNV_BLOCK_SIZE  = 16
)

type ec_vbnvcontext_op uint8

const (
	EC_VBNV_CONTEXT_OP_READ ec_vbnvcontext_op = iota
	EC_VBNV_CONTEXT_OP_WRITE
)

/* TYPE */
type ec_params_vbnvcontext struct {
	op    uint32
	block [EC_VBNV_BLOCK_SIZE]uint8
}

/* TYPE */
type ec_response_vbnvcontext struct {
	block [EC_VBNV_BLOCK_SIZE]uint8
}

/*****************************************************************************/
/* PWM commands */

/* Get fan target RPM */
const EC_CMD_PWM_GET_FAN_TARGET_RPM = 0x20

/* TYPE */
type ec_response_pwm_get_fan_rpm struct {
	rpm uint32
}

/* Set target fan RPM */
const EC_CMD_PWM_SET_FAN_TARGET_RPM = 0x21

/* TYPE */
/* Version 0 of input params */
type ec_params_pwm_set_fan_target_rpm_v0 struct {
	rpm uint32
}

/* TYPE */
/* Version 1 of input params */
type ec_params_pwm_set_fan_target_rpm_v1 struct {
	rpm     uint32
	fan_idx uint8
}

/* Get keyboard backlight */
const EC_CMD_PWM_GET_KEYBOARD_BACKLIGHT = 0x22

/* TYPE */
type ec_response_pwm_get_keyboard_backlight struct {
	percent uint8
	enabled uint8
}

/* Set keyboard backlight */
const EC_CMD_PWM_SET_KEYBOARD_BACKLIGHT = 0x23

/* TYPE */
type ec_params_pwm_set_keyboard_backlight struct {
	percent uint8
}

/* Set target fan PWM duty cycle */
const EC_CMD_PWM_SET_FAN_DUTY = 0x24

/* TYPE */
/* Version 0 of input params */
type ec_params_pwm_set_fan_duty_v0 struct {
	percent uint32
}

/* TYPE */
/* Version 1 of input params */
type ec_params_pwm_set_fan_duty_v1 struct {
	percent uint32
	fan_idx uint8
}

/*****************************************************************************/
/*
 * Lightbar commands. This looks worse than it is. Since we only use one HOST
 * command to say "talk to the lightbar", we put the "and tell it to do X" part
 * into a subcommand. We'll make separate structs for subcommands with
 * different input args, so that we know how much to expect.
 */
const EC_CMD_LIGHTBAR_CMD = 0x28

/* TYPE */
type rgb_s struct {
	r, g, b uint8
}

const LB_BATTERY_LEVELS = 4

/* TYPE */
/* List of tweakable parameters. NOTE: It's  so it can be sent in a
 * host command, but the alignment is the same regardless. Keep it that way.
 */
type lightbar_params_v0 struct {
	/* Timing */
	google_ramp_up   int32
	google_ramp_down int32
	s3s0_ramp_up     int32
	s0_tick_delay    [2]int32 /* AC=0/1 */
	s0a_tick_delay   [2]int32 /* AC=0/1 */
	s0s3_ramp_down   int32
	s3_sleep_for     int32
	s3_ramp_up       int32
	s3_ramp_down     int32

	/* Oscillation */
	new_s0  uint8
	osc_min [2]uint8 /* AC=0/1 */
	osc_max [2]uint8 /* AC=0/1 */
	w_ofs   [2]uint8 /* AC=0/1 */

	/* Brightness limits based on the backlight and AC. */
	bright_bl_off_fixed [2]uint8 /* AC=0/1 */
	bright_bl_on_min    [2]uint8 /* AC=0/1 */
	bright_bl_on_max    [2]uint8 /* AC=0/1 */

	/* Battery level thresholds */
	batteryhreshold [LB_BATTERY_LEVELS - 1]uint8

	/* Map [AC][battery_level] to color index */
	s0_idx [2][LB_BATTERY_LEVELS]uint8 /* AP is running */
	s3_idx [2][LB_BATTERY_LEVELS]uint8 /* AP is sleeping */

	/* Color palette */
	color [8]rgb_s /* 0-3 are Google colors */
}

/* TYPE */
type lightbar_params_v1 struct {
	/* Timing */
	google_ramp_up   int32
	google_ramp_down int32
	s3s0_ramp_up     int32
	s0_tick_delay    [2]int32 /* AC=0/1 */
	s0a_tick_delay   [2]int32 /* AC=0/1 */
	s0s3_ramp_down   int32
	s3_sleep_for     int32
	s3_ramp_up       int32
	s3_ramp_down     int32
	s5_ramp_up       int32
	s5_ramp_down     int32
	tap_tick_delay   int32
	tap_gate_delay   int32
	tap_display_time int32

	/* Tap-for-battery params */
	tap_pct_red    uint8
	tap_pct_green  uint8
	tap_seg_min_on uint8
	tap_seg_max_on uint8
	tap_seg_osc    uint8
	tap_idx        [3]uint8

	/* Oscillation */
	osc_min [2]uint8 /* AC=0/1 */
	osc_max [2]uint8 /* AC=0/1 */
	w_ofs   [2]uint8 /* AC=0/1 */

	/* Brightness limits based on the backlight and AC. */
	bright_bl_off_fixed [2]uint8 /* AC=0/1 */
	bright_bl_on_min    [2]uint8 /* AC=0/1 */
	bright_bl_on_max    [2]uint8 /* AC=0/1 */

	/* Battery level thresholds */
	batteryhreshold [LB_BATTERY_LEVELS - 1]uint8

	/* Map [AC][battery_level] to color index */
	s0_idx [2][LB_BATTERY_LEVELS]uint8 /* AP is running */
	s3_idx [2][LB_BATTERY_LEVELS]uint8 /* AP is sleeping */

	/* s5: single color pulse on inhibited power-up */
	s5_idx uint8

	/* Color palette */
	color [8]rgb_s /* 0-3 are Google colors */
}

/* TYPE */
/* Lightbar command params v2
 * crbug.com/467716
 *
 * lightbar_parms_v1 was too big for i2c, therefore in v2, we split them up by
 * logical groups to make it more manageable ( < 120 bytes).
 *
 * NOTE: Each of these groups must be less than 120 bytes.
 */

type lightbar_params_v2_timing struct {
	/* Timing */
	google_ramp_up   int32
	google_ramp_down int32
	s3s0_ramp_up     int32
	s0_tick_delay    [2]int32 /* AC=0/1 */
	s0a_tick_delay   [2]int32 /* AC=0/1 */
	s0s3_ramp_down   int32
	s3_sleep_for     int32
	s3_ramp_up       int32
	s3_ramp_down     int32
	s5_ramp_up       int32
	s5_ramp_down     int32
	tap_tick_delay   int32
	tap_gate_delay   int32
	tap_display_time int32
}

/* TYPE */
type lightbar_params_v2_tap struct {
	/* Tap-for-battery params */
	tap_pct_red    uint8
	tap_pct_green  uint8
	tap_seg_min_on uint8
	tap_seg_max_on uint8
	tap_seg_osc    uint8
	tap_idx        [3]uint8
}

/* TYPE */
type lightbar_params_v2_oscillation struct {
	/* Oscillation */
	osc_min [2]uint8 /* AC=0/1 */
	osc_max [2]uint8 /* AC=0/1 */
	w_ofs   [2]uint8 /* AC=0/1 */
}

/* TYPE */
type lightbar_params_v2_brightness struct {
	/* Brightness limits based on the backlight and AC. */
	bright_bl_off_fixed [2]uint8 /* AC=0/1 */
	bright_bl_on_min    [2]uint8 /* AC=0/1 */
	bright_bl_on_max    [2]uint8 /* AC=0/1 */
}

/* TYPE */
type lightbar_params_v2_thresholds struct {
	/* Battery level thresholds */
	batteryhreshold [LB_BATTERY_LEVELS - 1]uint8
}

/* TYPE */
type lightbar_params_v2_colors struct {
	/* Map [AC][battery_level] to color index */
	s0_idx [2][LB_BATTERY_LEVELS]uint8 /* AP is running */
	s3_idx [2][LB_BATTERY_LEVELS]uint8 /* AP is sleeping */

	/* s5: single color pulse on inhibited power-up */
	s5_idx uint8

	/* Color palette */
	color [8]rgb_s /* 0-3 are Google colors */
}

/* Lightbyte program. */
const EC_LB_PROG_LEN = 192

/* TYPE */
type lightbar_program struct {
	size uint8
	data [EC_LB_PROG_LEN]uint8
}

/* TYPE */
/* this one is messy
type ec_params_lightbar struct {
cmd uint8		      /* Command (see enum lightbar_command)
	union {
		struct {
			/* no args
		} dump, off, on, init, get_seq, get_params_v0, get_params_v1
			version, get_brightness, get_demo, suspend, resume
			get_params_v2_timing, get_params_v2_tap
			get_params_v2_osc, get_params_v2_bright
			get_params_v2_thlds, get_params_v2_colors;

		struct {
		uint8 num;
		} set_brightness, seq, demo;

		struct {
		uint8 ctrl, reg, value;
		} reg;

		struct {
		uint8 led, red, green, blue;
		} set_rgb;

		struct {
		uint8 led;
		} get_rgb;

		struct {
		uint8 enable;
		} manual_suspend_ctrl;

		struct lightbar_params_v0 set_params_v0;
		struct lightbar_params_v1 set_params_v1;

		struct lightbar_params_v2_timing set_v2par_timing;
		struct lightbar_params_v2_tap set_v2par_tap;
		struct lightbar_params_v2_oscillation set_v2par_osc;
		struct lightbar_params_v2_brightness set_v2par_bright;
		struct lightbar_params_v2_thresholds set_v2par_thlds;
		struct lightbar_params_v2_colors set_v2par_colors;

		struct lightbar_program set_program;
	};
} ;
*/
/* TYPE */
/*
type ec_response_lightbar struct {
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
		} get_seq, get_brightness, get_demo;

		struct lightbar_params_v0 get_params_v0;
		struct lightbar_params_v1 get_params_v1;


		struct lightbar_params_v2_timing get_params_v2_timing;
		struct lightbar_params_v2_tap get_params_v2_tap;
		struct lightbar_params_v2_oscillation get_params_v2_osc;
		struct lightbar_params_v2_brightness get_params_v2_bright;
		struct lightbar_params_v2_thresholds get_params_v2_thlds;
		struct lightbar_params_v2_colors get_params_v2_colors;

		struct {
		uint32 num;
		uint32 flags;
		} version;

		struct {
		uint8 red, green, blue;
		} get_rgb;

		struct {
			/* no return params *
		} off, on, init, set_brightness, seq, reg, set_rgb
			demo, set_params_v0, set_params_v1
			set_program, manual_suspend_ctrl, suspend, resume
			set_v2par_timing, set_v2par_tap
			set_v2par_osc, set_v2par_bright, set_v2par_thlds
			set_v2par_colors;
	};
} ;
*/
/* Lightbar commands */
type lightbar_command uint8

const (
	LIGHTBAR_CMD_DUMP                      lightbar_command = 0
	LIGHTBAR_CMD_OFF                                        = 1
	LIGHTBAR_CMD_ON                                         = 2
	LIGHTBAR_CMD_INIT                                       = 3
	LIGHTBAR_CMD_SET_BRIGHTNESS                             = 4
	LIGHTBAR_CMD_SEQ                                        = 5
	LIGHTBAR_CMD_REG                                        = 6
	LIGHTBAR_CMD_SET_RGB                                    = 7
	LIGHTBAR_CMD_GET_SEQ                                    = 8
	LIGHTBAR_CMD_DEMO                                       = 9
	LIGHTBAR_CMD_GET_PARAMS_V0                              = 10
	LIGHTBAR_CMD_SET_PARAMS_V0                              = 11
	LIGHTBAR_CMD_VERSION                                    = 12
	LIGHTBAR_CMD_GET_BRIGHTNESS                             = 13
	LIGHTBAR_CMD_GET_RGB                                    = 14
	LIGHTBAR_CMD_GET_DEMO                                   = 15
	LIGHTBAR_CMD_GET_PARAMS_V1                              = 16
	LIGHTBAR_CMD_SET_PARAMS_V1                              = 17
	LIGHTBAR_CMD_SET_PROGRAM                                = 18
	LIGHTBAR_CMD_MANUAL_SUSPEND_CTRL                        = 19
	LIGHTBAR_CMD_SUSPEND                                    = 20
	LIGHTBAR_CMD_RESUME                                     = 21
	LIGHTBAR_CMD_GET_PARAMS_V2_TIMING                       = 22
	LIGHTBAR_CMD_SET_PARAMS_V2_TIMING                       = 23
	LIGHTBAR_CMD_GET_PARAMS_V2_TAP                          = 24
	LIGHTBAR_CMD_SET_PARAMS_V2_TAP                          = 25
	LIGHTBAR_CMD_GET_PARAMS_V2_OSCILLATION                  = 26
	LIGHTBAR_CMD_SET_PARAMS_V2_OSCILLATION                  = 27
	LIGHTBAR_CMD_GET_PARAMS_V2_BRIGHTNESS                   = 28
	LIGHTBAR_CMD_SET_PARAMS_V2_BRIGHTNESS                   = 29
	LIGHTBAR_CMD_GET_PARAMS_V2_THRESHOLDS                   = 30
	LIGHTBAR_CMD_SET_PARAMS_V2_THRESHOLDS                   = 31
	LIGHTBAR_CMD_GET_PARAMS_V2_COLORS                       = 32
	LIGHTBAR_CMD_SET_PARAMS_V2_COLORS                       = 33
	LIGHTBAR_NUM_CMDS
)

/*****************************************************************************/
/* LED control commands */

const EC_CMD_LED_CONTROL = 0x29

type ec_led_id uint8

const (
	/* LED to indicate battery state of charge */
	EC_LED_ID_BATTERY_LED ec_led_id = 0
	/*
	 * LED to indicate system power state (on or in suspend).
	 * May be on power button or on C-panel.
	 */
	EC_LED_ID_POWER_LED
	/* LED on power adapter or its plug */
	EC_LED_ID_ADAPTER_LED

	EC_LED_ID_COUNT
)
const (
	/* LED control flags */
	EC_LED_FLAGS_QUERY = (1 << iota) /* Query LED capability only */
	EC_LED_FLAGS_AUTO                /* Switch LED back to automatic control */
)

type ec_led_colors uint8

const (
	EC_LED_COLOR_RED ec_led_colors = 0
	EC_LED_COLOR_GREEN
	EC_LED_COLOR_BLUE
	EC_LED_COLOR_YELLOW
	EC_LED_COLOR_WHITE

	EC_LED_COLOR_COUNT
)

/* TYPE */
type ec_params_led_control struct {
	led_id uint8 /* Which LED to control */
	flags  uint8 /* Control flags */

	brightness [EC_LED_COLOR_COUNT]uint8
}

/* TYPE */
type ec_response_led_control struct {
	/*
	 * Available brightness value range.
	 *
	 * Range 0 means color channel not present.
	 * Range 1 means on/off control.
	 * Other values means the LED is control by PWM.
	 */
	brightness_range [EC_LED_COLOR_COUNT]uint8
}

/*****************************************************************************/
/* Verified boot commands */

/*
 * Note: command code 0x29 version 0 was VBOOT_CMD in Link EVT; it may be
 * reused for other purposes with version > 0.
 */

/* Verified boot hash command */
const EC_CMD_VBOOT_HASH = 0x2A

/* TYPE */
type ec_params_vboot_hash struct {
	cmd        uint8     /* enum ec_vboot_hash_cmd */
	hashype    uint8     /* enum ec_vboot_hash_type */
	nonce_size uint8     /* Nonce size; may be 0 */
	reserved0  uint8     /* Reserved; set 0 */
	offset     uint32    /* Offset in flash to hash */
	size       uint32    /* Number of bytes to hash */
	nonce_data [64]uint8 /* Nonce data; ignored if nonce_size=0 */
}

/* TYPE */
type ec_response_vboot_hash struct {
	status      uint8     /* enum ec_vboot_hash_status */
	hashype     uint8     /* enum ec_vboot_hash_type */
	digest_size uint8     /* Size of hash digest in bytes */
	reserved0   uint8     /* Ignore; will be 0 */
	offset      uint32    /* Offset in flash which was hashed */
	size        uint32    /* Number of bytes hashed */
	hash_digest [64]uint8 /* Hash digest data */
}

type ec_vboot_hash_cmd uint8

const (
	EC_VBOOT_HASH_GET    ec_vboot_hash_cmd = 0 /* Get current hash status */
	EC_VBOOT_HASH_ABORT                    = 1 /* Abort calculating current hash */
	EC_VBOOT_HASH_START                    = 2 /* Start computing a new hash */
	EC_VBOOT_HASH_RECALC                   = 3 /* Synchronously compute a new hash */
)

type ec_vboot_hash_type uint8

const (
	EC_VBOOT_HASH_TYPE_SHA256 ec_vboot_hash_type = 0 /* SHA-256 */
)

type ec_vboot_hash_status uint8

const (
	EC_VBOOT_HASH_STATUS_NONE ec_vboot_hash_status = 0 /* No hash (not started, or aborted) */
	EC_VBOOT_HASH_STATUS_DONE                      = 1 /* Finished computing a hash */
	EC_VBOOT_HASH_STATUS_BUSY                      = 2 /* Busy computing a hash */
)

/*
 * Special values for offset for EC_VBOOT_HASH_START and EC_VBOOT_HASH_RECALC.
 * If one of these is specified, the EC will automatically update offset and
 * size to the correct values for the specified image (RO or RW).
 */
const EC_VBOOT_HASH_OFFSET_RO = 0xfffffffe
const EC_VBOOT_HASH_OFFSET_RW = 0xfffffffd

/*****************************************************************************/
/*
 * Motion sense commands. We'll make separate structs for sub-commands with
 * different input args, so that we know how much to expect.
 */
const EC_CMD_MOTION_SENSE_CMD = 0x2B

/* Motion sense commands */
type motionsense_command uint8

const (
	/* Dump command returns all motion sensor data including motion sense
	 * module flags and individual sensor flags.
	 */
	MOTIONSENSE_CMD_DUMP motionsense_command = iota

	/*
	 * Info command returns data describing the details of a given sensor
	 * including enum motionsensor_type, enum motionsensor_location, and
	 * enum motionsensor_chip.
	 */
	MOTIONSENSE_CMD_INFO

	/*
	 * EC Rate command is a setter/getter command for the EC sampling rate
	 * of all motion sensors in milliseconds.
	 */
	MOTIONSENSE_CMD_EC_RATE

	/*
	 * Sensor ODR command is a setter/getter command for the output data
	 * rate of a specific motion sensor in millihertz.
	 */
	MOTIONSENSE_CMD_SENSOR_ODR

	/*
	 * Sensor range command is a setter/getter command for the range of
	 * a specified motion sensor in +/-G's or +/- deg/s.
	 */
	MOTIONSENSE_CMD_SENSOR_RANGE

	/*
	 * Setter/getter command for the keyboard wake angle. When the lid
	 * angle is greater than this value, keyboard wake is disabled in S3
	 * and when the lid angle goes less than this value, keyboard wake is
	 * enabled. Note, the lid angle measurement is an approximate
	 * un-calibrated value, hence the wake angle isn't exact.
	 */
	MOTIONSENSE_CMD_KB_WAKE_ANGLE

	/* Number of motionsense sub-commands. */
	MOTIONSENSE_NUM_CMDS
)

/* List of motion sensor types. */
type motionsensor_type uint8

const (
	MOTIONSENSE_TYPE_ACCEL motionsensor_type = 0
	MOTIONSENSE_TYPE_GYRO                    = 1
)

/* List of motion sensor locations. */
type motionsensor_location uint8

const (
	MOTIONSENSE_LOC_BASE motionsensor_location = 0
	MOTIONSENSE_LOC_LID                        = 1
)

/* List of motion sensor chips. */
type motionsensor_chip uint8

const (
	MOTIONSENSE_CHIP_KXCJ9   motionsensor_chip = 0
	MOTIONSENSE_CHIP_LSM6DS0                   = 1
)

/* Module flag masks used for the dump sub-command. */
const MOTIONSENSE_MODULE_FLAG_ACTIVE = (1 << 0)

/* Sensor flag masks used for the dump sub-command. */
const MOTIONSENSE_SENSOR_FLAG_PRESENT = (1 << 0)

/*
 * Send this value for the data element to only perform a read. If you
 * send any other value, the EC will interpret it as data to set and will
 * return the actual value set.
 */
const EC_MOTION_SENSE_NO_VALUE = -1

/* some other time
type ec_params_motion_sense struct {
cmd uint8
	union {
		/* Used for MOTIONSENSE_CMD_DUMP * /
		struct {
			/*
			 * Maximal number of sensor the host is expecting.
			 * 0 means the host is only interested in the number
			 * of sensors controlled by the EC.
			 * /
		uint8 max_sensor_count;
		} dump;

		/*
		 * Used for MOTIONSENSE_CMD_EC_RATE and
		 * MOTIONSENSE_CMD_KB_WAKE_ANGLE.
		 * /
		struct {
			/* Data to set or EC_MOTION_SENSE_NO_VALUE to read. * /
			data int16
		} ec_rate, kb_wake_angle;

		/* Used for MOTIONSENSE_CMD_INFO. * /
		struct {
		uint8 sensor_num;
		} info;

		/*
		 * Used for MOTIONSENSE_CMD_SENSOR_ODR and
		 * MOTIONSENSE_CMD_SENSOR_RANGE.
		 * /
		struct {
		uint8 sensor_num;

			/* Rounding flag, true for round-up, false for down. * /
		uint8 roundup;

		uint16 reserved;

			/* Data to set or EC_MOTION_SENSE_NO_VALUE to read. * /
			data int32
		} sensor_odr, sensor_range;
	};
} ;
*/
/* TYPE */
type ec_response_motion_sensor_data struct {
	/* Flags for each sensor. */
	flags   uint8
	padding uint8

	/* Each sensor is up to 3-axis. */
	data [3]int16
}

/* TYPE */
/* some other time

type ec_response_motion_sense struct {
	union {
		/* Used for MOTIONSENSE_CMD_DUMP * /
		struct {
			/* Flags representing the motion sensor module. * /
		uint8 module_flags;

			/* Number of sensors managed directly by the EC * /
		uint8 sensor_count;

			/*
			 * sensor data is truncated if response_max is too small
			 * for holding all the data.
			 * /
			struct ec_response_motion_sensor_data sensor[0];
		} dump;

		/* Used for MOTIONSENSE_CMD_INFO. * /
		struct {
			/* Should be element of enum motionsensor_type. * /
		uint8 type;

			/* Should be element of enum motionsensor_location. * /
		uint8 location;

			/* Should be element of enum motionsensor_chip. * /
		uint8 chip;
		} info;

		/*
		 * Used for MOTIONSENSE_CMD_EC_RATE, MOTIONSENSE_CMD_SENSOR_ODR
		 * MOTIONSENSE_CMD_SENSOR_RANGE, and
		 * MOTIONSENSE_CMD_KB_WAKE_ANGLE.
		 * /
		struct {
			/* Current value of the parameter queried. * /
			ret int32
		} ec_rate, sensor_odr, sensor_range, kb_wake_angle;
	};
} ;
*/
/*****************************************************************************/
/* Force lid open command */

/* Make lid event always open */
const EC_CMD_FORCE_LID_OPEN = 0x2c

/* TYPE */
type ec_params_force_lid_open struct {
	enabled uint8
}

/*****************************************************************************/
/* USB charging control commands */

/* Set USB port charging mode */
const EC_CMD_USB_CHARGE_SET_MODE = 0x30

/* TYPE */
type ec_params_usb_charge_set_mode struct {
	usb_port_id uint8
	mode        uint8
}

/*****************************************************************************/
/* Persistent storage for host */

/* Maximum bytes that can be read/written in a single command */
const EC_PSTORE_SIZE_MAX = 64

/* Get persistent storage info */
const EC_CMD_PSTORE_INFO = 0x40

/* TYPE */
type ec_response_pstore_info struct {
	/* Persistent storage size, in bytes */
	pstore_size uint32
	/* Access size; read/write offset and size must be a multiple of this */
	access_size uint32
}

/*
 * Read persistent storage
 *
 * Response is params.size bytes of data.
 */
const EC_CMD_PSTORE_READ = 0x41

/* TYPE */
type ec_params_pstore_read struct {
	offset uint32 /* Byte offset to read */
	size   uint32 /* Size to read in bytes */
}

/* Write persistent storage */
const EC_CMD_PSTORE_WRITE = 0x42

/* TYPE */
type ec_params_pstore_write struct {
	offset uint32 /* Byte offset to write */
	size   uint32 /* Size to write in bytes */
	data   [EC_PSTORE_SIZE_MAX]uint8
}

/* TYPE */
/*****************************************************************************/
/* Real-time clock */

/* RTC params and response structures */
type ec_params_rtc struct {
	time uint32
}

/* TYPE */
type ec_response_rtc struct {
	time uint32
}

/* These use ec_response_rtc */
const EC_CMD_RTC_GET_VALUE = 0x44
const EC_CMD_RTC_GET_ALARM = 0x45

/* These all use ec_params_rtc */
const EC_CMD_RTC_SET_VALUE = 0x46
const EC_CMD_RTC_SET_ALARM = 0x47

/*****************************************************************************/
/* Port80 log access */

/* Maximum entries that can be read/written in a single command */
const EC_PORT80_SIZE_MAX = 32

/* Get last port80 code from previous boot */
const EC_CMD_PORT80_LAST_BOOT = 0x48
const EC_CMD_PORT80_READ = 0x48

type ec_port80_subcmd uint8

const (
	EC_PORT80_GET_INFO ec_port80_subcmd = 0
	EC_PORT80_READ_BUFFER
)

/* TYPE */
type ec_params_port80_read struct {
	subcmd      uint16
	offset      uint32
	num_entries uint32
}

/* TYPE */
type ec_response_port80_read struct {
	/*
		struct {
		uint32 writes;
		uint32 history_size;
		uint32 last_boot;
		} get_info;*/

	codes [EC_PORT80_SIZE_MAX]uint16
}

/* TYPE */
type ec_response_port80_last_boot struct {
	code uint16
}

/*****************************************************************************/
/* Thermal engine commands. Note that there are two implementations. We'll
 * reuse the command number, but the data and behavior is incompatible.
 * Version 0 is what originally shipped on Link.
 * Version 1 separates the CPU thermal limits from the fan control.
 */

const EC_CMD_THERMAL_SET_THRESHOLD = 0x50
const EC_CMD_THERMAL_GET_THRESHOLD = 0x51

/* TYPE */
/* The version 0 structs are opaque. You have to know what they are for
 * the get/set commands to make any sense.
 */

/* Version 0 - set */
type ec_params_thermal_set_threshold struct {
	sensorype    uint8
	threshold_id uint8
	value        uint16
}

/* TYPE */
/* Version 0 - get */
type ec_params_thermal_get_threshold struct {
	sensorype    uint8
	threshold_id uint8
}

/* TYPE */
type ec_response_thermal_get_threshold struct {
	value uint16
}

/* The version 1 structs are visible. */
type ec_temp_thresholds uint8

const (
	EC_TEMP_THRESH_WARN ec_temp_thresholds = 0
	EC_TEMP_THRESH_HIGH
	EC_TEMP_THRESH_HALT

	EC_TEMP_THRESH_COUNT
)

/* TYPE */
/* Thermal configuration for one temperature sensor. Temps are in degrees K.
 * Zero values will be silently ignored by the thermal task.
 */
type ec_thermal_config struct {
	temp_host    [EC_TEMP_THRESH_COUNT]uint32 /* levels of hotness */
	temp_fan_off uint32                       /* no active cooling needed */
	temp_fan_max uint32                       /* max active cooling needed */
}

/* TYPE */
/* Version 1 - get config for one sensor. */
type ec_params_thermal_get_threshold_v1 struct {
	sensor_num uint32
}

/* TYPE */
/* This returns a struct ec_thermal_config */

/* Version 1 - set config for one sensor.
 * Use read-modify-write for best results! */
type ec_params_thermal_set_threshold_v1 struct {
	sensor_num uint32
	cfg        ec_thermal_config
}

/* This returns no data */

/****************************************************************************/

/* Toggle automatic fan control */
const EC_CMD_THERMAL_AUTO_FAN_CTRL = 0x52

/* TYPE */
/* Version 1 of input params */
type ec_params_auto_fan_ctrl_v1 struct {
	fan_idx uint8
}

/* Get/Set TMP006 calibration data */
const EC_CMD_TMP006_GET_CALIBRATION = 0x53
const EC_CMD_TMP006_SET_CALIBRATION = 0x54

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
type ec_params_tmp006_get_calibration struct {
	index uint8
}

/* TYPE */
/* Version 0 */
type ec_response_tmp006_get_calibration_v0 struct {
	s0, b0, b1, bw float32
}

/* TYPE */
type ec_params_tmp006_set_calibration_v0 struct {
	index          uint8
	reserved       [3]uint8
	s0, b0, b1, b2 float32
}

/* TYPE */
/* Version 1 */
type ec_response_tmp006_get_calibration_v1 struct {
	algorithm  uint8
	num_params uint8
	reserved   [2]uint8
	val        []float32
}

/* TYPE */
type ec_params_tmp006_set_calibration_v1 struct {
	index      uint8
	algorithm  uint8
	num_params uint8
	reserved   uint8
	val        []float32
}

/* Read raw TMP006 data */
const EC_CMD_TMP006_GET_RAW = 0x55

/* TYPE */
type ec_params_tmp006_get_raw struct {
	index uint8
}

/* TYPE */
type ec_response_tmp006_get_raw struct {
	t int32 /* In 1/100 K */
	v int32 /* In nV */
}

/*****************************************************************************/
/* MKBP - Matrix KeyBoard Protocol */

/*
 * Read key state
 *
 * Returns raw data for keyboard cols; see ec_response_mkbp_info.cols for
 * expected response size.
 */
const EC_CMD_MKBP_STATE = 0x60

/* Provide information about the matrix : number of rows and columns */
const EC_CMD_MKBP_INFO = 0x61

/* TYPE */
type ec_response_mkbp_info struct {
	rows     uint32
	cols     uint32
	switches uint8
}

/* Simulate key press */
const EC_CMD_MKBP_SIMULATE_KEY = 0x62

/* TYPE */
type ec_params_mkbp_simulate_key struct {
	col     uint8
	row     uint8
	pressed uint8
}

/* Configure keyboard scanning */
const EC_CMD_MKBP_SET_CONFIG = 0x64
const EC_CMD_MKBP_GET_CONFIG = 0x65

/* flags */
type mkbp_config_flags uint8

const (
	EC_MKBP_FLAGS_ENABLE mkbp_config_flags = 1 /* Enable keyboard scanning */
)

type mkbp_config_valid uint8

const (
	EC_MKBP_VALID_SCAN_PERIOD         mkbp_config_valid = 1 << 0
	EC_MKBP_VALID_POLL_TIMEOUT                          = 1 << 1
	EC_MKBP_VALID_MIN_POST_SCAN_DELAY                   = 1 << 3
	EC_MKBP_VALID_OUTPUT_SETTLE                         = 1 << 4
	EC_MKBP_VALID_DEBOUNCE_DOWN                         = 1 << 5
	EC_MKBP_VALID_DEBOUNCE_UP                           = 1 << 6
	EC_MKBP_VALID_FIFO_MAX_DEPTH                        = 1 << 7
)

/* TYPE */
/* Configuration for our key scanning algorithm */
type ec_mkbp_config struct {
	valid_mask     uint32 /* valid fields */
	flags          uint8  /* some flags (enum mkbp_config_flags) */
	valid_flags    uint8  /* which flags are valid */
	scan_period_us uint16 /* period between start of scans */
	/* revert to interrupt mode after no activity for this long */
	pollimeout_us uint32
	/*
	 * minimum post-scan relax time. Once we finish a scan we check
	 * the time until we are due to start the next one. If this time is
	 * shorter this field, we use this instead.
	 */
	min_post_scan_delay_us uint16
	/* delay between setting up output and waiting for it to settle */
	output_settle_us uint16
	debounce_down_us uint16 /* time for debounce on key down */
	debounce_up_us   uint16 /* time for debounce on key up */
	/* maximum depth to allow for fifo (0 = no keyscan output) */
	fifo_max_depth uint8
}

/* TYPE */
type ec_params_mkbp_set_config struct {
	config ec_mkbp_config
}

/* TYPE */
type ec_response_mkbp_get_config struct {
	config ec_mkbp_config
}

/* Run the key scan emulation */
const EC_CMD_KEYSCAN_SEQ_CTRL = 0x66

type ec_keyscan_seq_cmd uint8

const (
	EC_KEYSCAN_SEQ_STATUS  ec_keyscan_seq_cmd = 0 /* Get status information */
	EC_KEYSCAN_SEQ_CLEAR                      = 1 /* Clear sequence */
	EC_KEYSCAN_SEQ_ADD                        = 2 /* Add item to sequence */
	EC_KEYSCAN_SEQ_START                      = 3 /* Start running sequence */
	EC_KEYSCAN_SEQ_COLLECT                    = 4 /* Collect sequence summary data */
)

type ec_collect_flags uint8

const (
	/* Indicates this scan was processed by the EC. Due to timing, some
	 * scans may be skipped.
	 */
	EC_KEYSCAN_SEQ_FLAG_DONE ec_collect_flags = 1 << iota
)

/* TYPE */
type ec_collect_item struct {
	flags uint8 /* some flags (enum ec_collect_flags) */
}

/* TYPE */
/* later
type ec_params_keyscan_seq_ctrl struct {
cmd uint8	/* Command to send (enum ec_keyscan_seq_cmd) * /
	union {
		struct {
		uint8 active;		/* still active * /
		uint8 num_items;	/* number of items * /
			/* Current item being presented * /
		uint8 cur_item;
		} status;
		struct {
			/*
			 * Absolute time for this scan, measured from the
			 * start of the sequence.
			 * /
		uint32 time_us;
		uint8 scan[0];	/* keyscan data * /
		} add;
		struct {
		uint8 start_item;	/* First item to return * /
		uint8 num_items;	/* Number of items to return * /
		} collect;
	};
} ;
*/
/* TYPE */
/* lter
type ec_result_keyscan_seq_ctrl struct {
	union {
		struct {
		uint8 num_items;	/* Number of items *
			/* Data for each item *
			struct ec_collect_item item[0]
		} collect;
	};
} ;
*/
/*
 * Get the next pending MKBP event.
 *
 * Returns EC_RES_UNAVAILABLE if there is no event pending.
 */
const EC_CMD_GET_NEXT_EVENT = 0x67

type ec_mkbp_event uint8

const (
	/* Keyboard matrix changed. The event data is the new matrix state. */
	EC_MKBP_EVENT_KEY_MATRIX = iota

	/* New host event. The event data is 4 bytes of host event flags. */
	EC_MKBP_EVENT_HOST_EVENT

	/* Number of MKBP events */
	EC_MKBP_EVENT_COUNT
)

/* TYPE */
type ec_response_get_next_event struct {
	eventype uint8
	/* Followed by event data if any */
}

/*****************************************************************************/
/* Temperature sensor commands */

/* Read temperature sensor info */
const EC_CMD_TEMP_SENSOR_GET_INFO = 0x70

/* TYPE */
type ec_params_temp_sensor_get_info struct {
	id uint8
}

/* TYPE */
type ec_response_temp_sensor_get_info struct {
	sensor_name [32]byte
	sensorype   uint8
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
type ec_params_host_event_mask struct {
	mask uint32
}

/* TYPE */
type ec_response_host_event_mask struct {
	mask uint32
}

/* These all use ec_response_host_event_mask */
const EC_CMD_HOST_EVENT_GET_B = 0x87
const EC_CMD_HOST_EVENT_GET_SMI_MASK = 0x88
const EC_CMD_HOST_EVENT_GET_SCI_MASK = 0x89
const EC_CMD_HOST_EVENT_GET_WAKE_MASK = 0x8d

/* These all use ec_params_host_event_mask */
const EC_CMD_HOST_EVENT_SET_SMI_MASK = 0x8a
const EC_CMD_HOST_EVENT_SET_SCI_MASK = 0x8b
const EC_CMD_HOST_EVENT_CLEAR = 0x8c
const EC_CMD_HOST_EVENT_SET_WAKE_MASK = 0x8e
const EC_CMD_HOST_EVENT_CLEAR_B = 0x8f

/*****************************************************************************/
/* Switch commands */

/* Enable/disable LCD backlight */
const EC_CMD_SWITCH_ENABLE_BKLIGHT = 0x90

/* TYPE */
type ec_params_switch_enable_backlight struct {
	enabled uint8
}

/* Enable/disable WLAN/Bluetooth */
const EC_CMD_SWITCH_ENABLE_WIRELESS = 0x91
const EC_VER_SWITCH_ENABLE_WIRELESS = 1

/* TYPE */
/* Version 0 params; no response */
type ec_params_switch_enable_wireless_v0 struct {
	enabled uint8
}

/* TYPE */
/* Version 1 params */
type ec_params_switch_enable_wireless_v1 struct {
	/* Flags to enable now */
	now_flags uint8

	/* Which flags to copy from now_flags */
	now_mask uint8

	/*
	 * Flags to leave enabled in S3, if they're on at the S0->S3
	 * transition.  (Other flags will be disabled by the S0->S3
	 * transition.)
	 */
	suspend_flags uint8

	/* Which flags to copy from suspend_flags */
	suspend_mask uint8
}

/* TYPE */
/* Version 1 response */
type ec_response_switch_enable_wireless_v1 struct {
	/* Flags to enable now */
	now_flags uint8

	/* Flags to leave enabled in S3 */
	suspend_flags uint8
}

/*****************************************************************************/
/* GPIO commands. Only available on EC if write protect has been disabled. */

/* Set GPIO output value */
const EC_CMD_GPIO_SET = 0x92

/* TYPE */
type ec_params_gpio_set struct {
	name [32]byte
	val  uint8
}

/* Get GPIO value */
const EC_CMD_GPIO_GET = 0x93

/* TYPE */
/* Version 0 of input params and response */
type ec_params_gpio_get struct {
	name [32]byte
}

/* TYPE */
type ec_response_gpio_get struct {
	val uint8
}

/* TYPE */
/* Version 1 of input params and response */
type ec_params_gpio_get_v1 struct {
	subcmd uint8
	data   [32]byte
}

/* TYPE */
/* later
type ec_response_gpio_get_v1 struct {
	union {
		struct {
		uint8 val;
		} get_value_by_name, get_count;
		struct {
		uint8 val;
			char name[32];
		uint32 flags;
		} get_info;
	};
} ;
*/
type gpio_get_subcmd uint8

const (
	EC_GPIO_GET_BY_NAME gpio_get_subcmd = 0
	EC_GPIO_GET_COUNT                   = 1
	EC_GPIO_GET_INFO                    = 2
)

/*****************************************************************************/
/* I2C commands. Only available when flash write protect is unlocked. */

/*
 * TODO(crosbug.com/p/23570): These commands are deprecated, and will be
 * removed soon.  Use EC_CMD_I2C_XFER instead.
 */

/* Read I2C bus */
const EC_CMD_I2C_READ = 0x94

/* TYPE */
type ec_params_i2c_read struct {
	addr      uint16 /* 8-bit address (7-bit shifted << 1) */
	read_size uint8  /* Either 8 or 16. */
	port      uint8
	offset    uint8
}

/* TYPE */
type ec_response_i2c_read struct {
	data uint16
}

/* Write I2C bus */
const EC_CMD_I2C_WRITE = 0x95

/* TYPE */
type ec_params_i2c_write struct {
	data       uint16
	addr       uint16 /* 8-bit address (7-bit shifted << 1) */
	write_size uint8  /* Either 8 or 16. */
	port       uint8
	offset     uint8
}

/*****************************************************************************/
/* Charge state commands. Only available when flash write protect unlocked. */

/* Force charge state machine to stop charging the battery or force it to
 * discharge the battery.
 */
const EC_CMD_CHARGE_CONTROL = 0x96
const EC_VER_CHARGE_CONTROL = 1

type ec_charge_control_mode uint8

const (
	CHARGE_CONTROL_NORMAL ec_charge_control_mode = 0
	CHARGE_CONTROL_IDLE
	CHARGE_CONTROL_DISCHARGE
)

/* TYPE */
type ec_params_charge_control struct {
	mode uint32 /* enum charge_control_mode */
}

/*****************************************************************************/
/* Console commands. Only available when flash write protect is unlocked. */

/* Snapshot console output buffer for use by EC_CMD_CONSOLE_READ. */
const EC_CMD_CONSOLE_SNAPSHOT = 0x97

/*
 * Read next chunk of data from saved snapshot.
 *
 * Response is null-terminated string.  Empty string, if there is no more
 * remaining output.
 */
const EC_CMD_CONSOLE_READ = 0x98

/*****************************************************************************/

/*
 * Cut off battery power immediately or after the host has shut down.
 *
 * return EC_RES_INVALID_COMMAND if unsupported by a board/battery.
 *	  EC_RES_SUCCESS if the command was successful.
 *	  EC_RES_ERROR if the cut off command failed.
 */
const EC_CMD_BATTERY_CUT_OFF = 0x99

const EC_BATTERY_CUTOFF_FLAG_AT_SHUTDOWN = (1 << 0)

/* TYPE */
type ec_params_battery_cutoff struct {
	flags uint8
}

/*****************************************************************************/
/* USB port mux control. */

/*
 * Switch USB mux or return to automatic switching.
 */
const EC_CMD_USB_MUX = 0x9a

/* TYPE */
type ec_params_usb_mux struct {
	mux uint8
}

/*****************************************************************************/
/* LDOs / FETs control. */

type ec_ldo_state uint8

const (
	EC_LDO_STATE_OFF ec_ldo_state = 0 /* the LDO / FET is shut down */
	EC_LDO_STATE_ON               = 1 /* the LDO / FET is ON / providing power */
)

/*
 * Switch on/off a LDO.
 */
const EC_CMD_LDO_SET = 0x9b

/* TYPE */
type ec_params_ldo_set struct {
	index uint8
	state uint8
}

/*
 * Get LDO state.
 */
const EC_CMD_LDO_GET = 0x9c

/* TYPE */
type ec_params_ldo_get struct {
	index uint8
}

/* TYPE */
type ec_response_ldo_get struct {
	state uint8
}

/*****************************************************************************/
/* Power info. */

/*
 * Get power info.
 */
const EC_CMD_POWER_INFO = 0x9d

/* TYPE */
type ec_response_power_info struct {
	usb_devype        uint32
	voltage_ac        uint16
	voltage_system    uint16
	current_system    uint16
	usb_current_limit uint16
}

/*****************************************************************************/
/* I2C passthru command */

const EC_CMD_I2C_PASSTHRU = 0x9e

/* Read data; if not present, message is a write */
const EC_I2C_FLAG_READ = (1 << 15)

/* Mask for address */
const EC_I2C_ADDR_MASK = 0x3ff

const EC_I2C_STATUS_NAK = (1 << 0)     /* Transfer was not acknowledged */
const EC_I2C_STATUS_TIMEOUT = (1 << 1) /* Timeout during transfer */

/* Any error */
const EC_I2C_STATUS_ERROR = (EC_I2C_STATUS_NAK | EC_I2C_STATUS_TIMEOUT)

/* TYPE */
type ec_params_i2c_passthru_msg struct {
	addr_flags uint16 /* I2C slave address (7 or 10 bits) and flags */
	len        uint16 /* Number of bytes to read or write */
}

/* TYPE */
type ec_params_i2c_passthru struct {
	port     uint8 /* I2C port number */
	num_msgs uint8 /* Number of messages */
	msg      []ec_params_i2c_passthru_msg
	/* Data to write for all messages is concatenated here */
}

/* TYPE */
type ec_response_i2c_passthru struct {
	i2c_status uint8   /* Status flags (EC_I2C_STATUS_...) */
	num_msgs   uint8   /* Number of messages processed */
	data       []uint8 /* Data read by messages concatenated here */
}

/*****************************************************************************/
/* Power button hang detect */

const EC_CMD_HANG_DETECT = 0x9f

/* Reasons to start hang detection timer */
/* Power button pressed */
const EC_HANG_START_ON_POWER_PRESS = (1 << 0)

/* Lid closed */
const EC_HANG_START_ON_LID_CLOSE = (1 << 1)

/* Lid opened */
const EC_HANG_START_ON_LID_OPEN = (1 << 2)

/* Start of AP S3->S0 transition (booting or resuming from suspend) */
const EC_HANG_START_ON_RESUME = (1 << 3)

/* Reasons to cancel hang detection */

/* Power button released */
const EC_HANG_STOP_ON_POWER_RELEASE = (1 << 8)

/* Any host command from AP received */
const EC_HANG_STOP_ON_HOST_COMMAND = (1 << 9)

/* Stop on end of AP S0->S3 transition (suspending or shutting down) */
const EC_HANG_STOP_ON_SUSPEND = (1 << 10)

/*
 * If this flag is set, all the other fields are ignored, and the hang detect
 * timer is started.  This provides the AP a way to start the hang timer
 * without reconfiguring any of the other hang detect settings.  Note that
 * you must previously have configured the timeouts.
 */
const EC_HANG_START_NOW = (1 << 30)

/*
 * If this flag is set, all the other fields are ignored (including
 * EC_HANG_START_NOW).  This provides the AP a way to stop the hang timer
 * without reconfiguring any of the other hang detect settings.
 */
const EC_HANG_STOP_NOW = (1 << 31)

/* TYPE */
type ec_params_hang_detect struct {
	/* Flags; see EC_HANG_* */
	flags uint32

	/* Timeout in msec before generating host event, if enabled */
	host_eventimeout_msec uint16

	/* Timeout in msec before generating warm reboot, if enabled */
	warm_rebootimeout_msec uint16
}

/*****************************************************************************/
/* Commands for battery charging */

/*
 * This is the single catch-all host command to exchange data regarding the
 * charge state machine (v2 and up).
 */
const EC_CMD_CHARGE_STATE = 0xa0

/* Subcommands for this host command */
type charge_state_command uint8

const (
	CHARGE_STATE_CMD_GET_STATE charge_state_command = iota
	CHARGE_STATE_CMD_GET_PARAM
	CHARGE_STATE_CMD_SET_PARAM
	CHARGE_STATE_NUM_CMDS
)

/*
 * Known param numbers are defined here. Ranges are reserved for board-specific
 * params, which are handled by the particular implementations.
 */
type charge_state_params uint8

const (
	CS_PARAM_CHG_VOLTAGE       charge_state_params = iota /* charger voltage limit */
	CS_PARAM_CHG_CURRENT                                  /* charger current limit */
	CS_PARAM_CHG_INPUT_CURRENT                            /* charger input current limit */
	CS_PARAM_CHG_STATUS                                   /* charger-specific status */
	CS_PARAM_CHG_OPTION                                   /* charger-specific options */
	/* How many so far? */
	CS_NUM_BASE_PARAMS

	/* Range for CONFIG_CHARGER_PROFILE_OVERRIDE params */
	CS_PARAM_CUSTOM_PROFILE_MIN = 0x10000
	CS_PARAM_CUSTOM_PROFILE_MAX = 0x1ffff

	/* Other custom param ranges go here... */
)

/* TYPE */
/* ler
type ec_params_charge_state struct {
cmd uint8				/* enum charge_state_command * /
	union {
		struct {
			/* no args * /
		} get_state;

		struct {
		uint32 param;		/* enum charge_state_param * /
		} get_param;

		struct {
		uint32 param;		/* param to set * /
		uint32 value;		/* value to set * /
		} set_param;
	};
} ;

/* TYPE */
/* later
type ec_response_charge_state struct {
	union {
		struct {
			int ac;
			int chg_voltage;
			int chg_current;
			int chg_input_current;
			int batt_state_of_charge;
		} get_state;

		struct {
		uint32 value;
		} get_param;
		struct {
			/* no return values *
		} set_param;
	};
} ;
*/

/*
 * Set maximum battery charging current.
 */
const EC_CMD_CHARGE_CURRENT_LIMIT = 0xa1

/* TYPE */
type ec_params_current_limit struct {
	limit uint32 /* in mA */
}

/*
 * Set maximum external power current.
 */
const EC_CMD_EXT_POWER_CURRENT_LIMIT = 0xa2

/* TYPE */
type ec_params_ext_power_current_limit struct {
	limit uint32 /* in mA */
}

/*****************************************************************************/
/* Smart battery pass-through */

/* Get / Set 16-bit smart battery registers */
const EC_CMD_SB_READ_WORD = 0xb0
const EC_CMD_SB_WRITE_WORD = 0xb1

/* Get / Set string smart battery parameters
 * formatted as SMBUS "block".
 */
const EC_CMD_SB_READ_BLOCK = 0xb2
const EC_CMD_SB_WRITE_BLOCK = 0xb3

/* TYPE */
type ec_params_sb_rd struct {
	reg uint8
}

/* TYPE */
type ec_response_sb_rd_word struct {
	value uint16
}

/* TYPE */
type ec_params_sb_wr_word struct {
	reg   uint8
	value uint16
}

/* TYPE */
type ec_response_sb_rd_block struct {
	data [32]uint8
}

/* TYPE */
type ec_params_sb_wr_block struct {
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

const EC_CMD_BATTERY_VENDOR_PARAM = 0xb4

type ec_battery_vendor_param_mode uint8

const (
	BATTERY_VENDOR_PARAM_MODE_GET ec_battery_vendor_param_mode = 0
	BATTERY_VENDOR_PARAM_MODE_SET
)

/* TYPE */
type ec_params_battery_vendor_param struct {
	param uint32
	value uint32
	mode  uint8
}

/* TYPE */
type ec_response_battery_vendor_param struct {
	value uint32
}

/*****************************************************************************/
/*
 * Smart Battery Firmware Update Commands
 */
const EC_CMD_SB_FW_UPDATE = 0xb5

type ec_sb_fw_update_subcmd uint8

const (
	EC_SB_FW_UPDATE_PREPARE ec_sb_fw_update_subcmd = 0x0
	EC_SB_FW_UPDATE_INFO                           = 0x1 /*query sb info */
	EC_SB_FW_UPDATE_BEGIN                          = 0x2 /*check if protected */
	EC_SB_FW_UPDATE_WRITE                          = 0x3 /*check if protected */
	EC_SB_FW_UPDATE_END                            = 0x4
	EC_SB_FW_UPDATE_STATUS                         = 0x5
	EC_SB_FW_UPDATE_PROTECT                        = 0x6
	EC_SB_FW_UPDATE_MAX                            = 0x7
)

const SB_FW_UPDATE_CMD_WRITE_BLOCK_SIZE = 32
const SB_FW_UPDATE_CMD_STATUS_SIZE = 2
const SB_FW_UPDATE_CMD_INFO_SIZE = 8

/* TYPE */
type ec_sb_fw_update_header struct {
	subcmd uint16 /* enum ec_sb_fw_update_subcmd */
	fw_id  uint16 /* firmware id */
}

/* TYPE */
type ec_params_sb_fw_update struct {
	hdr ec_sb_fw_update_header
	/* no args. */
	/* EC_SB_FW_UPDATE_PREPARE  = 0x0 */
	/* EC_SB_FW_UPDATE_INFO     = 0x1 */
	/* EC_SB_FW_UPDATE_BEGIN    = 0x2 */
	/* EC_SB_FW_UPDATE_END      = 0x4 */
	/* EC_SB_FW_UPDATE_STATUS   = 0x5 */
	/* EC_SB_FW_UPDATE_PROTECT  = 0x6 */
	/* or ... */
	/* EC_SB_FW_UPDATE_WRITE    = 0x3 */
	data [SB_FW_UPDATE_CMD_WRITE_BLOCK_SIZE]uint8
}

/* TYPE */
type ec_response_sb_fw_update struct {
	data []uint8
	/* EC_SB_FW_UPDATE_INFO     = 0x1 */
	//uint8 data[SB_FW_UPDATE_CMD_INFO_SIZE];
	/* EC_SB_FW_UPDATE_STATUS   = 0x5 */
	//uint8 data[SB_FW_UPDATE_CMD_STATUS_SIZE];
}

/*
 * Entering Verified Boot Mode Command
 * Default mode is VBOOT_MODE_NORMAL if EC did not receive this command.
 * Valid Modes are: normal, developer, and recovery.
 */
const EC_CMD_ENTERING_MODE = 0xb6

/* TYPE */
type ec_params_entering_mode struct {
	vboot_mode int
}

const VBOOT_MODE_NORMAL = 0
const VBOOT_MODE_DEVELOPER = 1
const VBOOT_MODE_RECOVERY = 2

/*****************************************************************************/
/* System commands */

/*
 * TODO(crosbug.com/p/23747): This is a confusing name, since it doesn't
 * necessarily reboot the EC.  Rename to "image" or something similar?
 */
const EC_CMD_REBOOT_EC = 0xd2

/* Command */
type ec_reboot_cmd uint8

const (
	EC_REBOOT_CANCEL  ec_reboot_cmd = 0 /* Cancel a pending reboot */
	EC_REBOOT_JUMP_RO               = 1 /* Jump to RO without rebooting */
	EC_REBOOT_JUMP_RW               = 2 /* Jump to RW without rebooting */
	/* (command 3 was jump to RW-B) */
	EC_REBOOT_COLD         = 4 /* Cold-reboot */
	EC_REBOOT_DISABLE_JUMP = 5 /* Disable jump until next reboot */
	EC_REBOOT_HIBERNATE    = 6 /* Hibernate EC */
)

/* Flags for ec_params_reboot_ec.reboot_flags */
const EC_REBOOT_FLAG_RESERVED0 = (1 << 0)      /* Was recovery request */
const EC_REBOOT_FLAG_ON_AP_SHUTDOWN = (1 << 1) /* Reboot after AP shutdown */

/* TYPE */
type ec_params_reboot_ec struct {
	cmd   uint8 /* enum ec_reboot_cmd */
	flags uint8 /* See EC_REBOOT_FLAG_* */
}

/*
 * Get information on last EC panic.
 *
 * Returns variable-length platform-dependent panic information.  See panic.h
 * for details.
 */
const EC_CMD_GET_PANIC_INFO = 0xd3

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
 * Use EC_CMD_REBOOT_EC to reboot the EC more politely.
 */
const EC_CMD_REBOOT = 0xd1 /* Think "die" */

/*
 * Resend last response (not supported on LPC).
 *
 * Returns EC_RES_UNAVAILABLE if there is no response available - for example
 * there was no previous command, or the previous command's response was too
 * big to save.
 */
const EC_CMD_RESEND_RESPONSE = 0xdb

/*
 * This header byte on a command indicate version 0. Any header byte less
 * than this means that we are talking to an old EC which doesn't support
 * versioning. In that case, we assume version 0.
 *
 * Header bytes greater than this indicate a later version. For example
 * EC_CMD_VERSION0 + 1 means we are using version 1.
 *
 * The old EC interface must not use commands 0xdc or higher.
 */
const EC_CMD_VERSION0 = 0xdc

/*****************************************************************************/
/*
 * PD commands
 *
 * These commands are for PD MCU communication.
 */

/* EC to PD MCU exchange status command */
const EC_CMD_PD_EXCHANGE_STATUS = 0x100

type pd_charge_state uint8

const (
	PD_CHARGE_NO_CHANGE pd_charge_state = 0 /* Don't change charge state */
	PD_CHARGE_NONE                          /* No charging allowed */
	PD_CHARGE_5V                            /* 5V charging only */
	PD_CHARGE_MAX                           /* Charge at max voltage */
)

/* TYPE */
/* Status of EC being sent to PD */
type ec_params_pd_status struct {
	batt_soc     int8  /* battery state of charge */
	charge_state uint8 /* charging state (from enum pd_charge_state) */
}

/* Status of PD being sent back to EC */
const PD_STATUS_HOST_EVENT = (1 << 0)      /* Forward host event to AP */
const PD_STATUS_IN_RW = (1 << 1)           /* Running RW image */
const PD_STATUS_JUMPED_TO_IMAGE = (1 << 2) /* Current image was jumped to */
/* TYPE */
type ec_response_pd_status struct {
	status             uint32 /* PD MCU status */
	curr_lim_ma        uint32 /* input current limit */
	active_charge_port int32  /* active charging port */
}

/* AP to PD MCU host event status command, cleared on read */
const EC_CMD_PD_HOST_EVENT_STATUS = 0x104

/* PD MCU host event status bits */
const PD_EVENT_UPDATE_DEVICE = (1 << 0)
const PD_EVENT_POWER_CHANGE = (1 << 1)
const PD_EVENT_IDENTITY_RECEIVED = (1 << 2)

/* TYPE */
type ec_response_host_event_status struct {
	status uint32 /* PD MCU host event status */
}

/* Set USB type-C port role and muxes */
const EC_CMD_USB_PD_CONTROL = 0x101

type usb_pd_control_role uint8

const (
	USB_PD_CTRL_ROLE_NO_CHANGE    usb_pd_control_role = 0
	USB_PD_CTRL_ROLE_TOGGLE_ON                        = 1 /* == AUTO */
	USB_PD_CTRL_ROLE_TOGGLE_OFF                       = 2
	USB_PD_CTRL_ROLE_FORCE_SINK                       = 3
	USB_PD_CTRL_ROLE_FORCE_SOURCE                     = 4
	USB_PD_CTRL_ROLE_COUNT
)

type usb_pd_control_mux uint8

const (
	USB_PD_CTRL_MUX_NO_CHANGE usb_pd_control_mux = 0
	USB_PD_CTRL_MUX_NONE                         = 1
	USB_PD_CTRL_MUX_USB                          = 2
	USB_PD_CTRL_MUX_DP                           = 3
	USB_PD_CTRL_MUX_DOCK                         = 4
	USB_PD_CTRL_MUX_AUTO                         = 5
	USB_PD_CTRL_MUX_COUNT
)

/* TYPE */
type ec_params_usb_pd_control struct {
	port uint8
	role uint8
	mux  uint8
}

/* TYPE */
type ec_response_usb_pd_control struct {
	enabled  uint8
	role     uint8
	polarity uint8
	state    uint8
}

/* TYPE */
type ec_response_usb_pd_control_v1 struct {
	enabled  uint8
	role     uint8 /* [0] power: 0=SNK/1=SRC [1] data: 0=UFP/1=DFP */
	polarity uint8
	state    [32]byte
}

const EC_CMD_USB_PD_PORTS = 0x102

/* TYPE */
type ec_response_usb_pd_ports struct {
	num_ports uint8
}

const EC_CMD_USB_PD_POWER_INFO = 0x103

const PD_POWER_CHARGING_PORT = 0xff

/* TYPE */
type ec_params_usb_pd_power_info struct {
	port uint8
}

type usb_chg_type uint8

const (
	USB_CHG_TYPE_NONE usb_chg_type = iota
	USB_CHG_TYPE_PD
	USB_CHG_TYPE_C
	USB_CHG_TYPE_PROPRIETARY
	USB_CHG_TYPE_BC12_DCP
	USB_CHG_TYPE_BC12_CDP
	USB_CHG_TYPE_BC12_SDP
	USB_CHG_TYPE_OTHER
	USB_CHG_TYPE_VBUS
	USB_CHG_TYPE_UNKNOWN
)

type usb_power_roles uint8

const (
	USB_PD_PORT_POWER_DISCONNECTED usb_power_roles = iota
	USB_PD_PORT_POWER_SOURCE
	USB_PD_PORT_POWER_SINK
	USB_PD_PORT_POWER_SINK_NOT_CHARGING
)

/* TYPE */
type usb_chg_measures struct {
	voltage_max uint16
	voltage_now uint16
	current_max uint16
	current_lim uint16
}

/* TYPE */
type ec_response_usb_pd_power_info struct {
	role      uint8
	etype     uint8
	dualrole  uint8
	reserved1 uint8
	meas      usb_chg_measures
	max_power uint32
}

/* Write USB-PD device FW */
const EC_CMD_USB_PD_FW_UPDATE = 0x110

type usb_pd_fw_update_cmds uint8

const (
	USB_PD_FW_REBOOT usb_pd_fw_update_cmds = iota
	USB_PD_FW_FLASH_ERASE
	USB_PD_FW_FLASH_WRITE
	USB_PD_FW_ERASE_SIG
)

/* TYPE */
type ec_params_usb_pd_fw_update struct {
	dev_id uint16
	cmd    uint8
	port   uint8
	size   uint32 /* Size to write in bytes */
	/* Followed by data to write */
}

/* Write USB-PD Accessory RW_HASH table entry */
const EC_CMD_USB_PD_RW_HASH_ENTRY = 0x111

/* RW hash is first 20 bytes of SHA-256 of RW section */
const PD_RW_HASH_SIZE = 20

/* TYPE */
type ec_params_usb_pd_rw_hash_entry struct {
	dev_id        uint16
	dev_rw_hash   [PD_RW_HASH_SIZE]uint8
	reserved      uint8  /* For alignment of current_image */
	current_image uint32 /* One of ec_current_image */
}

/* Read USB-PD Accessory info */
const EC_CMD_USB_PD_DEV_INFO = 0x112

/* TYPE */
type ec_params_usb_pd_info_request struct {
	port uint8
}

/* Read USB-PD Device discovery info */
const EC_CMD_USB_PD_DISCOVERY = 0x113

/* TYPE */
type ec_params_usb_pd_discovery_entry struct {
	vid   uint16 /* USB-IF VID */
	pid   uint16 /* USB-IF PID */
	ptype uint8  /* product type (hub,periph,cable,ama) */
}

/* Override default charge behavior */
const EC_CMD_PD_CHARGE_PORT_OVERRIDE = 0x114

/* Negative port parameters have special meaning */
type usb_pd_override_ports int8

const (
	OVERRIDE_DONT_CHARGE usb_pd_override_ports = -2
	OVERRIDE_OFF                               = -1
	/* [0, PD_PORT_COUNT): Port# */
)

/* TYPE */
type ec_params_charge_port_override struct {
	override_port int16 /* Override port# */
}

/* Read (and delete) one entry of PD event log */
const EC_CMD_PD_GET_LOG_ENTRY = 0x115

/* TYPE */
type ec_response_pd_log struct {
	timestamp uint32  /* relative timestamp in milliseconds */
	etype     uint8   /* event type : see PD_EVENT_xx below */
	size_port uint8   /* [7:5] port number [4:0] payload size in bytes */
	data      uint16  /* type-defined data payload */
	payload   []uint8 /* optional additional data payload: 0..16 bytes */
}

/* The timestamp is the microsecond counter shifted to get about a ms. */
const PD_LOG_TIMESTAMP_SHIFT = 10 /* 1 LSB = 1024us */

const PD_LOG_SIZE_MASK = 0x1F
const PD_LOG_PORT_MASK = 0xE0
const PD_LOG_PORT_SHIFT = 5

func pd_log_port_size(port, size uint8) uint8 {
	return (port << PD_LOG_PORT_SHIFT) | (size & PD_LOG_SIZE_MASK)
}

func pd_log_port(size_port uint8) uint8 {
	return size_port >> PD_LOG_PORT_SHIFT
}
func pd_log_size(size_port uint8) uint8 {
	return size_port & PD_LOG_SIZE_MASK
}

/* PD event log : entry types */
/* PD MCU events */
const PD_EVENT_MCU_BASE = 0x00
const PD_EVENT_MCU_CHARGE = (PD_EVENT_MCU_BASE + 0)
const PD_EVENT_MCU_CONNECT = (PD_EVENT_MCU_BASE + 1)

/* Reserved for custom board event */
const PD_EVENT_MCU_BOARD_CUSTOM = (PD_EVENT_MCU_BASE + 2)

/* PD generic accessory events */
const PD_EVENT_ACC_BASE = 0x20
const PD_EVENT_ACC_RW_FAIL = (PD_EVENT_ACC_BASE + 0)
const PD_EVENT_ACC_RW_ERASE = (PD_EVENT_ACC_BASE + 1)

/* PD power supply events */
const PD_EVENT_PS_BASE = 0x40
const PD_EVENT_PS_FAULT = (PD_EVENT_PS_BASE + 0)

/* PD video dongles events */
const PD_EVENT_VIDEO_BASE = 0x60
const PD_EVENT_VIDEO_DP_MODE = (PD_EVENT_VIDEO_BASE + 0)
const PD_EVENT_VIDEO_CODEC = (PD_EVENT_VIDEO_BASE + 1)

/* Returned in the "type" field, when there is no entry available */
const PD_EVENT_NO_ENTRY = 0xFF

/*
 * PD_EVENT_MCU_CHARGE event definition :
 * the payload is "struct usb_chg_measures"
 * the data field contains the port state flags as defined below :
 */
/* Port partner is a dual role device */
const CHARGE_FLAGS_DUAL_ROLE = (1 << 15)

/* Port is the pending override port */
const CHARGE_FLAGS_DELAYED_OVERRIDE = (1 << 14)

/* Port is the override port */
const CHARGE_FLAGS_OVERRIDE = (1 << 13)

/* Charger type */
const CHARGE_FLAGS_TYPE_SHIFT = 3
const CHARGE_FLAGS_TYPE_MASK = (0xF << CHARGE_FLAGS_TYPE_SHIFT)

/* Power delivery role */
const CHARGE_FLAGS_ROLE_MASK = (7 << 0)

/*
 * PD_EVENT_PS_FAULT data field flags definition :
 */
const PS_FAULT_OCP = 1
const PS_FAULT_FAST_OCP = 2
const PS_FAULT_OVP = 3
const PS_FAULT_DISCH = 4

/* TYPE */
/*
 * PD_EVENT_VIDEO_CODEC payload is "struct mcdp_info".
 */
type mcdp_version struct {
	major uint8
	minor uint8
	build uint16
}

/* TYPE */
type mcdp_info struct {
	family [2]uint8
	chipid [2]uint8
	irom   mcdp_version
	fw     mcdp_version
}

/* struct mcdp_info field decoding */
func mcdp_chipid(chipid []uint8) uint16 {
	return (uint16(chipid[0]) << 8) | uint16(chipid[1])
}

func mcdp_family(family []uint8) uint16 {
	return (uint16(family[0]) << 8) | uint16(family[1])
}

/* Get/Set USB-PD Alternate mode info */
const EC_CMD_USB_PD_GET_AMODE = 0x116

/* TYPE */
type ec_params_usb_pd_get_mode_request struct {
	svid_idx uint16 /* SVID index to get */
	port     uint8  /* port */
}

/* TYPE */
type ec_params_usb_pd_get_mode_response struct {
	svid uint16    /* SVID */
	opos uint16    /* Object Position */
	vdo  [6]uint32 /* Mode VDOs */
}

const EC_CMD_USB_PD_SET_AMODE = 0x117

type pd_mode_cmd uint8

const (
	PD_EXIT_MODE  pd_mode_cmd = 0
	PD_ENTER_MODE             = 1
	/* Not a command.  Do NOT remove. */
	PD_MODE_CMD_COUNT
)

/* TYPE */
type ec_params_usb_pd_set_mode_request struct {
	cmd  uint32 /* enum pd_mode_cmd */
	svid uint16 /* SVID to set */
	opos uint8  /* Object Position */
	port uint8  /* port */
}

/* Ask the PD MCU to record a log of a requested type */
const EC_CMD_PD_WRITE_LOG_ENTRY = 0x118

/* TYPE */
type ec_params_pd_write_log_entry struct {
	etype uint8 /* event type : see PD_EVENT_xx above */
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
func ec_cmd_passthru_offset(n uint) uint {
	return 0x4000 * n
}

func ec_cmd_passthru_max(n uint) uint {
	return ec_cmd_passthru_offset(n) + 0x3fff
}

/*****************************************************************************/
/*
 * Deprecated constants. These constants have been renamed for clarity. The
 * meaning and size has not changed. Programs that use the old names should
 * switch to the new names soon, as the old names may not be carried forward
 * forever.
 */
const EC_HOST_PARAM_SIZE = EC_PROTO2_MAX_PARAM_SIZE
const EC_LPC_ADDR_OLD_PARAM = EC_HOST_CMD_REGION1
const EC_OLD_PARAM_SIZE = EC_HOST_CMD_REGION_SIZE

func ec_ver_mask(version uint8) uint8 {
	/* Command version mask */
	return 1 << version
}
