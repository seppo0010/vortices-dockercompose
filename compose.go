package dockercompose

import (
	"errors"
	"path"

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

	*ComposeConfig

	Services map[string]*Service
	status   composeStatus
	exec     exec.Commander
	os       os.OS
}

func NewCompose(compose *ComposeConfig) *Compose {
	if compose.Version == "" {
		compose.Version = "2.1"
	}
	id := uuid.New().String()
	return &Compose{
		id:            id,
		ComposeConfig: compose,
		Services:      map[string]*Service{},
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

func (c *Compose) AddService(name string, serviceConfig *ServiceConfig) *Service {
	service := &Service{ServiceConfig: serviceConfig}
	if c.status != composeStatusSetup {
		panic("cannot register a service after started")
	}
	if _, found := c.Services[name]; found {
		panic("registering the same service twice")
	}
	c.Services[name] = service
	return service
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

	cmd := c.exec.New("docker-compose", "up", "-d")
	cmd.SetDir(c.getTmpDir())

	if err := cmd.Run(); err != nil {
		log.Errorf("failed to start docker compose: %s", err.Error())
		return errors.New("failed to start docker-compose")
	}

	return nil
}

func (c *Compose) Stop() error {
	if c.status != composeStatusRunning {
		return errors.New("cannot stop if status is not running")
	}
	c.status = composeStatusStopped

	cmd := c.exec.New("docker-compose", "down")
	cmd.SetDir(c.getTmpDir())

	if err := cmd.Run(); err != nil {
		log.Errorf("failed to stop docker compose: %s", err.Error())
		return errors.New("failed to stop docker-compose")
	}

	return nil
}
