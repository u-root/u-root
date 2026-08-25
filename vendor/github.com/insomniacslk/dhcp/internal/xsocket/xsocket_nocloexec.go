//go:build !(dragonfly || freebsd || linux || netbsd || openbsd)

package xsocket

import "golang.org/x/sys/unix"

func socketCloexec(domain, typ, proto int) (int, error) {
	return unix.Socket(domain, typ, proto)
}
