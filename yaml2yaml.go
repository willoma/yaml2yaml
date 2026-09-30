// Package yaml2yaml provides utilities for processing and transforming YAML data.
package yaml2yaml

import (
	oerrors "errors"
	"fmt"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
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
		return nil, oerrors.Join(pErr, fmt.Errorf("converting data to YAML: %w", err))
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
		return nil, oerrors.Join(pErr, fmt.Errorf("converting data to YAML: %w", err))
	}

	return result, pErr
}

// Process applies any processing to a YAML node, returning the processed node and any errors.
func Process(node *yaml.Node) (*yaml.Node, error) {
	return process(definitions.New(), "", node)
}

func process(defs *definitions.Definitions, key string, node *yaml.Node) (*yaml.Node, error) {
	result, err := processTransparent(defs, node)

	return result, errors.Locate(key, node.Line, node.Column, err)
}

// processTransparent processees a node, but does not locate the returned error..
func processTransparent(defs *definitions.Definitions, node *yaml.Node) (*yaml.Node, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		return processTransparent(defs, node.Content[0])
	case yaml.AliasNode:
		return process(defs, "alias:"+node.Alias.Anchor, node.Alias)
	case yaml.SequenceNode:
		return processSequence(defs, node)
	case yaml.MappingNode:
		return processMapping(defs, node)
	case yaml.ScalarNode:
		return processScalar(defs, node)
	case yaml.StreamNode:
		return nil, errors.ErrStreamNotSupported
	default:
		return nil, errors.ErrUnknownKind
	}
}
