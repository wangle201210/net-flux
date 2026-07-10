//go:build !unix

package logger

import "errors"

func dup2(oldfd int, newfd int) error {
	return errors.ErrUnsupported
}
