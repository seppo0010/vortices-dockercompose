package os

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealTempDir(t *testing.T) {
	tempdir := (&RealOS{}).TempDir()
	assert.Equal(t, tempdir, "/tmp")
}

func TestRealMakeDeleteDir(t *testing.T) {
	os := &RealOS{}
	tempdir := os.TempDir()
	dir := path.Join(tempdir, "a", "b")
	err := os.MkdirAll(dir, 0744)
	assert.Nil(t, err)
	err = os.RemoveAll(dir)
	assert.Nil(t, err)
	err = os.RemoveAll(dir)
	assert.Nil(t, err)
}

func TestRealMakeDeleteFile(t *testing.T) {
	os := &RealOS{}
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

func TestRealFileExists(t *testing.T) {
	os := &RealOS{}
	tempdir := os.TempDir()
	dir := path.Join(tempdir, "a")
	err := os.MkdirAll(dir, 0744)
	assert.Nil(t, err)
	assert.Equal(t, os.FileExists(dir), true)
	filePath := path.Join(dir, "c")
	f, err := os.Create(filePath)
	assert.Nil(t, err)
	f.Close()
	assert.Nil(t, err)
	assert.Equal(t, os.FileExists(filePath), true)
	err = os.RemoveAll(dir)
	assert.Nil(t, err)
	assert.Equal(t, os.FileExists(filePath), false)
}
