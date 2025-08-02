// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package blobs implements OpenBMC IPMI Blob Protocol commands.
//
// This file declares functions that implement the generic blob transfer
// interface detailed at https://github.com/openbmc/phosphor-ipmi-blobs
// with IPMI as a transport layer.
// See https://github.com/openbmc/google-ipmi-i2c for details on OEM
// commands.
package blobs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/ipmi"
)

// CRCOption is an option for sending/receiving CRCs.
type CRCOption string

// SessionID is a unique identifier for an open blob.
type SessionID uint16

// BlobStats contains statistics for a given blob.
type BlobStats struct {
	State       uint16
	Size        uint32
	MetadataLen uint8
	Metadata    []uint8
}

type blobStatsHdr struct {
	State       uint16
	Size        uint32
	MetadataLen uint8
}

type blobID []byte

// String returns the string representation of the blob ID, without the trailing NUL.
func (b blobID) String() string {
	return strings.TrimSuffix(string(b), "\000")
}

// blobIDFromString returns the IPMI blob-protocol representation of the given blob path (i.e., with one trailing NUL).
func blobIDFromString(blob string) blobID {
	terminated := strings.TrimSuffix(blob, "\000") + "\000"
	return blobIDFromBytes([]byte(terminated))
}

// blobIDFromBytes returns the IPMI blob-protocol representation of the given blob path from the given byte-array.
func blobIDFromBytes(bytes []byte) blobID {
	return blobID(bytes)
}

// BlobHandler provides an interface for the blob protocol. IT can be used
// to call all the blob transfer commands.
type BlobHandler struct {
	IPMI *ipmi.IPMI
}

const (
	_IPMI_GGL_NET_FN   ipmi.NetFn   = 46
	_IPMI_GGL_LUN                   = 0
	_IPMI_GGL_BLOB_CMD ipmi.Command = 128

	_OEN_LEN = 3
	_CRC_LEN = 2
)

// Blob transfer command codes.
const (
	_BMC_BLOB_CMD_CODE_GET_COUNT    = 0
	_BMC_BLOB_CMD_CODE_ENUMERATE    = 1
	_BMC_BLOB_CMD_CODE_OPEN         = 2
	_BMC_BLOB_CMD_CODE_READ         = 3
	_BMC_BLOB_CMD_CODE_WRITE        = 4
	_BMC_BLOB_CMD_CODE_COMMIT       = 5
	_BMC_BLOB_CMD_CODE_CLOSE        = 6
	_BMC_BLOB_CMD_CODE_DELETE       = 7
	_BMC_BLOB_CMD_CODE_STAT         = 8
	_BMC_BLOB_CMD_CODE_SESSION_STAT = 9
)

// Flags for blob open command.
const (
	BMC_BLOB_OPEN_FLAG_READ  = 1 << 0
	BMC_BLOB_OPEN_FLAG_WRITE = 1 << 1
	// Blob open: bit positions 2-7 are reserved for future protocol use.
	// Bit positions 8-15 are available for blob-specific definitions.
)

// Flags for blob state.
const (
	BMC_BLOB_STATE_OPEN_R       = 1 << 0
	BMC_BLOB_STATE_OPEN_W       = 1 << 1
	BMC_BLOB_STATE_COMMITTING   = 1 << 2
	BMC_BLOB_STATE_COMMITTED    = 1 << 3
	BMC_BLOB_STATE_COMMIT_ERROR = 1 << 4
	// Blob state: bit positions 5-7 are reserved for future protocol use.
	// Bit positions 8-16 are available for blob-specific definitions.
)

// CRC options
const (
	REQ_CRC     CRCOption = "REQ_CRC"
	RES_CRC     CRCOption = "RES_CRC"
	NO_CRC      CRCOption = "NO_CRC"
	REQ_RES_CRC CRCOption = "REQ_RES_CRC"
)

// OENMap maps OEM names to a 3 byte OEM number.
// OENs are typically serialized as the first 3 bytes of a request body.
var OENMap = map[string][3]uint8{
	"OpenBMC": {0xcf, 0xc2, 0x00},
}

