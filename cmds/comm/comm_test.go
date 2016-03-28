package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestComm(t *testing.T) {

	txt1 := []string{"line1", "line2", "Line3", "lIne4", "line"}
	txt2 := []string{"line1", "Line2", "line3", "lInes", "LINEZ"}

	want1 := []string{
		"\t\tline1",
		"\tLine2",
		"line2",
		"Line3",
		"lIne4",
		"line",
		"\tline3",
		"\tlInes",
		"\tLINEZ\n",
	}

	want2 := []string{
		"\t\tline1",
		"\t\tline2",
		"\t\tLine3",
		"lIne4",
		"\tlInes",
		"\tLINEZ",
		"line\n",
	}

	f1, err := ioutil.TempFile(os.TempDir(), "txt1")
	if err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	defer os.Remove(f1.Name())

	f2, err := ioutil.TempFile(os.TempDir(), "txt2")
	if err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	defer os.Remove(f2.Name())

	for i, _ := range txt1 {
		if _, err := f1.Write([]byte(txt1[i] + "\n")); err != nil {
			t.Fatalf("Can't write to file1: %v")
		}
		if _, err := f2.Write([]byte(txt2[i] + "\n")); err != nil {
			t.Fatalf("Can't write to file2: %v")
		}
	}

	t.Logf("Testing case sensitive")
	cmd := exec.Command("go", "run", "comm.go", f1.Name(), f2.Name())
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("can't get output of comm: %v", err)
	} else if string(output) != strings.Join(want1, "\n") {
		t.Fatalf("Fail: want\n %s\n got\n %s", strings.Join(want1, "\n"), output)
	}

	t.Logf("Testing case insensitive")
	cmd = exec.Command("go", "run", "comm.go", "-i", f1.Name(), f2.Name())
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("can't get output of comm: %v", err)
	} else if string(output) != strings.Join(want2, "\n") {
		t.Fatalf("Fail: want\n %s\n got\n %s", strings.Join(want2, "\n"), output)
	}
}
