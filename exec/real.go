package exec

import (
	"os"
	"os/exec"
)

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

func (r *RealCmd) Kill() error {
	return r.Process.Kill()
}

func (r *RealCmd) Signal(sig os.Signal) error {
	return r.Process.Signal(sig)
}

type RealCommander struct{}

func (*RealCommander) New(name string, arg ...string) Cmd {
	return &RealCmd{exec.Command(name, arg...)}
}
