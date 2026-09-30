package convert

import "go.yaml.in/yaml/v4"

func nodesEqual(left, right *yaml.Node) bool {
	if left == nil || right == nil {
		return left == right
	}

	switch left.Kind {
	case yaml.DocumentNode:
		return nodesEqual(left.Content[0], right)

	case yaml.AliasNode:
		return nodesEqual(left.Alias, right)

	case yaml.SequenceNode, yaml.MappingNode:
		return nodesSequenceOrMappingEqual(left, right)

	case yaml.ScalarNode:
		return nodesScalarEqual(left, right)

	default:
		return false

	}
}

func nodesSequenceOrMappingEqual(left, right *yaml.Node) bool {
	switch right.Kind {
	case yaml.DocumentNode:
		return nodesSequenceOrMappingEqual(left, right.Content[0])

	case yaml.AliasNode:
		return nodesSequenceOrMappingEqual(left, right.Alias)

	case left.Kind:
		if len(left.Content) != len(right.Content) {
			return false
		}

		for i := range len(left.Content) {
			if !nodesEqual(left.Content[i], right.Content[i]) {
				return false
			}
		}

		return true

	default:
		return false
	}
}

func nodesScalarEqual(left, right *yaml.Node) bool {
	switch right.Kind {
	case yaml.DocumentNode:
		return nodesScalarEqual(left, right.Content[0])

	case yaml.AliasNode:
		return nodesScalarEqual(left, right.Alias)

	case yaml.ScalarNode:
		return left.Tag == right.Tag && left.Value == right.Value

	default:
		return false
	}
}
