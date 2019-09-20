package os

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type FakeDir struct {
	Path string
	Perm os.FileMode
}

type FakeWrittenFile struct {
	Name     string
	Contents *bytes.Buffer
}

type FakeOS struct {
	dirs         []FakeDir
	writtenFiles []*FakeWrittenFile
}

func (f *FakeOS) MkdirAll(path string, perm os.FileMode) error {
	f.dirs = append(f.dirs, FakeDir{
		Path: path,
		Perm: perm,
	})
	return nil
}

func (f *FakeOS) TempDir() string {
	return "/tmp"
}

func (f *FakeOS) Create(name string) (io.WriteCloser, error) {
	r, w := io.Pipe()
	buffer := &bytes.Buffer{}
	go func() { io.Copy(buffer, r) }()
	f.writtenFiles = append(f.writtenFiles, &FakeWrittenFile{
		Name:     name,
		Contents: buffer,
	})
	return w, nil
}

func (f *FakeOS) RemoveAll(path string) error {
	i := 0 // output index
	for _, x := range f.dirs {
		if !strings.HasPrefix(x.Path, path) {
			// copy and increment index
			f.dirs[i] = x
			i++
		}
	}
	f.dirs = f.dirs[:i]
	return nil
}
