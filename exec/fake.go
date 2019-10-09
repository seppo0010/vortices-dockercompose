package exec

import (
	"io"
	"os"
)

type FakeCmd struct {
	Args []string
	Dir  string
	Path string

	fakeCommander *FakeCommander
	pipes         []io.Closer
	finished      chan error
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
	pipe, err := f.fakeCommander.StderrHandler(f)
	if pipe != nil {
		f.pipes = append(f.pipes, pipe)
	}
	return pipe, err
}

func (f *FakeCmd) StdinPipe() (io.WriteCloser, error) {
	pipe, err := f.fakeCommander.StdinHandler(f)
	if pipe != nil {
		f.pipes = append(f.pipes, pipe)
	}
	return pipe, err
}

func (f *FakeCmd) StdoutPipe() (io.ReadCloser, error) {
	pipe, err := f.fakeCommander.StdoutHandler(f)
	if pipe != nil {
		f.pipes = append(f.pipes, pipe)
	}
	return pipe, err
}

func (f *FakeCmd) Start() error {
	if f.fakeCommander.RunHandler != nil {
		go func() {
			f.finished <- f.fakeCommander.RunHandler(f)
		}()
	} else {
		go func() {
			f.finished <- nil
		}()
	}
	return nil
}

func (f *FakeCmd) Wait() error {
	return <-f.finished
}

func (f *FakeCmd) Run() error {
	if err := f.Start(); err != nil {
		return err
	}
	err := f.Wait()
	for _, p := range f.pipes {
		p.Close()
	}
	return err
}

func (f *FakeCmd) Kill() error {
	return nil
}

func (f *FakeCmd) Signal(sig os.Signal) error {
	return nil
}

type FakeCommander struct {
	StderrHandler func(*FakeCmd) (io.ReadCloser, error)
	StdoutHandler func(*FakeCmd) (io.ReadCloser, error)
	StdinHandler  func(*FakeCmd) (io.WriteCloser, error)
	RunHandler    func(*FakeCmd) error
}

func (f *FakeCommander) New(name string, arg ...string) Cmd {
	return &FakeCmd{
		Path:          name,
		Args:          arg,
		fakeCommander: f,
		finished:      make(chan error),
	}
}
