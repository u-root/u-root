package boot

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
)

type mockOSImage struct {
	packErr error
}

func newMockImage(*cpio.Archive) (OSImage, error) {
	return &mockOSImage{}, nil
}

func (mockOSImage) ExecutionInfo(log *log.Logger)  {}
func (mockOSImage) Execute() error                 { return nil }
func (m mockOSImage) Pack(sw *SigningWriter) error { return m.packErr }

func packageEqual(p1, p2 *Package) bool {
	li1, ok := p1.OSImage.(*LinuxImage)
	if !ok {
		return false
	}
	li2, ok := p2.OSImage.(*LinuxImage)
	if !ok {
		return false
	}
	return imageEqual(li1, li2) && reflect.DeepEqual(p1.Metadata, p2.Metadata)
}

func TestBootPackage(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Errorf("GenerateKey() = %v", err)
	}

	for _, tt := range []struct {
		pkg       *Package
		packErr   error
		unpackErr error
		signer    *rsa.PrivateKey
		verifier  *rsa.PublicKey
	}{
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  strings.NewReader("mcnulty"),
					Cmdline: "foo=bar",
				},
				Metadata: map[string]string{
					"stuff": "fooasdf",
				},
			},
			signer:   privateKey,
			verifier: &privateKey.PublicKey,
			packErr:  nil,
		},
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  strings.NewReader("mcnulty"),
					Cmdline: "foo=bar",
				},
				Metadata: map[string]string{
					"stuff": "fooasdf",
				},
			},
		},
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  strings.NewReader("mcnulty"),
					Cmdline: "foo=bar",
				},
				Metadata: map[string]string{
					"stuff": "fooasdf",
				},
			},
			verifier:  &privateKey.PublicKey,
			unpackErr: rsa.ErrVerification,
		},
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  nil,
					Cmdline: "foo=bar",
				},
				Metadata: map[string]string{},
			},
			signer:   privateKey,
			verifier: &privateKey.PublicKey,
			packErr:  nil,
		},
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  nil,
					Cmdline: "",
				},
				Metadata: map[string]string{},
			},
			signer:   privateKey,
			verifier: &privateKey.PublicKey,
			packErr:  nil,
		},
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  nil,
					Cmdline: "",
				},
				Metadata: map[string]string{
					"abc":      "def",
					"abcd/foo": "haha",
				},
			},
			signer:   privateKey,
			verifier: &privateKey.PublicKey,
			packErr:  nil,
		},
		{
			pkg: &Package{
				OSImage: &LinuxImage{
					Kernel:  strings.NewReader("lana"),
					Initrd:  nil,
					Cmdline: "",
				},
				Metadata: map[string]string{
					"abc":     "def",
					"abc/foo": "haha",
				},
			},
			signer:   privateKey,
			verifier: &privateKey.PublicKey,
			packErr:  nil,
		},
	} {
		a := cpio.InMemArchive()
		if err := tt.pkg.Pack(a, tt.signer); err != tt.packErr {
			t.Errorf("Pack(%v) = %v, want %v", tt.pkg, err, tt.packErr)
		}

		var p2 Package
		if err := (&p2).Unpack(a.Reader(), tt.verifier); err != tt.unpackErr {
			t.Errorf("Unpack() = %v, want %v", err, tt.unpackErr)
		} else if err == nil {
			if !packageEqual(tt.pkg, &p2) {
				t.Errorf("packages are not equal: got %v\nwant %v", p2, tt.pkg)
			}
		}
	}
}
