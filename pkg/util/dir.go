package util

import (
	"os"
	"path/filepath"
)

func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	return dir
}
