// Package yaml2yaml provides utilities for processing and transforming YAML data.
package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

const indentSize = 2

// ProcessBytes converts a YAML byte slice to a YAML byte slice, with any processing applied.
func ProcessBytes(src []byte) ([]byte, error) {
	var node yaml.Node

	if err := yaml.Unmarshal(src, &node); err != nil {
		return nil, fmt.Errorf("reading data as YAML: %w", err)
	}

	parsed, pErr := Process(&node)

	result, err := yaml.Dump(parsed, yaml.WithIndent(indentSize))
	if err != nil {
		return nil, joinErrors(pErr, fmt.Errorf("converting data to YAML: %w", err))
	}

	return result, pErr
}

// ProcessFile reads a YAML file from the given path and returns the processed YAML byte slice.
func ProcessFile(path string) ([]byte, error) {
	node, err := readFile(path)
	if err != nil {
		return nil, err
	}

	parsed, pErr := Process(node)

	result, err := yaml.Dump(parsed, yaml.WithIndent(indentSize))
	if err != nil {
		return nil, joinErrors(pErr, fmt.Errorf("converting data to YAML: %w", err))
	}

	return result, pErr
}

// Process applies any processing to a YAML node, returning the processed node and any errors.
func Process(src *yaml.Node) (*yaml.Node, error) {
	return newNode(src).process()
}
