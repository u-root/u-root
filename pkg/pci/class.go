// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

// Class definitions for PCI.
const (
	ClassNotDefined    = 0x0000
	ClassNotDefinedVGA = 0x0001

	ClassStorage       = 0x01
	ClassStorageSCSI   = 0x0100
	ClassStorageIDE    = 0x0101
	ClassStorageFLOPPY = 0x0102
	ClassStorageIPI    = 0x0103
	ClassStorageRAID   = 0x0104
	ClassStorageATA    = 0x0105
	ClassStorageSATA   = 0x0106
	ClassStorageSAS    = 0x0107
	ClassStorageOther  = 0x0180

	ClassNetwork         = 0x02
	ClassNetworkEthernet = 0x0200
	ClassNetworkOther    = 0x0280

	ClassDisplay      = 0x03
	ClassDisplayVGA   = 0x0300
	ClassDisplayXGA   = 0x0301
	ClassDisplay3D    = 0x0302
	ClassDisplayOther = 0x0380

	ClassMultimedia         = 0x04
	ClassMultimediaVideo    = 0x0400
	ClassMultimediaAudio    = 0x0401
	ClassMultimediaPhone    = 0x0402
	ClassMultimediaAudioDev = 0x0403
	ClassMultimediaOther    = 0x0480

	ClassMemory      = 0x05
	ClassMemoryRAM   = 0x0500
	ClassMemoryFLASH = 0x0501
	ClassMemoryOther = 0x0580

	ClassBridge        = 0x06
	ClassBridgeHost    = 0x0600
	ClassBridgeISA     = 0x0601
	ClassBridgeEISA    = 0x0602
	ClassBridgeMC      = 0x0603
	ClassBridgePCI     = 0x0604
	ClassBridgePCMCIA  = 0x0605
	ClassBridgeNUBUS   = 0x0606
	ClassBridgeCARDBUS = 0x0607
	ClassBridgeRACEWAY = 0x0608
	ClassBridgePCISemi = 0x0609
	ClassBridgeIBToPCI = 0x060a
	ClassBridgeOther   = 0x0680

	ClassCommunication         = 0x07
	ClassCommunicationSerial   = 0x0700
	ClassCommunicationParallel = 0x0701
	ClassCommunicationMSerial  = 0x0702
	ClassCommunicationModem    = 0x0703
	ClassCommunicationOther    = 0x0780

	ClassSystem           = 0x08
	ClassSystemPIC        = 0x0800
	ClassSystemDMA        = 0x0801
	ClassSystemTimer      = 0x0802
	ClassSystemRTC        = 0x0803
	ClassSystemPCIHotplug = 0x0804
	ClassSystemOther      = 0x0880

	ClassInput         = 0x09
	ClassInputKeyboard = 0x0900
	ClassInputPen      = 0x0901
	ClassInputMouse    = 0x0902
	ClassInputScanner  = 0x0903
	ClassInputGameport = 0x0904
	ClassInputOther    = 0x0980

	ClassDocking        = 0x0a
	ClassDockingGeneric = 0x0a00
	ClassDockingOther   = 0x0a80

	ClassProcessor        = 0x0b
	ClassProcessor386     = 0x0b00
	ClassProcessor486     = 0x0b01
	ClassProcessorPentium = 0x0b02
	ClassProcessorALPHA   = 0x0b10
	ClassProcessorPOWERPC = 0x0b20
	ClassProcessorMIPS    = 0x0b30
	ClassProcessorCO      = 0x0b40

	ClassSerial           = 0x0c
	ClassSerialFirewire   = 0x0c00
	ClassSerialAccess     = 0x0c01
	ClassSerialSSA        = 0x0c02
	ClassSerialUSB        = 0x0c03
	ClassSerialFIBER      = 0x0c04
	ClassSerialSMBUS      = 0x0c05
	ClassSerialINFINIBAND = 0x0c06

	ClassWireless           = 0x0d
	ClassWirelessIRDA       = 0x0d00
	ClassWirelessCONSUMERIR = 0x0d01
	ClassWirelessRF         = 0x0d10
	ClassWirelessOther      = 0x0d80

	ClassSatellite      = 0x0f
	ClassSatelliteTV    = 0x0f00
	ClassSatelliteAudio = 0x0f01
	ClassSatelliteVoice = 0x0f03
	ClassSatelliteData  = 0x0f04

	ClassCrypt              = 0x10
	ClassCryptNetwork       = 0x1000
	ClassCryptEntertainment = 0x1010
	ClassCryptOther         = 0x1080

	ClassSignal             = 0x11
	ClassSignalDPIO         = 0x1100
	ClassSignalPERFCTR      = 0x1101
	ClassSignalSynchronizer = 0x1110
	ClassSignalOther        = 0x1180

	ClassOthers = 0xff
)

