// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package txtlog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

func parseTcgBiosSpecIDEvent(handle io.Reader) (*TcgBiosSpecIDEvent, error) {
	var endianness binary.ByteOrder = binary.LittleEndian
	var biosSpecEvent TcgBiosSpecIDEvent

	if err := binary.Read(handle, endianness, &biosSpecEvent.signature); err != nil {
		return nil, err
	}

	identifier := string(bytes.Trim(biosSpecEvent.signature[:], "\x00"))
	if string(identifier) != TCGOldEfiFormatID {
		return nil, nil
	}

	if err := binary.Read(handle, endianness, &biosSpecEvent.platformClass); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &biosSpecEvent.specVersionMinor); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &biosSpecEvent.specVersionMajor); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &biosSpecEvent.specErrata); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &biosSpecEvent.uintnSize); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &biosSpecEvent.vendorInfoSize); err != nil {
		return nil, err
	}

	biosSpecEvent.vendorInfo = make([]byte, biosSpecEvent.vendorInfoSize)
	if err := binary.Read(handle, endianness, &biosSpecEvent.vendorInfo); err != nil {
		return nil, err
	}

	return &biosSpecEvent, nil
}

func parseEfiSpecEvent(handle io.Reader) (*TcgEfiSpecIDEvent, error) {
	var endianness binary.ByteOrder = binary.LittleEndian
	var efiSpecEvent TcgEfiSpecIDEvent

	if err := binary.Read(handle, endianness, &efiSpecEvent.signature); err != nil {
		return nil, err
	}

	identifier := string(bytes.Trim(efiSpecEvent.signature[:], "\x00"))
	if string(identifier) != TCGAgileEventFormatID {
		return nil, nil
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.platformClass); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.specVersionMinor); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.specVersionMajor); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.specErrata); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.uintnSize); err != nil {
		return nil, err
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.numberOfAlgorithms); err != nil {
		return nil, err
	}

	efiSpecEvent.digestSizes = make([]TcgEfiSpecIDEventAlgorithmSize, efiSpecEvent.numberOfAlgorithms)
	for i := uint32(0); i < efiSpecEvent.numberOfAlgorithms; i++ {
		if err := binary.Read(handle, endianness, &efiSpecEvent.digestSizes[i].algorithID); err != nil {
			return nil, err
		}
		if err := binary.Read(handle, endianness, &efiSpecEvent.digestSizes[i].digestSize); err != nil {
			return nil, err
		}
	}

	if err := binary.Read(handle, endianness, &efiSpecEvent.vendorInfoSize); err != nil {
		return nil, err
	}

	efiSpecEvent.vendorInfo = make([]byte, efiSpecEvent.vendorInfoSize)
	if err := binary.Read(handle, endianness, &efiSpecEvent.vendorInfo); err != nil {
		return nil, err
	}

	return &efiSpecEvent, nil
}

// TcgPcrEvent parser and PCREvent interface implementation
func parseTcgPcrEvent(handle io.Reader) (*TcgPcrEvent, error) {
	var endianness binary.ByteOrder = binary.LittleEndian
	var pcrEvent TcgPcrEvent

	if err := binary.Read(handle, endianness, &pcrEvent.pcrIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, endianness, &pcrEvent.eventType); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, endianness, &pcrEvent.digest); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, endianness, &pcrEvent.eventSize); err != nil {
		return nil, err
	}

	pcrEvent.event = make([]byte, pcrEvent.eventSize)
	if err := binary.Read(handle, endianness, &pcrEvent.event); err != nil {
		return nil, err
	}

	return &pcrEvent, nil
}

func (e *TcgPcrEvent) PcrIndex() int {
	return int(e.pcrIndex)
}

func (e *TcgPcrEvent) PcrEventType() uint32 {
	return e.eventType
}

func (e *TcgPcrEvent) PcrEventName() string {
	if BIOSLogTypes[BIOSLogID(e.eventType)] != "" {
		return BIOSLogTypes[BIOSLogID(e.eventType)]
	}
	if EFILogTypes[EFILogID(e.eventType)] != "" {
		return EFILogTypes[EFILogID(e.eventType)]
	}
	if TxtLogTypes[TxtLogID(e.eventType)] != "" {
		return TxtLogTypes[TxtLogID(e.eventType)]
	}

	return ""
}

func (e *TcgPcrEvent) PcrEventData() string {
	if BIOSLogID(e.eventType) == EvNoAction {
		return string(e.event)
	}

	eventDataString, _ := getEventDataString(e.eventType, e.event)
	if eventDataString != nil {
		return *eventDataString
	}

	return ""
}

func (e *TcgPcrEvent) Digests() *[]PCRDigestValue {
	d := make([]PCRDigestValue, 1)
	d[0].DigestAlg = TPMAlgSha
	d[0].Digest = make([]byte, TPMAlgShaSize)
	copy(d[0].Digest, e.digest[:])

	return &d
}

