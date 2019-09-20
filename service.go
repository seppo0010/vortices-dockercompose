package dockercompose

type ServiceConfig struct {
	Image      string
	Privileged bool
}

type Service struct {
	ServiceConfig `yaml:",inline"`
	name          string
	NetworkNames  []string `yaml:"networks"`
	networks      []*Network
}

func (s *Service) SetNetworks(networks []*Network) {
	s.networks = networks
	s.NetworkNames = make([]string, len(s.networks))
	for i, network := range networks {
		s.NetworkNames[i] = network.name
	}
}
