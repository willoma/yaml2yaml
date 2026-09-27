package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

func (n *node) processMapping() (*yaml.Node, error) {
	var errs []error

	children := make([]mappingChild, 0, len(n.Content)/2)

	for i := 0; i < len(n.Content); i += 2 {
		keyNode := n.Content[i]

		keyValue := keyNode.Value
		if n.Content[i].Kind != yaml.ScalarNode {
			keyValue = strconv.Itoa(i / 2)
		}

		key := n.keyedChild(keyNode.Tag, keyValue, keyNode)
		value := n.keyedChild(keyNode.Tag, keyValue, n.Content[i+1])

		if key.Kind == yaml.ScalarNode {
			thisChildren, err := n.mappingDirective(key, value)
			if err != nil {
				errs = append(errs, err)
			}

			children = append(children, thisChildren...)

			continue
		}

		keyResult, kErr := key.process()
		if kErr != nil {
			errs = append(errs, kErr)
		}

		valueResult, vErr := value.process()
		if vErr != nil {
			errs = append(errs, vErr)
		}

		children = append(children, mappingChild{path: key.path, key: keyResult, value: valueResult})
	}

	result, err := n.processMappingChildren(children)
	if err != nil {
		errs = append(errs, err)
	}

	return result, errors.Join(errs...)
}

func (n *node) processMappingChildren(children []mappingChild) (*yaml.Node, error) {
	var errs []error

	switch len(children) {
	case 0:
		// No children: pass definitions to parent and return nil.
		return nil, n.definitions.PassToParent()

	case 1:
		if children[0].key == nil {
			// Only one child, which is a single value: replace the whole mapping with the single evaluated value.
			return children[0].value, nil
		}

		return &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: []*yaml.Node{children[0].key, children[0].value},
		}, errors.Join(errs...)
	}

	if children[0].key == nil {
		// Hopefully, the children are all sequence nodes: flatten them into a single sequence node.
		childrenNodes := make([]*yaml.Node, 0, len(children))
		for _, child := range children {
			if child.key != nil || child.value.Kind != yaml.SequenceNode {
				errs = append(errs, child.locatedError(fmt.Errorf("%w: unkeyed value as mapping entry", ErrInvalid)))
			}

			childrenNodes = append(childrenNodes, child.value.Content...)
		}

		if len(childrenNodes) == 0 {
			return nil, errors.Join(errs...)
		}

		return &yaml.Node{Kind: yaml.SequenceNode, Content: childrenNodes}, errors.Join(errs...)
	}

	childrenNodes := make([]*yaml.Node, 0, 2*len(children))
	for _, child := range children {
		if child.key == nil {
			errs = append(errs, child.locatedError(fmt.Errorf("%w: single value as mapping entry", ErrInvalid)))
		}

		childrenNodes = append(childrenNodes, child.key, child.value)
	}

	return &yaml.Node{Kind: yaml.MappingNode, Content: childrenNodes}, errors.Join(errs...)
}
