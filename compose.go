package dockercompose

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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
	cmd      *exec.Cmd
}

func NewCompose(compose *ComposeConfig) *Compose {
	if compose.Version == "" {
		compose.Version = "2.1"
	}
	id := uuid.New().String()
	tmpDir := path.Join(os.TempDir(), "vortices-dockercompose", id)
	return &Compose{
		id:            id,
		tmpDir:        tmpDir,
		ComposeConfig: compose,
		Services:      map[string]*Service{},
		status:        composeStatusSetup,
	}
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
	err := os.MkdirAll(c.tmpDir, 0744)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(c.tmpDir, "docker-compose.yml"))
	if err != nil {
		return err
	}
	yml, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = f.Write(yml)
	if err != nil {
		return err
	}
	f.Close()

	c.cmd = exec.Command("docker-compose", "up", "-d")
	c.cmd.Dir = c.tmpDir

	if err := c.cmd.Run(); err != nil {
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

	c.cmd = exec.Command("docker-compose", "down")
	c.cmd.Dir = c.tmpDir

	if err := c.cmd.Run(); err != nil {
		log.Errorf("failed to stop docker compose: %s", err.Error())
		return errors.New("failed to stop docker-compose")
	}

	return nil
}
