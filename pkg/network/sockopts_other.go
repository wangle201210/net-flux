//go:build !unix && !windows

package network

import "errors"

func sysSetsockoptReuseAddr(fd uintptr, v int) error {
	return errors.ErrUnsupported
}

func sysSetsockoptReusePort(fd uintptr, v int) error {
	return errors.ErrUnsupported
}
