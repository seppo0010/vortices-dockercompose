package os

import (
	"io"
	"os"
)

type RealOS struct{}

func (*RealOS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (*RealOS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (*RealOS) TempDir() string {
	return os.TempDir()
}

func (*RealOS) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (*RealOS) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
