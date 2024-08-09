// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package gpt

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	equalHeaderError = "p.Signature(0x5452415020494646) != b.Signature(0x5452415020494645); p.Revision(65537) != b.Revision(65536); p.HeaderSize(93) != b.HeaderSize(92); p.CurrentLBA(0x2) != b.BackupLBA(0x1); p.FirstLBA(0x23) != b.FirstLBA(0x22); p.LastLBA(0x43cf9f) != b.LastLBA(0x43cf9e); p.DiskGUID({0xbad41e2d 0x93ef 0xb04a 0x856e2e3a6a2d73bf}) != b.DiskGUID({0xbad41e2d 0x93ef 0xb04a 0x846e2e3a6a2d73bf}); p.NPart(127) != b.NPart(128); p.PartSize(127) != b.PartSize(128)"
	equalPartsError  = "Partition 3: p.PartGUID({0xfe3a2a5d 0x4f32 0x41a7 0xb825accc3285a309}) != b.PartGUID({0xfe3a2a5d 0x4f32 0x41a7 0xb725accc3285a309}); Partition 8: p.UniqueGUID({0x513a98ed 0xc43e 0x144a 0x8399a47a7e6ae42c}) != b.UniqueGUID({0x513a98ed 0xc43e 0x144a 0x8398a47a7e6ae42c}); Partition 11: p.FirstLBA(0x3d001) != b.FirstLBA(0x3d000); Partition 21: p.LastLBA(0x1) != b.LastLBA(0x0); Partition 61: p.Name(0x000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000) != b.Name(0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000)"
)

var (
	header = Header{
		Signature:  Signature,
		Revision:   Revision,
		HeaderSize: HeaderSize,
		CRC:        0x22519292,
		Reserved:   0,
		CurrentLBA: 1,
		BackupLBA:  0x43cfbf,
		FirstLBA:   0x22,
		LastLBA:    0x43cf9e,
		DiskGUID:   GUID{L: 0xbad41e2d, W1: 0x93ef, W2: 0xb04a, B: [8]byte{0x84, 0x6e, 0x2e, 0x3a, 0x6a, 0x2d, 0x73, 0xbf}},
		PartStart:  2,
		NPart:      MaxNPart,
		PartSize:   0x80, // This is not constant, but was used for this chromeos disk.
		PartCRC:    0x8d728e57,
	}
	disk = make([]byte, 0x100000000)
)

func InstallGPT() {
	for i, d := range block {
		copy(disk[i:], d)
	}
}

// GPT is GUID Partition Table, so technically, this test name is
// Test Guid Partition Table Table. :-)
func TestGPTTable(t *testing.T) {
	tests := []struct {
		mangle int
		msg    string
	}{
		{-1, ""},
		{0x8, "Primary GPT revision (100ff) is not supported value (10000)"},
		{0x0, "Primary GPT signature invalid (54524150204946ff), needs to be 5452415020494645"},
		{0xf, "Primary GPT HeaderSize (ff00005c) is not supported value (5c)"},
		{0x51, "Primary GPT MaxNPart (ff80) is above maximum of 80"},
		{0x59, "Primary Partition CRC: Header {\n\t\"Signature\": 6075990659671082565,\n\t\"Revision\": 65536,\n\t\"HeaderSize\": 92,\n\t\"CRC\": 575771282,\n\t\"Reserved\": 0,\n\t\"CurrentLBA\": 1,\n\t\"BackupLBA\": 4444095,\n\t\"FirstLBA\": 34,\n\t\"LastLBA\": 4444062,\n\t\"DiskGUID\": {\n\t\t\"L\": 3134463533,\n\t\t\"W1\": 37871,\n\t\t\"W2\": 45130,\n\t\t\"B\": [\n\t\t\t132,\n\t\t\t110,\n\t\t\t46,\n\t\t\t58,\n\t\t\t106,\n\t\t\t45,\n\t\t\t115,\n\t\t\t191\n\t\t]\n\t},\n\t\"PartStart\": 2,\n\t\"NPart\": 128,\n\t\"PartSize\": 128,\n\t\"PartCRC\": 2373123927,\n\t\"Parts\": null\n}, computed checksum is 8d728e57, header has 8d72ff57"},
		{0x10, "Primary Header CRC: computed checksum is 22519292, header has 225192ff"},
	}

	for _, test := range tests {
		InstallGPT()
		if test.mangle > -1 {
			disk[BlockSize+test.mangle] = 0xff
		}
		r := bytes.NewReader(disk)
		g, err := Table(r, BlockSize)
		if err != nil {
			if err.Error() != test.msg {
				t.Errorf("New GPT: got %q, want %q", err, test.msg)
				continue
			}
			t.Logf("Got expected error %q", err)
			continue
		}

		if err == nil && test.msg != "" {
			t.Errorf("New GPT: got nil, want %s", test.msg)
			continue
		}

		if !reflect.DeepEqual(header, g.Header) {
			t.Errorf("Check GUID equality from\n%v to\n%v: got false, want true", header, g.Header)
			continue
		}
	}
}

// TestGPTTtables tests whether we can match the primary and backup
// or, if they differ, we catch that error.
// We know from other tests that the tables read fine.
// This test verifies that they match and that therefore we
// are able to read the backup table and test that it is ok.
func TestGPTTables(t *testing.T) {
	tests := []struct {
		mangle int
		what   string
	}{
		{-1, "No error test"},
		{0x10, "Should differ test"},
	}

	for _, test := range tests {
		InstallGPT()
		if test.mangle > -1 {
			disk[BlockSize+test.mangle] = 0xff
		}
		r := bytes.NewReader(disk)
		_, err := New(r)
		switch {
		case err != nil && test.mangle > -1:
			t.Logf("Got expected error %s", test.what)
		case err != nil && test.mangle == -1:
			t.Errorf("%s: got %s, want nil", test.what, err)
			continue
		case err == nil && test.mangle > -1:
			t.Errorf("%s: got nil, want err", test.what)
			continue
		}
		t.Logf("Passed %s", test.what)
	}
}

