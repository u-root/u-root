package cmd

import (
	"bytes"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-tpm-tools/internal"
	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

func makeTempFile(tb testing.TB, content []byte) string {
	tb.Helper()
	file, err := ioutil.TempFile("", "gotpm_test_*.txt")
	if err != nil {
		tb.Fatal(err)
	}
	defer file.Close()
	if content != nil {
		if _, err := file.Write(content); err != nil {
			tb.Fatal(err)
		}
	}
	return file.Name()
}

func TestSealPlain(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)
	ExternalTPM = rwc

	tests := []struct {
		name        string
		algo        string
		sealPCRs    string
		certifyPCRs string
	}{
		{"RSASeal", "rsa", "", ""},
		{"ECCSeal", "ecc", "", ""},
		{"RSASealWithPCR", "rsa", "7", ""},
		{"ECCSealWithPCR", "ecc", "7", ""},
		{"RSACertifyWithPCR", "rsa", "", "7"},
		{"ECCCertifyWithPCR", "ecc", "", "7"},
		{"RSASealAndCertifyWithPCR", "rsa", "7,8", "1"},
		{"ECCSealAndCertifyWithPCR", "ecc", "7", "7,23"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			secretIn := []byte("Hello")
			secretFile1 := makeTempFile(t, secretIn)
			defer os.Remove(secretFile1)
			sealedFile := makeTempFile(t, nil)
			defer os.Remove(sealedFile)
			secretFile2 := makeTempFile(t, nil)
			defer os.Remove(secretFile2)

			sealArgs := []string{"seal", "--quiet", "--input", secretFile1, "--output", sealedFile}
			if test.sealPCRs != "" {
				sealArgs = append(sealArgs, "--pcrs", test.sealPCRs)
			}
			if test.algo != "" {
				sealArgs = append(sealArgs, "--algo", test.algo)
			}
			RootCmd.SetArgs(sealArgs)
			if err := RootCmd.Execute(); err != nil {
				t.Error(err)
			}
			pcrs = []int{} // "flush" pcrs value in last Execute() cmd

			unsealArgs := []string{"unseal", "--quiet", "--input", sealedFile, "--output", secretFile2}
			if test.certifyPCRs != "" {
				unsealArgs = append(unsealArgs, "--pcrs", test.certifyPCRs)
			}
			RootCmd.SetArgs(unsealArgs)
			if err := RootCmd.Execute(); err != nil {
				t.Error(err)
			}
			secretOut, err := ioutil.ReadFile(secretFile2)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(secretIn, secretOut) {
				t.Errorf("Expected %s, got %s", secretIn, secretOut)
			}
		})
	}
}

func TestUnsealFail(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)
	ExternalTPM = rwc
	extension := bytes.Repeat([]byte{0xAA}, sha256.Size)

	tests := []struct {
		name        string
		sealPCRs    string
		certifyPCRs string
		pcrToExtend []int
	}{
		// TODO(joerichey): Add test that TPM2_Reset make unsealing fail
		{"ExtendPCRAndUnseal", "23", "", []int{23}},
		{"ExtendPCRAndCertify", "23", "7", []int{7}},
		{"ExtendPCRAndCertify2", "", "5", []int{5}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			secretIn := []byte("Hello")
			secretFile := makeTempFile(t, secretIn)
			defer os.Remove(secretFile)
			sealedFile := makeTempFile(t, nil)
			defer os.Remove(sealedFile)

			sealArgs := []string{"seal", "--quiet", "--input", secretFile, "--output", sealedFile}
			if test.sealPCRs != "" {
				sealArgs = append(sealArgs, "--pcrs", test.sealPCRs)
			}
			RootCmd.SetArgs(sealArgs)
			if err := RootCmd.Execute(); err != nil {
				t.Error(err)
			}
			pcrs = []int{} // "flush" pcrs value in last Execute() cmd

			for _, pcr := range test.pcrToExtend {
				pcrHandle := tpmutil.Handle(pcr)
				if err := tpm2.PCRExtend(rwc, pcrHandle, tpm2.AlgSHA256, extension, ""); err != nil {
					t.Fatal(err)
				}
			}

			unsealArgs := []string{"unseal", "--quiet", "--input", sealedFile, "--output", secretFile}
			if test.certifyPCRs != "" {
				unsealArgs = append(unsealArgs, "--pcrs", test.certifyPCRs)
			}
			RootCmd.SetArgs(unsealArgs)
			if RootCmd.Execute() == nil {
				t.Error("Unsealing should have failed")
			}
		})
	}
}
