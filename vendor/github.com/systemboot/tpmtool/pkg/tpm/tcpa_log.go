package tpm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf16"
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

// HashAlgoToSize is a map converter for hash to length
var HashAlgoToSize = map[IAlgHash]IAlgHashSize{
	TPMAlgSha:     TPMAlgShaSize,
	TPMAlgSha256:  TPMAlgSha256Size,
	TPMAlgSha384:  TPMAlgSha384Size,
	TPMAlgSha512:  TPMAlgSha512Size,
	TPMAlgSm3s256: TPMAlgSm3s256Size,
}

func stripControlSequences(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

func getTaggedEvent(eventData []byte) (*string, error) {
	var eventReader = bytes.NewReader(eventData)
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
	var eventReader = bytes.NewReader(eventData)
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

	eventInfo := fmt.Sprint("Tables: ")
	for _, table := range handoffTablePointers.tableEntry {
		guid := fmt.Sprintf("%x-%x-%x-%x-%x", table.vendorGUID.blockA, table.vendorGUID.blockB, table.vendorGUID.blockC, table.vendorGUID.blockD, table.vendorGUID.blockE)
		eventInfo += fmt.Sprintf("At address 0x%d with Guid %s", table.vendorTable, guid)
	}
	return &eventInfo, nil
}

func getPlatformFirmwareBlob(eventData []byte) (*string, error) {
	var eventReader = bytes.NewReader(eventData)
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
	var eventReader = bytes.NewReader(eventData)
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

	eventInfo := fmt.Sprint("Disk Guid - ")
	eventInfo += gptEvent.uefiPartitionHeader.DiskGUID.String()
	return &eventInfo, nil
}

func getImageLoadEventString(eventData []byte) (*string, error) {
	var eventReader = bytes.NewReader(eventData)
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
	var eventReader = bytes.NewReader(eventData)
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
	eventInfo += fmt.Sprintf("%s", string(utf16String))

	return &eventInfo, nil
}

func getEventDataString(eventType uint32, eventData []byte) (*string, error) {
	if eventType < uint32(EvEFIEventBase) {
		switch BIOSLogID(eventType) {
		case EvSeparator:
			eventInfo := fmt.Sprintf("%x", eventData)
			return &eventInfo, nil
		case EvAction:
			eventInfo := string(bytes.Trim(eventData, "\x00"))
			return &eventInfo, nil
		case EvOmitBootDeviceEvents:
			eventInfo := string("BOOT ATTEMPTS OMITTED")
			return &eventInfo, nil
		case EvPostCode:
			eventInfo := string(bytes.Trim(eventData, "\x00"))
			return &eventInfo, nil
		case EvEventTag:
			eventInfo, err := getTaggedEvent(eventData)
			if err != nil {
				return nil, err
			}
			return eventInfo, nil
		case EvSCRTMContents:
			eventInfo := string(bytes.Trim(eventData, "\x00"))
			return &eventInfo, nil
		case EvIPL:
			eventInfo := string(bytes.Trim(eventData, "\x00"))
			return &eventInfo, nil
		}
	} else {
		switch EFILogID(eventType) {
		case EvEFIHCRTMEvent:
			eventInfo := "HCRTM"
			return &eventInfo, nil
		case EvEFIAction:
			eventInfo := string(bytes.Trim(eventData, "\x00"))
			return &eventInfo, nil
		case EvEFIVariableDriverConfig, EvEFIVariableBoot, EvEFIVariableAuthority:
			eventInfo, err := getVariableDataString(eventData)
			if err != nil {
				return nil, err
			}
			return eventInfo, nil
		case EvEFIRuntimeServicesDriver, EvEFIBootServicesDriver, EvEFIBootServicesApplication:
			eventInfo, err := getImageLoadEventString(eventData)
			if err != nil {
				return nil, err
			}
			return eventInfo, nil
		case EvEFIGPTEvent:
			eventInfo, err := getGPTEventString(eventData)
			if err != nil {
				return nil, err
			}
			return eventInfo, nil
		case EvEFIPlatformFirmwareBlob:
			eventInfo, err := getPlatformFirmwareBlob(eventData)
			if err != nil {
				return nil, err
			}
			return eventInfo, nil
		case EvEFIHandoffTables:
			eventInfo, err := getHandoffTablePointers(eventData)
			if err != nil {
				return nil, err
			}
			return eventInfo, nil
		}
	}

	eventInfo := string(bytes.Trim(eventData, "\x00"))
	return &eventInfo, errors.New("Event type couldn't get parsed")
}

func readTPM2Log(firmware string) (*PCRLog, error) {
	var pcrLog PCRLog
	pcrLog.Firmware = firmware

	file, err := os.Open(DefaultTCPABinaryLog)
	if err != nil {
		return nil, err
	}

	var endianess binary.ByteOrder = binary.LittleEndian
	var pcrDigest PCRDigestInfo
	var pcrEvent TcgPcrEvent2
	for {
		if err := binary.Read(file, endianess, &pcrEvent.pcrIndex); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if err := binary.Read(file, endianess, &pcrEvent.eventType); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if BIOSLogID(pcrEvent.eventType) == EvNoAction {
			var efiSpecEvent TcgEfiSpecIDEvent
			if err := binary.Read(file, endianess, make([]byte, TPMAlgShaSize)); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &pcrEvent.eventSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			pcrEvent.event = make([]byte, pcrEvent.eventSize)

			if err := binary.Read(file, endianess, &efiSpecEvent.signature); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			identifier := string(bytes.Trim(efiSpecEvent.signature[:], "\x00"))
			if string(identifier) != TCGAgileEventFormatID {
				continue
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.platformClass); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.specVersionMinor); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.specVersionMajor); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.specErrata); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.uintnSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.numberOfAlgorithms); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			efiSpecEvent.digestSizes = make([]TcgEfiSpecIDEventAlgorithmSize, efiSpecEvent.numberOfAlgorithms)
			for i := uint32(0); i < efiSpecEvent.numberOfAlgorithms; i++ {
				if err := binary.Read(file, endianess, &efiSpecEvent.digestSizes[i].algorithID); err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}
				if err := binary.Read(file, endianess, &efiSpecEvent.digestSizes[i].digestSize); err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}
			}

			if err := binary.Read(file, endianess, &efiSpecEvent.vendorInfoSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			efiSpecEvent.vendorInfo = make([]byte, efiSpecEvent.vendorInfoSize)
			if err := binary.Read(file, endianess, &efiSpecEvent.vendorInfo); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			var in bytes.Buffer
			binary.Write(&in, endianess, efiSpecEvent)
			copy(pcrEvent.event, in.Bytes())

			if BIOSLogTypes[BIOSLogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = BIOSLogTypes[BIOSLogID(pcrEvent.eventType)]
			}
			if EFILogTypes[EFILogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = EFILogTypes[EFILogID(pcrEvent.eventType)]
			}

			pcrDigest.PcrIndex = int(pcrEvent.pcrIndex)
			pcrDigest.PcrEventData = string(pcrEvent.event)
			pcrLog.PcrList = append(pcrLog.PcrList, pcrDigest)
		} else {
			if err := binary.Read(file, endianess, &pcrEvent.digests.count); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			pcrEvent.digests.digests = make([]THA, pcrEvent.digests.count)
			for i := uint32(0); i < pcrEvent.digests.count; i++ {
				if err := binary.Read(file, endianess, &pcrEvent.digests.digests[i].hashAlg); err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}

				pcrEvent.digests.digests[i].digest.hash = make([]byte, HashAlgoToSize[pcrEvent.digests.digests[i].hashAlg])
				if err := binary.Read(file, endianess, &pcrEvent.digests.digests[i].digest.hash); err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}
			}

			if err := binary.Read(file, endianess, &pcrEvent.eventSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			// Placeholder
			pcrEvent.event = make([]byte, pcrEvent.eventSize)
			if err := binary.Read(file, endianess, &pcrEvent.event); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			pcrDigest.Digests = make([]PCRDigestValue, pcrEvent.digests.count)
			for i := uint32(0); i < pcrEvent.digests.count; i++ {
				pcrDigest.Digests[i].DigestAlg = pcrEvent.digests.digests[i].hashAlg
				pcrDigest.Digests[i].Digest = make([]byte, HashAlgoToSize[pcrEvent.digests.digests[i].hashAlg])
				copy(pcrDigest.Digests[i].Digest, pcrEvent.digests.digests[i].digest.hash)
			}

			if BIOSLogTypes[BIOSLogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = BIOSLogTypes[BIOSLogID(pcrEvent.eventType)]
			}
			if EFILogTypes[EFILogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = EFILogTypes[EFILogID(pcrEvent.eventType)]
			}

			pcrDigest.PcrIndex = int(pcrEvent.pcrIndex)
			eventDataString, _ := getEventDataString(pcrEvent.eventType, pcrEvent.event)
			if eventDataString != nil {
				pcrDigest.PcrEventData = *eventDataString
			}
			pcrLog.PcrList = append(pcrLog.PcrList, pcrDigest)
		}
	}
	file.Close()

	return &pcrLog, nil
}

