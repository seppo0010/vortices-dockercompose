package exec

import (
	"io"
	"os"
)

type Cmd interface {
	SetPath(path string)
	SetArgs(args []string)
	SetDir(dir string)
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	Run() error
	Kill() error
	Signal(sig os.Signal) error
}

type Commander interface {
	New(name string, arg ...string) Cmd
}
