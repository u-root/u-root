package tpm2tools

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/google/go-tpm-tools/internal"
	tpmpb "github.com/google/go-tpm-tools/proto"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var tests = []struct {
	inAlg        tpm2.Algorithm
	inExtensions [][]byte
}{
	{tpm2.AlgSHA1, nil},
	{tpm2.AlgSHA1, [][]byte{bytes.Repeat([]byte{0x00}, sha1.Size)}},
	{tpm2.AlgSHA1, [][]byte{bytes.Repeat([]byte{0x01}, sha1.Size)}},
	{tpm2.AlgSHA1, [][]byte{bytes.Repeat([]byte{0x02}, sha1.Size)}},

	{tpm2.AlgSHA256, nil},
	{tpm2.AlgSHA256, [][]byte{bytes.Repeat([]byte{0x00}, sha256.Size)}},
	{tpm2.AlgSHA256, [][]byte{bytes.Repeat([]byte{0x01}, sha256.Size)}},
	{tpm2.AlgSHA256, [][]byte{bytes.Repeat([]byte{0x02}, sha256.Size)}},

	{tpm2.AlgSHA384, nil},
	{tpm2.AlgSHA384, [][]byte{bytes.Repeat([]byte{0x00}, sha512.Size384)}},
	{tpm2.AlgSHA384, [][]byte{bytes.Repeat([]byte{0x01}, sha512.Size384)}},
	{tpm2.AlgSHA384, [][]byte{bytes.Repeat([]byte{0x02}, sha512.Size384)}},
}

func pcrExtend(alg tpm2.Algorithm, old, new []byte) ([]byte, error) {
	hCon, err := alg.Hash()
	if err != nil {
		return nil, fmt.Errorf("not a valid hash type: %v", alg)
	}
	h := hCon.New()
	h.Write(old)
	h.Write(new)
	return h.Sum(nil), nil
}

func TestReadPCRs(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	testPcrs := make(map[tpm2.Algorithm][]byte, 2)
	testPcrs[tpm2.AlgSHA1] = bytes.Repeat([]byte{0x00}, sha1.Size)
	testPcrs[tpm2.AlgSHA256] = bytes.Repeat([]byte{0x00}, sha256.Size)
	testPcrs[tpm2.AlgSHA384] = bytes.Repeat([]byte{0x00}, sha512.Size384)

	for _, test := range tests {
		for _, extension := range test.inExtensions {
			if err := tpm2.PCRExtend(rwc, tpmutil.Handle(0), test.inAlg, extension, ""); err != nil {
				t.Fatalf("failed to extend pcr for test %v", err)
			}

			pcrVal, err := pcrExtend(test.inAlg, testPcrs[test.inAlg], extension)
			if err != nil {
				t.Fatalf("could not extend pcr: %v", err)
			}
			testPcrs[test.inAlg] = pcrVal
		}

		sel := tpm2.PCRSelection{Hash: test.inAlg, PCRs: []int{0}}
		proto, err := ReadPCRs(rwc, sel)
		if err != nil {
			t.Fatalf("failed to read pcrs %v", err)
		}

		if !bytes.Equal(proto.Pcrs[0], testPcrs[test.inAlg]) {
			t.Errorf("%v not equal to expected %v", proto.Pcrs[0], testPcrs[test.inAlg])
		}
	}
}

func TestCheckContainedPCRs(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	sel := FullPcrSel(tpm2.AlgSHA256)
	baseline, err := ReadPCRs(rwc, sel)
	if err != nil {
		t.Fatalf("Failed to Read PCRs: %v", err)
	}

	toBeCertified, err := ReadPCRs(rwc, tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{1, 2, 3}})
	if err != nil {
		t.Fatalf("failed to read pcrs %v", err)
	}
	if err := checkContainedPCRs(toBeCertified, baseline); err != nil {
		t.Fatalf("Validation should pass: %v", err)
	}

	if err := tpm2.PCRExtend(rwc, tpmutil.Handle(2), tpm2.AlgSHA256, bytes.Repeat([]byte{0x00}, sha256.Size), ""); err != nil {
		t.Fatalf("failed to extend pcr for test %v", err)
	}

	toBeCertified, err = ReadPCRs(rwc, tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{1, 2, 3}})
	if err != nil {
		t.Fatalf("failed to read pcrs %v", err)
	}
	if err := checkContainedPCRs(toBeCertified, baseline); err == nil {
		t.Fatalf("validation should fail due to PCR 2 changed")
	}

	toBeCertified, err = ReadPCRs(rwc, tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{}})
	if err != nil {
		t.Fatalf("failed to read pcrs %v", err)
	}
	if err := checkContainedPCRs(toBeCertified, baseline); err != nil {
		t.Fatalf("empty pcrs is always validate")
	}
}

func TestHasSamePCRSelection(t *testing.T) {
	var tests = []struct {
		pcrs        *tpmpb.Pcrs
		pcrSel      tpm2.PCRSelection
		expectedRes bool
	}{
		{&tpmpb.Pcrs{}, tpm2.PCRSelection{}, true},
		{&tpmpb.Pcrs{Hash: tpmpb.HashAlgo(tpm2.AlgSHA256), Pcrs: map[uint32][]byte{1: []byte{}}},
			tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{1}}, true},
		{&tpmpb.Pcrs{Hash: tpmpb.HashAlgo(tpm2.AlgSHA256), Pcrs: map[uint32][]byte{}},
			tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{}}, true},
		{&tpmpb.Pcrs{Hash: tpmpb.HashAlgo(tpm2.AlgSHA256), Pcrs: map[uint32][]byte{1: []byte{}}},
			tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{4}}, false},
		{&tpmpb.Pcrs{Hash: tpmpb.HashAlgo(tpm2.AlgSHA256), Pcrs: map[uint32][]byte{1: []byte{}, 4: []byte{}}},
			tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{4}}, false},
		{&tpmpb.Pcrs{Hash: tpmpb.HashAlgo(tpm2.AlgSHA256), Pcrs: map[uint32][]byte{1: []byte{}, 2: []byte{}}},
			tpm2.PCRSelection{Hash: tpm2.AlgSHA1, PCRs: []int{1, 2}}, false},
	}
	for _, test := range tests {
		if HasSamePCRSelection(test.pcrs, test.pcrSel) != test.expectedRes {
			t.Errorf("HasSamePCRSelection result is not expected")
		}
	}
}
