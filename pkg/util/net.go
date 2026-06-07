package util

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func ParseProtocol(addr string) (string, error) {
	if strings.Contains(addr, "://") {
		protocol := strings.Split(addr, "://")[0]
		return protocol, nil
	}
	return "", errors.New("invalid address")
}

func ParseHostAndPort(addr string) (string, int, error) {
	if strings.Contains(addr, "://") {
		addr = strings.TrimPrefix(addr, "://")
	}
	if strings.Contains(addr, "/") {
		addr = strings.Split(addr, "/")[0]
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return host, p, nil
}

func ParseHost(addr string) (string, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return host, nil
}

func ParsePort(addr string) (int, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	return p, nil
}
