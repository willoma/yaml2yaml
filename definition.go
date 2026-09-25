package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

func (n *node) setDefinition(name string, value *yaml.Node) error {
	if n.definitions[name] != nil {
		return fmt.Errorf("%q %w", name, ErrAlreadyDefined)
	}

	n.definitions[name] = value

	return nil
}

func (n *node) applyDefinitions(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return n.setDefinition("_", value)
	}

	var errs []error

	for i := 0; i < len(value.Content); i += 2 {
		keyAny, err := nodeToAny(value.Content[i])
		if err != nil {
			errs = append(errs, err)
		}

		key := anyToString(keyAny)

		if err := n.setDefinition(key, value.Content[i+1]); err != nil {
			errs = append(errs, err)
		}
	}

	return joinErrors(errs...)
}

func (n *node) getDefinition(name string) (*yaml.Node, error) {
	if value, ok := n.definitions[name]; ok {
		return value, nil
	}

	if n.parent == nil {
		return nil, fmt.Errorf("%q is %w", name, ErrNotDefined)
	}

	return n.parent.getDefinition(name)
}

func (n *node) passDefinitionsToParent() error {
	var errs []error

	for name, value := range n.definitions {
		err := n.parent.setDefinition(name, value)
		if err != nil {
			errs = append(errs, n.locatedError(err))
		}
	}

	return joinErrors(errs...)
}
