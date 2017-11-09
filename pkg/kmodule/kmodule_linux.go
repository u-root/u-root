package kmodule

import (
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Flags to finit_module(2) / FileInit.
const (
	// Ignore symbol version hashes.
	MODULE_INIT_IGNORE_MODVERSIONS = 0x1

	// Ignore kernel version magic.
	MODULE_INIT_IGNORE_VERMAGIC = 0x2
)

// Init loads the kernel module given by image with the given options.
func Init(image []byte, opts string) error {
	optsNull, err := unix.BytePtrFromString(opts)
	if err != nil {
		return fmt.Errorf("kmodule.Init: could not convert %q to C string: %v", opts, err)
	}

	if _, _, e := unix.Syscall(unix.SYS_INIT_MODULE, uintptr(unsafe.Pointer(&image[0])), uintptr(len(image)), uintptr(unsafe.Pointer(optsNull))); e != 0 {
		return fmt.Errorf("init_module(%v, %q) failed with %v", image, opts, e)
	}

	return nil
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags.
//
// FileInit falls back to Init when the finit_module(2) syscall is not available.
func FileInit(f *os.File, opts string, flags uintptr) error {
	optsNull, err := unix.BytePtrFromString(opts)
	if err != nil {
		return fmt.Errorf("kmodule.Init: could not convert %q to C string: %v", opts, err)
	}

	if _, _, e := unix.Syscall(unix.SYS_FINIT_MODULE, f.Fd(), uintptr(unsafe.Pointer(optsNull)), flags); e == unix.ENOSYS {
		if flags != 0 {
			return fmt.Errorf("finit_module unavailable")
		}

		// Fall back to regular init_module(2).
		img, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("kmodule.FileInit: %v", err)
		}
		return Init(img, opts)
	} else if e != 0 {
		return fmt.Errorf("finit_module(%v, %q, %#x) failed with %v", f, opts, flags, e)
	}

	return nil
}

// Delete removes a kernel module.
func Delete(name string, flags uintptr) error {
	modnameptr, err := unix.BytePtrFromString(name)
	if err != nil {
		return fmt.Errorf("could not delete module %q: %v", name, err)
	}

	if _, _, e := unix.Syscall(unix.SYS_DELETE_MODULE, uintptr(unsafe.Pointer(modnameptr)), flags, 0); e != 0 {
		return fmt.Errorf("could not delete module %q: %v", name, e)
	}

	return nil
}
