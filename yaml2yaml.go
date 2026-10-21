// Package yaml2yaml provides utilities for processing and transforming YAML data.
package yaml2yaml

import (
	"fmt"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/read"
)

const (
	indentSize = 2
	lineWidth  = -1
)

// ProcessBytes converts a YAML byte slice to a YAML byte slice, with any processing applied.
func ProcessBytes(src []byte) ([]byte, error) {
	var srcNode yaml.Node

	if err := yaml.Unmarshal(src, &srcNode); err != nil {
		return nil, fmt.Errorf("reading data as YAML: %w", err)
	}

	return processToBytes("", &srcNode)
}

// ProcessFile reads a YAML file from the given path and returns the processed YAML byte slice.
func ProcessFile(path string) ([]byte, error) {
	src, err := read.File(path)
	if err != nil {
		return nil, err
	}

	return processToBytes(filepath.Dir(path), src)
}

func processToBytes(basePath string, src *yaml.Node) ([]byte, error) {
	parsed, err := Process(basePath, src)
	if err != nil {
		return nil, err
	}

	result, err := yaml.Dump(parsed, yaml.WithIndent(indentSize), yaml.WithLineWidth(lineWidth))
	if err != nil {
		return nil, fmt.Errorf("converting data to YAML: %w", err)
	}

	return result, nil
}

// Process applies any processing to a YAML node, returning the processed node and any errors.
func Process(basePath string, src *yaml.Node) (*yaml.Node, error) {
	sourceNode, err := parseNode(basePath, src)
	if err != nil {
		return nil, err
	}

	processed, err := sourceNode.Process(rootDefinitions{})
	if err != nil {
		return nil, err
	}

	return processed.ToYAML(), nil
}
