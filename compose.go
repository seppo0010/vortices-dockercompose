package dockercompose

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/seppo0010/vortices-dockercompose/exec"
	"github.com/seppo0010/vortices-dockercompose/os"
)

type composeStatus int

const (
	composeStatusSetup composeStatus = iota
	composeStatusRunning
	composeStatusStopped
)

type ComposeConfig struct {
	Version string
}

type Compose struct {
	id     string
	tmpDir string

	ComposeConfig `yaml:",inline"`

	Services map[string]*Service
	Networks map[string]*Network
	status   composeStatus
	exec     exec.Commander
	os       os.OS
}

func NewCompose(compose ComposeConfig) *Compose {
	if compose.Version == "" {
		compose.Version = "2.1"
	}
	id := uuid.New().String()
	return &Compose{
		id:            id,
		ComposeConfig: compose,
		Services:      map[string]*Service{},
		Networks:      map[string]*Network{},
		status:        composeStatusSetup,

		exec: &exec.RealCommander{},
		os:   &os.RealOS{},
	}
}

func (c *Compose) getTmpDir() string {
	if c.tmpDir == "" {
		c.tmpDir = path.Join(c.os.TempDir(), "vortices-dockercompose", c.id)
	}
	return c.tmpDir
}

func (c *Compose) AddService(name string, serviceConfig ServiceConfig, networks []ServiceNetworkConfig) *Service {
	if c.status != composeStatusSetup {
		panic("cannot register a service after started")
	}
	service := &Service{ServiceConfig: serviceConfig, ContainerName: name, name: name, compose: c}
	if networks != nil {
		service.SetNetworks(networks)
	}
	if _, found := c.Services[name]; found {
		panic("registering the same service twice")
	}
	c.Services[name] = service
	return service
}

func (c *Compose) AddNetwork(name string, networkConfig NetworkConfig) *Network {
	network := &Network{NetworkConfig: networkConfig, name: name}
	if c.status != composeStatusSetup {
		panic("cannot register a network after started")
	}
	if _, found := c.Networks[name]; found {
		panic("registering the same network twice")
	}
	c.Networks[name] = network
	return network
}

func (c *Compose) execOrFail(context, name string, arg ...string) ([]byte, error) {
	cmd := c.exec.New(name, arg...)
	cmd.SetDir(c.getTmpDir())
	stdoutPipe, err := cmd.StdoutPipe()
	stderrPipe, err := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		log.Errorf("failed to start %s: %s", context, err.Error())
		return nil, fmt.Errorf("failed to start %s: %s", context, err.Error())
	}
	stdout, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		log.Errorf("failed to read stdout %s: %s", context, err.Error())
		return nil, fmt.Errorf("failed to read stdout %s: %s", context, err.Error())
	}
	stderr, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		log.Errorf("failed to read stderr %s: %s", context, err.Error())
		return nil, fmt.Errorf("failed to read stderr %s: %s", context, err.Error())
	}
	if err = cmd.Wait(); err != nil {
		log.Errorf("failed to %s: %s\n%s", context, err.Error(), string(stderr))
		return nil, fmt.Errorf("failed to %s: %s\n%s", context, err.Error(), string(stderr))
	}
	return stdout, nil
}

func (c *Compose) Start() error {
	c.status = composeStatusRunning
	err := c.os.MkdirAll(c.getTmpDir(), 0744)
	if err != nil {
		return err
	}
	f, err := c.os.Create(path.Join(c.getTmpDir(), "docker-compose.yml"))
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	err = encoder.Close()
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	log.Infof("starting docker compose")
	defer log.Infof("finished starting docker compose")

	_, err = c.execOrFail("start docker compose", "docker-compose", "up", "-d")
	if err != nil {
		return errors.New("failed to start docker-compose")
	}

	return nil
}

func (c *Compose) Logs(machine ...string) (string, error) {
	logs, err := c.execOrFail("docker compose logs", "docker-compose", append([]string{"logs", "--no-color"}, machine...)...)
	if err != nil {
		return "", errors.New("failed to run docker-compose logs")
	}

	return string(logs), nil
}

func (c *Compose) Stop() error {
	if c.status != composeStatusRunning {
		return errors.New("cannot stop if status is not running")
	}
	c.status = composeStatusStopped

	log.Infof("stopping docker compose")
	defer log.Infof("finished stopping docker compose")

	_, err := c.execOrFail("stop docker compose", "docker-compose", "down")
	if err != nil {
		return errors.New("failed to stop docker-compose")
	}

	return nil
}

func (c *Compose) Clear() error {
	if c.status == composeStatusRunning {
		return errors.New("cannot clear if status is running")
	}

	log.Infof("clearing docker compose")
	defer log.Infof("finished clearing docker compose")

	if c.tmpDir == "" {
		return nil
	}

	return c.os.RemoveAll(c.tmpDir)
}

func (c *Compose) BuildDockerPath(name, path string) (string, error) {
	if !c.os.FileExists(path) {
		return "", fmt.Errorf("path %s does not exist", path)
	}

	log.Infof("starting to build docker image %s", name)
	defer log.Infof("finished building docker image %s", name)

	cmd := c.exec.New("docker", "build", path)
	outPipe, err := cmd.StdoutPipe()
	errPipe, err := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to build docker image at path %s: %s", path, err.Error())
	}
	out, err := ioutil.ReadAll(outPipe)
	if err != nil {
		return "", fmt.Errorf("failed to read docker image at path %s: %s", path, err.Error())
	}
	stderr, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return "", fmt.Errorf("failed to read docker stderr at path %s: %s", path, err.Error())
	}
	if err = cmd.Wait(); err != nil {
		return "", fmt.Errorf("failed to build docker image at path %s: %s\n%s", path, err.Error(), string(stderr))
	}
	submatches := regexp.MustCompile(`Successfully built ([a-fA-F0-9]*)`).FindStringSubmatch(string(out))
	if len(submatches) == 0 {
		return "", fmt.Errorf("could not find docker image tag. Full output:\n%s", string(out))
	}
	return submatches[1], nil
}

func (c *Compose) BuildDocker(name, script string) (string, error) {
	return c.buildDocker(name, script, uuid.New().String())
}

func (c *Compose) buildDocker(name, script, uuidString string) (string, error) {
	dirPath := path.Join(c.os.TempDir(), uuidString)
	err := c.os.MkdirAll(dirPath, 0744)
	if err != nil {
		return "", err
	}

	filePath := path.Join(dirPath, "Dockerfile")
	f, err := c.os.Create(filePath)
	if err != nil {
		return "", err
	}
	_, err = f.Write([]byte(script))
	if err != nil {
		return "", err
	}
	f.Close()
	defer c.os.RemoveAll(dirPath)

	return c.BuildDockerPath(name, dirPath)
}