// TestEqualHeader tests all variations of headers not being equal.
// We test to make sure it works, then break some aspect of the header
// and test that too.
func TestEqualHeader(t *testing.T) {
	InstallGPT()
	r := bytes.NewReader(disk)
	p, err := New(r)
	if err != nil {
		t.Fatalf("TestEqualHeader: Reading in gpt: got %v, want nil", err)
	}

	if err := EqualHeader(p.Primary.Header, p.Backup.Header); err != nil {
		t.Fatalf("TestEqualHeader: got %v, want nil", err)
	}
	// Yes, we assume a certain order, but it sure simplifies the test :-)
	p.Primary.Signature++
	p.Primary.Revision++
	p.Primary.HeaderSize++
	p.Primary.CurrentLBA++
	p.Primary.FirstLBA++
	p.Primary.LastLBA++
	p.Primary.DiskGUID.B[0]++
	p.Primary.NPart--
	p.Primary.PartSize--
	p.Primary.PartCRC++
	if err = EqualHeader(p.Primary.Header, p.Backup.Header); err == nil {
		t.Fatalf("TestEqualHeader: got %v, want nil", err)
	}
	t.Logf("TestEqualHeader: EqualHeader returns %v", err)

	if err.Error() != equalHeaderError {
		t.Fatalf("TestEqualHeader: got %v, want %v", err.Error(), equalHeaderError)
	}
}

func TestEqualParts(t *testing.T) {
	InstallGPT()
	r := bytes.NewReader(disk)
	p, err := New(r)
	if err != nil {
		t.Fatalf("TestEqualParts: Reading in gpt: got %v, want nil", err)
	}

	if err = EqualParts(p.Primary, p.Backup); err != nil {
		t.Fatalf("TestEqualParts: Checking equality: got %v, want nil", err)
	}
	// Test some equality things before we do the 'length is the same' test
	// Note that testing the NParts header variable is done in the HeaderTest
	p.Primary.Parts[3].PartGUID.B[0]++
	p.Primary.Parts[8].UniqueGUID.B[1]++
	p.Primary.Parts[11].FirstLBA++
	p.Primary.Parts[21].LastLBA++
	p.Primary.Parts[53].Attribute++
	p.Primary.Parts[61].Name[1]++
	if err = EqualParts(p.Primary, p.Backup); err == nil {
		t.Errorf("TestEqualParts: Checking equality: got nil, want '%v'", equalPartsError)
	}
	if err.Error() != equalPartsError {
		t.Errorf("TestEqualParts: Checking equality: got '%v', want '%v'", err, equalPartsError)
	}

	if err = EqualParts(p.Primary, p.Backup); err == nil {
		t.Errorf("TestEqualParts: Checking number of parts: got nil, want 'Primary Number of partitions (127) differs from Backup (128)'")
	}
}

// writeLog is a history of []byte written to the iodisk. Each write to iodisk creates a new writeLog entry.
type writeLog [][]byte

// iodisk is a fake disk that is used for testing.
// Each write is logged into the `writes` map.
// iodisk implements the WriterAt interface and can be passed to Write() for testing.
type iodisk struct {
	bytes []byte

	// mapping of address=>writes.
	// This is used for verifying that the correct data was written into the correct locations.
	writes map[int64]writeLog
}

func newIOdisk(size int) *iodisk {
	return &iodisk{
		bytes:  make([]byte, size),
		writes: make(map[int64]writeLog),
	}
}

func (d *iodisk) WriteAt(b []byte, offset int64) (int, error) {
	copy([]byte(d.bytes)[offset:], b)
	d.writes[offset] = append(d.writes[offset], b)
	return len(b), nil
}

func TestWrite(t *testing.T) {
	InstallGPT()
	r := bytes.NewReader(disk)
	p, err := New(r)
	if err != nil {
		t.Fatalf("Reading partitions: got %v, want nil", err)
	}
	targ := newIOdisk(len(disk))

	if err := Write(targ, p); err != nil {
		t.Fatalf("Writing: got %v, want nil", err)
	}
	if n, err := New(bytes.NewReader([]byte(targ.bytes))); err != nil {
		t.Logf("Old GPT: %s", p.Primary)
		var b bytes.Buffer
		w := hex.Dumper(&b)
		io.Copy(w, bytes.NewReader(disk[:4096]))
		t.Logf("%s\n", b.String())
		t.Fatalf("Reading back new header: new:%s\n%v", n, err)
	}

	tests := []struct {
		desc   string
		offset int64
		size   int64
	}{
		{
			desc:   "Protective MBR",
			offset: 0x00000000,
			size:   BlockSize,
		},
		{
			desc:   "Primary GPT header",
			offset: 0x00000200,
			size:   BlockSize,
		},
		{
			desc:   "Backup GPT header",
			offset: 0x879f7e00,
			size:   BlockSize,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			// Verify that there was exactly one write.
			if count := len(targ.writes[tc.offset]); count != 1 {
				t.Fatalf("Expected exactly 1 write to address 0x%08x, got %d", tc.offset, count)
			}
			// Verify that the contents were exactly as expected.
			if !cmp.Equal(targ.writes[tc.offset][0], disk[tc.offset:tc.offset+tc.size]) {
				t.Fatalf("Data written to 0x%08x does not match the source data", tc.offset)
			}
		})
	}
}
