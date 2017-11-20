package gpt

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const (
	equalHeaderError = "p.Signature(0x5452415020494646) != b.Signature(0x5452415020494645); p.Revision(65537) != b.Revision(65536); p.HeaderSize(93) != b.HeaderSize(92); p.CurrentLBA(0x2) != b.BackupLBA(0x1); p.FirstLBA(0x23) != b.FirstLBA(0x22); p.LastLBA(0x43cf9f) != b.LastLBA(0x43cf9e); p.DiskGUID(0x32653165643462612d656639332d346162302d383436652d326533613661326437336266) != b.DiskGUID(0x32643165643462612d656639332d346162302d383436652d326533613661326437336266); p.NPart(127) != b.NPart(128); p.PartSize(127) != b.PartSize(128)"
	equalPartsError  = "Partition 3: p.PartGUID(0x35653261336166652d333234662d613734312d623732352d616363633332383561333039) != b.PartGUID(0x35643261336166652d333234662d613734312d623732352d616363633332383561333039); Partition 8: p.UniqueGUID(0x65643939336135312d336563342d346131342d383339382d613437613765366165343263) != b.UniqueGUID(0x65643938336135312d336563342d346131342d383339382d613437613765366165343263); Partition 11: p.FirstLBA(0x3d001) != b.FirstLBA(0x3d000); Partition 21: p.LastLBA(0x1) != b.LastLBA(0x0); Partition 61: p.Name(0x000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000) != b.Name(0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000)"
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

// TestEqualHeader tests all variations of headers not being equal.
// We test to make sure it works, then break some aspect of the header
// and test that too.
func TestEqualHeader(t *testing.T) {
	InstallGPT()
	r := bytes.NewReader(disk)
	p, b, err := New(r)
	if err != nil {
		t.Fatalf("TestEqualHeader: Reading in gpt: got %v, want nil", err)
	}

	if err := EqualHeader(p.Header, b.Header); err != nil {
		t.Fatalf("TestEqualHeader: got %v, want nil", err)
	}
	// Yes, we assume a certain order, but it sure simplifies the test :-)
	p.Signature++
	p.Revision++
	p.HeaderSize++
	p.CurrentLBA++
	p.FirstLBA++
	p.LastLBA++
	p.DiskGUID[0]++
	p.NPart--
	p.PartSize--
	p.PartCRC++
	if err = EqualHeader(p.Header, b.Header); err == nil {
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
	p, b, err := New(r)
	if err != nil {
		t.Fatalf("TestEqualParts: Reading in gpt: got %v, want nil", err)
	}

	if err = EqualParts(p, b); err != nil {
		t.Fatalf("TestEqualParts: Checking equality: got %v, want nil", err)
	}
	// Test some equality things before we do the 'length is the same' test
	// Note that testing the NParts header variable is done in the HeaderTest
	p.Parts[3].PartGUID[0]++
	p.Parts[8].UniqueGUID[1]++
	p.Parts[11].FirstLBA++
	p.Parts[21].LastLBA++
	p.Parts[53].Attribute++
	p.Parts[61].Name[1]++
	if err = EqualParts(p, b); err == nil {
		t.Errorf("TestEqualParts: Checking equality: got nil, want '%v'", equalPartsError)
	}
	if err.Error() != equalPartsError {
		t.Errorf("TestEqualParts: Checking equality: got '%v', want '%v'", err, equalPartsError)
	}

	if err = EqualParts(p, b); err == nil {
		t.Errorf("TestEqualParts: Checking number of parts: got nil, want 'Primary Number of partitions (127) differs from Backup (128)'")
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
