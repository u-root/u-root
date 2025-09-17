//go:generate go run build.go

package vm

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/u-root/cpu/client"
)

//go:embed initramfs_linux_amd64.cpio.gz
var linux_amd64 []byte

//go:embed initramfs_linux_arm64.cpio.gz
var linux_arm64 []byte

//go:embed initramfs_linux_arm.cpio.gz
var linux_arm []byte

//go:embed initramfs_linux_riscv64.cpio.gz
var linux_riscv64 []byte

//go:embed kernel_linux_amd64.gz
var kernel_linux_amd64 []byte

//go:embed kernel_linux_arm64.gz
var kernel_linux_arm64 []byte

//go:embed kernel_linux_arm.gz
var kernel_linux_arm []byte

//go:embed kernel_linux_riscv64.gz
var kernel_linux_riscv64 []byte

// Image defines an image, including []byte for a kernel and initramfs;
// a []string for the Cmd and its args; and its environment;
// and a directory in which to run.
type Image struct {
	Kernel    []byte
	InitRAMFS []byte
	Cmd       []string
	Env       []string
	dir       string
	GOOS      string
	GOARCH    string
	Opts      []string
}

func (i *Image) String() string {
	return fmt.Sprintf("Image for %s:%s(%s): Cmd %s, Env %s", i.GOOS, i.GOARCH, i.Opts, i.Cmd, i.Env)
}

var images = map[string]Image{
	"linux_amd64":   {Kernel: kernel_linux_amd64, InitRAMFS: linux_amd64, Cmd: []string{"qemu-system-x86_64", "-m", "4G"}},
	"linux_arm64":   {Kernel: kernel_linux_arm64, InitRAMFS: linux_arm64, Cmd: []string{"qemu-system-aarch64", "-machine", "virt", "-cpu", "max", "-m", "1G"}},
	"linux_arm":     {Kernel: kernel_linux_arm, InitRAMFS: linux_arm, Cmd: []string{"qemu-system-arm", "-M", "virt,highmem=off"}, Opts: []string{"GOARM=5"}},
	"linux_riscv64": {Kernel: kernel_linux_riscv64, InitRAMFS: linux_riscv64, Cmd: []string{"qemu-system-riscv64", "-M", "virt", "-cpu", "rv64", "-m", "1G"}},
}

// New creates an Image, using the kernel and arch to select the Image.
// It will return an error if there is a problem uncompressing
// the kernel and initramfs.
func New(goos, arch string) (*Image, error) {
	common := []string{
		"-nographic",
		"-netdev", "user,id=net0,ipv4=on,hostfwd=tcp::17010-:17010",
		// required for mac. No idea why. Should work on linux. If not, we'll need a bit
		// more logic.
		"-device", "e1000-82545em,netdev=net0,id=net0,mac=52:54:00:c9:18:27",
		"-netdev", "user,id=net1",
		"-device", "e1000-82545em,netdev=net1,id=net1,mac=52:54:00:c9:18:28",
		// No password needed, you're just a guest vm ...
		// The kernel may not understand ip=dhcp, in which case it just ends up
		// in init's environment.
		// To add cpud debugging, you can add -d to the append string.
		"--append", "ip=dhcp init=/init -pk \"\"",
	}
	env := []string{}
	n := fmt.Sprintf("%s_%s", goos, arch)
	im, ok := images[n]
	if !ok {
		return nil, fmt.Errorf("(%q,%q): %w", goos, arch, os.ErrNotExist)
	}

	r, err := gzip.NewReader(bufio.NewReader(bytes.NewBuffer(im.InitRAMFS)))
	if err != nil {
		return nil, err
	}

	i, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unzipped %d bytes: %w", len(i), err)
	}

	if r, err = gzip.NewReader(bufio.NewReader(bytes.NewBuffer(im.Kernel))); err != nil {
		return nil, err
	}

	k, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unzipped %d bytes: %w", len(k), err)
	}
	return &Image{Kernel: k, InitRAMFS: i, Cmd: append(im.Cmd, common...), Env: append(env, im.Env...), GOOS: goos, GOARCH: arch, Opts: im.Opts}, nil
}

// Uroot builds a uroot cpio into the a directory.
// It returns the full path of the file, or an error.
func Uroot(d, GOOS, GOARCH string, opts ...string) (string, error) {
	out := filepath.Join(d, GOOS+"_"+GOARCH+".cpio")
	c := exec.Command("u-root", "-o", out)
	c.Env = append(append(os.Environ(), "CGO_ENABLED=0", "GOARCH="+GOARCH, "GOOS="+GOOS), opts...)
	if out, err := c.CombinedOutput(); err != nil {
		return "", fmt.Errorf("u-root initramfs:%q:%w", out, err)
	}
	return out, nil

}

// Uroot builds a uroot cpio for an Image
func (i *Image) Uroot(d string) (string, error) {
	return Uroot(d, i.GOOS, i.GOARCH, i.Opts...)
}

// CommandContext starts qemu, given a context, directory in which to
// run the command. The variadic arguments are a set of cpios which will
// be merged into Image.InitRAMFS. A typical use of the extra arguments
// will be to extend the initramfs; they are usually not needed.
func (image *Image) CommandContext(ctx context.Context, d string, extra ...string) (*exec.Cmd, error) {
	image.dir = d
	i, k := filepath.Join(d, "initramfs"), filepath.Join(d, "kernel")

	if err := os.WriteFile(k, image.Kernel, 0644); err != nil {
		return nil, err
	}

	ir, err := os.OpenFile(i, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	defer ir.Close()
	for _, n := range extra {
		f, err := os.Open(n)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if _, err := io.Copy(ir, f); err != nil {
			return nil, err
		}
	}
	if _, err := ir.Write(image.InitRAMFS); err != nil {
		return nil, err
	}

	ir.Close()
	c := exec.CommandContext(ctx, image.Cmd[0], append(image.Cmd[1:], "-kernel", k, "-initrd", i)...)
	c.Env = append(os.Environ(), c.Env...)
	c.Dir = d
	return c, nil
}

// StartVM is used to start a VM (or in fact any exec.Cmd).
// Once cmd.Start is called, StartVM delays for one second.
// This time has been experimentally determined as the minimum
// required for the guest network to be ready.
func (*Image) StartVM(c *exec.Cmd) error {
	if err := c.Start(); err != nil {
		return fmt.Errorf("starting VM: %w", err)
	}
	time.Sleep(time.Second)
	return nil
}

// CPUCommand runs a command in a guest running cpud.
// It is similar to exec.Command, in that it accepts an arg and
// a set of optional args. It differs in that it can return an error.
// If there are no errors, it returns a client.Cmd.
// The returned client.Cmd can be called with CombinedOutput, to make
// it easier to scan output from a command for error messages.
func (i *Image) CPUCommand(arg string, args ...string) (*client.Cmd, error) {
	cpu := client.Command("127.0.0.1", append([]string{arg}, args...)...)
	cpu.Env = os.Environ()

	if err := cpu.SetOptions(
		client.WithDisablePrivateKey(true),
		client.WithPort("17010"),
		client.WithRoot(i.dir),
		client.WithNameSpace("/"),
		client.With9P(true),
		client.WithTimeout("5s"),
	); err != nil {
		return nil, err
	}

	if err := cpu.Dial(); err != nil {
		return nil, err
	}
	cpu.Env = append(cpu.Env, "PATH=/bbin", "PWD=/", "SHELL=/bbin/x")
	return cpu, nil
}
