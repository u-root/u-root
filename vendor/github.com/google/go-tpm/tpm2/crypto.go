package tpm2

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"math/big"
)

// RSAPub converts a TPM RSA public key into one recognized by the rsa package.
func RSAPub(parms *TPMSRSAParms, pub *TPM2BPublicKeyRSA) (*rsa.PublicKey, error) {
	result := rsa.PublicKey{
		N: big.NewInt(0).SetBytes(pub.Buffer),
		E: int(parms.Exponent),
	}
	// TPM considers 65537 to be the default RSA public exponent, and 0 in
	// the parms
	// indicates so.
	if result.E == 0 {
		result.E = 65537
	}
	return &result, nil
}

// ECDHPubKey converts a TPM ECC public key into one recognized by the ecdh package
func ECDHPubKey(curve ecdh.Curve, pub *TPMSECCPoint) (*ecdh.PublicKey, error) {

	var c elliptic.Curve
	switch curve {
	case ecdh.P256():
		c = elliptic.P256()
	case ecdh.P384():
		c = elliptic.P384()
	case ecdh.P521():
		c = elliptic.P521()
	default:
		return nil, fmt.Errorf("unknown curve: %v", curve)
	}

	pubKey := ecdsa.PublicKey{
		Curve: c,
		X:     big.NewInt(0).SetBytes(pub.X.Buffer),
		Y:     big.NewInt(0).SetBytes(pub.Y.Buffer),
	}

	return pubKey.ECDH()
}

// ECCPoint returns an uncompressed ECC Point
func ECCPoint(pubKey *ecdh.PublicKey) (*big.Int, *big.Int, error) {
	b := pubKey.Bytes()
	size, err := elementLength(pubKey.Curve())
	if err != nil {
		return nil, nil, fmt.Errorf("ECCPoint: %w", err)
	}
	return big.NewInt(0).SetBytes(b[1 : size+1]),
		big.NewInt(0).SetBytes(b[size+1:]), nil
}

func elementLength(c ecdh.Curve) (int, error) {
	switch c {
	case ecdh.P256():
		// crypto/internal/nistec/fiat.p256ElementLen
		return 32, nil
	case ecdh.P384():
		// crypto/internal/nistec/fiat.p384ElementLen
		return 48, nil
	case ecdh.P521():
		// crypto/internal/nistec/fiat.p521ElementLen
		return 66, nil
	default:
		return 0, fmt.Errorf("unknown element length for curve: %v", c)
	}
}
