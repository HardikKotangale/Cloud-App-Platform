package spec

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Resources struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type AppSpec struct {
	Name      string    `yaml:"name"`
	Namespace string    `yaml:"namespace"`
	Image     string    `yaml:"image"`
	Port      int       `yaml:"port"`
	Replicas  int       `yaml:"replicas"`
	Resources Resources `yaml:"resources"`
}

func LoadFromFile(path string) (*AppSpec, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spec: %w", err)
	}

	var app AppSpec
	if err := yaml.Unmarshal(b, &app); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	return &app, nil
}
