// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"bytes"
	"reflect"
	"testing"
)

func TestOpcode_String(t *testing.T) {
	cases := []struct {
		code opcode

		expected string
	}{
		{
			code:     opCodeRRQ,
			expected: "READ_REQUEST",
		},
		{
			code:     opCodeWRQ,
			expected: "WRITE_REQUEST",
		},
		{
			code:     opCodeDATA,
			expected: "DATA",
		},
		{
			code:     opCodeACK,
			expected: "ACK",
		},
		{
			code:     opCodeERROR,
			expected: "ERROR",
		},
		{
			code:     opCodeOACK,
			expected: "OPTION_ACK",
		},
		{
			code:     13,
			expected: "UNKNOWN_OPCODE_13",
		},
	}

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			if c.code.String() != c.expected {
				t.Errorf("Expected opcode(%d).String() to be %q, but it was %q", c.code, c.expected, c.code.String())
			}
		})
	}
}

func TestErrorCode_String(t *testing.T) {
	cases := []struct {
		code ErrorCode

		expected string
	}{
		{
			code:     ErrCodeNotDefined,
			expected: "NOT_DEFINED",
		},
		{
			code:     ErrCodeFileNotFound,
			expected: "FILE_NOT_FOUND",
		},
		{
			code:     ErrCodeAccessViolation,
			expected: "ACCESS_VIOLATION",
		},
		{
			code:     ErrCodeDiskFull,
			expected: "DISK_FULL",
		},
		{
			code:     ErrCodeIllegalOperation,
			expected: "ILLEGAL_OPERATION",
		},
		{
			code:     ErrCodeUnknownTransferID,
			expected: "UNKNOWN_TRANSFER_ID",
		},
		{
			code:     ErrCodeFileAlreadyExists,
			expected: "FILE_ALREADY_EXISTS",
		},
		{
			code:     ErrCodeNoSuchUser,
			expected: "NO_SUCH_USER",
		},
		{
			code:     13,
			expected: "UNKNOWN_ERROR_13",
		},
	}

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			if c.code.String() != c.expected {
				t.Errorf("Expected errCode(%d).String() to be %q, but it was %q", c.code, c.expected, c.code.String())
			}
		})
	}
}

func TestDatagram_String(t *testing.T) {
	cases := []struct {
		name string
		dg   datagram

		expected string
	}{
		{
			name: "RRQ",
			dg: func() datagram {
				d := datagram{}
				d.writeReadReq("readFile", ModeNetASCII, options{"first": "option"})
				return d
			}(),
			expected: `READ_REQUEST[Filename: "readFile"; Mode: "netascii"; Options: {"first": "option"}]`,
		},
		{
			name: "WRQ",
			dg: func() datagram {
				d := datagram{}
				d.writeWriteReq("readFile", ModeNetASCII, options{})
				return d
			}(),
			expected: `WRITE_REQUEST[Filename: "readFile"; Mode: "netascii"; Options: {}]`,
		},
		{
			name: "DATA",
			dg: func() datagram {
				d := datagram{}
				d.writeData(678, []byte("the data"))
				return d
			}(),
			expected: `DATA[Block: 678; Data Length: 8]`,
		},
		{
			name: "OACK",
			dg: func() datagram {
				d := datagram{}
				d.writeOptionAck(options{"first": "option"})
				return d
			}(),
			expected: `OPTION_ACK[Options: {"first": "option"}]`,
		},
		{
			name: "ACK",
			dg: func() datagram {
				d := datagram{}
				d.writeAck(65000)
				return d
			}(),
			expected: `ACK[Block: 65000]`,
		},
		{
			name: "ERROR",
			dg: func() datagram {
				d := datagram{}
				d.writeError(ErrCodeDiskFull, "my error")
				return d
			}(),
			expected: `ERROR[Code: DISK_FULL; Message: "my error"]`,
		},
		{
			name:     "Bad Datagram",
			dg:       datagram{},
			expected: `INVALID_DATAGRAM[Error: "Datagram has no opcode"]`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.dg.String() != c.expected {
				t.Errorf("expected to be %q, but it was %q", c.expected, c.dg.String())
			}
		})
	}
}

