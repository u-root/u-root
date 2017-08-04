// Copyright (c) 2014, Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

// packedSize computes the size of a sequence of types that can be passed to
// binary.Read or binary.Write.
func packedSize(elts []interface{}) int {
	var size int
	for _, e := range elts {
		v := reflect.ValueOf(e)
		switch v.Kind() {
		case reflect.Ptr:
			s := packedSize([]interface{}{reflect.Indirect(v).Interface()})
			if s < 0 {
				return s
			}

			size += s
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				s := packedSize([]interface{}{v.Field(i).Interface()})
				if s < 0 {
					return s
				}

				size += s
			}
		case reflect.Slice:
			b, ok := e.([]byte)
			if !ok {
				return -1
			}

			size += 4 + len(b)
		default:
			s := binary.Size(e)
			if s < 0 {
				return s
			}

			size += s
		}
	}

	return size
}

// packWithHeader takes a header and a sequence of elements that are either of
// fixed length or slices of fixed-length types and packs them into a single
// byte array using binary.Write. It updates the CommandHeader to have the right
// length.
func packWithHeader(ch commandHeader, cmd []interface{}) ([]byte, error) {
	hdrSize := binary.Size(ch)
	bodySize := packedSize(cmd)
	if bodySize < 0 {
		return nil, errors.New("couldn't compute packed size for message body")
	}

	ch.Size = uint32(hdrSize + bodySize)

	in := []interface{}{ch}
	in = append(in, cmd...)
	return pack(in)
}

// pack encodes a set of elements into a single byte array, using
// encoding/binary. This means that all the elements must be encodeable
// according to the rules of encoding/binary. It has one difference from
// encoding/binary: it encodes byte slices with a prepended uint32 length, to
// match how the TPM encodes variable-length arrays.
func pack(elts []interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := packType(buf, elts); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// packType recursively packs types the same way that encoding/binary does under
// binary.BigEndian, but with one difference: it packs a byte slice as a uint32
// size followed by the bytes. The function unpackType performs the inverse
// operation of unpacking slices stored in this manner and using encoding/binary
// for everything else.
func packType(buf io.Writer, elts []interface{}) error {
	for _, e := range elts {
		v := reflect.ValueOf(e)
		switch v.Kind() {
		case reflect.Ptr:
			if err := packType(buf, []interface{}{reflect.Indirect(v).Interface()}); err != nil {
				return err
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				if err := packType(buf, []interface{}{v.Field(i).Interface()}); err != nil {
					return err
				}
			}
		case reflect.Slice:
			b, ok := e.([]byte)
			if !ok {
				return errors.New("can't pack slices of non-byte values")
			}

			if err := binary.Write(buf, binary.BigEndian, uint32(len(b))); err != nil {
				return err
			}

			if err := binary.Write(buf, binary.BigEndian, b); err != nil {
				return err
			}
		default:
			if err := binary.Write(buf, binary.BigEndian, e); err != nil {
				return err
			}
		}

	}

	return nil
}

func unpackKeyHandleList(b []byte) ([]Handle, error) {
	// TODO(kwalsh): handle pack/unpack of TPM_KEY_HANDLE_LIST more gracefully
	buf := bytes.NewBuffer(b)
	var n uint16
	if err := unpackType(buf, []interface{}{&n}); err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	h := make([]Handle, n)
	for i := range h {
		if err := unpackType(buf, []interface{}{&h[i]}); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// unpack performs the inverse operation from pack.
func unpack(b []byte, elts []interface{}) error {
	buf := bytes.NewBuffer(b)
	return unpackType(buf, elts)
}

// resizeBytes changes the size of the byte slice according to the second param.
func resizeBytes(b *[]byte, size uint32) {
	// Append to the slice if it's too small and shrink it if it's too large.
	l := len(*b)
	ss := int(size)
	if l > ss {
		*b = (*b)[:ss]
	} else if l < ss {
		*b = append(*b, make([]byte, ss-l)...)
	}
}

// unpackType recursively unpacks types from a reader just as encoding/binary
// does under binary.BigEndian, but with one difference: it unpacks a byte slice
// by first reading a uint32, then reading that many bytes. It assumes that
// incoming values are pointers to values so that, e.g., underlying slices can
// be resized as needed.
func unpackType(buf io.Reader, elts []interface{}) error {
	for _, e := range elts {
		v := reflect.ValueOf(e)
		k := v.Kind()
		if k != reflect.Ptr {
			return errors.New("all values passed to unpack must be pointers")
		}

		if v.IsNil() {
			return errors.New("can't fill a nil pointer")
		}

		iv := reflect.Indirect(v)
		switch iv.Kind() {
		case reflect.Struct:
			// Decompose the struct and copy over the values.
			for i := 0; i < iv.NumField(); i++ {
				if err := unpackType(buf, []interface{}{iv.Field(i).Addr().Interface()}); err != nil {
					return err
				}
			}
		case reflect.Slice:
			// Read a uint32 and resize the byte array as needed
			var size uint32
			if err := binary.Read(buf, binary.BigEndian, &size); err != nil {
				return err
			}

			// A zero size is used by the TPM to signal that certain elements
			// are not present.
			if size == 0 {
				continue
			}

			b, ok := e.(*[]byte)
			if !ok {
				return errors.New("can't fill pointers to slices of non-byte values")
			}

			resizeBytes(b, size)
			if err := binary.Read(buf, binary.BigEndian, e); err != nil {
				return err
			}
		default:
			if err := binary.Read(buf, binary.BigEndian, e); err != nil {
				return err
			}
		}

	}

	return nil
}
