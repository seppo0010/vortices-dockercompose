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
	Dirs         []FakeDir
	WrittenFiles []*FakeWrittenFile
}

func (f *FakeOS) MkdirAll(path string, perm os.FileMode) error {
	f.Dirs = append(f.Dirs, FakeDir{
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
	f.WrittenFiles = append(f.WrittenFiles, &FakeWrittenFile{
		Name:     name,
		Contents: buffer,
	})
	return w, nil
}

func (f *FakeOS) RemoveAll(path string) error {
	i := 0 // output index
	for _, x := range f.Dirs {
		if !strings.HasPrefix(x.Path, path) {
			// copy and increment index
			f.Dirs[i] = x
			i++
		}
	}
	f.Dirs = f.Dirs[:i]
	i = 0
	for _, x := range f.WrittenFiles {
		if !strings.HasPrefix(x.Name, path) {
			// copy and increment index
			f.WrittenFiles[i] = x
			i++
		}
	}
	f.WrittenFiles = f.WrittenFiles[:i]
	return nil
}

func (f *FakeOS) FileExists(path string) bool {
	for _, dir := range f.Dirs {
		if strings.HasPrefix(dir.Path, path) {
			return true
		}
	}
	for _, f := range f.WrittenFiles {
		if f.Name == path {
			return true
		}
	}
	return false
}
