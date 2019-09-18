package executor

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealCommandStdoutPipe(t *testing.T) {
	t.Parallel()
	cmd := (&RealCommander{}).New("/bin/echo", "1", "2")
	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)
	cmd.Start()
	stdout, err := ioutil.ReadAll(stdoutPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "1 2\n")
}

func TestRealCommandStdoutPipeFinished(t *testing.T) {
	t.Parallel()
	cmd := (&RealCommander{}).New("echo", "1", "2")
	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)
	cmd.Run()
	_, err = ioutil.ReadAll(stdoutPipe)
	assert.NotNil(t, err)
}

func TestRealCommandSetDir(t *testing.T) {
	t.Parallel()
	cmd := (&RealCommander{}).New("pwd")
	cmd.SetDir("/")
	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)
	cmd.Start()
	stdout, err := ioutil.ReadAll(stdoutPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "/\n")
}

func TestRealCommandStderrPipe(t *testing.T) {
	t.Parallel()
	expectedStderr := "ls: cannot access '/fake': No such file or directory\n"
	cmd := (&RealCommander{}).New("ls", "/fake")
	stderrPipe, err := cmd.StderrPipe()
	assert.Nil(t, err)
	cmd.Start()
	stderr, err := ioutil.ReadAll(stderrPipe)
	assert.Nil(t, err)
	assert.Equal(t, string(stderr), expectedStderr)
}

func TestRealCommandStdinPipe(t *testing.T) {
	t.Parallel()
	cmd := (&RealCommander{}).New("wc", "-c")
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
