package tss

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"io"

	proto "github.com/google/go-tpm-tools/proto"
	tpm2tools "github.com/google/go-tpm-tools/tpm2tools"
	tpm1 "github.com/google/go-tpm/tpm"
	tpm2 "github.com/google/go-tpm/tpm2"
	tpmutil "github.com/google/go-tpm/tpmutil"
)

func defaultSymScheme() *tpm2.SymScheme {
	return &tpm2.SymScheme{
		Alg:     tpm2.AlgAES,
		KeyBits: 128,
		Mode:    tpm2.AlgCFB,
	}
}

func defaultECCParams() *tpm2.ECCParams {
	return &tpm2.ECCParams{
		Symmetric: defaultSymScheme(),
		CurveID:   tpm2.CurveNISTP256,
		Point: tpm2.ECPoint{
			XRaw: make([]byte, 32),
			YRaw: make([]byte, 32),
		},
	}
}

func loadSRK20(rwc io.ReadWriteCloser, srkPW string) (*tpm2tools.Key, error) {
	var srkAuth tpmutil.U16Bytes
	var hash [32]byte

	if srkPW != "" {
		hash = sha256.Sum256([]byte(srkPW))
	}

	srkPWBytes := bytes.NewBuffer(hash[:])
	err := srkAuth.TPMMarshal(srkPWBytes)
	if err != nil {
		return nil, err
	}
	srkTemplate := tpm2.Public{
		Type:          tpm2.AlgECC,
		NameAlg:       tpm2.AlgSHA256,
		Attributes:    tpm2.FlagStorageDefault,
		ECCParameters: defaultECCParams(),
		AuthPolicy:    srkAuth,
	}

	return tpm2tools.NewCachedKey(rwc, tpm2.HandleOwner, srkTemplate, tpm2tools.SRKECCReservedHandle)
}

func seal12(rwc io.ReadWriteCloser, srkPW string, pcrs []int, data []byte) ([]byte, error) {
	var srkAuth [20]byte

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}

	return tpm1.Seal(rwc, tpm1.Locality(1), pcrs, data, srkAuth[:])
}

func reseal12(rwc io.ReadWriteCloser, srkPW string, pcrs map[uint32][]byte, data []byte) ([]byte, error) {
	var srkAuth [20]byte
	var pcrMap map[int][]byte

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}
	for k, v := range pcrs {
		pcrMap[int(k)] = v
	}

	return tpm1.Reseal(rwc, tpm1.Locality(1), pcrMap, data, srkAuth[:])
}

func unseal12(rwc io.ReadWriteCloser, srkPW string, sealed []byte) ([]byte, error) {
	var srkAuth [20]byte

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}

	return tpm1.Unseal(rwc, sealed, srkAuth[:])
}

func seal20(rwc io.ReadWriteCloser, srkPW string, pcrs []int, data []byte) (*proto.SealedBytes, error) {
	key, err := loadSRK20(rwc, srkPW)
	if err != nil {
		return nil, err
	}
	sOpt := tpm2tools.SealCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrs,
		},
	}

	return key.Seal(data, sOpt)
}

func unseal20(rwc io.ReadWriteCloser, srkPW string, pcrs []int, sealed *proto.SealedBytes) ([]byte, error) {
	key, err := loadSRK20(rwc, srkPW)
	if err != nil {
		return nil, err
	}
	cOpt := tpm2tools.CertifyCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrs,
		},
	}

	return key.Unseal(sealed, cOpt)
}

func reseal20(rwc io.ReadWriteCloser, srkPW string, pcrs map[uint32][]byte, sealed *proto.SealedBytes) (*proto.SealedBytes, error) {
	key, err := loadSRK20(rwc, srkPW)
	if err != nil {
		return nil, err
	}
	keys := make([]int, 0, len(pcrs))
	for k := range pcrs {
		keys = append(keys, int(k))
	}
	sel := tpm2.PCRSelection{
		Hash: tpm2.AlgSHA256,
		PCRs: keys,
	}
	cOpt := tpm2tools.CertifyCurrent{
		PCRSelection: sel,
	}
	predictedPcrs := proto.Pcrs{
		Hash: proto.HashAlgo(sel.Hash),
		Pcrs: pcrs,
	}
	sOpt := tpm2tools.SealTarget{
		Pcrs: &predictedPcrs,
	}

	return key.Reseal(sealed, cOpt, sOpt)
}
