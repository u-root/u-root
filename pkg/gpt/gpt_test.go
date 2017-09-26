package gpt

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"

	"github.com/google/uuid"
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
		DiskGUID:   uuid.Must(uuid.Parse("2d1ed4ba-ef93-4ab0-846e-2e3a6a2d73bf")),
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

func TestGPTTable(t *testing.T) {
	var tests = []struct {
		mangle int
		msg    string
	}{
		{-1, ""},
		{0x8, "Primary GPT revision (100ff) is not supported value (10000)"},
		{0x0, "Primary GPT signature invalid (54524150204946ff), needs to be 5452415020494645"},
		{0xf, "Primary GPT HeaderSize (ff00005c) is not supported value (5c)"},
		{0x51, "Primary GPT MaxNPart (ff80) is above maximum of 80"},
		{0x59, "Primary Partition CRC: Header {\n\t\"Signature\": 6075990659671082565,\n\t\"Revision\": 65536,\n\t\"HeaderSize\": 92,\n\t\"CRC\": 575771282,\n\t\"Reserved\": 0,\n\t\"CurrentLBA\": 1,\n\t\"BackupLBA\": 4444095,\n\t\"FirstLBA\": 34,\n\t\"LastLBA\": 4444062,\n\t\"DiskGUID\": \"2d1ed4ba-ef93-4ab0-846e-2e3a6a2d73bf\",\n\t\"PartStart\": 2,\n\t\"NPart\": 128,\n\t\"PartSize\": 128,\n\t\"PartCRC\": 2373123927,\n\t\"Parts\": null\n}, computed checksum is 8d728e57, header has 8d72ff57"},
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

		t.Logf("New GPT: %s", g)
		if !reflect.DeepEqual(header, g.Header) {
			t.Errorf("Check UUID equality from\n%v to\n%v: got false, want true", header, g.Header)
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
	var tests = []struct {
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
		_, _, err := New(r)
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

type iodisk []byte

func (d *iodisk) WriteAt(b []byte, off int64) (int, error) {
	copy([]byte(*d)[off:], b)
	return len(b), nil
}

func TestWrite(t *testing.T) {
	InstallGPT()
	r := bytes.NewReader(disk)
	g, err := Table(r, BlockSize)
	if err != nil {
		t.Fatalf("Reading table: got %v, want nil", err)
	}
	var targ = make(iodisk, len(disk))

	if err := Write(&targ, g); err != nil {
		t.Fatalf("Writing: got %v, want nil", err)
	}
	if n, err := Table(bytes.NewReader([]byte(targ)), BlockSize); err != nil {
		t.Logf("Old GPT: %s", g)
		var b bytes.Buffer
		w := hex.Dumper(&b)
		io.Copy(w, bytes.NewReader(disk[:4096]))
		t.Logf("%s\n", string(b.Bytes()))
		t.Fatalf("Reading back new header: new:%s\n%v", n, err)
	}
}
