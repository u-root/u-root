package integration

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if len(os.Getenv("UROOT_KERNEL")) == 0 {
		log.Fatalf("Failed to run tests: no kernel provided")
	}
	if len(os.Getenv("UROOT_QEMU")) == 0 {
		log.Fatalf("Failed to run tests: no QEMU binary provided")
	}

	log.Printf("Starting generic tests...")

	os.Exit(m.Run())
}