// NewBlobHandler takes an IPMI struct, which provides a reference to the IPMI
// device driver, and returns a BlobHandler.
func NewBlobHandler(i *ipmi.IPMI) *BlobHandler {
	return &BlobHandler{IPMI: i}
}

// sendBlobCmd takes a command code, data given in little endian format, and
// an option for cyclic redundancy checks (CRC). It constructs the request
// and sends the command over IPMI. It receives the response, validates it,
// and then returns the response body.
func (h *BlobHandler) sendBlobCmd(code uint8, data []uint8, crcOpt CRCOption) ([]byte, error) {
	i := h.IPMI
	// Initialize a buffer with the correct OEN and code.
	oen, ok := OENMap["OpenBMC"]
	if !ok {
		return nil, fmt.Errorf("couldn't find OEN for OpenBMC")
	}

	buf := []uint8{oen[0], oen[1], oen[2], code}

	// If the request should have a CRC, derive a CRC based on the request body.
	if crcOpt == REQ_CRC || crcOpt == REQ_RES_CRC {
		crc := new(bytes.Buffer)
		if err := binary.Write(crc, binary.LittleEndian, genCRC(data)); err != nil {
			return nil, fmt.Errorf("failed to generate request CRC: %w", err)
		}
		buf = append(buf, crc.Bytes()...)
	}

	buf = append(buf, data...)

	// The request buffer should now be as follows:
	// - 3-byte OEN
	// - 1-byte subcommand code
	// - (optionally) 2-byte CRC over request body in little endian format
	// - request body in little endian format

	res, err := i.SendRecv(_IPMI_GGL_NET_FN, _IPMI_GGL_BLOB_CMD, buf)
	if err != nil {
		return nil, err
	}
	// Response always has a leading 0, so we ignore it.
	res = res[1:]

	// The response buffer is expected to be as follows:
	// - 3-byte OEN
	// - 2-byte CRC over response body in little endian format
	// - response body in little endian format
	// We verify that the OEN and CRC match the expected values.

	if len(res) < _OEN_LEN {
		return nil, fmt.Errorf("response too small: %d < size of OEN", len(res))
	}
	resOen, resBody := res[0:_OEN_LEN], res[_OEN_LEN:]

	// if oen[0] != resOen[0] || oen[1] != resOen[1] || oen[2] != resOen[2] {
	if !bytes.Equal(oen[0:3], resOen) {
		return nil, fmt.Errorf("response OEN incorrect: got %v, expected %v", resOen, oen)
	}

	// If the response should have a CRC, validate the CRC for the response body.
	if crcOpt == RES_CRC || crcOpt == REQ_RES_CRC {
		if err := verifyCRC(resBody); err != nil {
			return nil, fmt.Errorf("failed to verify response CRC: %w", err)
		}
		resBody = resBody[_CRC_LEN:]
	}

	return resBody, nil
}

// Sets a CCITT CRC based on the contents of the buffer.
// TODO(plaud): this is right now a copied implementation. Better to get and use
// a functional library (tried snksoft_crc but didn't work?)
func genCRC(data []uint8) uint16 {
	var kPoly uint16 = 0x1021
	var kLeftBit uint16 = 0x8000
	var crc uint16 = 0xFFFF
	kExtraRounds := 2

	for i := 0; i < len(data)+kExtraRounds; i++ {
		for j := range 8 {
			xorFlag := false
			if (crc & kLeftBit) != 0 {
				xorFlag = true
			}
			crc = crc << 1
			// If this isn't an extra round and the current byte's j'th bit from the
			// left is set, increment the CRC.
			if i < len(data) && (data[i]&(1<<(7-j))) != 0 {
				crc = crc + 1
			}
			if xorFlag {
				crc = crc ^ kPoly
			}
		}
	}

	return crc
}

