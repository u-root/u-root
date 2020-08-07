// Package internal provides helper methods for testing. It should never be
// included in non-test libraries/binaries.
package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"io"
	"testing"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

// LoadRandomExternalKey loads a randomly generated external key into the
// TPM simulator and returns its' handle. If any errors occur, calls Fatal()
// on the passed testing.TB.
func LoadRandomExternalKey(tb testing.TB, rw io.ReadWriter) tpmutil.Handle {
	tb.Helper()
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		tb.Fatal(err)
	}
	public := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA1,
		Attributes: tpm2.FlagSign | tpm2.FlagUserWithAuth,
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA1,
			},
			KeyBits:     2048,
			ExponentRaw: uint32(pk.PublicKey.E),
			ModulusRaw:  pk.PublicKey.N.Bytes(),
		},
	}
	private := tpm2.Private{
		Type:      tpm2.AlgRSA,
		Sensitive: pk.Primes[0].Bytes(),
	}
	handle, _, err := tpm2.LoadExternal(rw, public, private, tpm2.HandleNull)
	if err != nil {
		tb.Error(err)
	}
	return handle
}
