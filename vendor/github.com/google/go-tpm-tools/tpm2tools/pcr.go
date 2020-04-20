package tpm2tools

import (
	"bytes"
	"crypto"
	"fmt"
	"io"

	tpmpb "github.com/google/go-tpm-tools/proto"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

// NumPCRs is set to the spec minimum of 24, as that's all go-tpm supports.
const NumPCRs = 24

// We hard-code SHA256 as the policy session hash algorithms. Note that this
// differs from the PCR hash algorithm (which selects the bank of PCRs to use)
// and the Public area Name algorithm. We also chose this for compatibility with
// github.com/google/go-tpm/tpm2, as it hardcodes the nameAlg as SHA256 in
// several places. Two constants are used to avoid repeated conversions.
const sessionHashAlg = crypto.SHA256
const sessionHashAlgTpm = tpm2.AlgSHA256

// CertifyHashAlgTpm is the hard-coded algorithm used in certify PCRs.
const CertifyHashAlgTpm = tpm2.AlgSHA256

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ReadPCRs fetches all the PCR values specified in sel, making multiple calls
// to the TPM if necessary.
func ReadPCRs(rw io.ReadWriter, sel tpm2.PCRSelection) (*tpmpb.Pcrs, error) {
	pl := tpmpb.Pcrs{
		Hash: tpmpb.HashAlgo(sel.Hash),
		Pcrs: map[uint32][]byte{},
	}

	for i := 0; i < len(sel.PCRs); i += 8 {
		end := min(i+8, len(sel.PCRs))
		pcrSel := tpm2.PCRSelection{
			Hash: sel.Hash,
			PCRs: sel.PCRs[i:end],
		}

		pcrMap, err := tpm2.ReadPCRs(rw, pcrSel)
		if err != nil {
			return nil, err
		}

		for pcr, val := range pcrMap {
			pl.Pcrs[uint32(pcr)] = val
		}
	}

	return &pl, nil
}

// SealCurrent seals data to the current specified PCR selection.
type SealCurrent struct{ tpm2.PCRSelection }

// SealTarget predicatively seals data to the given specified PCR values.
type SealTarget struct{ Pcrs *tpmpb.Pcrs }

// SealOpt specifies the PCR values that should be used for Seal().
type SealOpt interface {
	PCRsForSealing(rw io.ReadWriter) (*tpmpb.Pcrs, error)
}

// PCRsForSealing read from TPM and return the selected PCRs.
func (p SealCurrent) PCRsForSealing(rw io.ReadWriter) (*tpmpb.Pcrs, error) {
	if len(p.PCRSelection.PCRs) == 0 {
		panic("SealCurrent contains 0 PCRs")
	}
	return ReadPCRs(rw, p.PCRSelection)
}

// PCRsForSealing return the target PCRs.
func (p SealTarget) PCRsForSealing(_ io.ReadWriter) (*tpmpb.Pcrs, error) {
	if len(p.Pcrs.GetPcrs()) == 0 {
		panic("SealTaget contains 0 PCRs")
	}
	return p.Pcrs, nil
}

// CertifyCurrent certifies that a selection of current PCRs have the same value when sealing.
// Hash Algorithm in the selection should be CertifyHashAlgTpm.
type CertifyCurrent struct{ tpm2.PCRSelection }

// CertifyExpected certifies that the TPM had a specific set of PCR values when sealing.
// Hash Algorithm in the PCR proto should be CertifyHashAlgTpm.
type CertifyExpected struct{ Pcrs *tpmpb.Pcrs }

// CertifyOpt determines if the given PCR value can pass certification in Unseal().
type CertifyOpt interface {
	CertifyPCRs(rw io.ReadWriter, certified *tpmpb.Pcrs) error
}

// CertifyPCRs from CurrentPCRs will read PCR values from TPM and compare the digest.
func (p CertifyCurrent) CertifyPCRs(rw io.ReadWriter, pcrs *tpmpb.Pcrs) error {
	if len(p.PCRSelection.PCRs) == 0 {
		panic("CertifyCurrent contains 0 PCRs")
	}
	current, err := ReadPCRs(rw, p.PCRSelection)
	if err != nil {
		return err
	}
	return checkContainedPCRs(current, pcrs)
}

// CertifyPCRs will compare the digest with given expected PCRs values.
func (p CertifyExpected) CertifyPCRs(_ io.ReadWriter, pcrs *tpmpb.Pcrs) error {
	if len(p.Pcrs.GetPcrs()) == 0 {
		panic("CertifyExpected contains 0 PCRs")
	}
	return checkContainedPCRs(p.Pcrs, pcrs)
}

// Check if the "superset" PCRs contain a valid "subset" PCRs, the PCR value must match
// If there is one or more PCRs in subset which don't exist in superset, will return
// an error with the first missing PCR.
// If there is one or more PCRs value mismatch with the superset, will return an error
// with the first mismatched PCR numbers.
func checkContainedPCRs(subset *tpmpb.Pcrs, superset *tpmpb.Pcrs) error {
	if subset.GetHash() != superset.GetHash() {
		return fmt.Errorf("PCR hash algo not matching: %v, %v", subset.GetHash(), superset.GetHash())
	}
	for pcrNum, pcrVal := range subset.GetPcrs() {
		if expectedVal, ok := superset.GetPcrs()[pcrNum]; ok {
			if !bytes.Equal(expectedVal, pcrVal) {
				return fmt.Errorf("PCR %d mismatch: expected %v, got %v", pcrNum, expectedVal, pcrVal)
			}
		} else {
			return fmt.Errorf("PCR %d mismatch: value missing from the superset PCRs", pcrNum)
		}
	}
	return nil
}

// PCRSelection returns the corresponding tpm2.PCRSelection for a tpmpb.Pcrs
func PCRSelection(pcrs *tpmpb.Pcrs) tpm2.PCRSelection {
	sel := tpm2.PCRSelection{Hash: tpm2.Algorithm(pcrs.GetHash())}

	for pcrNum := range pcrs.GetPcrs() {
		sel.PCRs = append(sel.PCRs, int(pcrNum))
	}
	return sel
}

// HasSamePCRSelection checks the given tpmpb.Pcrs has the same PCRSelection as the
// given tpm2.PCRSelection (including the hash algorithm).
func HasSamePCRSelection(pcrs *tpmpb.Pcrs, pcrSel tpm2.PCRSelection) bool {
	if tpm2.Algorithm(pcrs.Hash) != pcrSel.Hash {
		return false
	}
	if len(pcrs.GetPcrs()) != len(pcrSel.PCRs) {
		return false
	}
	for _, p := range pcrSel.PCRs {
		if _, ok := pcrs.Pcrs[uint32(p)]; !ok {
			return false
		}
	}
	return true
}

// FullPcrSel will return a full PCR selection based on the total PCR number
// of the TPM with the given hash algo.
func FullPcrSel(hash tpm2.Algorithm) tpm2.PCRSelection {
	sel := tpm2.PCRSelection{Hash: hash}
	for i := 0; i < NumPCRs; i++ {
		sel.PCRs = append(sel.PCRs, int(i))
	}
	return sel
}

// ComputePCRSessionAuth calculates the authorization value for the given PCRs.
func ComputePCRSessionAuth(pcrs *tpmpb.Pcrs) []byte {
	// Start with all zeros, we only use a single policy command on our session.
	oldDigest := make([]byte, sessionHashAlg.Size())
	ccPolicyPCR, _ := tpmutil.Pack(tpm2.CmdPolicyPCR)

	// Extend the policy digest, see TPM2_PolicyPCR in Part 3 of the spec.
	hash := sessionHashAlg.New()
	hash.Write(oldDigest)
	hash.Write(ccPolicyPCR)
	hash.Write(encodePCRSelection(PCRSelection(pcrs)))
	hash.Write(computePCRDigest(pcrs))
	newDigest := hash.Sum(nil)
	return newDigest[:]
}

// ComputePCRDigest will take in a PCR proto and compute the digest based on the
// given PCR proto.
func computePCRDigest(pcrs *tpmpb.Pcrs) []byte {
	hash := sessionHashAlg.New()
	for i := 0; i < 24; i++ {
		if pcrValue, exists := pcrs.Pcrs[uint32(i)]; exists {
			hash.Write(pcrValue)
		}
	}
	return hash.Sum(nil)
}

// Encode a tpm2.PCRSelection as if it were a TPML_PCR_SELECTION
func encodePCRSelection(sel tpm2.PCRSelection) []byte {
	// Encode count, pcrSelections.hash and pcrSelections.sizeofSelect fields
	buf, _ := tpmutil.Pack(uint32(1), sel.Hash, byte(3))
	// Encode pcrSelect bitmask
	pcrBits := make([]byte, 3)
	for _, pcr := range sel.PCRs {
		byteNum := pcr / 8
		bytePos := 1 << uint(pcr%8)
		pcrBits[byteNum] |= byte(bytePos)
	}

	return append(buf, pcrBits...)
}
