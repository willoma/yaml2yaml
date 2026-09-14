package yaml2yaml

import "go.yaml.in/yaml/v4"

func kindString(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "document"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.MappingNode:
		return "mapping"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	case yaml.StreamNode:
		return "stream"
	default:
		return "unknown"
	}
}
