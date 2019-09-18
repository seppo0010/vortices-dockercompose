package executor

import "io"

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
}
