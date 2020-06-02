// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boottest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/curl"
)

func module(r io.ReaderAt) map[string]interface{} {
	m := make(map[string]interface{})
	if f, ok := r.(curl.File); ok {
		m["url"] = f.URL().String()
	} else if f, ok := r.(fmt.Stringer); ok {
		m["stringer"] = f.String()
	}
	return m
}

// CompareImagesToJSON compares the names, cmdlines, and file URLs in imgs to
// the ones stored in jsonEncoded.
func CompareImagesToJSON(imgs []boot.OSImage, jsonEncoded []byte) error {
	var want interface{}
	if err := json.Unmarshal(jsonEncoded, &want); err != nil {
		return fmt.Errorf("failed to unmarshall test json %q: %v", jsonEncoded, err)
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
	return ioutil.WriteFile(filename, enc, 0644)
}

// ImagesToJSONLike spits out a json-convertible reproducible representation of
// the given boot images. This can be used in configuration parser tests (when
// the content of the images doesn't matter, but the file URLs, cmdlines,
// names, etc.)
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
	if li.Kernel != nil {
		m["kernel"] = module(li.Kernel)
	}
	if li.Initrd != nil {
		m["initrd"] = module(li.Initrd)
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
	if mi.Kernel != nil {
		m["kernel"] = module(mi.Kernel)
	}

	var modules []interface{}
	for _, mod := range mi.Modules {
		mmod := module(mod.Module)
		mmod["cmdline"] = mod.CmdLine
		mmod["name"] = mod.Name
		modules = append(modules, mmod)
	}
	m["modules"] = modules
	return m
}
