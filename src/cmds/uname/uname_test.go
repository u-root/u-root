package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestUname(t *testing.T) {
	exec.Command("go", "build", "uname.go").Run()
	defer os.Remove("uname")

	options := []string{"-m", "-n", "-r", "-s", "-v"}

	for i := 0; i < len(options); i++ {
		want, err := exec.Command("./uname", options[i]).Output()
		if err != nil {
			t.Errorf("Can't exec ./uname %s: %v", options[i], err)
		}

		got, err := exec.Command("uname", options[i]).Output()
		if err != nil {
			t.Errorf("Can't exec uname %s (from unix): %v", options[i], err)
		}

		if string(want) != string(got) {
			t.Errorf("Fail while trying option %s: want %s got %s", options[i], want, got)
		}
	}
}
