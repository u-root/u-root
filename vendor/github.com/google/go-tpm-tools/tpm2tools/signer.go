package tpm2tools

import (
	"crypto"
	"crypto/rsa"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/google/go-tpm/tpm2"
)

// Global mutex to protect against concurrent TPM access.
var signerMutex sync.Mutex

type tpmSigner struct {
	Key  *Key
	Hash crypto.Hash
}

// Public returns the tpmSigners public key.
func (signer *tpmSigner) Public() crypto.PublicKey {
	return signer.Key.PublicKey()
}

// Sign uses the TPM key to sign the digest.
// The digest must be hashed from the same hash algorithm as the keys scheme.
// The opts hash function must also match the keys scheme (or be nil).
// Concurrent use of Sign is thread safe, but it is not safe to access the TPM
// from other sources while Sign is executing.
// For RSAPSS signatures, you cannot specify custom salt lengths. The salt
// length will be (keyBits/8) - digestSize - 2, unless that is less than the
// digestSize in which case, saltLen will be digestSize. The only normal case
// where saltLen is not digestSize is when using 1024 keyBits with SHA512.
func (signer *tpmSigner) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if pssOpts, ok := opts.(*rsa.PSSOptions); ok {
		if signer.Key.pubArea.RSAParameters.Sign.Alg != tpm2.AlgRSAPSS {
			return nil, fmt.Errorf("invalid options: PSSOptions cannot be used with signing alg: %v", signer.Key.pubArea.RSAParameters.Sign.Alg)
		}
		if pssOpts.SaltLength != rsa.PSSSaltLengthAuto {
			return nil, fmt.Errorf("salt length must be rsa.PSSSaltLengthAuto")
		}
	}
	if opts != nil && opts.HashFunc() != signer.Hash {
		return nil, fmt.Errorf("hash algorithm: got %v, want %v", opts.HashFunc(), signer.Hash)
	}
	if len(digest) != signer.Hash.Size() {
		return nil, fmt.Errorf("digest length: got %d, want %d", digest, signer.Hash.Size())
	}

	signerMutex.Lock()
	defer signerMutex.Unlock()

	auth, err := signer.Key.session.Auth()
	if err != nil {
		return nil, err
	}

	sig, err := tpm2.SignWithSession(signer.Key.rw, auth.Session, signer.Key.handle, "", digest, nil, nil)
	if err != nil {
		return nil, err
	}

	switch sig.Alg {
	case tpm2.AlgRSASSA:
		return sig.RSA.Signature, nil
	case tpm2.AlgRSAPSS:
		return sig.RSA.Signature, nil
	case tpm2.AlgECDSA:
		sigStruct := struct{ R, S *big.Int }{sig.ECC.R, sig.ECC.S}
		return asn1.Marshal(sigStruct)
	default:
		panic("unsupported signing algorithm")
	}
}

// GetSigner returns a crypto.Signer wrapping the loaded TPM Key.
// Concurrent use of one or more Signers is thread safe, but it is not safe to
// access the TPM from other sources while using a Signer.
// The returned Signer lasts the lifetime of the Key, and will no longer work
// once the Key has been closed.
func (k *Key) GetSigner() (crypto.Signer, error) {
	if k.hasAttribute(tpm2.FlagRestricted) {
		return nil, fmt.Errorf("restricted keys are not supported")
	}
	if !k.hasAttribute(tpm2.FlagSign) {
		return nil, fmt.Errorf("non-signing key used with GetSigner()")
	}

	var sigScheme *tpm2.SigScheme

	switch k.pubArea.Type {
	case tpm2.AlgRSA:
		sigScheme = k.pubArea.RSAParameters.Sign
		if sigScheme.Alg != tpm2.AlgRSAPSS && sigScheme.Alg != tpm2.AlgRSASSA {
			return nil, fmt.Errorf("unsupported signing algorithm: %v", sigScheme.Alg)
		}
	case tpm2.AlgECC:
		sigScheme = k.pubArea.ECCParameters.Sign
		if sigScheme.Alg != tpm2.AlgECDSA {
			return nil, fmt.Errorf("unsupported signing algorithm: %v", sigScheme.Alg)
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %v", k.pubArea.Type)
	}
	hash, err := sigScheme.Hash.Hash()
	if err != nil {
		return nil, err
	}
	return &tpmSigner{k, hash}, nil
}
