package main

import (
	"os/exec"
	"testing"
)

func TestDmesg(t *testing.T) {
	
	out, err := exec.Command("go", "run", "dmesg.go", "-c").Output()
	if err != nil {
		t.Fatalf("can't run dmesg: %v", err)
	}

	out, err = exec.Command("go", "run", "dmesg.go").Output()
	if err != nil {
		t.Fatalf("can't run dmesg: %v", err)
	}

	if len(out) > 0 {
		t.Fatalf("The log wasn't cleared, got %v", out)
	}
}
