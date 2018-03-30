package boot

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
)

func TestSigningWriterWriteFile(t *testing.T) {
	m := cpio.InMemArchive()
	s := NewSigningWriter(m)
	digest := &bytes.Buffer{}

	for _, tt := range []struct {
		r       cpio.Record
		err     error
		measure bool
	}{
		{
			r:   cpio.Directory("foobar", 0777),
			err: nil,
		},
		{
			r:   cpio.Directory("signature", 0777),
			err: fmt.Errorf("cannot write signature or signature_algo files"),
		},
		{
			r:   cpio.Directory("signature_algo", 0777),
			err: fmt.Errorf("cannot write signature or signature_algo files"),
		},
		{
			r:   cpio.StaticFile("signature", "foobar", 0700),
			err: fmt.Errorf("cannot write signature or signature_algo files"),
		},
		{
			r:   cpio.StaticFile("signature_algo", "foobar", 0700),
			err: fmt.Errorf("cannot write signature or signature_algo files"),
		},
		{
			r:       cpio.StaticFile("modules/foo/kernel", "barfoo", 0700),
			err:     nil,
			measure: true,
		},
	} {
		if err := s.WriteRecord(tt.r); err != tt.err && err.Error() != tt.err.Error() {
			t.Errorf("WriteFile(%v) = %v, want %v", tt.r.Name, err, tt.err)
		} else if err == nil {
			if !m.Contains(tt.r) {
				t.Errorf("Archive should contain %q but doesn't", tt.r.Name)
			}
			if tt.measure {
				digest.WriteString(tt.r.Name)
				digest.ReadFrom(uio.Reader(tt.r))
			}
		} else if err != nil && m.Contains(tt.r) {
			t.Errorf("Archive contains file %q but shouldn't", tt.r.Name)
		}
	}

	if len(digest.Bytes()) == 0 {
		t.Errorf("digest should contain something")
	}
	if sha1.Sum(digest.Bytes()) != s.SHA1Sum() {
		t.Errorf("sha1 differs")
	}
}

func TestWriterAndReader(t *testing.T) {
	m := cpio.InMemArchive()
	s := NewSigningWriter(m)
	digest := &bytes.Buffer{}

	records := []cpio.Record{
		cpio.Directory("modules", 0700),
		cpio.Directory("modules/foo", 0700),
		cpio.Directory("metadata", 0700),
		cpio.StaticFile("modules/foo/kernel", "foobar", 0700),
		cpio.StaticFile("metadata/hahaha", "arrgh", 0700),
	}
	if err := cpio.WriteRecords(s, records); err != nil {
		t.Errorf("WriteRecords() = %v, want nil", err)
	}
	digest.WriteString("modules/foo/kernel")
	digest.WriteString("foobar")
	digest.WriteString("metadata/hahaha")
	digest.WriteString("arrgh")

	if len(digest.Bytes()) == 0 {
		t.Errorf("digest should contain something")
	}
	if sha1.Sum(digest.Bytes()) != s.SHA1Sum() {
		t.Errorf("sha1 differs")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Errorf("rsa GenerateKey() = %v", err)
	}

	if err := s.WriteSignature(privateKey); err != nil {
		t.Errorf("WriteSignature() = %v, want nil", err)
	}
	if err := cpio.WriteTrailer(s); err != nil {
		t.Errorf("WriteTrailer() = %v, want nil", err)
	}

	want := []cpio.Record{
		cpio.Directory("modules", 0700),
		cpio.Directory("modules/foo", 0700),
		cpio.Directory("metadata", 0700),
		cpio.StaticFile("modules/foo/kernel", "foobar", 0700),
		cpio.StaticFile("metadata/hahaha", "arrgh", 0700),
	}

	r := NewMeasuringReader(m.Reader())
	got, err := cpio.ReadAllRecords(r)
	if err != nil {
		t.Errorf("ReadAllRecords() = %v, want nil", err)
	}
	if !cpio.AllEqual(got, want) {
		t.Errorf("ReadAllRecords() = \n%v, want \n%v", got, want)
	}

	if err := r.Verify(&privateKey.PublicKey); err != nil {
		t.Errorf("Verify() = %v, want nil", err)
	}
}
