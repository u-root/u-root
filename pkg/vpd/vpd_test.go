// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vpd

import (
	"bytes"
	"reflect"
	"testing"
)

func TestGetReadOnly(t *testing.T) {
	r := NewReader()
	r.VpdDir = "./tests"
	value, err := r.Get("key1", true)
	if err != nil || !bytes.Equal(value, []byte("value1\n")) {
		t.Errorf(`r.Get("key1", true) = %v, %v, want nil, %v`, value, err, []byte("value1\n"))
	}
}

func TestGetReadWrite(t *testing.T) {
	r := NewReader()
	r.VpdDir = "./tests"
	value, err := r.Get("mysecretpassword", false)
	if err != nil || !bytes.Equal(value, []byte("passw0rd\n")) {
		t.Errorf(`r.Get("mysecretpassword", false) = %v, %v, want nil, %v`, value, err, []byte("passwor0rd\n"))
	}
}

func TestGetReadBinary(t *testing.T) {
	r := NewReader()
	r.VpdDir = "./tests"
	value, err := r.Get("binary1", true)
	if err != nil || !bytes.Equal(value, []byte("some\x00binary\ndata")) {
		t.Errorf(`r.Get("binary1", true) = %v, %v, want nil, %v`, value, err, []byte("some\x00binary\ndata"))
	}
}

func TestGetAllReadOnly(t *testing.T) {
	r := NewReader()
	r.VpdDir = "./tests"
	expected := map[string][]byte{
		"binary1": []byte("some\x00binary\ndata"),
		"key1":    []byte("value1\n"),
	}
	vpdMap, err := r.GetAll(true)
	if err != nil || !reflect.DeepEqual(vpdMap, expected) {
		t.Errorf(`r.GetAll(true) = %v, %v, want nil, %v`, vpdMap, err, expected)
		t.FailNow()
	}
}
