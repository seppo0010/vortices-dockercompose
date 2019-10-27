package dockercompose

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/seppo0010/vortices-dockercompose/exec"
	"github.com/stretchr/testify/assert"
)

func TestCIDR(t *testing.T) {
	ranCommands := []*exec.FakeCmd{}
	compose, fakeExec, _ := mockCompose()
	fakeExec.StdoutHandler = func(f *exec.FakeCmd) (io.ReadCloser, error) {
		if len(f.Args) == 4 && f.Args[0] == "inspect" && f.Args[1] == "-f" && f.Args[2] == "{{(index .IPAM.Config 0).Subnet}}" && f.Args[3] == fmt.Sprintf("%s_network1", strings.Replace(compose.id, "-", "", -1)) {
			r, w := io.Pipe()
			go func() {
				w.Write([]byte("1.2.3.4/5"))
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
	cidr, err := network1.GetCIDR()
	assert.Nil(t, err)
	assert.Equal(t, cidr, "1.2.3.4/5")
}