func TestDatagram(t *testing.T) {
	cases := []struct {
		name string
		dg   datagram

		valid      bool
		len        int
		data       []byte
		offset     int
		code       opcode
		block      uint16
		filename   *string
		mode       *TransferMode
		opts       options
		errCode    *ErrorCode
		errMessage *string
	}{
		{
			name: "ack",
			dg: func() datagram {
				dg := datagram{}
				dg.writeAck(3)
				return dg
			}(),

			valid:  true,
			len:    4,
			data:   []byte{},
			offset: 4,
			code:   opCodeACK,
			block:  3,
		},
		{
			name: "data",
			dg: func() datagram {
				dg := datagram{}
				dg.writeData(314, []byte("this is the data"))
				return dg
			}(),

			valid:  true,
			len:    20,
			offset: 20,
			code:   opCodeDATA,
		},
		{
			name: "RRQ",
			dg: func() datagram {
				dg := datagram{}
				dg.writeReadReq("the file", ModeNetASCII, options{})
				return dg
			}(),

			valid:    true,
			len:      20,
			offset:   20,
			code:     opCodeRRQ,
			filename: ptrString("the file"),
			mode:     ptrMode(ModeNetASCII),
			opts:     options{},
		},
		{
			name: "WRQ",
			dg: func() datagram {
				dg := datagram{}
				dg.writeWriteReq("a file", ModeOctet, options{})
				return dg
			}(),

			valid:    true,
			len:      15,
			offset:   15,
			code:     opCodeWRQ,
			filename: ptrString("a file"),
			mode:     ptrMode(ModeOctet),
			opts:     options{},
		},
		{
			name: "OACK, no options",
			dg: func() datagram {
				dg := datagram{}
				dg.writeOptionAck(options{})
				return dg
			}(),

			valid: false,
		},
		{
			name: "OACK",
			dg: func() datagram {
				dg := datagram{}
				dg.writeOptionAck(options{optBlocksize: "345"})
				return dg
			}(),

			valid:  true,
			len:    14,
			offset: 14,
			code:   opCodeOACK,
			opts:   options{optBlocksize: "345"},
		},
		{
			name: "error",
			dg: func() datagram {
				dg := datagram{}
				dg.writeError(ErrCodeDiskFull, "the message")
				return dg
			}(),

			valid:      true,
			len:        16,
			offset:     16,
			code:       opCodeERROR,
			errCode:    ptrErrCode(ErrCodeDiskFull),
			errMessage: ptrString("the message"),
		},
		{
			name: "no opcode",
			dg: func() datagram {
				dg := datagram{}
				return dg
			}(),

			valid: false,
		},
		{
			name: "invalid opcode",
			dg: func() datagram {
				dg := datagram{}
				dg.reset(2)
				dg.writeUint16(13)
				return dg
			}(),

			valid: false,
		},
		{
			name: "empty filename",
			dg: func() datagram {
				dg := datagram{}
				dg.writeReadReq("", ModeOctet, options{})
				dg.buf[dg.offset-1] = 'x'
				return dg
			}(),

			valid: false,
		},
		{
			name: "request doesn't end with null",
			dg: func() datagram {
				dg := datagram{}
				dg.writeReadReq("file", ModeOctet, options{})
				dg.buf[dg.offset-1] = 'x'
				return dg
			}(),

			valid: false,
		},
		{
			name: "request has odd number of null",
			dg: func() datagram {
				dg := datagram{}
				dg.writeReadReq("file\x00name", ModeOctet, options{})
				return dg
			}(),

			valid: false,
		},
		{
			name: "mail",
			dg: func() datagram {
				dg := datagram{}
				dg.writeReadReq("file", modeMail, options{})
				return dg
			}(),

			valid: false,
		},
		{
			name: "invalid mode",
			dg: func() datagram {
				dg := datagram{}
				dg.writeReadReq("file", "fast", options{})
				return dg
			}(),

			valid: false,
		},
		{
			name: "corrupt block #",
			dg: func() datagram {
				dg := datagram{}
				dg.writeData(133, []byte("data"))
				dg.offset = 3
				return dg
			}(),

			valid: false,
		},
		{
			name: "corrupt error",
			dg: func() datagram {
				dg := datagram{}
				dg.reset(4)
				dg.writeUint16(uint16(opCodeERROR))
				dg.writeUint16(uint16(ErrCodeAccessViolation))
				return dg
			}(),

			valid: false,
		},
		{
			name: "error doesn't end with null",
			dg: func() datagram {
				dg := datagram{}
				dg.reset(8)
				dg.writeUint16(uint16(opCodeERROR))
				dg.writeUint16(uint16(ErrCodeAccessViolation))
				dg.writeString("data")
				return dg
			}(),

			valid: false,
		},
		{
			name: "error has more than one null",
			dg: func() datagram {
				dg := datagram{}
				dg.reset(8)
				dg.writeError(ErrCodeDiskFull, "the\x00data")
				return dg
			}(),

			valid: false,
		},
		{
			name: "corrupt options",
			dg: func() datagram {
				dg := datagram{}
				dg.reset(10)
				dg.writeUint16(uint16(opCodeOACK))
				dg.writeString(optBlocksize)
				dg.writeNull()
				return dg
			}(),

			valid: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Valid
			if err := c.dg.validate(); (err == nil) != c.valid {
				t.Errorf("expected %s to be valid %t, but it wasn't: %s", c.dg, c.valid, err)
			}
			if !c.valid {
				return // No point in checking an invalid datagram
			}

			// Len
			if len(c.dg.buf) != c.len {
				t.Errorf("expected %s to have len %d, but it was %d", c.dg, c.len, len(c.dg.buf))
			}

			// Data
			if c.data != nil && !bytes.Equal(c.dg.data(), c.data) {
				t.Errorf("expected %s, to have data %q, but it was %q", c.dg, c.data, c.dg.data())
			}

			// Offset
			if c.offset != c.dg.offset {
				t.Errorf("expected %s to have offset %d, but it was %d", c.dg, c.offset, c.dg.offset)
			}

			// Code
			if c.code != c.dg.opcode() {
				t.Errorf("expected %s to have code %d, but it was %d", c.dg, c.code, c.dg.opcode())
			}

			// Filename
			if c.filename != nil && *c.filename != c.dg.filename() {
				t.Errorf("expected %s to have filename %q, but it was %q", c.dg, *c.filename, c.dg.filename())
			}

			// Mode
			if c.mode != nil && *c.mode != c.dg.mode() {
				t.Errorf("expected %s to have mode %q, but it was %q", c.dg, *c.mode, c.dg.mode())
			}

			// Options
			if c.opts != nil && !reflect.DeepEqual(c.opts, c.dg.options()) {
				t.Errorf("expected %s to have options %q, but it was %q", c.dg, c.opts, c.dg.options())
			}

			// Error Code
			if c.errCode != nil && *c.errCode != c.dg.errorCode() {
				t.Errorf("expected %s to have error code %d, but it was %d", c.dg, *c.errCode, c.dg.errorCode())
			}

			// Error Message
			if c.errMessage != nil && *c.errMessage != c.dg.errMsg() {
				t.Errorf("expected %s to have error message %q, but it was %q", c.dg, *c.errMessage, c.dg.errMsg())
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}

func ptrMode(s TransferMode) *TransferMode {
	return &s
}

func ptrErrCode(e ErrorCode) *ErrorCode {
	return &e
}
