package boot

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

// Package is a netboot21 boot package.
//
// It contains an OSImage to boot as well as arbitrary metadata.
type Package struct {
	OSImage

	// Metadata is a map of relative archive paths -> arbitrary metadata
	// content.
	Metadata map[string]string
}

// NewPackage returns a new package based on the given OSImage.
func NewPackage(osi OSImage) *Package {
	return &Package{
		OSImage:  osi,
		Metadata: make(map[string]string),
	}
}

// AddMetadata adds metadata at a relative path.
func (p *Package) AddMetadata(relPath string, content string) {
	p.Metadata[relPath] = content
}

// Pack writes the boot package into archive w.
//
// TODO(hugelgupf): use a generic private key interface. No idea if we intend
// to keep using RSA here. Make usable with TPM.
func (p *Package) Pack(w cpio.RecordWriter, signer *rsa.PrivateKey) error {
	sw := NewSigningWriter(w)

	if len(p.Metadata) > 0 {
		if err := sw.WriteRecord(cpio.Directory("metadata", 0700)); err != nil {
			return err
		}

		for name, r := range p.Metadata {
			if err := sw.WriteRecord(cpio.StaticFile(path.Join("metadata", name), r, 0700)); err != nil {
				return err
			}
		}
	}

	if err := p.OSImage.Pack(sw); err != nil {
		return err
	}

	if signer != nil {
		return sw.WriteSignature(signer)
	}
	return nil
}

// Unpack unpacks a boot package in rr to p.
//
// TODO(hugelgupf): RSA? Generalize.
func (p *Package) Unpack(rr cpio.RecordReader, pk *rsa.PublicKey) error {
	*p = Package{
		Metadata: make(map[string]string),
	}

	recs := NewMeasuringReader(rr)
	a, err := cpio.ReadArchive(recs)
	if err != nil {
		return err
	}
	if pk != nil {
		if err := recs.Verify(pk); err != nil {
			return err
		}
	}

	for pth, content := range a.Files {
		s := strings.Split(pth, "/")
		if s[0] == "metadata" && content.Info.Mode&unix.S_IFMT == unix.S_IFREG {
			c, err := uio.ReadAll(content)
			if err != nil {
				return err
			}
			p.Metadata[path.Join(s[1:]...)] = string(c)
		}
	}

	typFile, ok := a.Files["package_type"]
	if !ok {
		return errors.New("file 'package_type' missing from boot package")
	}

	tb, err := uio.ReadAll(typFile)
	if err != nil {
		return err
	}

	pkgType := strings.TrimSpace(string(tb))
	imager, ok := osimageMap[pkgType]
	if !ok {
		return fmt.Errorf("invalid package type %q not supported", pkgType)
	}
	img, err := imager(a)
	if err != nil {
		return err
	}
	p.OSImage = img
	return nil
}
