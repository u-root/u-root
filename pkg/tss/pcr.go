// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"fmt"
	"io"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpmutil"
)

func extendPCR12(rwc io.ReadWriter, pcrIndex uint32, hash [20]byte) error {
	if _, err := tpm.PcrExtend(rwc, pcrIndex, hash); err != nil {
		return err
	}
	return nil
}

func extendPCR20(rwc io.ReadWriter, pcrIndex uint32, hash []byte) error {
	if err := tpm2.PCRExtend(rwc, tpmutil.Handle(pcrIndex), tpm2.AlgSHA256, hash, ""); err != nil {
		return err
	}
	return nil
}

func readAllPCRs20(tpm io.ReadWriter, alg tpm2.Algorithm) (map[uint32][]byte, error) {
	numPCRs := 24
	out := map[uint32][]byte{}

	// The TPM 2.0 spec says that the TPM can partially fulfill the
	// request. As such, we repeat the command up to 8 times to get all
	// 24 PCRs.
	for i := 0; i < numPCRs; i++ {
		// Build a selection structure, specifying all PCRs we do
		// not have the value for.
		sel := tpm2.PCRSelection{Hash: alg}
		for pcr := 0; pcr < numPCRs; pcr++ {
			if _, present := out[uint32(pcr)]; !present {
				sel.PCRs = append(sel.PCRs, pcr)
			}
		}

		// Ask the TPM for those PCR values.
		ret, err := tpm2.ReadPCRs(tpm, sel)
		if err != nil {
			return nil, fmt.Errorf("tpm2.ReadPCRs(%+v) failed with err: %w", sel, err)
		}
		// Keep track of the PCRs we were actually given.
		for pcr, digest := range ret {
			out[uint32(pcr)] = digest
		}
		if len(out) == numPCRs {
			break
		}
	}

	if len(out) != numPCRs {
		return nil, fmt.Errorf("failed to read all PCRs, only read %d", len(out))
	}

	return out, nil
}

func readAllPCRs12(rwc io.ReadWriter) (map[uint32][]byte, error) {
	numPCRs := 24
	out := map[uint32][]byte{}

	for i := 0; i < numPCRs; i++ {
		// Ask the TPM for those PCR values.
		pcr, err := tpm.ReadPCR(rwc, uint32(i))
		if err != nil {
			return nil, fmt.Errorf("tpm.ReadPCR(%d) failed with err: %w", i, err)
		}
		out[uint32(i)] = pcr
		if len(out) == numPCRs {
			break
		}
	}

	if len(out) != numPCRs {
		return nil, fmt.Errorf("failed to read all PCRs, only read %d", len(out))
	}

	return out, nil
}

func readPCR12(rwc io.ReadWriter, pcrIndex uint32) ([]byte, error) {
	return tpm.ReadPCR(rwc, pcrIndex)
}

func readPCR20(rwc io.ReadWriter, pcrIndex uint32) ([]byte, error) {
	return tpm2.ReadPCR(rwc, int(pcrIndex), tpm2.AlgSHA256)
}
