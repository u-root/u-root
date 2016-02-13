package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMount(t *testing.T) {

	if err := exec.Command("dd", "if=/dev/zero", "of=test.bin", "count=10k").Run(); err != nil {
		t.Fatalf("Failed to create block device: %v", err)
	}
	defer os.Remove("test.bin")

	if err := exec.Command("losetup", "/dev/loop0", "test.bin").Run(); err != nil {
		t.Fatalf("Failed to connect loop0 to file: %v", err)
	}

	if err := exec.Command("mkfs.ext2", "/dev/loop0").Run(); err != nil {
		t.Fatalf("Failed to create test file system: %v", err)
	}

	device := "/dev/loop0"
	target := "/mnt"
	cmd := exec.Command("go", "run", "mount.go", "-t", "ext2", device, target)

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run mount: %v", err)
	}

	output, err := exec.Command("stat", "-f", "-c", "'%T'", target).Output()
	if err != nil {
		t.Fatalf("Failed to exec stat: %v")
	}

	if string(output) != "ext2/ext3" {
		t.Errorf("Seems %s is not a mounted device. want ext2/ext3, got %s", target, output)
	}

}
