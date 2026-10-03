package read

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

func File(path string) (*yaml.Node, error) {
	//nolint:gosec // Yes, file inclusion via variable: it's the goal here.
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var node yaml.Node

	if err := yaml.Unmarshal(src, &node); err != nil {
		return nil, fmt.Errorf("reading data as YAML: %w", err)
	}

	return &node, nil
}
