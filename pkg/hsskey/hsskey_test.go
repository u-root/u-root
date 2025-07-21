// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hsskey

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/ipmi/blobs"
)

type mockBlobReader struct {
	data []uint8
}

func (h *mockBlobReader) BlobOpen(id string, flags int16) (blobs.SessionID, error) {
	return 0, nil
}

func (h *mockBlobReader) BlobRead(sid blobs.SessionID, offset, size uint32) ([]uint8, error) {
	return h.data, nil
}

func (h *mockBlobReader) BlobClose(sid blobs.SessionID) error {
	return nil
}

func TestGenPassword(t *testing.T) {
	mockHss := []byte{
		0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55,
		0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55,
		0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55,
		0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55,
	}

	expected := []byte{
		175, 97, 131, 232, 89, 61, 152, 29, 245,
		45, 164, 141, 98, 78, 7, 243, 120, 96, 179, 166, 18,
		59, 22, 172, 16, 151, 191, 99, 141, 25, 35, 246,
	}

	key, err := GenPassword(mockHss, DefaultPasswordSalt, "a", "b")

	if err != nil || !bytes.Equal(key, expected) {
		t.Fatalf("GenPassword generated a wrong key, expected:\n%v\nbut returned:\n%v", expected, key)
	}
}

func TestReadHssBlobEmpty(t *testing.T) {
	h := mockBlobReader{[]uint8{}}
	_, err := readHssBlob("", &h)

	if err == nil {
		t.Fatalf("Expected invalid length failure")
	}
}

func TestReadHssBlob(t *testing.T) {
	data := [hostSecretSeedLen]uint8{}
	h := mockBlobReader{data: data[:]}
	_, err := readHssBlob("", &h)
	if err != nil {
		t.Fatalf("Expected success, got err: %v", err)
	}
}
