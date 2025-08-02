// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path"
	"testing"
)

func TestTar(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	filePath := path.Join(tmpDir, "file")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	content := "hello from tar"
	_, err = f.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	create, err := command(params{
		file:        "file.tar",
		create:      true,
		extract:     false,
		list:        false,
		noRecursion: false,
		verbose:     false,
	}, []string{"file"})
	if err != nil {
		t.Fatal(err)
	}

	err = create.run()
	if err != nil {
		t.Fatal(err)
	}

	archPath := path.Join(tmpDir, "file.tar")
	_, err = os.Stat(archPath)
	if err != nil {
		t.Fatal(err)
	}

	list, err := command(params{
		file:        "file.tar",
		create:      false,
		extract:     false,
		list:        true,
		noRecursion: false,
		verbose:     true,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = list.run()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	extract, err := command(params{
		file:        "file.tar",
		create:      false,
		extract:     true,
		list:        false,
		noRecursion: false,
		verbose:     false,
	}, []string{"."})
	if err != nil {
		t.Fatal(err)
	}

	err = extract.run()
	if err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != content {
		t.Errorf("expected %q, got %q", content, string(b))
	}
}

func TestCommandErrors(t *testing.T) {
	tests := []struct {
		err  error
		args []string
		p    params
	}{
		{
			err: errCreateAndExtract,
			p:   params{create: true, extract: true},
		},
		{
			err: errCreateAndList,
			p:   params{create: true, list: true},
		},
		{
			err: errExtractAndList,
			p:   params{extract: true, list: true},
		},
		{
			err:  errExtractArgsLen,
			p:    params{extract: true},
			args: []string{"1", "2"},
		},
		{
			err: errMissingMandatoryFlag,
		},
		{
			err:  errEmptyFile,
			p:    params{extract: true, file: ""},
			args: []string{"1"},
		},
	}

	for _, tt := range tests {
		_, err := command(tt.p, tt.args)
		if err != tt.err {
			t.Errorf("expected %v, got %v", tt.err, err)
		}
	}
}
