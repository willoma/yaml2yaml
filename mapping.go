package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"
)

func (n *node) mappingLength() int {
	return len(n.Content) / 2
}

func (n *node) processMapping() (*yaml.Node, error) {
	var errs []error

	children := make([]mappingChild, 0, n.mappingLength())

	for i := 0; i < len(n.Content); i += 2 {
		keyNode := n.Content[i]

		keyValue := keyNode.Value
		if n.Content[i].Kind != yaml.ScalarNode {
			keyValue = strconv.Itoa(i / 2)
		}

		keyedKey := n.keyedChild(keyNode.Tag, keyValue, keyNode)

		result, err := n.evaluateMappingChild(
			keyedKey,
			n.keyedChild(keyNode.Tag, keyValue, n.Content[i+1]),
		)
		if err != nil {
			errs = append(errs, err)
		}

		result.path = keyedKey.path

		children = append(children, result)
	}

	if len(children) == 1 && children[0].value != nil {
		// Only one child, which is a single value: replace the whole mapping with the single evaluated value.
		return children[0].value, nil
	}

	childrenNodes := make([]*yaml.Node, 0, 2*len(children))
	for _, child := range children {
		if child.value != nil {
			errs = append(errs, child.valueLocatedError(fmt.Errorf("%w: single value as mapping entry", ErrInvalid)))
		}

		childrenNodes = append(childrenNodes, child.kvs...)
	}

	if len(childrenNodes) == 0 {
		errs = append(errs, n.passDefinitionsToParent())
	}

	return &yaml.Node{Kind: yaml.MappingNode, Content: childrenNodes}, joinErrors(errs...)
}

func (n *node) toMappingChild() (mappingChild, error) {
	resultValue, vErr := n.process()
	if resultValue.Kind == yaml.MappingNode {
		return mappingChild{kvs: resultValue.Content}, vErr
	}

	return mappingChild{value: resultValue}, vErr
}
