package os

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeTempDir(t *testing.T) {
	tempdir := (&FakeOS{}).TempDir()
	assert.Equal(t, tempdir, "/tmp")
}

func TestFakeMakeDeleteDir(t *testing.T) {
	os := &FakeOS{}
	tempdir := os.TempDir()
	dir := path.Join(tempdir, "a", "b")
	err := os.MkdirAll(dir, 0744)
	assert.Nil(t, err)
	err = os.RemoveAll(dir)
	assert.Nil(t, err)
	err = os.RemoveAll(dir)
	assert.Nil(t, err)
}

func TestFakeMakeDeleteFile(t *testing.T) {
	os := &FakeOS{}
	tempdir := os.TempDir()
	dir := path.Join(tempdir, "a")
	err := os.MkdirAll(dir, 0744)
	assert.Nil(t, err)
	f, err := os.Create(path.Join(dir, "c"))
	assert.Nil(t, err)
	written, err := f.Write([]byte{1, 2, 3})
	assert.Nil(t, err)
	assert.Equal(t, written, 3)
	f.Close()
	assert.Nil(t, err)
	err = os.RemoveAll(dir)
	assert.Nil(t, err)
}
