package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
)

// Returns temporary directory and QEMU instance.
func testWithQEMU(t *testing.T, uinitPkg string) (string, *qemu.QEMU) {
	if _, ok := os.LookupEnv("UROOT_QEMU"); !ok {
		t.Skip("test is skipped unless UROOT_QEMU is set")
	}

	// TempDir
	tmpDir, err := ioutil.TempDir("", "uroot-integration")
	if err != nil {
		t.Fatal(err)
	}

	// Env
	env := golang.Default()
	env.CgoEnabled = false

	// Builder
	builder, err := uroot.GetBuilder("bb")
	if err != nil {
		t.Fatal(err)
	}

	// Packages
	pkgs, err := uroot.DefaultPackageImports(env)
	if err != nil {
		t.Fatal(err)
	}
	pkgs = append(pkgs, uinitPkg)

	// Archiver
	archiver, err := uroot.GetArchiver("cpio")
	if err != nil {
		t.Fatal(err)
	}

	// OutputFile
	outputFile := filepath.Join(tmpDir, fmt.Sprintf("initramfs.%s_%s.cpio", env.GOOS, env.GOARCH))
	w, err := archiver.OpenWriter(outputFile, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Build u-root
	opts := uroot.Opts{
		TempDir: tmpDir,
		Env:     env,
		Commands: []uroot.Commands{
			{
				Builder:  builder,
				Packages: pkgs,
			},
		},
		Archiver:   archiver,
		OutputFile: w,
	}
	if err := uroot.CreateInitramfs(opts); err != nil {
		t.Fatal(err)
	}

	// Start QEMU
	q := &qemu.QEMU{
		InitRAMFS: outputFile,
		Kernel:    "testdata/bzImage_amd64",
	}
	t.Logf("command line:\n%s", q.CmdLineQuoted())
	if err := q.Start(); err != nil {
		t.Fatal("could not spawn QEMU: ", err)
	}
	return tmpDir, q
}

func cleanup(t *testing.T, tmpDir string, q *qemu.QEMU) {
	q.Close()
	if t.Failed() {
		t.Log("Temp dir: ", tmpDir)
	} else {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temporary directory %s", tmpDir)
		}
	}
}

// TestHelloWorld runs an init which prints the string "HELLO WORLD" and exits.
func TestHelloWorld(t *testing.T) {
	// Create the CPIO and start QEMU.
	tmpDir, q := testWithQEMU(t, "github.com/u-root/u-root/integration/testdata/helloworld/uinit")
	defer cleanup(t, tmpDir, q)

	if err := q.Expect("HELLO WORLD"); err != nil {
		t.Fatal(err)
	}
}

// TestHelloWorldNegative runs an init which does not print the string "GOODBYE WORLD".
func TestHelloWorldNegative(t *testing.T) {
	// Create the CPIO and start QEMU.
	tmpDir, q := testWithQEMU(t, "github.com/u-root/u-root/integration/testdata/helloworld/uinit")
	defer cleanup(t, tmpDir, q)

	if err := q.Expect("GOODBYE WORLD"); err == nil {
		t.Fatal(err)
	}
}
