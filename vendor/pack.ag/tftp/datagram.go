// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type opcode uint16

func (o opcode) String() string {
	name, ok := opcodeStrings[o]
	if ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN_OPCODE_%v", uint16(o))
}

// ErrorCode is a TFTP error code as defined in RFC 1350
type ErrorCode uint16

func (e ErrorCode) String() string {
	name, ok := errorStrings[e]
	if ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN_ERROR_%v", uint16(e))
}

const (
	opCodeRRQ   opcode = 0x1 // Read Request
	opCodeWRQ   opcode = 0x2 // Write Request
	opCodeDATA  opcode = 0x3 // Data
	opCodeACK   opcode = 0x4 // Acknowledgement
	opCodeERROR opcode = 0x5 // Error
	opCodeOACK  opcode = 0x6 // Option Acknowledgement

	// ErrCodeNotDefined - Not defined, see error message (if any).
	ErrCodeNotDefined ErrorCode = 0x0
	// ErrCodeFileNotFound - File not found.
	ErrCodeFileNotFound ErrorCode = 0x1
	// ErrCodeAccessViolation - Access violation.
	ErrCodeAccessViolation ErrorCode = 0x2
	// ErrCodeDiskFull - Disk full or allocation exceeded.
	ErrCodeDiskFull ErrorCode = 0x3
	// ErrCodeIllegalOperation - Illegal TFTP operation.
	ErrCodeIllegalOperation ErrorCode = 0x4
	// ErrCodeUnknownTransferID - Unknown transfer ID.
	ErrCodeUnknownTransferID ErrorCode = 0x5
	// ErrCodeFileAlreadyExists - File already exists.
	ErrCodeFileAlreadyExists ErrorCode = 0x6
	// ErrCodeNoSuchUser - No such user.
	ErrCodeNoSuchUser ErrorCode = 0x7

	// ModeNetASCII is the string for netascii transfer mode
	ModeNetASCII TransferMode = "netascii"
	// ModeOctet is the string for octet/binary transfer mode
	ModeOctet TransferMode = "octet"
	modeMail  TransferMode = "mail"

	optBlocksize    = "blksize"
	optTimeout      = "timeout"
	optTransferSize = "tsize"
	optWindowSize   = "windowsize"
)

// TransferMode is a TFTP transer mode
type TransferMode string

var (
	errorStrings = map[ErrorCode]string{
		ErrCodeNotDefined:        "NOT_DEFINED",
		ErrCodeFileNotFound:      "FILE_NOT_FOUND",
		ErrCodeAccessViolation:   "ACCESS_VIOLATION",
		ErrCodeDiskFull:          "DISK_FULL",
		ErrCodeIllegalOperation:  "ILLEGAL_OPERATION",
		ErrCodeUnknownTransferID: "UNKNOWN_TRANSFER_ID",
		ErrCodeFileAlreadyExists: "FILE_ALREADY_EXISTS",
		ErrCodeNoSuchUser:        "NO_SUCH_USER",
	}
	opcodeStrings = map[opcode]string{
		opCodeRRQ:   "READ_REQUEST",
		opCodeWRQ:   "WRITE_REQUEST",
		opCodeDATA:  "DATA",
		opCodeACK:   "ACK",
		opCodeERROR: "ERROR",
		opCodeOACK:  "OPTION_ACK",
	}
)

type datagram struct {
	buf    []byte
	offset int
}

func (d datagram) String() string {
	if err := d.validate(); err != nil {
		return fmt.Sprintf("INVALID_DATAGRAM[Error: %q]", err.Error())
	}

	switch o := d.opcode(); o {
	case opCodeRRQ, opCodeWRQ:
		return fmt.Sprintf("%s[Filename: %q; Mode: %q; Options: %s]", o, d.filename(), d.mode(), d.options())
	case opCodeDATA:
		return fmt.Sprintf("%s[Block: %d; Data Length: %d]", o, d.block(), len(d.data()))
	case opCodeOACK:
		return fmt.Sprintf("%s[Options: %s]", o, d.options())
	case opCodeACK:
		return fmt.Sprintf("%s[Block: %d]", o, d.block())
	case opCodeERROR:
		return fmt.Sprintf("%s[Code: %s; Message: %q]", o, d.errorCode(), d.errMsg())
	default:
		return o.String()
	}
}

// Sets the buffer from raw bytes
func (d *datagram) setBytes(b []byte) {
	d.buf = b
	d.offset = len(b)
}

// Returns the allocated bytes
func (d *datagram) bytes() []byte {
	return d.buf[:d.offset]
}

// Resets the byte buffer.
// If requested size is larger than allocated the buffer is reallocated.
func (d *datagram) reset(size int) {
	if len(d.buf) < size {
		d.buf = make([]byte, size)
	}
	d.offset = 0
}

// DATAGRAM CONSTRUCTORS
func (d *datagram) writeAck(block uint16) {
	d.reset(2 + 2)

	d.writeUint16(uint16(opCodeACK))
	d.writeUint16(block)
}

func (d *datagram) writeData(block uint16, data []byte) {
	d.reset(2 + 2 + len(data))

	d.writeUint16(uint16(opCodeDATA))
	d.writeUint16(block)
	d.writeBytes(data)
}

func (d *datagram) writeError(code ErrorCode, msg string) {
	d.reset(2 + 2 + len(msg) + 1)

	d.writeUint16(uint16(opCodeERROR))
	d.writeUint16(uint16(code))
	d.writeString(msg)
	d.writeNull()
}

func (d *datagram) writeReadReq(filename string, mode TransferMode, options map[string]string) {
	d.writeReq(opCodeRRQ, filename, mode, options)
}

