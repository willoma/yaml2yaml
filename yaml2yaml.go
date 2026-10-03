// Package yaml2yaml provides utilities for processing and transforming YAML data.
package yaml2yaml

import (
	"fmt"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
	"github.com/willoma/yaml2yaml/internal/read"
)

const indentSize = 2

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
	var errs []error

	parsed, err := Process(basePath, src)
	if err != nil {
		errs = append(errs, err)
	}

	result, err := yaml.Dump(parsed, yaml.WithIndent(indentSize))
	if err != nil {
		errs = append(errs, fmt.Errorf("converting data to YAML: %w", err))
	}

	return result, errors.Join(errs...)
}

// Process applies any processing to a YAML node, returning the processed node and any errors.
func Process(basePath string, src *yaml.Node) (*yaml.Node, error) {
	var errs []error

	sourceNode := parseNode(basePath, src)

	if sourceNode == nil {
		return nil, errors.ErrEmptyNode
	}

	processed, err := sourceNode.Process(node.NewDefinitions())
	if err != nil {
		errs = append(errs, err)
	}

	return processed.ToYAML(), errors.Join(errs...)
}
