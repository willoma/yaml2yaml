// Package yaml2yaml provides utilities for processing and transforming YAML data.
package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

const indentSize = 2

// Bytes converts a YAML byte slice to a YAML byte slice, with any processing applied.
func Bytes(src []byte) ([]byte, error) {
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

// Process applies any processing to a YAML node, returning the processed node and any errors.
func Process(src *yaml.Node) (*yaml.Node, error) {
	return newNode(src).process()
}
