package dockercompose

import (
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

type NetworkConfig struct {
}

type Network struct {
	NetworkConfig `yaml:",inline"`
	name          string
	compose       *Compose
}

func (n *Network) GetCIDR() (string, error) {
	networkID := fmt.Sprintf("%s_%s", strings.Replace(n.compose.id, "-", "", -1), n.name)
	networksExec := n.compose.exec.New("docker", "inspect", "-f", "{{(index .IPAM.Config 0).Subnet}}", networkID)
	stdout, err := networksExec.StdoutPipe()
	if err != nil {
		log.Errorf("failed to pipe network settings stdout: %s", err.Error())
		return "", err
	}
	if err := networksExec.Start(); err != nil {
		log.Errorf("failed to inspect network settings: %s", err.Error())
		return "", err
	}
	stdoutStr, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Errorf("failed to read network settings json: %s", err.Error())
		return "", err
	}
	if err := networksExec.Wait(); err != nil {
		log.Errorf("failed to wait inspect network settings: %s", err.Error())
		return "", err
	}
	return string(stdoutStr), nil
}
