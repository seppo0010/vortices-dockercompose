package exec

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeCommandStdoutPipe(t *testing.T) {
	t.Parallel()
	cmd := (&FakeCommander{
		StdoutHandler: func(f *FakeCmd) (io.ReadCloser, error) {
			if f.Path == "echo" {
				r, w := io.Pipe()
				go func() {
					w.Write([]byte(strings.Join(f.Args, " ") + "\n"))
					w.Close()
				}()
				return r, nil
			}
			panic("unexpected")
		},
	}).New("echo", "1", "2")
	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)
	cmd.Start()
	stdout, err := ioutil.ReadAll(stdoutPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "1 2\n")
}

func TestFakeCommandStdoutPipeFinished(t *testing.T) {
	t.Parallel()
	cmd := (&FakeCommander{
		StdoutHandler: func(f *FakeCmd) (io.ReadCloser, error) {
			if f.Path == "echo" {
				r, w := io.Pipe()
				go func() {
					w.Write([]byte(strings.Join(f.Args, " ") + "\n"))
					w.Close()
				}()
				return r, nil
			}
			panic("unexpected")
		},
	}).New("echo", "1", "2")
	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)
	cmd.Run()
	_, err = ioutil.ReadAll(stdoutPipe)
	assert.NotNil(t, err)
}

func TestFakeCommandStderrPipe(t *testing.T) {
	t.Parallel()
	expectedStderr := "ls: cannot access '/fake': No such file or directory\n"
	cmd := (&FakeCommander{
		StderrHandler: func(f *FakeCmd) (io.ReadCloser, error) {
			if f.Path == "ls" {
				r, w := io.Pipe()
				go func() {
					w.Write([]byte(expectedStderr))
					w.Close()
				}()
				return r, nil
			}
			panic("unexpected")
		},
	}).New("ls", "/fake")
	stderrPipe, err := cmd.StderrPipe()
	assert.Nil(t, err)
	cmd.Start()
	stderr, err := ioutil.ReadAll(stderrPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stderr), expectedStderr)
}

func TestFakeCommandStdinPipe(t *testing.T) {
	t.Parallel()
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	cmd := (&FakeCommander{
		StdinHandler: func(f *FakeCmd) (io.WriteCloser, error) {
			if f.Path == "wc" && len(f.Args) == 1 && f.Args[0] == "-c" {
				return stdinW, nil
			}
			panic("unexpected")
		},
		StdoutHandler: func(f *FakeCmd) (io.ReadCloser, error) {
			if f.Path == "wc" && len(f.Args) == 1 && f.Args[0] == "-c" {
				return stdoutR, nil
			}
			panic("unexpected")
		},
		RunHandler: func(f *FakeCmd) error {
			if f.Path == "wc" && len(f.Args) == 1 && f.Args[0] == "-c" {
				read, err := ioutil.ReadAll(stdinR)
				if err != nil {
					return err
				}
				stdoutW.Write([]byte(fmt.Sprintf("%d\n", len(read))))
				stdoutW.Close()
				return nil
			}
			panic("unexpected")
		},
	}).New("wc", "-c")
	stdinPipe, err := cmd.StdinPipe()
	assert.Nil(t, err)
	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)

	cmd.Start()
	n, err := stdinPipe.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, n, 5)
	err = stdinPipe.Close()
	assert.Nil(t, err)

	stdout, err := ioutil.ReadAll(stdoutPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "5\n")
}
