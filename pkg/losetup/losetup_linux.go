package losetup

import (
	"fmt"
	"os"
	"syscall"
)

const (
	/*
	 * IOCTL commands --- we will commandeer 0x4C ('L')
	 */
	LOOP_SET_CAPACITY = 0x4C07
	LOOP_CHANGE_FD    = 0x4C06
	LOOP_GET_STATUS64 = 0x4C05
	LOOP_SET_STATUS64 = 0x4C04
	LOOP_GET_STATUS   = 0x4C03
	LOOP_SET_STATUS   = 0x4C02
	LOOP_CLR_FD       = 0x4C01
	LOOP_SET_FD       = 0x4C00
	LO_NAME_SIZE      = 64
	LO_KEY_SIZE       = 32
	/* /dev/loop-control interface */
	LOOP_CTL_ADD      = 0x4C80
	LOOP_CTL_REMOVE   = 0x4C81
	LOOP_CTL_GET_FREE = 0x4C82

	SYS_IOCTL         = 16
)

// FindLoopDevice finds an unused loop device.
func FindLoopDevice() (name string, err error) {
	cfd, err := syscall.Open("/dev/loop-control", syscall.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer syscall.Close(cfd)

	if number, err := LoopCtlGetFree(uintptr(cfd)); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("/dev/loop%d", number), nil
	}

}

// LoopClearFd clears the loop device associated with filedescriptor fd.
func LoopClearFd(fd uintptr) error {
	if _, _, err := syscall.Syscall(SYS_IOCTL, fd, LOOP_CLR_FD, 0); err != 0 {
		return err
	}

	return nil
}

// LoopCtlGetFree finds a free loop device querying the loop control device pointed
// by fd. It returns the number of the free loop device /dev/loopX
func LoopCtlGetFree(fd uintptr) (uintptr, error) {
	number, _, err := syscall.Syscall(SYS_IOCTL, fd, LOOP_CTL_GET_FREE, 0)
	if err != 0 {
		return 0, err
	}
	return number, nil
}

// LoopSetFd associates a loop device pointed by lfd with a regular file pointed by ffd.
func LoopSetFd(lfd, ffd uintptr) error {
	_, _, err := syscall.Syscall(SYS_IOCTL, lfd, LOOP_SET_FD, ffd)
	if err != 0 {
		return err
	}

	return nil
}

// LoopSetFdFiles associates loop device "devicename" with regular file "filename"
func LoopSetFdFiles(devicename, filename string) error {
	mode := os.O_RDWR
	file, err := os.OpenFile(filename, mode, 0644)
	if err != nil {
		mode = os.O_RDONLY
		file, err = os.OpenFile(filename, mode, 0644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	device, err := os.OpenFile(devicename, mode, 0644)
	if err != nil {
		return err
	}
	defer device.Close()

	return LoopSetFd(device.Fd(), file.Fd())
}

// LoopClearFdFile clears the loop device "devicename"
func LoopClearFdFile(devicename string) error {
	device, err := os.Open(devicename)
	if err != nil {
		return err
	}
	defer device.Close()

	return LoopClearFd(device.Fd())
}
