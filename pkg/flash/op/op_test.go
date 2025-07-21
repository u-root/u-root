// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package op_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/flash/op"
)

func TestString(t *testing.T) {
	bad := op.OpCode(0xff)
	tests := []struct {
		opcode   op.OpCode
		expected string
	}{
		{op.PageProgram, "PageProgram"},
		{op.Read, "Read"},
		{op.WriteDisable, "WriteDisable"},
		{op.ReadStatus, "ReadStatus"},
		{op.WriteEnable, "WriteEnable"},
		{op.SectorErase, "SectorErase"},
		{op.ReadSFDP, "ReadSFDP"},
		{op.ReadJEDECID, "ReadJEDECID"},
		{op.PRDRES, "PRDRES"},
		{op.Enter4BA, "Enter4BA"},
		{op.BlockErase, "BlockErase"},
		{op.Exit4BA, "Exit4BA"},
		{bad, fmt.Sprintf("Unknown(%02x)", byte(bad))},
	}

	for _, tt := range tests {
		actual := tt.opcode.String()
		if actual != tt.expected {
			t.Errorf("String() for %v: expected %s, got %s", tt.opcode, tt.expected, actual)
		}
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		opcode   op.OpCode
		expected []byte
	}{
		{op.PageProgram, []byte{byte(op.PageProgram)}},
	}

	for _, tt := range tests {
		actual := tt.opcode.Bytes()
		if !bytes.Equal(actual, tt.expected) {
			t.Errorf("Bytes() for %v: expected %v, got %v", tt.opcode, tt.expected, actual)
		}
	}
}

func TestStatus(t *testing.T) {
	for i, tt := range []struct {
		name string
		val  op.Status
		out  string
	}{
		{name: "nothing", out: ""},
		{name: "write enabled", val: op.WriteEnabled, out: "WriteEnabled"},
		{name: "busy, write enabled", val: op.WriteBusy | op.WriteEnabled, out: "WriteBusy|WriteEnabled"},
	} {
		s := tt.val.String()
		if s != tt.out {
			t.Errorf("%d:%#x: got %s, want %s", i, tt.name, s, tt.out)
		}
	}
}

func TestStatusBus(t *testing.T) {
	for i, tt := range []struct {
		name string
		val  op.Status
		out  bool
	}{
		{name: "idle"},
		{name: "busy", val: op.WriteBusy, out: true},
		{name: "other than busy", val: op.WriteEnabled, out: false},
	} {
		busy := tt.val.Busy()
		if busy != tt.out {
			t.Errorf("%d:%#x: got %v, want %v", i, tt.name, busy, tt.out)
		}
	}
}
