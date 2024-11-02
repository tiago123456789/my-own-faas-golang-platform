package types

type Function struct {
	Trigger map[string]map[string]string `yaml:"trigger"`
}

type Config struct {
	Function Function          `yaml:"function"`
	Name     string            `yaml:"name"`
	Runtime  string            `yaml:"runtime"`
	Envs     map[string]string `yaml:"envs"`
	Cpu      string            `yaml:"cpu"`
	Memory   string            `yaml:"memory"`
}
