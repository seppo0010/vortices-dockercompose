package dockercompose

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/seppo0010/vortices-dockercompose/exec"
)

type ServiceConfig struct {
	Image      string
	Command    []string `yaml:"command,omitempty"`
	Privileged bool
}

type Service struct {
	ServiceConfig         `yaml:",inline"`
	ContainerName         string `yaml:"container_name"`
	name                  string
	Networks              map[string]ServiceNetworkConfig `yaml:"networks"`
	compose               *Compose
	serviceNetworksConfig []ServiceNetworkConfig
}

type ServiceNetworkConfig struct {
	Network *Network `yaml:"-"`
	Aliases []string `yaml:"aliases,omitempty"`
}

func (s *Service) SetNetworks(serviceNetworksConfig []ServiceNetworkConfig) {
	s.serviceNetworksConfig = serviceNetworksConfig
	s.Networks = make(map[string]ServiceNetworkConfig, len(s.serviceNetworksConfig))
	for _, network := range s.serviceNetworksConfig {
		s.Networks[network.Network.name] = network
	}
}

func (s *Service) Exec(path string, args ...string) exec.Cmd {
	args = append([]string{"exec", "-T", s.name, path}, args...)
	cmd := s.compose.exec.New("docker-compose", args...)
	cmd.SetDir(s.compose.getTmpDir())
	return cmd
}

func (s *Service) SudoExec(path string, args ...string) exec.Cmd {
	args = append([]string{"exec", "-T", "--privileged", s.name, path}, args...)
	cmd := s.compose.exec.New("docker-compose", args...)
	cmd.SetDir(s.compose.getTmpDir())
	return cmd
}

func (s *Service) GetIPAddressForNetwork(network *Network) (string, error) {
	networksExec := s.compose.exec.New("docker", "inspect", "-f", "{{json .NetworkSettings.Networks}}", s.name)
	stdout, err := networksExec.StdoutPipe()
	if err != nil {
		log.Errorf("failed to pipe network settings stdout: %s", err.Error())
		return "", err
	}
	if err := networksExec.Start(); err != nil {
		log.Errorf("failed to inspect network settings: %s", err.Error())
		return "", err
	}
	var networks map[string]map[string]interface{}
	err = json.NewDecoder(stdout).Decode(&networks)
	if err != nil {
		log.Errorf("failed to decode network settings json: %s", err.Error())
		return "", err
	}
	if err := networksExec.Wait(); err != nil {
		log.Errorf("failed to wait inspect network settings: %s", err.Error())
		return "", err
	}

	for network_id, data := range networks {
		networkLabelExec := s.compose.exec.New("docker", "inspect", "-f", "{{range $key, $value := .Labels}}{{if eq $key \"com.docker.compose.network\"}}{{$value}}{{end}}{{end}}", network_id)
		stdout, err := networkLabelExec.StdoutPipe()
		if err != nil {
			log.Errorf("failed to pipe network label settings stdout: %s", err.Error())
			return "", err
		}

		if err = networkLabelExec.Start(); err != nil {
			log.Errorf("failed to run network label settings: %s", err.Error())
			return "", err
		}

		stdoutBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Errorf("failed to read network label settings stdout: %s", err.Error())
			return "", err
		}
		if err = networkLabelExec.Wait(); err != nil {
			log.Errorf("failed to wait network label settings stdout: %s", err.Error())
			return "", err
		}

		if strings.Trim(string(stdoutBytes), " \n") == network.name {
			ip, found := data["IPAddress"]
			if !found {
				log.Errorf("ip address not found")
				return "", fmt.Errorf("ip address not found")
			}
			ipStr, ok := ip.(string)
			if !ok {
				log.Errorf("invalid ip address, got %T", ip)
				return "", fmt.Errorf("invalid ip address, got %T", ip)
			}
			return ipStr, nil
		}
	}
	log.Errorf("could not find ip address for %s in network %s", s.name, network.name)
	return "", fmt.Errorf("could not find ip address for %s in network %s", s.name, network.name)
}
