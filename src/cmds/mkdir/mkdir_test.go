//test3 might require root permission

package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func Test_mkdir(t *testing.T) {
	flag.Parse()
	os.Chdir(os.TempDir()) //hack to remove the test3 without deleting /tmp
	
	test1 := []string{os.TempDir() + "/dir1"}
	test2 := []string{os.TempDir() + "/dir2"}
	test3 := []string{os.TempDir() + "/long/path/to/dir3"}
	test4 := []string{os.TempDir() + "/dir4", os.TempDir() + "/dir5", os.TempDir() + "/dir6"}

	fmt.Printf("test1: defaults\n")
	mkdir(test1)
	defer os.Remove(test1[0])
	info, err := os.Stat(test1[0])
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Fatal("Directory not created?")
	}
	if info.Mode()&os.ModePerm != 0755 {
		t.Fatal("permission bits don't match")
	}
	fmt.Printf("test1 ok!\n")
	
	fmt.Println("test2: change permission")
	*mode = 0644
	mkdir(test2)
	defer os.Remove(test2[0])
	info, err = os.Stat(test2[0])
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Fatal("Directory not created?")
	}
	if info.Mode()&os.ModePerm != 0644 {
		t.Fatal("permission bits don't match")
	}
	fmt.Printf("test2 ok!\n")

	fmt.Println("test3: create entire path")
	*mkall = true
	mkdir(test3)
	defer os.RemoveAll("long")
	info, err = os.Stat(test3[0])
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Fatal("Directory not created?")
	}
	if info.Mode()&os.ModePerm != 0644 {
		t.Fatal("permission bits don't match")
	}
	fmt.Printf("test3 ok!\n")

	fmt.Println("test4: verbose mode")
	*verbose = true
	mkdir(test4)
	for i := 0; i < len(test4); i++ {
		defer os.Remove(test4[i])
	}
	for i := 0; i < len(test4); i++ {
		info, err = os.Stat(test4[i])
		if err != nil {
			t.Error(err)
		}
		if !info.IsDir() {
			t.Fatal("Directory not created?")
		}
		if info.Mode()&os.ModePerm != 0644 {
			t.Fatal("permission bits don't match")
		}
	}
	
	fmt.Printf("test4 ok!\n")
	
}
