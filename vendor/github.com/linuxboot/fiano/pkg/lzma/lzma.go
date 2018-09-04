// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lzma implements reading and writing of LZMA compressed files.
//
// This package is specifically designed for the LZMA format used popular UEFI
// implementations. It requires the `lzma` and `unlzma` programs to be
// installed and on the path.
package lzma

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/ulikunitz/xz/lzma"
)

// Mapping from compression level to dictionary size.
var lzmaDictCapExps = []uint{18, 20, 21, 22, 22, 23, 23, 24, 25, 26}
var compressionLevel = 7

// Decode decodes a byte slice of LZMA data.
func Decode(encodedData []byte) ([]byte, error) {
	r, err := lzma.NewReader(bytes.NewBuffer(encodedData))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

// Encode encodes a byte slice with LZMA.
func Encode(decodedData []byte) ([]byte, error) {
	// These options are supported by the xz's LZMA command and EDK2's LZMA.
	// TODO: This does not support the f86 feature used in EDK2.
	wc := lzma.WriterConfig{
		SizeInHeader: true,
		Size:         int64(len(decodedData)),
		EOSMarker:    false,
		Properties:   &lzma.Properties{LC: 3, LP: 0, PB: 2},
		DictCap:      1 << lzmaDictCapExps[compressionLevel],
	}
	if err := wc.Verify(); err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	w, err := wc.NewWriter(buf)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, bytes.NewBuffer(decodedData)); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
