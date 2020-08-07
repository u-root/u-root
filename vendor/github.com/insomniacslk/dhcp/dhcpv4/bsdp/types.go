package bsdp

import (
	"fmt"
)

// DefaultMacOSVendorClassIdentifier is a default vendor class identifier used
// on non-darwin hosts where the vendor class identifier cannot be determined.
// It should mostly be used for debugging if testing BSDP on a non-darwin
// system.
const DefaultMacOSVendorClassIdentifier = AppleVendorID + "/i386/MacMini6,1"

// optionCode are BSDP option codes.
//
// optionCode implements the dhcpv4.OptionCode interface.
type optionCode uint8

func (o optionCode) Code() uint8 {
	return uint8(o)
}

func (o optionCode) String() string {
	if s, ok := optionCodeToString[o]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", o)
}

// Options (occur as sub-options of DHCP option 43).
const (
	OptionMessageType                   optionCode = 1
	OptionVersion                       optionCode = 2
	OptionServerIdentifier              optionCode = 3
	OptionServerPriority                optionCode = 4
	OptionReplyPort                     optionCode = 5
	OptionBootImageListPath             optionCode = 6 // Not used
	OptionDefaultBootImageID            optionCode = 7
	OptionSelectedBootImageID           optionCode = 8
	OptionBootImageList                 optionCode = 9
	OptionNetboot1_0Firmware            optionCode = 10
	OptionBootImageAttributesFilterList optionCode = 11
	OptionShadowMountPath               optionCode = 128
	OptionShadowFilePath                optionCode = 129
	OptionMachineName                   optionCode = 130
)

// optionCodeToString maps BSDP OptionCodes to human-readable strings
// describing what they are.
var optionCodeToString = map[optionCode]string{
	OptionMessageType:                   "BSDP Message Type",
	OptionVersion:                       "BSDP Version",
	OptionServerIdentifier:              "BSDP Server Identifier",
	OptionServerPriority:                "BSDP Server Priority",
	OptionReplyPort:                     "BSDP Reply Port",
	OptionBootImageListPath:             "", // Not used
	OptionDefaultBootImageID:            "BSDP Default Boot Image ID",
	OptionSelectedBootImageID:           "BSDP Selected Boot Image ID",
	OptionBootImageList:                 "BSDP Boot Image List",
	OptionNetboot1_0Firmware:            "BSDP Netboot 1.0 Firmware",
	OptionBootImageAttributesFilterList: "BSDP Boot Image Attributes Filter List",
	OptionShadowMountPath:               "BSDP Shadow Mount Path",
	OptionShadowFilePath:                "BSDP Shadow File Path",
	OptionMachineName:                   "BSDP Machine Name",
}
