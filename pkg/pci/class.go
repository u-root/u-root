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

	ClassOtherS = 0xff
)

// ClassNames maps class names from PCI sysfs to a name.
var ClassNames = map[string]string{
	"000000": "NotDefined",
	"000100": "NotDefinedVGA",

	"01":     "Storage",
	"010000": "StorageSCSI",
	"010100": "StorageIDE",
	"010200": "StorageFLOPPY",
	"010300": "StorageIPI",
	"010400": "StorageRAID",
	"010500": "StorageATA",
	"010600": "StorageSATA",
	"010700": "StorageSAS",
	"018000": "StorageOther",

	"02":     "Network",
	"020000": "NetworkEthernet",
	"028000": "NetworkOther",

	"03":     "Display",
	"030000": "DisplayVGA",
	"030100": "DisplayXGA",
	"030200": "Display3D",
	"038000": "DisplayOther",

	"04":     "Multimedia",
	"040000": "MultimediaVideo",
	"040100": "MultimediaAudio",
	"040200": "MultimediaPhone",
	"040300": "MultimediaAudioDev",
	"048000": "MultimediaOther",

	"05":     "Memory",
	"050000": "MemoryRAM",
	"050100": "MemoryFLASH",
	"058000": "MemoryOther",

	"06":     "Bridge",
	"060000": "BridgeHost",
	"060100": "BridgeISA",
	"060200": "BridgeEISA",
	"060300": "BridgeMC",
	"060400": "BridgePCI",
	"060500": "BridgePCMCIA",
	"060600": "BridgeNUBUS",
	"060700": "BridgeCARDBUS",
	"060800": "BridgeRACEWAY",
	"060900": "BridgePCISemi",
	"060a00": "BridgeIBToPCI",
	"068000": "BridgeOther",

	"07":     "Communication",
	"070000": "CommunicationSerial",
	"070100": "CommunicationParallel",
	"070200": "CommunicationMSerial",
	"070300": "CommunicationModem",
	"078000": "CommunicationOther",

	"08":     "System",
	"080000": "SystemPIC",
	"080100": "SystemDMA",
	"080200": "SystemTimer",
	"080300": "SystemRTC",
	"080400": "SystemPCIHotplug",
	"088000": "SystemOther",

	"09":     "Input",
	"090000": "InputKeyboard",
	"090100": "InputPen",
	"090200": "InputMouse",
	"090300": "InputScanner",
	"090400": "InputGameport",
	"098000": "InputOther",

	"0a":     "Docking",
	"0a0000": "DockingGeneric",
	"0a8000": "DockingOther",

	"0b":     "Processor",
	"0b0000": "Processor386",
	"0b0100": "Processor486",
	"0b0200": "ProcessorPentium",
	"0b1000": "ProcessorALPHA",
	"0b2000": "ProcessorPOWERPC",
	"0b3000": "ProcessorMIPS",
	"0b4000": "ProcessorCO",

	"0c":     "Serial",
	"0c0000": "SerialFirewire",
	"0c0100": "SerialAccess",
	"0c0200": "SerialSSA",
	"0c0300": "SerialUSB",
	"0c0400": "SerialFIBER",
	"0c0500": "SerialSMBUS",
	"0c0600": "SerialINFINIBAND",

	"0d":     "Wireless",
	"0d0000": "WirelessIRDA",
	"0d0100": "WirelessCONSUMERIR",
	"0d1000": "WirelessRF",
	"0d8000": "WirelessOther",

	"0f":     "Satellite",
	"0f0000": "SatelliteTV",
	"0f0100": "SatelliteAudio",
	"0f0300": "SatelliteVoice",
	"0f0400": "SatelliteData",

	"10":     "Crypt",
	"100000": "CryptNetwork",
	"101000": "CryptEntertainment",
	"108000": "CryptOther",

	"11":     "Signal",
	"110000": "SignalDPIO",
	"110100": "SignalPERFCTR",
	"111000": "SignalSynchronizer",
	"118000": "SignalOther",

	"ff": "OtherS",
}