func (d *datagram) writeWriteReq(filename string, mode TransferMode, options map[string]string) {
	d.writeReq(opCodeWRQ, filename, mode, options)
}

func (d *datagram) writeOptionAck(options map[string]string) {
	optLen := 0
	for opt, val := range options {
		optLen += len(opt) + 1 + len(val) + 1
	}
	d.reset(2 + optLen)

	d.writeUint16(uint16(opCodeOACK))

	for opt, val := range options {
		d.writeOption(opt, val)
	}
}

// Combines duplicate logic from RRQ and WRQ
func (d *datagram) writeReq(o opcode, filename string, mode TransferMode, options map[string]string) {
	// This is ugly, could just set buf to 512
	// or use a bytes buffer. Intend to switch to bytes buffer
	// after implementing all RFCs so that perf can be compared
	// with a reasonable block and window size
	optLen := 0
	for opt, val := range options {
		optLen += len(opt) + 1 + len(val) + 1
	}
	d.reset(2 + len(filename) + 1 + len(mode) + 1 + optLen)

	d.writeUint16(uint16(o))
	d.writeString(filename)
	d.writeNull()
	d.writeString(string(mode))
	d.writeNull()

	for opt, val := range options {
		d.writeOption(opt, val)
	}
}

// FIELD ACCESSORS

// Block # from DATA and ACK datagrams
func (d *datagram) block() uint16 {
	return binary.BigEndian.Uint16(d.buf[2:4])
}

// Data from DATA datagram
func (d *datagram) data() []byte {
	return d.buf[4:d.offset]
}

// ErrorCode from ERROR datagram
func (d *datagram) errorCode() ErrorCode {
	return ErrorCode(binary.BigEndian.Uint16(d.buf[2:4]))
}

// ErrMsg from ERROR datagram
func (d *datagram) errMsg() string {
	end := d.offset - 1
	return string(d.buf[4:end])
}

// Filename from RRQ and WRQ datagrams
func (d *datagram) filename() string {
	offset := bytes.IndexByte(d.buf[2:], 0x0) + 2
	return string(d.buf[2:offset])
}

// Mode from RRQ and WRQ datagrams
func (d *datagram) mode() TransferMode {
	fields := bytes.Split(d.buf[2:], []byte{0x0})
	return TransferMode(fields[1])
}

// Opcode from all datagrams
func (d *datagram) opcode() opcode {
	return opcode(binary.BigEndian.Uint16(d.buf[:2]))
}

type options map[string]string

func (o options) String() string {
	opts := make([]string, 0, len(o))
	for k, v := range o {
		opts = append(opts, fmt.Sprintf("%q: %q", k, v))
	}

	return "{" + strings.Join(opts, "; ") + "}"
}

func (d *datagram) options() options {
	options := make(options)

	optSlice := bytes.Split(d.buf[2:d.offset-1], []byte{0x0}) // d.buf[2:d.offset-1] = file -> just before final NULL
	if op := d.opcode(); op == opCodeRRQ || op == opCodeWRQ {
		optSlice = optSlice[2:] // Remove filename, mode
	}

	for i := 0; i < len(optSlice); i += 2 {
		options[string(optSlice[i])] = string(optSlice[i+1])
	}
	return options
}

// BUFFER WRITING FUNCTIONS
func (d *datagram) writeBytes(b []byte) {
	copy(d.buf[d.offset:], b)
	d.offset += len(b)
}

func (d *datagram) writeNull() {
	d.buf[d.offset] = 0x0
	d.offset++
}

func (d *datagram) writeString(str string) {
	d.writeBytes([]byte(str))
}

func (d *datagram) writeUint16(i uint16) {
	binary.BigEndian.PutUint16(d.buf[d.offset:], i)
	d.offset += 2
}

func (d *datagram) writeOption(o string, v string) {
	d.writeString(o)
	d.writeNull()
	d.writeString(v)
	d.writeNull()
}

// VALIDATION

func (d *datagram) validate() error {
	switch {
	case d.offset < 2:
		return errors.New("Datagram has no opcode")
	case d.opcode() > 6:
		return errors.New("Invalid opcode")
	}

	switch d.opcode() {
	case opCodeRRQ, opCodeWRQ:
		switch {
		case len(d.filename()) < 1:
			return errors.New("No filename provided")
		case d.buf[d.offset-1] != 0x0: // End with NULL
			return fmt.Errorf("Corrupt %v datagram", d.opcode())
		case bytes.Count(d.buf[2:d.offset], []byte{0x0})%2 != 0: // Number of NULL chars is not even
			return fmt.Errorf("Corrupt %v datagram", d.opcode())
		default:
			switch d.mode() {
			case ModeNetASCII, ModeOctet:
				break
			case modeMail:
				return errors.New("MAIL transfer mode is unsupported")
			default:
				return errors.New("Invalid transfer mode")
			}
		}
	case opCodeACK, opCodeDATA:
		if d.offset < 4 {
			return errors.New("Corrupt block number")
		}
	case opCodeERROR:
		switch {
		case d.offset < 5:
			return errors.New("Corrupt ERROR datagram")
		case d.buf[d.offset-1] != 0x0:
			return errors.New("Corrupt ERROR datagram")
		case bytes.Count(d.buf[4:d.offset], []byte{0x0}) > 1:
			return errors.New("Corrupt ERROR datagram")
		}
	case opCodeOACK:
		switch {
		case d.buf[d.offset-1] != 0x0:
			return errors.New("Corrupt OACK datagram")
		case bytes.Count(d.buf[2:d.offset], []byte{0x0})%2 != 0: // Number of NULL chars is not even
			return errors.New("Corrupt OACK datagram")
		}
	}

	return nil
}
