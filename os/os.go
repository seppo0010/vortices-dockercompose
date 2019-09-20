package os

import (
	"io"
	"os"
)

type OS interface {
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
	TempDir() string
	Create(name string) (io.WriteCloser, error)
}
