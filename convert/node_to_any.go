package convert

import (
	"encoding/base64"
	"strconv"

	"go.yaml.in/yaml/v4"
)

// NodeToAny converts a YAML node to any value.
// It does not return errors on [yaml.StreamNode] or default kinds,
// because it runs on processed YAML nodes only
// (ie. those kinds have already been excluded).
func NodeToAny(node *yaml.Node) any {
	switch node.Kind {
	case yaml.DocumentNode:
		return NodeToAny(node.Content[0])
	case yaml.AliasNode:
		return NodeToAny(node.Alias)
	case yaml.SequenceNode:
		return sequenceToAny(node)
	case yaml.MappingNode:
		return mappingToAny(node)
	case yaml.ScalarNode:
		return scalarToAny(node)
	default:
		return nil
	}
}

func sequenceToAny(node *yaml.Node) any {
	children := make([]any, len(node.Content))
	for i, child := range node.Content {
		children[i] = NodeToAny(child)
	}

	return children
}

func mappingToAny(node *yaml.Node) any {
	children := make(map[string]any, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		children[NodeToString(node.Content[i])] = NodeToAny(node.Content[i+1])
	}

	return children
}

// scalarToAny converts a YAML scalar node to any value.
// If the value cannot be converted (unknown tag, or content not matching tag),
// it returns the raw string value.
func scalarToAny(node *yaml.Node) any {
	switch node.Tag {
	case "!!bool":
		result, err := strconv.ParseBool(node.Value)
		if err != nil {
			return node.Value
		}

		return result
	case "!!int":
		result, err := strconv.Atoi(node.Value)
		if err != nil {
			return node.Value
		}

		return result
	case "!!float":
		result, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return node.Value
		}

		return result
	case "!!str":
		return node.Value
	case "!!binary":
		result, err := base64.StdEncoding.DecodeString(node.Value)
		if err != nil {
			return node.Value
		}

		return result
	default:
		return node.Value
	}
}
