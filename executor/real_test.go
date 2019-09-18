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
