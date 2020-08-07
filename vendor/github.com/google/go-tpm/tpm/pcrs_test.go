// Copyright (c) 2014, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpm

import (
	"testing"
)

func TestPCRMask(t *testing.T) {
	var mask pcrMask
	if err := mask.setPCR(-1); err == nil {
		t.Fatal("Incorrectly allowed non-existent PCR -1 to be set")
	}

	if err := mask.setPCR(24); err == nil {
		t.Fatal("Incorrectly allowed non-existent PCR 24 to be set")
	}

	if err := mask.setPCR(0); err != nil {
		t.Fatal("Couldn't set PCR 0 in the mask:", err)
	}

	set, err := mask.isPCRSet(0)
	if err != nil {
		t.Fatal("Couldn't check to see if PCR 0 was set:", err)
	}

	if !set {
		t.Fatal("Incorrectly said PCR wasn't set when it should have been")
	}

	if err := mask.setPCR(18); err != nil {
		t.Fatal("Couldn't set PCR 18 in the mask:", err)
	}

	set, err = mask.isPCRSet(18)
	if err != nil {
		t.Fatal("Couldn't check to see if PCR 18 was set:", err)
	}

	if !set {
		t.Fatal("Incorrectly said PCR wasn't set when it should have been")
	}

	if _, err := mask.isPCRSet(-1); err == nil {
		t.Fatal("Incorrectly permitted a check for PCR -1")
	}

	if _, err := mask.isPCRSet(400); err == nil {
		t.Fatal("Incorrectly permitted a check for PCR 400")
	}
}

func TestNewPCRSelection(t *testing.T) {
	pcrs, err := newPCRSelection([]int{17, 18})
	if err != nil {
		t.Fatal("Couldn't set up a PCR selection with PCRs 17 and 18")
	}

	if pcrs.Size != 3 {
		t.Fatal("Incorrectly size in a PCR selection")
	}

	set, err := pcrs.Mask.isPCRSet(17)
	if err != nil {
		t.Fatal("Couldn't check a PCR on a mask in a PCR selection")
	}

	if !set {
		t.Fatal("PCR 17 wasn't set in a PCR selection after setting it")
	}

	set, err = pcrs.Mask.isPCRSet(20)
	if err != nil {
		t.Fatal("Couldn't check an unset PCR on a mask in a PCR selection")
	}

	if set {
		t.Fatal("PCR 20 was incorrectly set in a PCR mask in a PCR selection")
	}
}

func TestIncorrectCreatePCRComposite(t *testing.T) {
	pcrs, err := newPCRSelection([]int{17, 18})
	if err != nil {
		t.Fatal("Couldn't set up a PCR selection with PCRs 17 and 18")
	}

	// This byte array is far too long and isn't a multiple of PCRSize, since it
	// is of prime size.
	pcrValues := make([]byte, 541)
	if _, err := createPCRComposite(pcrs.Mask, pcrValues); err == nil {
		t.Fatal("Incorrectly created a PCR composite with wrong PCR length")
	}
}

func TestWrongCreatePCRInfoLong(t *testing.T) {
	pcrs, err := newPCRSelection([]int{17, 18})
	if err != nil {
		t.Fatal("Couldn't set up a PCR selection with PCRs 17 and 18")
	}

	// This byte array is far too long and isn't a multiple of PCRSize, since it
	// is of prime size.
	pcrValues := make([]byte, 541)
	if _, err := createPCRInfoLong(0, pcrs.Mask, pcrValues); err == nil {
		t.Fatal("Incorrectly created a PCR composite with wrong PCR length")
	}
}

func TestWrongNewPCRInfoLong(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	if _, err := newPCRInfoLong(rwc, 0, []int{400}); err == nil {
		t.Fatal("Incorrectly created a pcrInfoLong for PCR 400")
	}

	// This case uses a reasonable PCR value but a nil file.
	if _, err := newPCRInfoLong(nil, 0, []int{17}); err == nil {
		t.Fatal("Incorrectly created a pcrInfoLong using a nil file")
	}
}

func TestNewPCRInfoLongWithHashes(t *testing.T) {
	pcrMap := make(map[int][]byte)
	pcrMap[23] = make([]byte, 20)
	pcrMap[16] = make([]byte, 20)

	if _, err := newPCRInfoLongWithHashes(LocZero, pcrMap); err != nil {
		t.Fatal("Couldn't create pcrInfoLong structure")
	}
}
