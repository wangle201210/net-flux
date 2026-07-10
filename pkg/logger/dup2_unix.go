//go:build unix

package logger

import "golang.org/x/sys/unix"

func dup2(oldfd int, newfd int) error {
	return unix.Dup2(oldfd, newfd)
}
