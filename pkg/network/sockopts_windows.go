//go:build windows

package network

import (
	"errors"
	"syscall"
)

func sysSetsockoptReuseAddr(fd uintptr, v int) error {
	return syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, v)
}

func sysSetsockoptReusePort(fd uintptr, v int) error {
	return errors.ErrUnsupported
}
