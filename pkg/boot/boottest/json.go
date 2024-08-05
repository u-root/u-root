// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package boottest contains methods for comparing boot.OSImages to each other
// and to JSON representations of themselves for use in tests.
//
// The JSON representation for boot.OSImages is special because the built-in
// json.Marshal function cannot marshal interfaces such as io.ReaderAt nicely,
// especially when the underlying members in structs used (such as *os.File or
// curl.lazyFile) are not exported.
//
// They are not json.Marshalers as part of boot.OSImage itself because they're
// not a fully accurate representation of an OSImage, not including file
// contents and depending for example on the current working directory of the
// calling process.
package boottest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/curl"
)

// makeFileRel converts file:/// schemes to a path relative to the current
// working directory in order to make a hermetic comparison for the test. Note
// that a "relative file scheme" does not technically exist since all file
// schemes are absolute.
//
// For example, if the URL is
// "file:///home/ryan/go/src/github.com/u-root/u-root/pkg/boot/grub/testdata_new/a/b/c",
// and the working directory is in grub, the URL is converted to
// "file:///testdata_new/a/b/c".
// It assumes the testdata is in a sub-directory of the test.
func makeFileRel(u *url.URL) (*url.URL, error) {
	if u.Scheme != "file" {
		return u, nil
	}
	relU, err := url.Parse(u.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing %v: %w", u.String(), err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error from os.GetWd(): %w", err)
	}
	relU.Path = strings.TrimPrefix(relU.Path, wd)
	return relU, nil
}

func module(r io.ReaderAt) map[string]interface{} {
	m := make(map[string]interface{})
	if f, ok := r.(curl.File); ok {
		url, err := makeFileRel(f.URL())
		if err != nil {
			// Including the error message is sufficient to cause
			// the test to fail and simplifies error handling.
			m["error"] = err
		}
		m["url"] = url.String()
	} else if f, ok := r.(fmt.Stringer); ok {
		m["stringer"] = f.String()
	} else if f, ok := r.(*os.File); ok {
		m["name"] = f.Name()
	}
	return m
}

// CompareImagesToJSON compares the names, cmdlines, and file URLs in imgs to
// the ones stored in jsonEncoded.
//
// You can obtain such a JSON encoding with ToJSONFile.
func CompareImagesToJSON(imgs []boot.OSImage, jsonEncoded []byte) error {
	var want interface{}
	if err := json.Unmarshal(jsonEncoded, &want); err != nil {
		return fmt.Errorf("failed to unmarshall test json %q: %w", jsonEncoded, err)
	}

	got := ImagesToJSONLike(imgs)
	if !cmp.Equal(want, got) {
		return fmt.Errorf("mismatch(-want, +got):\n%s", cmp.Diff(want, got))
	}
	return nil
}

// ToJSONFile can be used to generate JSON-comparable files for use with
// CompareImagesToJSON in tests.
func ToJSONFile(imgs []boot.OSImage, filename string) error {
	enc, err := json.MarshalIndent(ImagesToJSONLike(imgs), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, enc, 0o644)
}

// ImagesToJSONLike spits out a json-convertible reproducible representation of
// the given boot images. This can be used in configuration parser tests (when
// the content of the images doesn't matter, but the file URLs, cmdlines,
// names, etc.)
//
// The JSON representation for boot.OSImages is special because the built-in
// json.Marshal function cannot marshal interfaces such as io.ReaderAt nicely,
// especially when the underlying structs used are not exported.
func ImagesToJSONLike(imgs []boot.OSImage) []interface{} {
	var infs []interface{}
	for _, img := range imgs {
		if l, ok := img.(*boot.LinuxImage); ok {
			infs = append(infs, LinuxImageToJSON(l))
		}
		if m, ok := img.(*boot.MultibootImage); ok {
			infs = append(infs, MultibootImageToJSON(m))
		}
	}
	return infs
}

// LinuxImageToJSON is implemented only in order to compare LinuxImages in
// tests.
//
// It should be json-encodable and decodable.
func LinuxImageToJSON(li *boot.LinuxImage) map[string]interface{} {
	m := make(map[string]interface{})
	m["image_type"] = "linux"
	m["name"] = li.Name
	m["cmdline"] = li.Cmdline
	m["rank"] = strconv.Itoa(li.BootRank)
	if li.Kernel != nil {
		m["kernel"] = module(li.Kernel)
	}
	if li.Initrd != nil {
		m["initrd"] = module(li.Initrd)
	}
	if li.DTB != nil {
		m["dtb"] = module(li.DTB)
	}
	return m
}

// MultibootImageToJSON is implemented only in order to compare MultibootImages
// in tests.
//
// It should be json-encodable and decodable.
func MultibootImageToJSON(mi *boot.MultibootImage) map[string]interface{} {
	m := make(map[string]interface{})
	m["image_type"] = "multiboot"
	m["name"] = mi.Name
	m["cmdline"] = mi.Cmdline
	m["rank"] = strconv.Itoa(mi.BootRank)
	if mi.Kernel != nil {
		m["kernel"] = module(mi.Kernel)
	}

	var modules []interface{}
	for _, mod := range mi.Modules {
		mmod := module(mod.Module)
		mmod["cmdline"] = mod.Cmdline
		mmod["name"] = mod.Name()
		modules = append(modules, mmod)
	}
	m["modules"] = modules
	return m
}
