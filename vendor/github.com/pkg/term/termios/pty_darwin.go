package termios

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/unix"
)

func open_pty_master() (uintptr, error) {
	return open_device("/dev/ptmx")
}

const (
	_IOC_PARAM_SHIFT = 13
	_IOC_PARAM_MASK  = (1 << _IOC_PARAM_SHIFT) - 1
)

func _IOC_PARM_LEN(ioctl uintptr) uintptr {
	return (ioctl >> 16) & _IOC_PARAM_MASK
}

func Ptsname(fd uintptr) (string, error) {
	n := make([]byte, _IOC_PARM_LEN(unix.TIOCPTYGNAME))

	err := ioctl(fd, unix.TIOCPTYGNAME, uintptr(unsafe.Pointer(&n[0])))
	if err != nil {
		return "", err
	}

	for i, c := range n {
		if c == 0 {
			return string(n[:i]), nil
		}
	}
	return "", errors.New("TIOCPTYGNAME string not NUL-terminated")
}

func grantpt(fd uintptr) error {
	return unix.IoctlSetInt(int(fd), unix.TIOCPTYGRANT, 0)
}

func unlockpt(fd uintptr) error {
	return unix.IoctlSetInt(int(fd), unix.TIOCPTYUNLK, 0)
}
