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

type makeit struct {
	name string      // name of the file.
	mode os.FileMode // mode of creation of a file by the OS.
}

var tests = []makeit{
	{
		"file1.txt",
		0777,
	},
	{
		"file2.txt",
		0777,
	},
	{
		"file3.txt",
		0777,
	},
}

func setup(data []byte) (string, error) {
	fmt.Println(":: Creating simulation data. ")
	dir, err := ioutil.TempDir(os.TempDir(), "cat.dir")
	if err != nil {
		return "", err
	}

	for i := range tests {
		if err := ioutil.WriteFile(path.Join(dir, tests[i].name), data, tests[i].mode); err != nil {
			return "", err
		}
	}

	return dir, nil
}

// Test_cat_1 test cat function against 3 files
func Test_cat_1(t *testing.T) {
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir, err := setup(someData)
	if err != nil {
		t.Fatal("Setup failed. Check errors.")
	}
	defer os.RemoveAll(dir)

	files := []string{path.Join(dir, "file1.txt"), path.Join(dir, "file2.txt"), path.Join(dir, "file3.txt")}

	var b bytes.Buffer
	if err := cat(b, files); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b.Bytes, someData) {
		t.Fatal("The output expected is not equal to the bytes provided in setup.")
	}
}
