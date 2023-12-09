// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package txtlog provides reading/parsing of Intel TXT logs.
// Huge parts were taken from 9elements/tpmtool
package txtlog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf16"

	tss "github.com/u-root/u-root/pkg/tss"
)

/*
[1] TCG EFI Platform Specification For TPM Family 1.1 or 1.2
https://trustedcomputinggroup.org/wp-content/uploads/TCG_EFI_Platform_1_22_Final_-v15.pdf

[2] TCG PC Client Specific Implementation Specification for Conventional BIOS", version 1.21
https://trustedcomputinggroup.org/wp-content/uploads/TCG_PCClientImplementation_1-21_1_00.pdf

[3] TCG EFI Protocol Specification, Family "2.0"
https://trustedcomputinggroup.org/wp-content/uploads/EFI-Protocol-Specification-rev13-160330final.pdf

[4] TCG PC Client Platform Firmware Profile Specification
https://trustedcomputinggroup.org/wp-content/uploads/PC-ClientSpecific_Platform_Profile_for_TPM_2p0_Systems_v51.pdf
*/
var (
	// DefaultTCPABinaryLog log file where the TCPA log is stored
	DefaultTCPABinaryLog = "/sys/kernel/security/tpm0/binary_bios_measurements"
)

var HashAlgoToSize = map[IAlgHash]IAlgHashSize{
	TPMAlgSha:     TPMAlgShaSize,
	TPMAlgSha256:  TPMAlgSha256Size,
	TPMAlgSha384:  TPMAlgSha384Size,
	TPMAlgSha512:  TPMAlgSha512Size,
	TPMAlgSm3s256: TPMAlgSm3s256Size,
}

func ParseLog(firmware FirmwareType, tpmSpec tss.TPMVersion) (*PCRLog, error) {
	var pcrLog *PCRLog
	var err error

	switch tpmSpec {
	case tss.TPMVersion12:
		pcrLog, err = readTPM1Log(firmware)
		if err != nil {
			return nil, err
		}
	case tss.TPMVersion20:
		pcrLog, err = readTPM2Log(firmware)
		if err != nil {
			// Kernel eventlog workaround does not export agile measurement log..
			pcrLog, err = readTPM1Log(firmware)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, errors.New("no valid TPM specification found")
	}

	return pcrLog, nil
}

func DumpLog(tcpaLog *PCRLog) error {
	for _, pcr := range tcpaLog.PcrList {
		fmt.Printf("%s\n", pcr)

		fmt.Println()
	}

	return nil
}

func readTPM1Log(firmware FirmwareType) (*PCRLog, error) {
	var pcrLog PCRLog

	file, err := os.Open(DefaultTCPABinaryLog)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pcrLog.Firmware = firmware

	if firmware == "TXT" {
		var pcrLog PCRLog

		container, err := readTxtEventLogContainer(file)
		if err != nil {
			return nil, err
		}

		// seek to first PCR event
		if _, err := file.Seek(int64(container.PcrEventsOffset), io.SeekStart); err != nil {
			return nil, err
		}

		for {
			offset, err := file.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, err
			}

			if offset >= int64(container.NextEventOffset) {
				break
			}

			pcrEvent, err := parseTcgPcrEvent(file)
			if err != nil {
				// NB: error out even for EOF because it should
				//     not be seen before NextEventOffset
				return nil, err
			}

			pcrLog.PcrList = append(pcrLog.PcrList, pcrEvent)
		}
	} else {
		for {
			pcrEvent, err := parseTcgPcrEvent(file)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			pcrLog.PcrList = append(pcrLog.PcrList, pcrEvent)
		}
	}

	return &pcrLog, nil
}

func readTPM2Log(firmware FirmwareType) (*PCRLog, error) {
	var pcrLog PCRLog
	var pcrEvent *TcgPcrEvent

	file, err := os.Open(DefaultTCPABinaryLog)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pcrLog.Firmware = firmware

	if pcrEvent, err = parseTcgPcrEvent(file); err != nil {
		return nil, err
	}
	if efiSpecID, err := parseEfiSpecEvent(bytes.NewBuffer(pcrEvent.event)); efiSpecID == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("first event was not an EFI SpecID Event")
	}

	pcrLog.PcrList = append(pcrLog.PcrList, pcrEvent)

	for {
		pcrEvent, err := parseTcgPcrEvent2(file)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// There may be times when give part of the buffer past the last event,
		// when that is the case just check to see if the event type is zero (reserved)
		if pcrEvent.eventType == 0 {
			break
		}
		pcrLog.PcrList = append(pcrLog.PcrList, pcrEvent)
	}

	return &pcrLog, nil
}

func getTaggedEvent(eventData []byte) (*string, error) {
	eventReader := bytes.NewReader(eventData)
	var taggedEvent TCGPCClientTaggedEvent

	if err := binary.Read(eventReader, binary.LittleEndian, &taggedEvent.taggedEventID); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &taggedEvent.taggedEventDataSize); err != nil {
		return nil, err
	}

	taggedEvent.taggedEventData = make([]byte, taggedEvent.taggedEventDataSize)
	if err := binary.Read(eventReader, binary.LittleEndian, &taggedEvent.taggedEventData); err != nil {
		return nil, err
	}

	eventInfo := fmt.Sprintf("Tag ID - %d - %s", taggedEvent.taggedEventID, string(taggedEvent.taggedEventData))
	return &eventInfo, nil
}

