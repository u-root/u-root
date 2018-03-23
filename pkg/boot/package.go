package boot

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
)

// Package is a netboot21 boot package.
//
// It contains an OSImage to boot as well as arbitrary metadata.
type Package struct {
	OSImage

	Metadata map[string]string
}

func NewPackage(osi OSImage) *Package {
	return &Package{
		Metadata: make(map[string]string),
		OSImage:  osi,
	}
}

func (p *Package) AddMetadata(relPath string, content string) {
	p.Metadata[relPath] = content
}

// Pack writes the boot package into archive w.
//
// TODO(hugelgupf): use a generic private key interface.
func (p *Package) Pack(w cpio.RecordWriter, signer *rsa.PrivateKey) error {
	sw := NewSigningWriter(w)

	if len(p.Metadata) > 0 {
		if err := sw.WriteRecord(cpio.Directory("metadata", 0700)); err != nil {
			return err
		}

		for name, r := range p.Metadata {
			if err := sw.WriteFile(path.Join("metadata", name), r); err != nil {
				return err
			}
		}
	}

	if err := p.OSImage.Pack(sw); err != nil {
		return err
	}
	return sw.WriteSignature(signer)
}

// archive is an archive of files expected to be a netboot21 boot package.
type archive struct {
	originalArchive io.ReaderAt
	Files           map[string]cpio.Record
}

// TODO(hugelgupf): stop using x509.
func (p *Package) Unpack(arch io.ReaderAt, pk *rsa.PublicKey) error {
	p = &Package{
		Metadata: make(map[string]string),
	}
	a := &archive{
		originalArchive: arch,
		Files:           make(map[string]cpio.Record),
	}

	recs := NewMeasuringReader(cpio.Newc.Reader(arch))
	for {
		r, err := recs.ReadRecord()
		if err == io.EOF {
			break
		}

		dir, name := path.Split(r.Name)
		if dir == "metadata" {
			b, err := ioutil.ReadAll(uio.Reader(r))
			if err != nil {
				return err
			}

			p.Metadata[name] = string(b)
		} else {
			a.Files[r.Name] = r
		}
	}
	if pk != nil {
		if err := recs.Verify(pk); err != nil {
			return err
		}
	}

	typFile, ok := a.Files["package_type"]
	if !ok {
		return errors.New("file 'package_type' missing from boot package")
	}

	tb, err := ioutil.ReadAll(uio.Reader(typFile.ReaderAt))
	if err != nil {
		return err
	}

	pkgType := strings.TrimSpace(string(tb))
	switch pkgType {
	case "linux":
		img, err := newLinuxImageFromArchive(a)
		if err != nil {
			return err
		}
		p.OSImage = img

	case "multiboot":
		return fmt.Errorf("multiboot image support not yet implemented")

	default:
		return fmt.Errorf("invalid package type %q not supported", pkgType)
	}

	return nil
}
