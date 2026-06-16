//go:build unix

package network

import "golang.org/x/sys/unix"

func sysSetsockoptReuseAddr(fd uintptr, v int) error {
	return unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, v)
}

func sysSetsockoptReusePort(fd uintptr, v int) error {
	return unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, v)
}
