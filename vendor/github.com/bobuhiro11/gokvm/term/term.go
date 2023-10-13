package term

import (
	"syscall"
	"unsafe"
)

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [32]uint8
	Pad    [3]byte
	Ispeed uint32
	Ospeed uint32
}

func read(fd int) (termios, error) {
	var t termios

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(fd), 0x5401,
		uintptr(unsafe.Pointer(&t)))

	if errno != 0 {
		return t, errno
	}

	return t, nil
}

func write(fd int, t termios) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(fd), 0x5402,
		uintptr(unsafe.Pointer(&t)))

	if errno != 0 {
		return errno
	}

	return nil
}

func IsTerminal() bool {
	_, err := read(0)

	return err == nil
}

func SetRawMode() (func(), error) {
	t, err := read(0)
	if err != nil {
		return func() {}, err
	}

	oldTermios := t

	t.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	t.Oflag &^= syscall.OPOST
	t.Cflag &^= syscall.CSIZE | syscall.PARENB
	t.Cflag |= syscall.CS8
	t.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0

	return func() {
		_ = write(0, oldTermios)
	}, write(0, t)
}
