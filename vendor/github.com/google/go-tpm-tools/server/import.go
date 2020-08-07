// Package server contains functions to be ran on a server (no TPM needed), as oppose to a client (with TPM).
package server

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash"
	"io"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"

	tpmpb "github.com/google/go-tpm-tools/proto"
	"github.com/google/go-tpm-tools/tpm2tools"
)

// CreateImportBlob uses the provided public EK to encrypt the sensitive data.
// The returned ImportBlob can then be decrypted and imported using the
// tpm2tools Key.Import() method. A non-nil pcrs parameter adds a requirement
// that the TPM must have specific PCR values for Import() to succeed.
func CreateImportBlob(ekPub crypto.PublicKey, sensitive []byte, pcrs *tpmpb.Pcrs) (*tpmpb.ImportBlob, error) {
	ek, err := CreateEKPublicAreaFromKey(ekPub)
	if err != nil {
		return nil, err
	}
	private := createPrivate(sensitive)
	public := createPublic(private)

	return createImportBlobHelper(ek, public, private, pcrs)
}

// CreateSigningKeyImportBlob uses the provided public EK to encrypt the signing
// key into import blob format. The returned import blob can be used to import
// the signing key into the TPM associated with the provided EK without exposing
// the private area to the TPM's OS using the tpm2tools Key.ImportSigningKey()
// method. A non-nil pcrs parameter adds a requirement that the TPM must have
// specific PCR values to use the signing key.
func CreateSigningKeyImportBlob(ekPub crypto.PublicKey, signingKey crypto.PrivateKey, pcrs *tpmpb.Pcrs) (*tpmpb.ImportBlob, error) {
	ek, err := CreateEKPublicAreaFromKey(ekPub)
	if err != nil {
		return nil, err
	}
	public, private, err := createPublicPrivateSign(signingKey)
	if err != nil {
		return nil, err
	}

	return createImportBlobHelper(ek, public, private, pcrs)
}

func createImportBlobHelper(ek, public tpm2.Public, private tpm2.Private, pcrs *tpmpb.Pcrs) (*tpmpb.ImportBlob, error) {
	setPublicAuth(&public, pcrs)

	var seed, encryptedSeed []byte
	var err error
	switch ek.Type {
	case tpm2.AlgRSA:
		seed, encryptedSeed, err = createRSASeed(ek)
		if err != nil {
			return nil, err
		}
	case tpm2.AlgECC:
		seed, encryptedSeed, err = createECCSeed(ek)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported EK type: %v", ek.Type)
	}
	duplicate, err := createDuplicate(private, seed, public, ek)
	if err != nil {
		return nil, err
	}
	pubEncoded, err := public.Encode()
	if err != nil {
		return nil, err
	}

	return &tpmpb.ImportBlob{
		Duplicate:     duplicate,
		EncryptedSeed: encryptedSeed,
		PublicArea:    pubEncoded,
		Pcrs:          pcrs,
	}, nil
}

func setPublicAuth(public *tpm2.Public, pcrs *tpmpb.Pcrs) {
	if len(pcrs.GetPcrs()) == 0 {
		// Allow password authorization so we can use a nil AuthPolicy.
		public.AuthPolicy = nil
		public.Attributes |= tpm2.FlagUserWithAuth
	} else {
		public.AuthPolicy = tpm2tools.ComputePCRSessionAuth(pcrs)
		public.Attributes |= tpm2.FlagAdminWithPolicy
	}
}

func createRSASeed(ek tpm2.Public) (seed, encryptedSeed []byte, err error) {
	seedSize := ek.RSAParameters.Symmetric.KeyBits / 8
	seed = make([]byte, seedSize)
	if _, err := io.ReadFull(rand.Reader, seed); err != nil {
		panic(err)
	}

	ekPub, err := ek.Key()
	if err != nil {
		return nil, nil, err
	}
	encryptedSeed, err = rsa.EncryptOAEP(
		getHash(ek.NameAlg),
		rand.Reader,
		ekPub.(*rsa.PublicKey),
		seed,
		[]byte("DUPLICATE\x00"))
	if err != nil {
		return nil, nil, err
	}
	encryptedSeed, err = tpmutil.Pack(encryptedSeed)
	return seed, encryptedSeed, err
}

