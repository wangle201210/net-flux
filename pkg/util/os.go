package util

import (
	"net"
	"os"
	"os/user"
	"path/filepath"
)

// get the current working directory
//
// @return: the current working directory
// @example:
// - Pwd() -> "/Users/user/project"
func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	return dir
}

// get the home directory
//
// @return: the home directory
// @example:
// - HomeDir() -> "/Users/user"
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// get the temporary directory
//
// @return: the temporary directory
// @example:
// - TempDir() -> "/tmp"
func TempDir() string {
	return os.TempDir()
}

// get the current user name
//
// @return: the current user name
// @example:
// - UserName() -> "user"
func UserName() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

// get the hostname
//
// @return: the hostname
// @example:
// - Hostname() -> "localhost"
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// get the host IP
//
// @return: the host IP
// @example:
// - HostIP() -> ["127.0.0.1", "::1"]
func HostIP() ([]string, error) {
	ip, err := net.LookupIP(Hostname())
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, ip := range ip {
		ips = append(ips, ip.String())
	}
	return ips, nil
}
