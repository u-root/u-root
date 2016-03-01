package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type makeit struct {
	n string		   // name
	m os.FileMode		   // mode
	s string		   // for symlinks or content
}

var tests = []makeit{
	{
		n: "hi1.txt",
		m: 0666,
		s: "",
	},
	{
		n: "hi2.txt",
		m: 0777,
		s: "",
	},
}

func setup() (string, error) {
	d, err := ioutil.TempDir(os.TempDir(), "hi.dir")
	if err != nil {
		return "", err
	}

	tmpdir := path.Join(d, "hi.sub.dir")
	if err := os.Mkdir(tmpdir, 0777); err != nil {
		return "", err
	}

	for i:= range tests {
		if err := ioutil.WriteFile(path.Join(d, tests[i].n), []byte("hi"), tests[i].m); err != nil {
			return "", err
		}
	}
	
	return d, nil
}

func Test_mv_1(t *testing.T) {
	d, err := setup()
	if err != nil {
		t.Fatal("err")
	}
	defer os.RemoveAll(d)

	fmt.Println("Renaming file...")
	files1 := []string{path.Join(d, "hi1.txt"), path.Join(d, "hi4.txt")}
	if err := mv(files1, false); err != nil { 
		t.Error(err)
	}

	dsub := path.Join(d, "hi.sub.dir")

	fmt.Println("Moving files to directory...")
	files2 := []string{path.Join(d, "hi2.txt"), path.Join(d, "hi4.txt"), dsub}
	if err := mv(files2, true); err != nil {
		t.Error(err)
	}
}
