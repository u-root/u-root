package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func init() {
	// test in verbose mode so that it covers more code
	*verbose = true
}

func TestRunNoConfigFileArg(t *testing.T) {
	*configFile = ""
	if run() != 1 {
		t.Errorf("Expected run() to return 1 on no configFile specified")
	}
}

func TestRunConfigFileDoesntExist(t *testing.T) {
	*configFile = "/non_existent_file/asdasdf91238109234"
	if run() != 1 {
		t.Errorf("Expected run() to return 1 on non-existent configFile")
	}
}

func TestRunConfigFileInvalidJSON(t *testing.T) {
	*configFile = tempFile("[").Name()
	defer os.Remove(*configFile)
	if run() != 1 {
		t.Errorf("Expected run() to return 1 on invalid JSON")
	}
}

func TestRunConfigFileEmpty(t *testing.T) {
	originalStdout := os.Stdout
	os.Stdout = tempFile("")
	defer func() {
		os.Stdout = originalStdout
	}()

	expectedOut := "Checker Results: []\n"

	*configFile = tempFile("[]").Name()
	defer os.Remove(*configFile)
	if run() != 0 {
		t.Errorf("Expected run() to return 0 on empty checklist")
	}

	os.Stdout.Seek(0, 0)
	out, _ := ioutil.ReadAll(os.Stdout)
	if string(out) != expectedOut {
		t.Errorf("Expected run() to write %#v, not %#v", expectedOut, string(out))
	}
}

func tempFile(contents string) *os.File {
	file, err := ioutil.TempFile("", "configFile")
	if err != nil {
		log.Fatalf("Could not create temporary config file: %v", err)
	}
	_, err = file.WriteString(contents)
	if err != nil {
		log.Fatalf("Could not write to temporary file: %v", err)
	}
	return file
}
