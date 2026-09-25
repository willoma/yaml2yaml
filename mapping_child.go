package yaml2yaml

import (
	"go.yaml.in/yaml/v4"
)

// mappingChild is the result of evaluating a mapping child.
//
//   - kvs is not nil if child type is mapping entry or entries
//   - value is not nil if child type is single value
//   - both are nil if child is to be ignored
type mappingChild struct {
	path  string
	key   *yaml.Node
	value *yaml.Node
}

func nodeToMappingChild(path string, node *yaml.Node) ([]mappingChild, error) {
	if node == nil {
		return nil, nil
	}

	if node.Kind != yaml.MappingNode {
		return []mappingChild{{path: path, value: node}}, nil
	}

	var errs []error

	children := make([]mappingChild, 0, len(node.Content)/2)

	for i := 0; i < len(node.Content); i += 2 {
		keyAny, err := nodeToAny(node.Content[i])
		if err != nil {
			errs = append(errs, err)
		}

		children = append(children, mappingChild{
			path:  childPath(path, anyToString(keyAny)),
			key:   node.Content[i],
			value: node.Content[i+1],
		})
	}

	return children, joinErrors(errs...)
}
