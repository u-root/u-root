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
	"strings"
	"syscall"
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
	sock      string
	GOOS      string
	GOARCH    string
	Opts      []string
}

func (i *Image) String() string {
	return fmt.Sprintf("Image for %s:%s(%s) at %s: Cmd %s, Env %s", i.GOOS, i.GOARCH, i.Opts, i.endpoint(), i.Cmd, i.Env)
}

// endpoint names where cpud can be reached on the host.
func (i *Image) endpoint() string {
	if i.sock != "" {
		return i.sock
	}
	return client.DefaultPort
}

// Socket is the host unix socket forwarded to cpud, or "" when the default
// port is forwarded instead.
func (i *Image) Socket() string {
	return i.sock
}

// An Option configures an Image as New creates it.
type Option func(*Image) error

// WithSocket forwards a host unix socket to cpud rather than the default
// host port.
//
// There is deliberately no option for a different port. A port is a global
// name, so picking a free one is a race: nothing stops another process taking
// it between the look and the bind, and the failure that follows is a
// confusing one. A socket path is a name the caller already owns, so
// concurrent VMs cannot collide and nothing has to be reserved. qemu creates
// the socket; the directory must exist.
//
// The path must be absolute. qemu refuses a relative one outright -- "unix
// socket path must be absolute" -- and CommandContext runs it with Dir set to
// the VM's directory, so a relative path would not mean the same thing to
// qemu as it does to the dial here. Refusing it names the rule instead.
//
// unix(7) caps sun_path, so the path must fit in maxSockPath bytes. A
// directory from os.MkdirTemp with a short prefix leaves room.
func WithSocket(path string) Option {
	return func(i *Image) error {
		if path == "" {
			return fmt.Errorf("socket path is empty: %w", os.ErrInvalid)
		}
		if !filepath.IsAbs(path) {
			return fmt.Errorf("socket path %q is relative, and qemu requires an absolute one: %w", path, os.ErrInvalid)
		}
		if strings.ContainsRune(path, 0) {
			return fmt.Errorf("socket path %q contains a NUL: %w", path, os.ErrInvalid)
		}
		if len(path) > maxSockPath {
			return fmt.Errorf("socket path %q is %d bytes, over the %d unix(7) allows: %w", path, len(path), maxSockPath, os.ErrInvalid)
		}
		i.sock = path
		return nil
	}
}

// maxSockPath is the longest usable unix socket path: sun_path less the
// terminator.
//
// Taken from the structure rather than written down, because it is not the
// same everywhere -- 107 on Linux, 103 on Darwin and FreeBSD, which this
// package also builds for. A hardcoded 107 would pass a path macOS cannot
// bind, and macOS $TMPDIR is long enough to reach that on its own. len of an
// array is a constant expression, so this is still a compile-time constant.
const maxSockPath = len(syscall.RawSockaddrUnix{}.Path) - 1

var images = map[string]Image{
	"linux_amd64":   {Kernel: kernel_linux_amd64, InitRAMFS: linux_amd64, Cmd: []string{"qemu-system-x86_64", "-m", "4G"}},
	"linux_arm64":   {Kernel: kernel_linux_arm64, InitRAMFS: linux_arm64, Cmd: []string{"qemu-system-aarch64", "-machine", "virt", "-cpu", "max", "-m", "1G"}},
	"linux_arm":     {Kernel: kernel_linux_arm, InitRAMFS: linux_arm, Cmd: []string{"qemu-system-arm", "-M", "virt,highmem=off"}, Opts: []string{"GOARM=5"}},
	"linux_riscv64": {Kernel: kernel_linux_riscv64, InitRAMFS: linux_riscv64, Cmd: []string{"qemu-system-riscv64", "-M", "virt", "-cpu", "rv64", "-m", "1G"}},
}

// New creates an Image, using the kernel and arch to select the Image.
// It will return an error if there is a problem uncompressing
// the kernel and initramfs.
func New(goos, arch string, opts ...Option) (*Image, error) {
	i := &Image{}
	for _, o := range opts {
		if err := o(i); err != nil {
			return nil, err
		}
	}

	fwd := fmt.Sprintf("hostfwd=tcp::%s-:%s", client.DefaultPort, client.DefaultPort)
	if i.sock != "" {
		// -netdev is a keyval list: it splits on "," and reads ",," as a
		// literal comma. Only the qemu argument is escaped; i.sock stays
		// as the caller gave it, because that is what Socket reports and
		// what net.Dial is handed.
		fwd = fmt.Sprintf("hostfwd=unix:%s-:%s", strings.ReplaceAll(i.sock, ",", ",,"), client.DefaultPort)
	}

	common := []string{
		"-nographic",
		"-netdev", "user,id=net0,ipv4=on," + fwd,
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

	ramfs, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unzipped %d bytes: %w", len(ramfs), err)
	}

	if r, err = gzip.NewReader(bufio.NewReader(bytes.NewBuffer(im.Kernel))); err != nil {
		return nil, err
	}

	k, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unzipped %d bytes: %w", len(k), err)
	}

	i.Kernel, i.InitRAMFS = k, ramfs
	i.Cmd, i.Env = append(im.Cmd, common...), append(env, im.Env...)
	i.GOOS, i.GOARCH, i.Opts = goos, arch, im.Opts
	return i, nil
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
	cpu, err := i.cpuCmd(arg, args...)
	if err != nil {
		return nil, err
	}
	if err := cpu.Dial(); err != nil {
		return nil, err
	}
	cpu.Env = append(cpu.Env, "PATH=/bbin", "PWD=/", "SHELL=/bbin/x")
	return cpu, nil
}

// cpuCmd builds the client command and points it at this Image's endpoint.
//
// Split from CPUCommand, which dials, so that a test can check where the dial
// would be made without needing a VM to dial to.
func (i *Image) cpuCmd(arg string, args ...string) (*client.Cmd, error) {
	// Always "127.0.0.1" here, never the socket path, because Command runs
	// the name it is given through GetHostUser before any option applies.
	// That cuts on "@", so a socket under a directory such as
	// "workspace@2" would lose everything before it, and it then asks
	// ssh_config for a HostName, where a wildcard rule matches a path --
	// "*" matches "/" -- and rewrites it to a host. Either way the dial
	// would be made against something other than the socket.
	cpu := client.Command("127.0.0.1", append([]string{arg}, args...)...)
	cpu.Env = os.Environ()

	// A unix socket is named by its path and has no port; a TCP forward is
	// reached on the loopback address. Set after Command, so the path is
	// never a candidate for the rewrites above.
	network := "tcp"
	if i.sock != "" {
		cpu.Host, cpu.HostName, network = i.sock, i.sock, "unix"
	}

	if err := cpu.SetOptions(
		client.WithDisablePrivateKey(true),
		client.WithNetwork(network),
		// The port is ignored on a unix socket, which is named by its
		// path.
		client.WithPort(client.DefaultPort),
		client.WithRoot(i.dir),
		client.WithNameSpace("/"),
		client.With9P(true),
		client.WithTimeout("5s"),
	); err != nil {
		return nil, err
	}
	return cpu, nil
}
