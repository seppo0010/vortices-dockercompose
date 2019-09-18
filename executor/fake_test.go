package executor

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeCommandStdoutPipe(t *testing.T) {
	t.Parallel()
	cmd := (&FakeCommander{
		stdoutHandler: func(f *FakeCmd) (*ClosableBuffer, error) {
			if f.Path == "echo" {
				return &ClosableBuffer{
					Buffer: bytes.NewBufferString(strings.Join(f.Args, " ") + "\n"),
				}, nil
			}
			return nil, nil
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
		stdoutHandler: func(f *FakeCmd) (*ClosableBuffer, error) {
			if f.Path == "echo" {
				return &ClosableBuffer{
					Buffer: bytes.NewBufferString(strings.Join(f.Args, " ") + "\n"),
				}, nil
			}
			return nil, nil
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
		stderrHandler: func(f *FakeCmd) (*ClosableBuffer, error) {
			if f.Path == "ls" {
				return &ClosableBuffer{
					Buffer: bytes.NewBufferString(expectedStderr),
				}, nil
			}
			return nil, nil
		},
	}).New("ls", "/fake")
	stderrPipe, err := cmd.StderrPipe()
	assert.Nil(t, err)
	cmd.Start()
	stderr, err := ioutil.ReadAll(stderrPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stderr), expectedStderr)
}
