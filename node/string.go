package node

import (
	"strings"

	"go.yaml.in/yaml/v4"
)

func ToString(node *yaml.Node) string {
	switch node.Kind {
	case yaml.DocumentNode:
		return ToString(node.Content[0])
	case yaml.AliasNode:
		return ToString(node.Alias)
	case yaml.SequenceNode:
		return sequenceToString(node)
	case yaml.MappingNode:
		return mappingToString(node)
	case yaml.ScalarNode:
		return scalarToString(node)
	case yaml.StreamNode:
		return "stream"
	default:
		return "unknown"
	}
}

func sequenceToString(node *yaml.Node) string {
	content := make([]string, len(node.Content))
	for i, item := range node.Content {
		content[i] = ToString(item)
	}

	return "[" + strings.Join(content, ",") + "]"
}

func mappingToString(node *yaml.Node) string {
	content := make([]string, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		content[i/2] = ToString(node.Content[i]) + ":" + ToString(node.Content[i+1])
	}

	return "{" + strings.Join(content, ",") + "}"
}

func scalarToString(node *yaml.Node) string {
	return node.Value
}