// Verifies the CRC in the buffer, which must be the first two bytes. The CRC
// is validated against all data that follows it.
func verifyCRC(buf []uint8) error {
	if len(buf) < _CRC_LEN {
		return fmt.Errorf("response too small")
	}

	var respCrc uint16
	if err := binary.Read(bytes.NewReader(buf[0:_CRC_LEN]), binary.LittleEndian, &respCrc); err != nil {
		return fmt.Errorf("failed to read response CRC: %w", err)
	}

	expCrc := genCRC(buf[_CRC_LEN:])

	if expCrc != respCrc {
		return fmt.Errorf("CRC error: generated 0x%04X, got 0x%04X", expCrc, respCrc)
	}
	return nil
}

// Convert all args to little endian format and append to the given buffer.
func appendLittleEndian(buf []uint8, args ...any) ([]uint8, error) {
	for _, arg := range args {
		data := new(bytes.Buffer)
		err := binary.Write(data, binary.LittleEndian, arg)
		if err != nil {
			return nil, err
		}
		buf = append(buf, data.Bytes()...)
	}

	return buf, nil
}

// BlobGetCount returns the number of enumerable blobs available.
func (h *BlobHandler) BlobGetCount() (int, error) {
	data, err := h.sendBlobCmd(_BMC_BLOB_CMD_CODE_GET_COUNT, []uint8{}, RES_CRC)
	if err != nil {
		return 0, err
	}

	buf := bytes.NewReader(data)
	var blobCount int32

	if err := binary.Read(buf, binary.LittleEndian, &blobCount); err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	return (int)(blobCount), nil
}

// BlobEnumerate returns the blob identifier for the given index.
//
// Note that the index for a given blob ID is not expected to be stable long
// term. Callers are expected to call BlobGetCount, followed by N calls to
// BlobEnumerate, to collect all blob IDs.
func (h *BlobHandler) BlobEnumerate(index int) (string, error) {
	req, err := appendLittleEndian([]uint8{}, (int32)(index))
	if err != nil {
		return "", fmt.Errorf("failed to create data buffer: %w", err)
	}

	data, err := h.sendBlobCmd(_BMC_BLOB_CMD_CODE_ENUMERATE, req, REQ_RES_CRC)
	if err != nil {
		return "", err
	}

	return blobIDFromBytes(data).String(), nil
}

// BlobOpen opens a blob referred to by |id| with the given |flags|, and returns
// a unique session identifier.
//
// The BMC allocates a unique session identifier, and internally maps it
// to the blob identifier. The sessionId should be used by the rest of the
// session based commands to operate on the blob.
// NOTE: the new blob is not serialized and stored until BlobCommit is called.
func (h *BlobHandler) BlobOpen(id string, flags int16) (SessionID, error) {
	req, err := appendLittleEndian([]uint8{}, flags, blobIDFromString(id))
	if err != nil {
		return 0, fmt.Errorf("failed to create data buffer: %w", err)
	}

	data, err := h.sendBlobCmd(_BMC_BLOB_CMD_CODE_OPEN, req, REQ_RES_CRC)
	if err != nil {
		return 0, err
	}

	buf := bytes.NewReader(data)
	var sid SessionID

	if err := binary.Read(buf, binary.LittleEndian, &sid); err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	return sid, nil
}

// BlobRead reads and return the blob data.
//
// |sessionID| returned from BlobOpen gives us the open blob.
// The byte sequence starts at |offset|, and |size| bytes are read.
// If there are not enough bytes, return the bytes that are available.
func (h *BlobHandler) BlobRead(sid SessionID, offset, size uint32) ([]uint8, error) {
	req, err := appendLittleEndian([]uint8{}, sid, offset, size)
	if err != nil {
		return nil, fmt.Errorf("failed to create data buffer: %w", err)
	}

	data, err := h.sendBlobCmd(_BMC_BLOB_CMD_CODE_READ, req, REQ_RES_CRC)
	if err != nil {
		return nil, err
	}

	return ([]uint8)(data), nil
}

