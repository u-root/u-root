package boot

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/uio"
)

// ErrKernelMissing is returned by LinuxImage.Pack if no kernel is given.
var ErrKernelMissing = errors.New("must have non-nil kernel")

// LinuxImage implements OSImage for a Linux kernel + initramfs.
type LinuxImage struct {
	Kernel  io.ReaderAt
	Initrd  io.ReaderAt
	Cmdline string
}

var _ OSImage = &LinuxImage{}

// NewLinuxImageFromArchive reads a netboot21 Linux OSImage from a CPIO file
// archive.
func NewLinuxImageFromArchive(a *cpio.Archive) (*LinuxImage, error) {
	kernel, ok := a.Files["modules/kernel/content"]
	if !ok {
		return nil, fmt.Errorf("kernel missing from archive")
	}

	li := &LinuxImage{}
	li.Kernel = kernel

	if params, ok := a.Files["modules/kernel/params"]; ok {
		b, err := uio.ReadAll(params)
		if err != nil {
			return nil, err
		}
		li.Cmdline = string(b)
	}

	if initrd, ok := a.Files["modules/initrd/content"]; ok {
		li.Initrd = initrd
	}
	return li, nil
}

// Pack implements OSImage.Pack and writes all necessary files to the modules
// directory of `sw`.
func (li *LinuxImage) Pack(sw cpio.RecordWriter) error {
	if err := sw.WriteRecord(cpio.Directory("modules", 0700)); err != nil {
		return err
	}
	if err := sw.WriteRecord(cpio.Directory("modules/kernel", 0700)); err != nil {
		return err
	}
	if li.Kernel == nil {
		return ErrKernelMissing
	}
	kernel, err := uio.ReadAll(li.Kernel)
	if err != nil {
		return err
	}
	if err := sw.WriteRecord(cpio.StaticFile("modules/kernel/content", string(kernel), 0700)); err != nil {
		return err
	}
	if err := sw.WriteRecord(cpio.StaticFile("modules/kernel/params", li.Cmdline, 0700)); err != nil {
		return err
	}

	if li.Initrd != nil {
		if err := sw.WriteRecord(cpio.Directory("modules/initrd", 0700)); err != nil {
			return err
		}
		initrd, err := uio.ReadAll(li.Initrd)
		if err != nil {
			return err
		}
		if err := sw.WriteRecord(cpio.StaticFile("modules/initrd/content", string(initrd), 0700)); err != nil {
			return err
		}
	}

	return sw.WriteRecord(cpio.StaticFile("package_type", "linux", 0700))
}

func copyToFile(r io.Reader) (*os.File, error) {
	f, err := ioutil.TempFile("", "nerf-netboot")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	if err := f.Sync(); err != nil {
		return nil, err
	}

	readOnlyF, err := os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	return readOnlyF, nil
}

// ExecutionInfo implements OSImage.ExecutionInfo.
func (li *LinuxImage) ExecutionInfo(l *log.Logger) {
	k, err := copyToFile(uio.Reader(li.Kernel))
	if err != nil {
		l.Printf("Copying kernel to file: %v", err)
	}
	defer k.Close()

	var i *os.File
	if li.Initrd != nil {
		i, err = copyToFile(uio.Reader(li.Initrd))
		if err != nil {
			l.Printf("Copying initrd to file: %v", err)
		}
		defer i.Close()
	}

	l.Printf("Kernel: %s", k.Name())
	if i != nil {
		l.Printf("Initrd: %s", i.Name())
	}
	l.Printf("Command line: %s", li.Cmdline)
}

// Execute implements OSImage.Execute and kexec's the kernel with its initramfs.
func (li *LinuxImage) Execute() error {
	k, err := copyToFile(uio.Reader(li.Kernel))
	if err != nil {
		return err
	}
	defer k.Close()

	var i *os.File
	if li.Initrd != nil {
		i, err = copyToFile(uio.Reader(li.Initrd))
		if err != nil {
			return err
		}
		defer i.Close()
	}

	if err := kexec.FileLoad(k, i, li.Cmdline); err != nil {
		return err
	}
	return kexec.Reboot()
}