func readTPM1Log(firmware string) (*PCRLog, error) {
	var pcrLog PCRLog
	pcrLog.Firmware = firmware

	file, err := os.Open(DefaultTCPABinaryLog)
	if err != nil {
		return nil, err
	}

	var endianess binary.ByteOrder = binary.LittleEndian
	var pcrDigest PCRDigestInfo
	var pcrEvent TcgPcrEvent
	for {
		if err := binary.Read(file, endianess, &pcrEvent.pcrIndex); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if err := binary.Read(file, endianess, &pcrEvent.eventType); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if err := binary.Read(file, endianess, &pcrEvent.digest); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if err := binary.Read(file, endianess, &pcrEvent.eventSize); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		pcrDigest.Digests = make([]PCRDigestValue, 1)
		pcrDigest.Digests[0].DigestAlg = TPMAlgSha
		if BIOSLogID(pcrEvent.eventType) == EvNoAction {
			var biosSpecEvent TcgBiosSpecIDEvent
			if err := binary.Read(file, endianess, make([]byte, TPMAlgShaSize)); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &pcrEvent.eventSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			pcrEvent.event = make([]byte, pcrEvent.eventSize)

			if err := binary.Read(file, endianess, &biosSpecEvent.signature); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			identifier := string(bytes.Trim(biosSpecEvent.signature[:], "\x00"))
			if string(identifier) != TCGOldEfiFormatID {
				continue
			}

			if err := binary.Read(file, endianess, &biosSpecEvent.platformClass); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &biosSpecEvent.specVersionMinor); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &biosSpecEvent.specVersionMajor); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &biosSpecEvent.specErrata); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &biosSpecEvent.uintnSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if err := binary.Read(file, endianess, &biosSpecEvent.vendorInfoSize); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			biosSpecEvent.vendorInfo = make([]byte, biosSpecEvent.vendorInfoSize)
			if err := binary.Read(file, endianess, &biosSpecEvent.vendorInfo); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			var in bytes.Buffer
			binary.Write(&in, endianess, biosSpecEvent)
			copy(pcrEvent.event, in.Bytes())

			if BIOSLogTypes[BIOSLogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = BIOSLogTypes[BIOSLogID(pcrEvent.eventType)]
			}
			if EFILogTypes[EFILogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = EFILogTypes[EFILogID(pcrEvent.eventType)]
			}

			pcrDigest.PcrIndex = int(pcrEvent.pcrIndex)
			pcrDigest.PcrEventData = string(pcrEvent.event)
			pcrLog.PcrList = append(pcrLog.PcrList, pcrDigest)
		} else {
			// Placeholder
			pcrEvent.event = make([]byte, pcrEvent.eventSize)
			if err := binary.Read(file, endianess, &pcrEvent.event); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			pcrDigest.Digests[0].Digest = make([]byte, TPMAlgShaSize)
			copy(pcrDigest.Digests[0].Digest, pcrEvent.digest[:])

			if BIOSLogTypes[BIOSLogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = BIOSLogTypes[BIOSLogID(pcrEvent.eventType)]
			}
			if EFILogTypes[EFILogID(pcrEvent.eventType)] != "" {
				pcrDigest.PcrEventName = EFILogTypes[EFILogID(pcrEvent.eventType)]
			}

			eventDataString, _ := getEventDataString(pcrEvent.eventType, pcrEvent.event)
			if eventDataString != nil {
				pcrDigest.PcrEventData = *eventDataString
			}

			pcrDigest.PcrIndex = int(pcrEvent.pcrIndex)
			pcrLog.PcrList = append(pcrLog.PcrList, pcrDigest)
		}
	}
	file.Close()

	return &pcrLog, nil
}

