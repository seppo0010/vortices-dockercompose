package dockercompose

import (
	"path"
	"testing"

	"github.com/seppo0010/vortices-dockercompose/exec"
	"github.com/seppo0010/vortices-dockercompose/os"
	"github.com/stretchr/testify/assert"
)

func mockCompose() (*Compose, *exec.FakeCommander, *os.FakeOS) {
	compose := NewCompose(ComposeConfig{})
	fakeExec := &exec.FakeCommander{}
	fakeOS := &os.FakeOS{}
	compose.exec = fakeExec
	compose.os = fakeOS
	return compose, fakeExec, fakeOS
}

func TestStartStopIntegration(t *testing.T) {
	compose := NewCompose(ComposeConfig{})
	compose.AddService("test-service", ServiceConfig{
		Image: "ubuntu",
	}, nil)
	err := compose.Start()
	assert.Nil(t, err)
	err = compose.Stop()
	assert.Nil(t, err)
}

func TestStart(t *testing.T) {
	ranCommands := []*exec.FakeCmd{}
	compose, fakeExec, fakeOS := mockCompose()
	fakeExec.RunHandler = func(cmd *exec.FakeCmd) error {
		ranCommands = append(ranCommands, cmd)
		return nil
	}
	network1 := compose.AddNetwork("test-network1", NetworkConfig{})
	network2 := compose.AddNetwork("test-network2", NetworkConfig{})
	compose.AddService("test-service", ServiceConfig{
		Image: "ubuntu",
	}, []*Network{network1, network2})
	err := compose.Start()
	assert.Nil(t, err)

	assert.Equal(t, len(ranCommands), 1)
	assert.Equal(t, len(fakeOS.WrittenFiles), 1)

	assert.Equal(t, ranCommands[0].Path, "docker-compose")
	assert.Equal(t, ranCommands[0].Args, []string{"up", "-d"})
	assert.Equal(t, ranCommands[0].Dir, path.Dir(fakeOS.WrittenFiles[0].Name))

	assert.Equal(t, string(fakeOS.WrittenFiles[0].Contents.Bytes()),
		`version: "2.1"
services:
  test-service:
    image: ubuntu
    privileged: false
    networks:
    - test-network1
    - test-network2
networks:
  test-network1: {}
  test-network2: {}
`)
}

func TestStop(t *testing.T) {
	ranCommands := []*exec.FakeCmd{}
	compose, fakeExec, fakeOS := mockCompose()
	fakeExec.RunHandler = func(cmd *exec.FakeCmd) error {
		ranCommands = append(ranCommands, cmd)
		return nil
	}
	compose.AddService("test-service", ServiceConfig{
		Image: "ubuntu",
	}, nil)
	err := compose.Start()
	assert.Nil(t, err)
	err = compose.Stop()
	assert.Nil(t, err)

	assert.Equal(t, len(ranCommands), 2)

	assert.Equal(t, ranCommands[1].Path, "docker-compose")
	assert.Equal(t, ranCommands[1].Args, []string{"down"})
	assert.Equal(t, ranCommands[1].Dir, path.Dir(fakeOS.WrittenFiles[0].Name))
}

func TestNetworkIntegration(t *testing.T) {
	compose := NewCompose(ComposeConfig{})
	compose.AddService("test-service", ServiceConfig{
		Image: "ubuntu",
	}, nil)
	compose.AddNetwork("test-network", NetworkConfig{})
	err := compose.Start()
	assert.Nil(t, err)
	err = compose.Stop()
	assert.Nil(t, err)
}
