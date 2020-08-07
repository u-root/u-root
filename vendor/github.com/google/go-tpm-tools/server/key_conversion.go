package server

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"

	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
)

var defaultNameAlg = tpm2tools.DefaultEKTemplateRSA().NameAlg

// CreateEKPublicAreaFromKey creates a public area from a go interface PublicKey.
// Supports RSA and ECC keys.
func CreateEKPublicAreaFromKey(k crypto.PublicKey) (tpm2.Public, error) {
	switch key := k.(type) {
	case *rsa.PublicKey:
		return createEKPublicRSA(key)
	case *ecdsa.PublicKey:
		return createEKPublicECC(key)
	default:
		return tpm2.Public{}, fmt.Errorf("unsupported public key type: %T", k)
	}
}

func createEKPublicRSA(rsaKey *rsa.PublicKey) (tpm2.Public, error) {
	public := tpm2tools.DefaultEKTemplateRSA()
	if rsaKey.N.BitLen() != int(public.RSAParameters.KeyBits) {
		return tpm2.Public{}, fmt.Errorf("unexpected RSA modulus size: %d bits", rsaKey.N.BitLen())
	}
	if rsaKey.E != int(public.RSAParameters.Exponent()) {
		return tpm2.Public{}, fmt.Errorf("unexpected RSA exponent: %d", rsaKey.E)
	}
	public.RSAParameters.ModulusRaw = rsaKey.N.Bytes()
	return public, nil
}

func createEKPublicECC(eccKey *ecdsa.PublicKey) (public tpm2.Public, err error) {
	public = tpm2tools.DefaultEKTemplateECC()
	public.ECCParameters.Point = tpm2.ECPoint{
		XRaw: eccIntToBytes(eccKey.Curve, eccKey.X),
		YRaw: eccIntToBytes(eccKey.Curve, eccKey.Y),
	}
	public.ECCParameters.CurveID, err = goCurveToCurveID(eccKey.Curve)
	return public, err
}

func createPublic(private tpm2.Private) tpm2.Public {
	publicHash := getHash(defaultNameAlg)
	publicHash.Write(private.SeedValue)
	publicHash.Write(private.Sensitive)
	return tpm2.Public{
		Type:    tpm2.AlgKeyedHash,
		NameAlg: defaultNameAlg,
		KeyedHashParameters: &tpm2.KeyedHashParams{
			Alg:    tpm2.AlgNull,
			Unique: publicHash.Sum(nil),
		},
	}
}

func createPrivate(sensitive []byte) tpm2.Private {
	private := tpm2.Private{
		Type:      tpm2.AlgKeyedHash,
		AuthValue: nil,
		SeedValue: make([]byte, getHash(defaultNameAlg).Size()),
		Sensitive: sensitive,
	}
	if _, err := io.ReadFull(rand.Reader, private.SeedValue); err != nil {
		panic(err)
	}
	return private
}

func createPublicPrivateSign(signingKey crypto.PrivateKey) (tpm2.Public, tpm2.Private, error) {
	rsaPriv, ok := signingKey.(*rsa.PrivateKey)
	if !ok {
		return tpm2.Public{}, tpm2.Private{}, fmt.Errorf("unsupported signing key type: %T", signingKey)
	}

	rsaPub := rsaPriv.PublicKey
	public := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    defaultNameAlg,
		Attributes: tpm2.FlagSign,
		RSAParameters: &tpm2.RSAParams{
			KeyBits:     uint16(rsaPub.N.BitLen()),
			ExponentRaw: uint32(rsaPub.E),
			ModulusRaw:  rsaPub.N.Bytes(),
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA256,
			},
		},
	}
	private := tpm2.Private{
		Type:      tpm2.AlgRSA,
		AuthValue: nil,
		SeedValue: nil, // Only Storage Keys need a seed value. See part 3 TPM2_CREATE b.3.
		Sensitive: rsaPriv.Primes[0].Bytes(),
	}

	return public, private, nil
}
