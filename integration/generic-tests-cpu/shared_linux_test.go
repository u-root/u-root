// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/u-root/cpu/client"
	"github.com/u-root/cpu/vm"
)

// The expensive part of a cpu test is building the u-root image, not running
// the VM: measured on linux/amd64, the build is 124s of CPU and the boot plus
// a handful of commands is 3s. So the image is built once per architecture and
// every test boots its own VM from it.
//
// Each VM is reached over its own unix socket, so the VMs run concurrently and
// tests can call t.Parallel. A socket path is a name this test owns, where a
// port is global and reserving one races with every other process. cpud still
// listens on the usual port inside the guest; only the host side moves.
// readyTimeout bounds the wait for cpud to answer. A guest of the host's own
// architecture answers in about 6s. An emulated one is far slower: the guest
// spends most of that time in init's DHCPv6 solicit, which is a fixed number
// of retries rather than a fixed delay, so it scales with emulation speed.
func readyTimeout(arch string) time.Duration {
	if arch == runtime.GOARCH {
		return 2 * time.Minute
	}
	return 10 * time.Minute
}

var images struct {
	sync.Mutex
	dir  string            // holds every built image
	path map[string]string // arch -> built cpio
}

func TestMain(m *testing.M) {
	// Set once here rather than per command: SetVerbose is global, and these
	// tests run in parallel.
	if _, spew := os.LookupEnv("VMTEST_SPEW"); spew {
		client.SetVerbose(log.Printf)
	}
	code := m.Run()
	if images.dir != "" {
		os.RemoveAll(images.dir)
	}
	os.Exit(code)
}

// qemuFor names the qemu binary an architecture needs. A missing binary skips
// rather than fails, which is what the console-based tests do through
// qemu.SkipIfNotArch.
var qemuFor = map[string]string{
	"amd64":   "qemu-system-x86_64",
	"arm":     "qemu-system-arm",
	"arm64":   "qemu-system-aarch64",
	"riscv64": "qemu-system-riscv64",
}

// imageFor builds arch's u-root image, or returns the one already built.
func imageFor(t *testing.T, arch string) string {
	t.Helper()

	images.Lock()
	defer images.Unlock()

	if p, ok := images.path[arch]; ok {
		return p
	}
	if images.dir == "" {
		d, err := os.MkdirTemp("", "cpu-images")
		if err != nil {
			t.Fatalf("image directory: %v", err)
		}
		images.dir = d
		images.path = map[string]string{}
	}

	i, err := vm.New("linux", arch)
	if err != nil {
		t.Fatalf("vm.New(linux, %s): %v", arch, err)
	}
	start := time.Now()
	p, err := i.Uroot(images.dir)
	if err != nil {
		t.Skipf("skipping: cannot build a %s u-root image: %v", arch, err)
	}
	t.Logf("built the %s image in %v, shared from here on", arch, time.Since(start).Round(time.Millisecond))

	images.path[arch] = p
	return p
}

// guest is a booted VM running cpud, reachable on its own host socket.
type guest struct {
	image *vm.Image // knows where to dial, and dials it
	dir   string    // the host directory the guest sees at /tmp/cpu
	sock  string    // the host socket cpud is reached on
	arch  string
}

// bootGuest boots a VM for arch from the shared image and waits for cpud to
// answer. The VM is killed when the test ends.
func bootGuest(t *testing.T, arch string) *guest {
	t.Helper()

	if _, err := exec.LookPath(qemuFor[arch]); err != nil {
		t.Skipf("skipping %s: %s is not installed", arch, qemuFor[arch])
	}
	image := imageFor(t, arch)

	// Reach this guest over a socket of its own rather than the shared port,
	// so guests can run at once. A path is a name this test owns, where a
	// port is global and picking a free one races with every other process.
	// unix(7) caps the path, and WithSocket refuses one that is too long
	// rather than leaving qemu to fail later, so keep the directory short.
	sockDir, err := os.MkdirTemp("", "s")
	if err != nil {
		t.Fatalf("socket directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(sockDir) })
	sock := filepath.Join(sockDir, "c")

	i, err := vm.New("linux", arch, vm.WithSocket(sock))
	if err != nil {
		t.Fatalf("vm.New(linux, %s, WithSocket(%s)): %v", arch, sock, err)
	}

	d := t.TempDir()
	c, err := i.CommandContext(t.Context(), d, image)
	if err != nil {
		t.Fatalf("building the VM command: %v", err)
	}
	if err := i.StartVM(c); err != nil {
		t.Fatalf("starting the VM: %v", err)
	}
	t.Cleanup(func() {
		if c.Process != nil {
			c.Process.Kill()
			c.Wait()
		}
	})

	g := &guest{image: i, sock: sock, dir: d, arch: arch}
	g.waitReady(t)
	return g
}

// waitReady polls until cpud answers. StartVM's fixed one second sleep is not
// enough on its own: without this the first command of a test absorbs the rest
// of the boot, which measured between 7 and 16 seconds.
func (g *guest) waitReady(t *testing.T) {
	t.Helper()

	limit := readyTimeout(g.arch)
	start := time.Now()
	for {
		if _, err := g.run("/bbin/date"); err == nil {
			t.Logf("%s guest ready at %s after %v", g.arch, g.sock, time.Since(start).Round(time.Millisecond))
			return
		}
		if time.Since(start) > limit {
			t.Fatalf("%s guest did not answer at %s within %v", g.arch, g.sock, limit)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Dir is the host directory the guest sees at /tmp/cpu.
func (g *guest) Dir() string { return g.dir }

// Run runs one command in the guest and returns its combined output. The
// output comes back over 9p rather than off the serial console, so there is
// nothing to parse and nothing to match against a prompt.
func (g *guest) Run(name string, args ...string) ([]byte, error) {
	return g.run(name, args...)
}

func (g *guest) run(name string, args ...string) ([]byte, error) {
	// CPUCommand rather than a hand-rolled client.Command: the Image knows
	// its own endpoint, and it knows not to pass a socket path to
	// client.Command, which would cut it at any "@" and then run it through
	// ssh_config, where a wildcard Host rule matches a path.
	cpu, err := g.image.CPUCommand(name, args...)
	if err != nil {
		return nil, err
	}
	// Close restores the terminal and shuts the session down. Running the
	// compiled test binary from a tty leaves stdin in raw mode without it.
	defer cpu.Close()
	return cpu.CombinedOutput()
}
