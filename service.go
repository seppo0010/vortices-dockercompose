package dockercompose

type ServiceConfig struct {
	Image      string
	Privileged bool
}

type Service struct {
	*ServiceConfig
}
