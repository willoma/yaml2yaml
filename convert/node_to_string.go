package convert

import (
	"strings"

	"go.yaml.in/yaml/v4"
)

// NodeToString converts a YAML node to a string.
// It does not return errors on [yaml.StreamNode] or default kinds,
// because it runs on processed YAML nodes only
// (ie. those kinds have already been excluded).
func NodeToString(node *yaml.Node) string {
	switch node.Kind {
	case yaml.DocumentNode:
		return NodeToString(node.Content[0])
	case yaml.AliasNode:
		return NodeToString(node.Alias)
	case yaml.SequenceNode:
		return sequenceToString(node)
	case yaml.MappingNode:
		return mappingToString(node)
	case yaml.ScalarNode:
		return scalarToString(node)
	default:
		return ""
	}
}

func sequenceToString(node *yaml.Node) string {
	content := make([]string, len(node.Content))
	for i, item := range node.Content {
		content[i] = NodeToString(item)
	}

	return "[" + strings.Join(content, ",") + "]"
}

func mappingToString(node *yaml.Node) string {
	content := make([]string, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		content[i/2] = NodeToString(node.Content[i]) + ":" + NodeToString(node.Content[i+1])
	}

	return "{" + strings.Join(content, ",") + "}"
}

func scalarToString(node *yaml.Node) string {
	return node.Value
}
