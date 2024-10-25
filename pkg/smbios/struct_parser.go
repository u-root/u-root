// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type fieldParser interface {
	ParseField(t *Table, off int) (int, error)
}

var (
	fieldTagKey              = "smbios" // Tag key for annotations.
	fieldParserInterfaceType = reflect.TypeOf((*fieldParser)(nil)).Elem()
)

func parseStruct(t *Table, off int, complete bool, sp interface{}) (int, error) {
	var err error
	var ok bool
	var sv reflect.Value
	if sv, ok = sp.(reflect.Value); !ok {
		sv = reflect.Indirect(reflect.ValueOf(sp)) // must be a pointer to struct then, dereference it
	}
	svtn := sv.Type().Name()
	// fmt.Printf("t %s\n", svtn)
	i := 0
	for ; i < sv.NumField() && off < t.Len(); i++ {
		f := sv.Type().Field(i)
		fv := sv.Field(i)
		ft := fv.Type()
		tags := f.Tag.Get(fieldTagKey)
		// fmt.Printf("XX %02Xh f %s t %s k %s %s\n", off, f.Name, f.Type.Name(), fv.Kind(), tags)
		// Check tags first
		ignore := false
		for _, tag := range strings.Split(tags, ",") {
			tp := strings.Split(tag, "=")
			switch tp[0] {
			case "-":
				ignore = true
			case "skip":
				numBytes, _ := strconv.Atoi(tp[1])
				off += numBytes
			}
		}
		if ignore {
			continue
		}
		var verr error
		switch fv.Kind() {
		case reflect.Uint8:
			v, err := t.GetByteAt(off)
			fv.SetUint(uint64(v))
			verr = err
			off++
		case reflect.Uint16:
			v, err := t.GetWordAt(off)
			fv.SetUint(uint64(v))
			verr = err
			off += 2
		case reflect.Uint32:
			v, err := t.GetDWordAt(off)
			fv.SetUint(uint64(v))
			verr = err
			off += 4
		case reflect.Uint64:
			v, err := t.GetQWordAt(off)
			fv.SetUint(v)
			verr = err
			off += 8
		case reflect.String:
			v, _ := t.GetStringAt(off)
			fv.SetString(v)
			off++
		default:
			if reflect.PtrTo(ft).Implements(fieldParserInterfaceType) {
				off, err = fv.Addr().Interface().(fieldParser).ParseField(t, off)
				if err != nil {
					return off, err
				}
				break
			}
			// If it's a struct, just invoke parseStruct recursively.
			if fv.Kind() == reflect.Struct {
				off, err = parseStruct(t, off, true /* complete */, fv)
				if err != nil {
					return off, err
				}
				break
			}
			return off, fmt.Errorf("%s.%s: unsupported type %s", svtn, f.Name, fv.Kind())
		}
		if verr != nil {
			return off, fmt.Errorf("failed to parse %s.%s: %w", svtn, f.Name, verr)
		}
	}
	if complete && i < sv.NumField() {
		return off, fmt.Errorf("%s incomplete, got %d of %d fields", svtn, i, sv.NumField())
	}
	return off, nil
}