// BlobWrite writes bytes to the requested blob offset, and returns number of
// bytes written if success.
//
// |sessionID| returned from BlobOpen gives us the open blob.
// |data| is bounded by max size of an IPMI packet, which is platform-dependent.
// If not all of the bytes can be written, this operation will fail.
func (h *BlobHandler) BlobWrite(sid SessionID, offset uint32, data []uint8) error {
	req, err := appendLittleEndian([]uint8{}, sid, offset, data)
	if err != nil {
		return fmt.Errorf("failed to create data buffer: %w", err)
	}

	_, err = h.sendBlobCmd(_BMC_BLOB_CMD_CODE_WRITE, req, REQ_CRC)
	return err
}

// BlobCommit commits the blob.
//
// Each blob defines its own commit behavior. Optional blob-specific commit data
// can be provided with |data|.
// This function does not poll for BMC_BLOB_STATE_COMMITTED.
func (h *BlobHandler) BlobCommit(sid SessionID, data []uint8) error {
	req, err := appendLittleEndian([]uint8{}, sid, (uint8)(len(data)), data)
	if err != nil {
		return fmt.Errorf("failed to create data buffer: %w", err)
	}

	_, err = h.sendBlobCmd(_BMC_BLOB_CMD_CODE_COMMIT, req, REQ_CRC)
	return err
}

// BlobClose has the BMC mark the specified blob as closed.
//
// It must be called after commit-polling has finished, regardless of the result.
func (h *BlobHandler) BlobClose(sid SessionID) error {
	req, err := appendLittleEndian([]uint8{}, sid)
	if err != nil {
		return fmt.Errorf("failed to create data buffer: %w", err)
	}

	_, err = h.sendBlobCmd(_BMC_BLOB_CMD_CODE_CLOSE, req, REQ_CRC)
	return err
}

// BlobDelete deletes a blob if the operation is supported.
//
// This command will fail if there are open sessions for the blob.
func (h *BlobHandler) BlobDelete(id string) error {
	req, err := appendLittleEndian([]uint8{}, blobIDFromString(id))
	if err != nil {
		return fmt.Errorf("failed to create data buffer: %w", err)
	}

	_, err = h.sendBlobCmd(_BMC_BLOB_CMD_CODE_DELETE, req, REQ_CRC)
	return err
}

// BlobStat returns statistics about a blob.
//
// |size| is the size of blob in bytes. This may be zero if the blob does not
// support reading.
// |state| will be set with OPEN_R, OPEN_W, and/or COMMITTED as appropriate
// |metadata| is optional blob-specific bytes
func (h *BlobHandler) BlobStat(id string) (*BlobStats, error) {
	req, err := appendLittleEndian([]uint8{}, blobIDFromString(id))
	if err != nil {
		return nil, fmt.Errorf("failed to create data buffer: %w", err)
	}

	data, err := h.sendBlobCmd(_BMC_BLOB_CMD_CODE_STAT, req, REQ_RES_CRC)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	var statsHdr blobStatsHdr

	if err := binary.Read(buf, binary.LittleEndian, &statsHdr); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &BlobStats{
		State:       statsHdr.State,
		Size:        statsHdr.Size,
		MetadataLen: statsHdr.MetadataLen,
		Metadata:    buf.Bytes(),
	}, nil
}

// BlobSessionStat command returns the same data as BmcBlobStat.
//
// However, this command operates on sessions, rather than blob IDs. Not all
// blobs must support this command; this is only useful when session semantics
// are more useful than raw blob IDs.
func (h *BlobHandler) BlobSessionStat(sid SessionID) (*BlobStats, error) {
	req, err := appendLittleEndian([]uint8{}, sid)
	if err != nil {
		return nil, fmt.Errorf("failed to create data buffer: %w", err)
	}

	data, err := h.sendBlobCmd(_BMC_BLOB_CMD_CODE_SESSION_STAT, req, REQ_RES_CRC)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	var statsHdr blobStatsHdr

	if err := binary.Read(buf, binary.LittleEndian, &statsHdr); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &BlobStats{
		State:       statsHdr.State,
		Size:        statsHdr.Size,
		MetadataLen: statsHdr.MetadataLen,
		Metadata:    buf.Bytes(),
	}, nil
}
