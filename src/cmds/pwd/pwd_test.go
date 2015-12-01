package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_pwd(t *testing.T) {
	//create fake directories with symlinks
	tempdir_path, _ := ioutil.TempDir("", "testdir")
	tempdir2_path, _ := ioutil.TempDir("", "testdir2")
	os.Chdir(tempdir_path)
	os.Symlink(tempdir2_path, "testlink")
	linkpath := tempdir_path + "/testlink"
	os.Chdir(linkpath)

	fmt.Println(wd)
	if wd != tempdir2_path {
		t.Fail()
	}

	os.Remove(tempdir_path)
	os.Remove(tempdir2_path)
}
