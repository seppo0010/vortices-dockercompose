package dockercompose

import (
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ComposeConfig struct {
	Version string
}

type Compose struct {
	id     string
	tmpDir string

	*ComposeConfig

	Services map[string]*Service
	started  bool
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
	}
}

func (c *Compose) AddService(name string, serviceConfig *ServiceConfig) *Service {
	service := &Service{ServiceConfig: serviceConfig}
	if c.started {
		panic("cannot register a service after started")
	}
	if _, found := c.Services[name]; found {
		panic("registering the same service twice")
	}
	c.Services[name] = service
	return service
}

func (c *Compose) Start() error {
	c.started = true
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

	if err := c.cmd.Start(); err != nil {
		log.Fatalf("failed to start docker-compose")
	}

	return nil
}
