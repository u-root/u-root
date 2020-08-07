// Copyright (c) 2018, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpm2

import (
	"crypto/sha256"
	"testing"

	. "github.com/google/go-tpm/tpm2"
)

func BenchmarkRSA2048Signing(b *testing.B) {
	b.StopTimer()
	rw := openTPM(b)
	defer rw.Close()

	pub := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}

	signerHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, pub)
	if err != nil {
		b.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, signerHandle)

	digest := sha256.Sum256([]byte("randomString"))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Sign(rw, signerHandle, defaultPassword, digest[:], nil, pub.RSAParameters.Sign); err != nil {
			b.Fatalf("Signing failed: %v", err)
		}
	}
}

func BenchmarkECCNISTP256Signing(b *testing.B) {
	b.StopTimer()
	rw := openTPM(b)
	defer rw.Close()
	skipOnUnsupportedAlg(b, rw, AlgECC)

	pub := Public{
		Type:       AlgECC,
		NameAlg:    AlgSHA256,
		Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
		ECCParameters: &ECCParams{
			Sign: &SigScheme{
				Alg:  AlgECDSA,
				Hash: AlgSHA256,
			},
			CurveID: CurveNISTP256,
		},
	}

	signerHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, pub)
	if err != nil {
		b.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, signerHandle)

	digest := sha256.Sum256([]byte("randomString"))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Sign(rw, signerHandle, defaultPassword, digest[:], nil, pub.ECCParameters.Sign); err != nil {
			b.Fatalf("Signing failed: %v", err)
		}
	}
}