func getHandoffTablePointers(eventData []byte) (*string, error) {
	eventReader := bytes.NewReader(eventData)
	var handoffTablePointers EFIHandoffTablePointers

	if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.numberOfTables); err != nil {
		return nil, err
	}

	handoffTablePointers.tableEntry = make([]EFIConfigurationTable, handoffTablePointers.numberOfTables)
	for i := uint64(0); i < handoffTablePointers.numberOfTables; i++ {
		if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.tableEntry[i].vendorGUID.blockA); err != nil {
			return nil, err
		}

		if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.tableEntry[i].vendorGUID.blockB); err != nil {
			return nil, err
		}

		if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.tableEntry[i].vendorGUID.blockC); err != nil {
			return nil, err
		}

		if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.tableEntry[i].vendorGUID.blockD); err != nil {
			return nil, err
		}

		if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.tableEntry[i].vendorGUID.blockE); err != nil {
			return nil, err
		}

		if err := binary.Read(eventReader, binary.LittleEndian, &handoffTablePointers.tableEntry[i].vendorTable); err != nil {
			return nil, err
		}
	}

	eventInfo := "Tables: "
	for _, table := range handoffTablePointers.tableEntry {
		guid := fmt.Sprintf("%x-%x-%x-%x-%x", table.vendorGUID.blockA, table.vendorGUID.blockB, table.vendorGUID.blockC, table.vendorGUID.blockD, table.vendorGUID.blockE)
		eventInfo += fmt.Sprintf("At address 0x%d with Guid %s", table.vendorTable, guid)
	}
	return &eventInfo, nil
}

func getPlatformFirmwareBlob(eventData []byte) (*string, error) {
	eventReader := bytes.NewReader(eventData)
	var platformFirmwareBlob EFIPlatformFirmwareBlob

	if err := binary.Read(eventReader, binary.LittleEndian, &platformFirmwareBlob.blobBase); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &platformFirmwareBlob.blobLength); err != nil {
		return nil, err
	}

	eventInfo := fmt.Sprintf("Blob address - 0x%d - with size - %db", platformFirmwareBlob.blobBase, platformFirmwareBlob.blobLength)
	return &eventInfo, nil
}

func getGPTEventString(eventData []byte) (*string, error) {
	eventReader := bytes.NewReader(eventData)
	var gptEvent EFIGptData

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.Signature); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.Revision); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.Size); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.CRC); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.HeaderStartLBA); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.HeaderCopyStartLBA); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.FirstUsableLBA); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.LastUsableLBA); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &gptEvent.uefiPartitionHeader.DiskGUID); err != nil {
		return nil, err
	}

	// Stop here we only want to know which device was used here.

	eventInfo := "Disk Guid - "
	eventInfo += gptEvent.uefiPartitionHeader.DiskGUID.String()
	return &eventInfo, nil
}

func getImageLoadEventString(eventData []byte) (*string, error) {
	eventReader := bytes.NewReader(eventData)
	var imageLoadEvent EFIImageLoadEvent

	if err := binary.Read(eventReader, binary.LittleEndian, &imageLoadEvent.imageLocationInMemory); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &imageLoadEvent.imageLengthInMemory); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &imageLoadEvent.imageLinkTimeAddress); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &imageLoadEvent.lengthOfDevicePath); err != nil {
		return nil, err
	}

	// Stop here we only want to know which device was used here.

	eventInfo := fmt.Sprintf("Image loaded at address 0x%d ", imageLoadEvent.imageLocationInMemory)
	eventInfo += fmt.Sprintf("with %db", imageLoadEvent.imageLengthInMemory)

	return &eventInfo, nil
}

func getVariableDataString(eventData []byte) (*string, error) {
	eventReader := bytes.NewReader(eventData)
	var variableData EFIVariableData

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableName.blockA); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableName.blockB); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableName.blockC); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableName.blockD); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableName.blockE); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.unicodeNameLength); err != nil {
		return nil, err
	}

	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableDataLength); err != nil {
		return nil, err
	}

	variableData.unicodeName = make([]uint16, variableData.unicodeNameLength)
	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.unicodeName); err != nil {
		return nil, err
	}

	variableData.variableData = make([]byte, variableData.variableDataLength)
	if err := binary.Read(eventReader, binary.LittleEndian, &variableData.variableData); err != nil {
		return nil, err
	}

	guid := fmt.Sprintf("Variable - %x-%x-%x-%x-%x - ", variableData.variableName.blockA, variableData.variableName.blockB, variableData.variableName.blockC, variableData.variableName.blockD, variableData.variableName.blockE)
	eventInfo := guid
	utf16String := utf16.Decode(variableData.unicodeName)
	eventInfo += string(utf16String)

	return &eventInfo, nil
}
