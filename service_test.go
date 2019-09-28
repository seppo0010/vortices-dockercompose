package dockercompose

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/seppo0010/vortices-dockercompose/exec"
	"github.com/stretchr/testify/assert"
)

func TestIPAddressForNetwork(t *testing.T) {
	ranCommands := []*exec.FakeCmd{}
	compose, fakeExec, _ := mockCompose()
	fakeExec.StdoutHandler = func(f *exec.FakeCmd) (io.ReadCloser, error) {
		if len(f.Args) == 4 && f.Args[0] == "inspect" && f.Args[1] == "-f" && f.Args[2] == "{{json .NetworkSettings.Networks}}" && f.Args[3] == "test-service" {
			r, w := io.Pipe()
			go func() {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"network1": map[string]interface{}{"IPAddress": "1.2.3.4"},
					"network2": map[string]interface{}{"IPAddress": "5.6.7.8"},
				})
				w.Close()
			}()
			return r, nil
		}
		if len(f.Args) == 4 && f.Args[0] == "inspect" && f.Args[1] == "-f" && f.Args[2] == "{{range $key, $value := .Labels}}{{if eq $key \"com.docker.compose.network\"}}{{$value}}{{end}}{{end}}" && f.Args[3] == "network1" {
			r, w := io.Pipe()
			go func() {
				w.Write([]byte("network1 \n"))
				w.Close()
			}()
			return r, nil
		}
		if len(f.Args) == 4 && f.Args[0] == "inspect" && f.Args[1] == "-f" && f.Args[2] == "{{range $key, $value := .Labels}}{{if eq $key \"com.docker.compose.network\"}}{{$value}}{{end}}{{end}}" && f.Args[3] == "network2" {
			r, w := io.Pipe()
			go func() {
				w.Write([]byte("network2 \n"))
				w.Close()
			}()
			return r, nil
		}

		panic("unexpected stdout handler call")
	}

	fakeExec.RunHandler = func(cmd *exec.FakeCmd) error {
		ranCommands = append(ranCommands, cmd)
		return nil
	}

	network1 := compose.AddNetwork("network1", NetworkConfig{})
	network2 := compose.AddNetwork("network2", NetworkConfig{})
	service := compose.AddService("test-service", ServiceConfig{}, []ServiceNetworkConfig{
		ServiceNetworkConfig{Network: network1},
		ServiceNetworkConfig{Network: network2},
	})
	ipAddress, err := service.GetIPAddressForNetwork(network2)
	assert.Nil(t, err)
	assert.Equal(t, ipAddress, "5.6.7.8")
}

func TestExec(t *testing.T) {
	ranCommands := []*exec.FakeCmd{}
	compose, fakeExec, _ := mockCompose()
	fakeExec.RunHandler = func(cmd *exec.FakeCmd) error {
		ranCommands = append(ranCommands, cmd)
		return nil
	}

	service := compose.AddService("test-service", ServiceConfig{}, []ServiceNetworkConfig{})
	cmd := service.Exec("ping", "google.com")
	err := cmd.Run()
	assert.Nil(t, err)
	assert.Equal(t, len(ranCommands), 1)
	assert.Equal(t, ranCommands[0].Path, "docker-compose")
	assert.Equal(t, ranCommands[0].Args, []string{"exec", "test-service", "ping", "google.com"})
	assert.Equal(t, ranCommands[0].Dir, fmt.Sprintf("/tmp/vortices-dockercompose/%s", compose.id))
}

func TestSudoExec(t *testing.T) {
	ranCommands := []*exec.FakeCmd{}
	compose, fakeExec, _ := mockCompose()
	fakeExec.RunHandler = func(cmd *exec.FakeCmd) error {
		ranCommands = append(ranCommands, cmd)
		return nil
	}

	service := compose.AddService("test-service", ServiceConfig{}, []ServiceNetworkConfig{})
	cmd := service.SudoExec("ping", "google.com")
	err := cmd.Run()
	assert.Nil(t, err)
	assert.Equal(t, len(ranCommands), 1)
	assert.Equal(t, ranCommands[0].Path, "docker-compose")
	assert.Equal(t, ranCommands[0].Args, []string{"exec", "--privileged", "test-service", "ping", "google.com"})
	assert.Equal(t, ranCommands[0].Dir, fmt.Sprintf("/tmp/vortices-dockercompose/%s", compose.id))
}
