package tpm2tools

import (
	"io"
	"reflect"
	"testing"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"

	"github.com/google/go-tpm-tools/internal"
)

func TestNameMatchesPublicArea(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)
	ek, err := EndorsementKeyRSA(rwc)
	if err != nil {
		t.Fatal(err)
	}
	defer ek.Close()

	matches, err := ek.Name().MatchesPublic(ek.pubArea)
	if err != nil {
		t.Fatal(err)
	}
	if !matches {
		t.Fatal("Returned name and computed name do not match")
	}
}

func TestCreateSigningKeysInHierarchies(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)
	template := AIKTemplateRSA()

	// We are not authorized to create keys in the Platform Hierarchy
	for _, hierarchy := range []tpmutil.Handle{tpm2.HandleOwner, tpm2.HandleEndorsement, tpm2.HandleNull} {
		key, err := NewKey(rwc, hierarchy, template)
		if err != nil {
			t.Errorf("Hierarchy %+v: %s", hierarchy, err)
		} else {
			key.Close()
		}
	}
}

func TestCachedRSAKeys(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)
	tests := []struct {
		name   string
		getKey func(io.ReadWriter) (*Key, error)
	}{
		{"SRK", StorageRootKeyRSA},
		{"EK", EndorsementKeyRSA},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Get the key the first time and persist
			srk, err := test.getKey(rwc)
			if err != nil {
				t.Fatal(err)
			}
			defer srk.Close()

			pub := srk.PublicKey()
			if tpm2.FlushContext(rwc, srk.Handle()) == nil {
				t.Error("Trying to flush persistent keys should fail.")
			}

			// Get the cached key (should be the same)
			srk, err = test.getKey(rwc)
			if err != nil {
				t.Fatal(err)
			}
			defer srk.Close()

			if !reflect.DeepEqual(srk.PublicKey(), pub) {
				t.Errorf("Expected pub key: %v got: %v", pub, srk.PublicKey())
			}

			// We should still get the same key if we evict the handle
			if err := tpm2.EvictControl(rwc, "", tpm2.HandleOwner, srk.Handle(), srk.Handle()); err != nil {
				t.Errorf("Evicting control failed: %v", err)
			}
			srk, err = test.getKey(rwc)
			if err != nil {
				t.Fatal(err)
			}
			defer srk.Close()

			if !reflect.DeepEqual(srk.PublicKey(), pub) {
				t.Errorf("Expected pub key: %v got: %v", pub, srk.PublicKey())
			}
		})
	}
}

func TestKeyCreation(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	tests := []struct {
		name   string
		getKey func(io.ReadWriter) (*Key, error)
	}{
		{"SRK-ECC", StorageRootKeyECC},
		{"EK-ECC", EndorsementKeyECC},
		{"AIK-ECC", AttestationIdentityKeyECC},
		{"SRK-RSA", StorageRootKeyRSA},
		{"EK-RSA", EndorsementKeyRSA},
		{"AIK-RSA", AttestationIdentityKeyRSA},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := test.getKey(rwc)
			if err != nil {
				t.Fatal(err)
			}
			key.Close()
		})
	}
}

func BenchmarkKeyCreation(b *testing.B) {
	rwc := internal.GetTPM(b)
	defer CheckedClose(b, rwc)

	benchmarks := []struct {
		name   string
		getKey func(io.ReadWriter) (*Key, error)
	}{
		{"SRK-ECC-Cached", StorageRootKeyECC},
		{"EK-ECC-Cached", EndorsementKeyECC},
		{"AIK-ECC-Cached", AttestationIdentityKeyECC},

		{"SRK-ECC", func(rw io.ReadWriter) (*Key, error) {
			return NewKey(rw, tpm2.HandleOwner, SRKTemplateECC())
		}},
		{"EK-ECC", func(rw io.ReadWriter) (*Key, error) {
			return NewKey(rw, tpm2.HandleEndorsement, DefaultEKTemplateECC())
		}},
		{"AIK-ECC", func(rw io.ReadWriter) (*Key, error) {
			return NewKey(rw, tpm2.HandleOwner, AIKTemplateECC())
		}},

		{"SRK-RSA-Cached", StorageRootKeyRSA},
		{"EK-RSA-Cached", EndorsementKeyRSA},
		{"AIK-RSA-Cached", AttestationIdentityKeyRSA},

		{"SRK-RSA", func(rw io.ReadWriter) (*Key, error) {
			return NewKey(rw, tpm2.HandleEndorsement, SRKTemplateRSA())
		}},
		{"EK-RSA", func(rw io.ReadWriter) (*Key, error) {
			return NewKey(rw, tpm2.HandleOwner, DefaultEKTemplateRSA())
		}},
		{"AIK-RSA", func(rw io.ReadWriter) (*Key, error) {
			return NewKey(rw, tpm2.HandleOwner, AIKTemplateRSA())
		}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Don't count time to populate the cache
			b.StopTimer()
			key, err := bm.getKey(rwc)
			if err != nil {
				b.Fatal(err)
			}
			key.Close()
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				key, err := bm.getKey(rwc)
				if err != nil {
					b.Fatal(err)
				}
				key.Close()
			}
		})
	}
}
