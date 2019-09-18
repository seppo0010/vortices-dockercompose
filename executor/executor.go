package executor

import (
	"bytes"
	"io"
)

type Cmd interface {
	SetPath(path string)
	SetArgs(args []string)
	SetDir(dir string)
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Run() error
}

type ClosableBuffer struct {
	*bytes.Buffer
	closed bool
}

func (c *ClosableBuffer) Read(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	return c.Buffer.Read(p)
}

func (c *ClosableBuffer) Write(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	return c.Buffer.Write(p)
}

func (c *ClosableBuffer) Close() error {
	if c.closed {
		return io.ErrClosedPipe
	}
	c.closed = true
	return nil
}
