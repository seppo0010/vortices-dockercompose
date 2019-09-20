package dockercompose

type NetworkConfig struct {
}

type Network struct {
	NetworkConfig `yaml:",inline"`
	name          string
}
