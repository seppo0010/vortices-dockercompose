package executor

import (
	"io"
)

type FakeCmd struct {
	Args []string
	Dir  string
	Path string

	fakeCommander *FakeCommander
	pipes         []io.Closer
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
	pipe, err := f.fakeCommander.stderrHandler(f)
	if pipe != nil {
		f.pipes = append(f.pipes, pipe)
	}
	return pipe, err
}

func (f *FakeCmd) StdinPipe() (io.WriteCloser, error) {
	pipe, err := f.fakeCommander.stdinHandler(f)
	if pipe != nil {
		f.pipes = append(f.pipes, pipe)
	}
	return pipe, err
}

func (f *FakeCmd) StdoutPipe() (io.ReadCloser, error) {
	pipe, err := f.fakeCommander.stdoutHandler(f)
	if pipe != nil {
		f.pipes = append(f.pipes, pipe)
	}
	return pipe, err
}

func (f *FakeCmd) Start() error {
	return nil
}

func (f *FakeCmd) Run() error {
	for _, p := range f.pipes {
		p.Close()
	}
	return nil
}

type FakeCommander struct {
	stderrHandler func(*FakeCmd) (*ClosableBuffer, error)
	stdoutHandler func(*FakeCmd) (*ClosableBuffer, error)
	stdinHandler  func(*FakeCmd) (*ClosableBuffer, error)
}

func (f *FakeCommander) New(name string, arg ...string) Cmd {
	return &FakeCmd{
		Path:          name,
		Args:          arg,
		fakeCommander: f,
	}
}
