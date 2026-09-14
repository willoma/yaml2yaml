package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

const (
	indentSize = 2
)

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

func Process(src *yaml.Node) (*yaml.Node, error) {
	return newNode(src).process()
}

// specificTags are tags that are processed specially by yaml2yaml, in [node.evaluateMappingChild].
var specificTags = []string{"!define", "!apply", "!if", "!for"}
