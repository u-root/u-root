/* Copyright 2012 the u-root Authors. All rights reserved
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * created by Rafael Campos Nunes <rafaelnunes@engineer.com>
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

// setup writes a set of files, putting 1 byte in each file.
func setup(data []byte) (string, error) {
	fmt.Println(":: Creating simulation data. ")
	dir, err := ioutil.TempDir(os.TempDir(), "cat.dir")
	if err != nil {
		return "", err
	}

	for i := range data {
		n := fmt.Sprintf("%v%d", path.Join(dir, "file"), i)
		if err := ioutil.WriteFile(n, []byte{data[i]}, 0666); err != nil {
			return "", err
		}
	}

	return dir, nil
}

// Test_cat_1 test cat function against 3 files
func Test_cat_1(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir, err := setup(someData)
	if err != nil {
		t.Fatal("setup has failed, %v", err)
	}
	defer os.RemoveAll(dir)

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", path.Join(dir, "file"), i))
	}

	var b bytes.Buffer
	if err := cat(&b, files); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(b.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", b.Bytes(), someData)
	}
}
