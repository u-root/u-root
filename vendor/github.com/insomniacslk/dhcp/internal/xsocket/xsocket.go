package xsocket

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// CloexecSocket creates a new socket with the close-on-exec flag set.
//
// If the OS doesn't support the close-on-exec flag, this function will try a workaround.
func CloexecSocket(domain, typ, proto int) (int, error) {
	fd, err := socketCloexec(domain, typ, proto)
	if err == nil {
		return fd, nil
	}

	if err == unix.EINVAL || err == unix.EPROTONOSUPPORT {
		// SOCK_CLOEXEC is not supported, try without it, but avoid racing with fork/exec
		syscall.ForkLock.RLock()
		defer syscall.ForkLock.RUnlock()

		fd, err = unix.Socket(domain, typ, proto)
		if err != nil {
			return -1, err
		}

		unix.CloseOnExec(fd)

		return fd, nil
	}

	return fd, err
}
