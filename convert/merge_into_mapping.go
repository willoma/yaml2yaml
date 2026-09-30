package convert

import (
	"slices"

	"go.yaml.in/yaml/v4"
)

func MergeIntoMapping(mapping, key, value *yaml.Node) error {
	knownKey := slices.IndexFunc(mapping.Content, func(node *yaml.Node) bool {
		return nodesEqual(key, node)
	})

	if knownKey == -1 {
		mapping.Content = append(mapping.Content, key, value)

		return nil
	}

	var err error

	mapping.Content[knownKey+1], err = MergeNodes(mapping.Content[knownKey+1], value)

	return err
}