// ParseLog is a ,..
func ParseLog(firmware string, tpmSpec string) (*PCRLog, error) {
	var pcrLog *PCRLog
	var err error

	switch tpmSpec {
	case TPM12:
		pcrLog, err = readTPM1Log(firmware)
		if err != nil {
			return nil, err
		}
	case TPM20:
		pcrLog, err = readTPM2Log(firmware)
		if err != nil {
			// Kernel eventlog workaround does not export agile measurement log..
			pcrLog, err = readTPM1Log(firmware)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, errors.New("No valid TPM specification found")
	}

	return pcrLog, nil
}

// DumpLog dumps the evenlog on stdio
func DumpLog(tcpaLog *PCRLog) error {
	for _, pcr := range tcpaLog.PcrList {
		fmt.Printf("PCR: %d\n", pcr.PcrIndex)
		fmt.Printf("Event Name: %s\n", pcr.PcrEventName)
		fmt.Printf("Event Data: %s\n", stripControlSequences(pcr.PcrEventData))

		for _, digest := range pcr.Digests {
			var algoName string
			switch digest.DigestAlg {
			case TPMAlgSha:
				algoName = "SHA1"
			case TPMAlgSha256:
				algoName = "SHA256"
			case TPMAlgSha384:
				algoName = "SHA384"
			case TPMAlgSha512:
				algoName = "SHA512"
			case TPMAlgSm3s256:
				algoName = "SM3"
			}

			fmt.Printf("%s Digest: %x\n", algoName, digest.Digest)
		}

		fmt.Println()
	}

	return nil
}