func createECCSeed(ek tpm2.Public) (seed, encryptedSeed []byte, err error) {
	curve, err := curveIDToGoCurve(ek.ECCParameters.CurveID)
	if err != nil {
		return nil, nil, err
	}
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	ekPoint := ek.ECCParameters.Point
	z, _ := curve.ScalarMult(ekPoint.X(), ekPoint.Y(), priv)
	xBytes := eccIntToBytes(curve, x)

	seed, err = tpm2.KDFe(
		ek.NameAlg,
		eccIntToBytes(curve, z),
		"DUPLICATE",
		xBytes,
		eccIntToBytes(curve, ekPoint.X()),
		getHash(ek.NameAlg).Size()*8)
	if err != nil {
		return nil, nil, err
	}
	encryptedSeed, err = tpmutil.Pack(tpmutil.U16Bytes(xBytes), tpmutil.U16Bytes(eccIntToBytes(curve, y)))
	return seed, encryptedSeed, err
}

func createDuplicate(private tpm2.Private, seed []byte, public, ek tpm2.Public) ([]byte, error) {
	nameEncoded, err := getEncodedName(public)
	if err != nil {
		return nil, err
	}
	secret, err := private.Encode()
	if err != nil {
		return nil, err
	}
	packedSecret, err := tpmutil.Pack(tpmutil.U16Bytes(secret))
	if err != nil {
		return nil, err
	}
	encryptedSecret, err := encryptSecret(packedSecret, seed, nameEncoded, ek)
	if err != nil {
		return nil, err
	}
	macSum, err := createHMAC(encryptedSecret, nameEncoded, seed, ek.NameAlg)
	if err != nil {
		return nil, err
	}
	return tpmutil.Pack(tpm2.IDObject{
		IntegrityHMAC: macSum,
		EncIdentity:   encryptedSecret,
	})
}

func getEncodedName(public tpm2.Public) ([]byte, error) {
	name, err := public.Name()
	if err != nil {
		return nil, err
	}
	return name.Digest.Encode()
}

func encryptSecret(secret, seed, nameEncoded []byte, ek tpm2.Public) ([]byte, error) {
	var symSize int
	switch ek.Type {
	case tpm2.AlgRSA:
		symSize = int(ek.RSAParameters.Symmetric.KeyBits)
	case tpm2.AlgECC:
		symSize = int(ek.ECCParameters.Symmetric.KeyBits)
	default:
		return nil, fmt.Errorf("unsupported EK type: %v", ek.Type)
	}

	symmetricKey, err := tpm2.KDFa(
		ek.NameAlg,
		seed,
		"STORAGE",
		nameEncoded,
		/*contextV=*/ nil,
		symSize)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, err
	}
	encSecret := make([]byte, len(secret))
	// The TPM spec requires an all-zero IV.
	iv := make([]byte, len(symmetricKey))
	cipher.NewCFBEncrypter(c, iv).XORKeyStream(encSecret, secret)
	return encSecret, nil
}

func createHMAC(encryptedSecret, nameEncoded, seed []byte, hashAlg tpm2.Algorithm) ([]byte, error) {
	macKey, err := tpm2.KDFa(
		hashAlg,
		seed,
		"INTEGRITY",
		/*contextU=*/ nil,
		/*contextV=*/ nil,
		getHash(hashAlg).Size()*8)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(func() hash.Hash { return getHash(hashAlg) }, macKey)
	mac.Write(encryptedSecret)
	mac.Write(nameEncoded)
	return mac.Sum(nil), nil
}

func getHash(hashAlg tpm2.Algorithm) hash.Hash {
	create, err := hashAlg.Hash()
	if err != nil {
		panic(err)
	}
	return create.New()
}
