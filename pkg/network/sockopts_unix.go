//go:build unix

package network

import "syscall"

func sysSetsockoptReuseAddr(fd uintptr, v int) error {
	return syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, v)
}

func sysSetsockoptReusePort(fd uintptr, v int) error {
	return syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, v)
}