func (e *TcgPcrEvent) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "PCR: %d\n", e.PcrIndex())
	fmt.Fprintf(&b, "Event Name: %s\n", e.PcrEventName())
	fmt.Fprintf(&b, "Event Data: %s\n", stripControlSequences(e.PcrEventData()))
	fmt.Fprintf(&b, "SHA1 Digest: %x", e.digest)

	return b.String()
}

// TcgPcrEvent2 parser and PCREvent interface implementation
func parseTcgPcrEvent2(handle io.Reader) (*TcgPcrEvent2, error) {
	var endianness binary.ByteOrder = binary.LittleEndian
	var pcrEvent TcgPcrEvent2

	if err := binary.Read(handle, endianness, &pcrEvent.pcrIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, endianness, &pcrEvent.eventType); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, endianness, &pcrEvent.digests.count); err != nil {
		return nil, err
	}

	pcrEvent.digests.digests = make([]THA, pcrEvent.digests.count)
	for i := uint32(0); i < pcrEvent.digests.count; i++ {
		if err := binary.Read(handle, endianness, &pcrEvent.digests.digests[i].hashAlg); err != nil {
			return nil, err
		}

		pcrEvent.digests.digests[i].digest.hash = make([]byte, HashAlgoToSize[pcrEvent.digests.digests[i].hashAlg])
		if err := binary.Read(handle, endianness, &pcrEvent.digests.digests[i].digest.hash); err != nil {
			return nil, err
		}
	}

	if err := binary.Read(handle, endianness, &pcrEvent.eventSize); err != nil {
		return nil, err
	}

	pcrEvent.event = make([]byte, pcrEvent.eventSize)
	if err := binary.Read(handle, endianness, &pcrEvent.event); err != nil {
		return nil, err
	}

	return &pcrEvent, nil
}

func (e *TcgPcrEvent2) PcrIndex() int {
	return int(e.pcrIndex)
}

func (e *TcgPcrEvent2) PcrEventType() uint32 {
	return e.eventType
}

func (e *TcgPcrEvent2) PcrEventName() string {
	if BIOSLogTypes[BIOSLogID(e.eventType)] != "" {
		return BIOSLogTypes[BIOSLogID(e.eventType)]
	}
	if EFILogTypes[EFILogID(e.eventType)] != "" {
		return EFILogTypes[EFILogID(e.eventType)]
	}
	if TxtLogTypes[TxtLogID(e.eventType)] != "" {
		return TxtLogTypes[TxtLogID(e.eventType)]
	}

	return ""
}

func (e *TcgPcrEvent2) PcrEventData() string {
	if BIOSLogID(e.eventType) == EvNoAction {
		return string(e.event)
	}
	eventDataString, _ := getEventDataString(e.eventType, e.event)
	if eventDataString != nil {
		return *eventDataString
	}

	return ""
}

func (e *TcgPcrEvent2) Digests() *[]PCRDigestValue {
	d := make([]PCRDigestValue, e.digests.count)
	for i := uint32(0); i < e.digests.count; i++ {
		d[i].DigestAlg = e.digests.digests[i].hashAlg
		d[i].Digest = make([]byte, HashAlgoToSize[e.digests.digests[i].hashAlg])
		copy(d[i].Digest, e.digests.digests[i].digest.hash)
	}
	return &d
}

func (e *TcgPcrEvent2) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "PCR: %d\n", e.PcrIndex())
	fmt.Fprintf(&b, "Event Name: %s\n", e.PcrEventName())
	fmt.Fprintf(&b, "Event Data: %s\n", stripControlSequences(e.PcrEventData()))
	for i := uint32(0); i < e.digests.count; i++ {
		d := &e.digests.digests[i]
		switch d.hashAlg {
		case TPMAlgSha:
			b.WriteString("SHA1 Digest: ")
		case TPMAlgSha256:
			b.WriteString("SHA256 Digest: ")
		case TPMAlgSha384:
			b.WriteString("SHA384 Digest: ")
		case TPMAlgSha512:
			b.WriteString("SHA512 Digest: ")
		case TPMAlgSm3s256:
			b.WriteString("SM3 Digest: ")
		}

		fmt.Fprintf(&b, "%x\n", d.digest.hash)
	}

	return b.String()
}

func readTxtEventLogContainer(handle io.Reader) (*TxtEventLogContainer, error) {
	var container TxtEventLogContainer

	if err := binary.Read(handle, binary.LittleEndian, &container.Signature); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.Reserved); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.ContainerVerMajor); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.ContainerVerMinor); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.PcrEventVerMajor); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.PcrEventVerMinor); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.Size); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.PcrEventsOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(handle, binary.LittleEndian, &container.NextEventOffset); err != nil {
		return nil, err
	}

	return &container, nil
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
	return &eventInfo, errors.New("event type couldn't get parsed")
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
