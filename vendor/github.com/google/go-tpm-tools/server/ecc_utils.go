package server

import (
	"crypto/elliptic"
	"fmt"
	"math/big"

	"github.com/google/go-tpm/tpm2"
)

// ECC coordinates need to maintain a specific size based on the curve, so we pad the front with zeros.
// This is particularly an issue for NIST-P521 coordinates, as they are frequently missing their first byte.
func eccIntToBytes(curve elliptic.Curve, i *big.Int) []byte {
	bytes := i.Bytes()
	curveBytes := (curve.Params().BitSize + 7) / 8
	return append(make([]byte, curveBytes-len(bytes)), bytes...)
}

func curveIDToGoCurve(curve tpm2.EllipticCurve) (elliptic.Curve, error) {
	switch curve {
	case tpm2.CurveNISTP224:
		return elliptic.P224(), nil
	case tpm2.CurveNISTP256:
		return elliptic.P256(), nil
	case tpm2.CurveNISTP384:
		return elliptic.P384(), nil
	case tpm2.CurveNISTP521:
		return elliptic.P521(), nil
	default:
		return nil, fmt.Errorf("unsupported TPM2 curve: %v", curve)
	}
}

func goCurveToCurveID(curve elliptic.Curve) (tpm2.EllipticCurve, error) {
	switch curve.Params().Name {
	case elliptic.P224().Params().Name:
		return tpm2.CurveNISTP224, nil
	case elliptic.P256().Params().Name:
		return tpm2.CurveNISTP256, nil
	case elliptic.P384().Params().Name:
		return tpm2.CurveNISTP384, nil
	case elliptic.P521().Params().Name:
		return tpm2.CurveNISTP521, nil
	default:
		return 0, fmt.Errorf("unsupported Go curve: %v", curve.Params().Name)
	}
}
