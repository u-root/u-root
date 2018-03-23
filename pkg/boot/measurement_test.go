package boot

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
)

func containsFile(a *cpio.Archive, path string, content string) bool {
	if r, ok := a.Files[path]; ok {
		return cpio.ReaderAtEqual(r.ReaderAt, strings.NewReader(content))
	}
	return false
}

func TestSigningWriterWriteRecord(t *testing.T) {
	m := cpio.InMemArchive()
	s := NewSigningWriter(m)

	for _, tt := range []struct {
		r   cpio.Record
		err error
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
			r:   cpio.StaticFile("haha", "boo", 0777),
			err: fmt.Errorf("use SigningWriter.WriteFile for file named \"haha\""),
		},
	} {
		if err := s.WriteRecord(tt.r); err != tt.err && err.Error() != tt.err.Error() {
			t.Errorf("WriteRecord(%v) = %v, want %v", tt.r, err, tt.err)
		} else if err == nil && !m.Contains(tt.r) {
			t.Errorf("Archive should contain record %q, but doesn't", tt.r)
		} else if err != nil && m.Contains(tt.r) {
			t.Errorf("Archive should not contain record %q, but it does", tt.r)
		}
	}
}

func TestSigningWriterWriteFile(t *testing.T) {
	m := cpio.InMemArchive()
	s := NewSigningWriter(m)
	digest := &bytes.Buffer{}

	for _, tt := range []struct {
		path string
		r    string
		err  error
	}{
		{
			path: "signature",
			r:    "foobar",
			err:  fmt.Errorf("cannot write signature or signature_algo files"),
		},
		{
			path: "signature_algo",
			r:    "foobar",
			err:  fmt.Errorf("cannot write signature or signature_algo files"),
		},
		{
			path: "modules/foo/kernel",
			r:    "barfoo",
			err:  nil,
		},
	} {
		if err := s.WriteFile(tt.path, tt.r); err != tt.err && err.Error() != tt.err.Error() {
			t.Errorf("WriteFile(%v) = %v, want %v", tt.path, err, tt.err)
		} else if err == nil {
			if !containsFile(m, tt.path, tt.r) {
				t.Errorf("Archive should contain %q but doesn't", tt.path)
			}
			digest.WriteString(tt.path)
			digest.WriteString(tt.r)
		} else if err != nil && containsFile(m, tt.path, tt.r) {
			t.Errorf("Archive contains file %q but shouldn't", tt.path)
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
	}
	if err := cpio.WriteRecords(s, records); err != nil {
		t.Errorf("WriteRecords() = %v, want nil", err)
	}

	for path, content := range map[string]string{
		"modules/foo/kernel": "foobar",
		"metadata/hahaha":    "arrrrrgh",
	} {
		if err := s.WriteFile(path, content); err != nil {
			t.Errorf("WriteFile(%q, %s) = %v, want nil", path, content, err)
		}
		digest.WriteString(path)
		digest.WriteString(content)
	}

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
		cpio.StaticFile("metadata/hahaha", "arrrrrgh", 0700),
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
