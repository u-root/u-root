// Package monoimporter provides a monorepo-compatible types.Importer for Go
// packages.
package monoimporter

import (
	"archive/zip"
	"fmt"
	"go/build"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/gcexportdata"
)

type finder interface {
	findAndOpen(pkg string) io.ReadCloser
}

func find(finders []finder, pkg string) io.ReadCloser {
	for _, f := range finders {
		if file := f.findAndOpen(pkg); file != nil {
			return file
		}
	}
	return nil
}

type zipReader struct {
	ctxt   build.Context
	stdlib *zip.Reader
	files  map[string]*zip.File
}

func newZipReader(stdlib *zip.Reader, ctxt build.Context) *zipReader {
	z := &zipReader{
		stdlib: stdlib,
		files:  make(map[string]*zip.File),
		ctxt:   ctxt,
	}
	for _, file := range z.stdlib.File {
		z.files[file.Name] = file
	}
	return z
}

// thatOneString is the Go build context directory name used by
// blaze/bazel/buck/Go.
//
// GOOS_GOARCH[_InstallSuffix], e.g. linux_amd64 or linux_amd64_pure.
func thatOneString(ctxt build.Context) string {
	var suffix string
	if len(ctxt.InstallSuffix) > 0 {
		suffix = fmt.Sprintf("_%s", ctxt.InstallSuffix)
	}
	return fmt.Sprintf("%s_%s%s", ctxt.GOOS, ctxt.GOARCH, suffix)
}

func (z *zipReader) findAndOpen(pkg string) io.ReadCloser {
	name := fmt.Sprintf("%s/%s.a", thatOneString(z.ctxt), pkg)
	f, ok := z.files[name]
	if !ok {
		return nil
	}
	rc, err := f.Open()
	if err != nil {
		return nil
	}
	return rc
}

type archives []string

func (a archives) findAndOpen(pkg string) io.ReadCloser {
	archive := fmt.Sprintf("/%s.a", pkg)
	for _, s := range a {
		if strings.HasSuffix(s, archive) {
			ar, err := os.Open(s)
			if err != nil {
				return nil
			}
			return ar
		}
	}
	return nil
}

// Importer implements a go/types.Importer for bazel-like monorepo build
// systems for Go packages.
//
// While open source Go depends on GOPATH and GOROOT to find packages,
// bazel-like build systems such as blaze or buck rely on a monorepo-style
// package search instead of using GOPATH and standard library packages are
// found in a zipped archive instead of GOROOT.
type Importer struct {
	fset *token.FileSet

	// imports is a cache of imported packages.
	imports map[string]*types.Package

	// archives is a list of paths to compiled Go package archives.
	archives archives

	// stdlib is an archive reader for standard library package object
	// files.
	stdlib *zipReader
}

// NewFromZips returns a new monorepo importer, using the build context to pick
// the desired standard library zip archive.
//
// zips refers to zip file paths with Go standard library object files.
//
// archives refers to directories in which to find compiled Go package object files.
func NewFromZips(ctxt build.Context, archives []string, zips []string) (*Importer, error) {
	var stdlib *zip.Reader
	want := fmt.Sprintf("%s.a.zip", thatOneString(ctxt))
	for _, dir := range zips {
		if filepath.Base(dir) == want {
			stdlibZ, err := zip.OpenReader(dir)
			if err != nil {
				return nil, err
			}
			stdlib = &stdlibZ.Reader
			break
		}
	}
	return New(ctxt, archives, stdlib), nil
}

// New returns a new monorepo importer.
func New(ctxt build.Context, archs []string, stdlib *zip.Reader) *Importer {
	i := &Importer{
		imports: map[string]*types.Package{
			"unsafe": types.Unsafe,
		},
		fset:     token.NewFileSet(),
		archives: archives(archs),
	}
	if stdlib != nil {
		i.stdlib = newZipReader(stdlib, ctxt)
	}
	return i
}

// Import implements types.Importer.Import.
func (i *Importer) Import(importPath string) (*types.Package, error) {
	if pkg, ok := i.imports[importPath]; ok && pkg.Complete() {
		return pkg, nil
	}

	pkg := strings.TrimPrefix(importPath, "google3/")
	finders := []finder{i.archives}
	if i.stdlib != nil {
		finders = append(finders, i.stdlib)
	}
	file := find(finders, pkg)
	if file == nil {
		return nil, fmt.Errorf("package %q not found", importPath)
	}
	defer file.Close()

	r, err := gcexportdata.NewReader(file)
	if err != nil {
		return nil, err
	}
	return gcexportdata.Read(r, i.fset, i.imports, importPath)
}
