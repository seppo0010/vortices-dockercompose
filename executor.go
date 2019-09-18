package dockercompose

import (
	"bytes"
	"io"
	"os/exec"
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
	bytes.Buffer
	closed bool
}

func (c ClosableBuffer) Read(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	return c.Buffer.Read(p)
}

func (c ClosableBuffer) Write(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	return c.Buffer.Write(p)
}

func (c ClosableBuffer) Close() error {
	return nil
}

type FakeCmd struct {
	Args   []string
	Dir    string
	Path   string
	Stderr ClosableBuffer
	Stdout ClosableBuffer
	Stdin  ClosableBuffer
}

func (f *FakeCmd) SetArgs(args []string) {
	f.Args = args
}

func (f *FakeCmd) SetPath(path string) {
	f.Path = path
}

func (f *FakeCmd) SetDir(dir string) {
	f.Dir = dir
}

func (f *FakeCmd) StderrPipe() (io.ReadCloser, error) {
	return f.Stderr, nil
}

func (f *FakeCmd) StdinPipe() (io.WriteCloser, error) {
	return f.Stdin, nil
}

func (f *FakeCmd) StdoutPipe() (io.ReadCloser, error) {
	return f.Stdout, nil
}

func (f *FakeCmd) Start() error {
	return nil
}

func (f *FakeCmd) Run() error {
	return nil
}

type RealCmd struct {
	*exec.Cmd
}

func (r *RealCmd) SetPath(path string) {
	r.Cmd.Path = path
}

func (r *RealCmd) SetArgs(args []string) {
	r.Cmd.Args = args
}

func (r *RealCmd) SetDir(dir string) {
	r.Cmd.Dir = dir
}

type Commander interface {
	New(name string, arg ...string) Cmd
}

type RealCommander struct{}

func (*RealCommander) New(name string, arg ...string) Cmd {
	return &RealCmd{exec.Command(name, arg...)}
}

type FakeCommander struct{}

func (*FakeCommander) New(name string, arg ...string) Cmd {
	return &FakeCmd{
		Path: name,
		Args: arg,
	}
}