// ClassNames maps class names from PCI sysfs to a name.
var ClassNames = map[uint32]string{
	0x000000: "NotDefined",
	0x000100: "NotDefinedVGA",

	0x01:     "Storage",
	0x010000: "StorageSCSI",
	0x010100: "StorageIDE",
	0x010200: "StorageFLOPPY",
	0x010300: "StorageIPI",
	0x010400: "StorageRAID",
	0x010500: "StorageATA",
	0x010600: "StorageSATA",
	0x010601: "StorageAHCI",
	0x010700: "StorageSAS",
	0x010802: "StorageNVMHCI",
	0x018000: "StorageOther",

	0x02:     "Network",
	0x020000: "NetworkEthernet",
	0x028000: "NetworkOther",

	0x03:     "Display",
	0x030000: "DisplayVGA",
	0x030100: "DisplayXGA",
	0x030200: "Display3D",
	0x038000: "DisplayOther",

	0x04:     "Multimedia",
	0x040000: "MultimediaVideo",
	0x040100: "MultimediaAudio",
	0x040200: "MultimediaPhone",
	0x040300: "MultimediaAudioDev",
	0x048000: "MultimediaOther",

	0x05:     "Memory",
	0x050000: "MemoryRAM",
	0x050100: "MemoryFLASH",
	0x058000: "MemoryOther",

	0x06:     "Bridge",
	0x060000: "BridgeHost",
	0x060100: "BridgeISA",
	0x060200: "BridgeEISA",
	0x060300: "BridgeMC",
	0x060400: "BridgePCI",
	0x060500: "BridgePCMCIA",
	0x060600: "BridgeNUBUS",
	0x060700: "BridgeCARDBUS",
	0x060800: "BridgeRACEWAY",
	0x060900: "BridgePCISemi",
	0x060a00: "BridgeIBToPCI",
	0x068000: "BridgeOther",

	0x07:     "Communication",
	0x070000: "CommunicationSerial",
	0x070100: "CommunicationParallel",
	0x070200: "CommunicationMSerial",
	0x070300: "CommunicationModem",
	0x078000: "CommunicationOther",

	0x08:     "System",
	0x080000: "SystemPIC",
	0x080100: "SystemDMA",
	0x080200: "SystemTimer",
	0x080300: "SystemRTC",
	0x080400: "SystemPCIHotplug",
	0x088000: "SystemOther",

	0x09:     "Input",
	0x090000: "InputKeyboard",
	0x090100: "InputPen",
	0x090200: "InputMouse",
	0x090300: "InputScanner",
	0x090400: "InputGameport",
	0x098000: "InputOther",

	0x0a:     "Docking",
	0x0a0000: "DockingGeneric",
	0x0a8000: "DockingOther",

	0x0b:     "Processor",
	0x0b0000: "Processor386",
	0x0b0100: "Processor486",
	0x0b0200: "ProcessorPentium",
	0x0b1000: "ProcessorALPHA",
	0x0b2000: "ProcessorPOWERPC",
	0x0b3000: "ProcessorMIPS",
	0x0b4000: "ProcessorCO",

	0x0c:     "Serial",
	0x0c0000: "SerialFirewire",
	0x0c0100: "SerialAccess",
	0x0c0200: "SerialSSA",
	0x0c0300: "SerialUSB",
	0x0c0400: "SerialFIBER",
	0x0c0500: "SerialSMBUS",
	0x0c0600: "SerialINFINIBAND",

	0x0d:     "Wireless",
	0x0d0000: "WirelessIRDA",
	0x0d0100: "WirelessCONSUMERIR",
	0x0d1000: "WirelessRF",
	0x0d8000: "WirelessOther",

	0x0f:     "Satellite",
	0x0f0000: "SatelliteTV",
	0x0f0100: "SatelliteAudio",
	0x0f0300: "SatelliteVoice",
	0x0f0400: "SatelliteData",

	0x10:     "Crypt",
	0x100000: "CryptNetwork",
	0x101000: "CryptEntertainment",
	0x108000: "CryptOther",

	0x11:     "Signal",
	0x110000: "SignalDPIO",
	0x110100: "SignalPERFCTR",
	0x111000: "SignalSynchronizer",
	0x118000: "SignalOther",

	0xff: "Others",
}
