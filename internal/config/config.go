package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level driftwatch configuration.
type Config struct {
	Version  string    `yaml:"version"`
	Services []Service `yaml:"services"`
}

// Service describes a single declared service and its expected state.
type Service struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Replicas    int               `yaml:"replicas"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate performs basic sanity checks on the loaded configuration.
func (c *Config) Validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("config must declare at least one service")
	}
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service[%d]: name is required", i)
		}
		if svc.Image == "" {
			return fmt.Errorf("service %q: image is required", svc.Name)
		}
		if svc.Replicas < 0 {
			return fmt.Errorf("service %q: replicas must be >= 0", svc.Name)
		}
	}
	return nil
}
