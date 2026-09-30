package yaml2yaml

import (
	"slices"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func processMapping(defs *definitions.Definitions, node *yaml.Node) (*yaml.Node, error) {
	errs := errors.Multi()
	children := make([]mappingChild, 0, len(node.Content)/2)

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Kind == yaml.ScalarNode {
			thisChildren, err := processMappingChildWithScalarKey(defs, node.Content[i], node.Content[i+1])
			children = append(children, thisChildren...)

			errs.Add(err)

			continue
		}

		thisChild, err := processMappingChildWithNonScalarKey(defs, i/2, node.Content[i], node.Content[i+1])
		children = append(children, thisChild)

		errs.Add(err)
	}

	result, err := postprocessMappingChildren(defs, children)
	errs.Add(err)

	return result, errs
}

type mappingChild struct {
	pathKey string
	key     *yaml.Node
	value   *yaml.Node
}

func processMappingChildWithNonScalarKey(
	defs *definitions.Definitions, index int, key, value *yaml.Node,
) (mappingChild, error) {
	errs := errors.Multi()
	childKey := strconv.Itoa(index)

	keyResult, err := process(defs.Child(), childKey, key)

	errs.Add(err)

	valueResult, err := process(defs.Child(), childKey, value)

	errs.Add(err)

	return mappingChild{pathKey: childKey, key: keyResult, value: valueResult}, errs
}

func processMappingChildWithScalarKey(defs *definitions.Definitions, key, value *yaml.Node) ([]mappingChild, error) {
	result, processed, err := directiveMapping(defs, key, value)
	if processed {
		return result, err
	}

	errs := errors.Multi()

	keyResult, err := process(defs.Child(), key.Value, key)

	errs.Add(err)

	valueResult, err := process(defs.Child(), key.Value, value)

	errs.Add(err)

	return []mappingChild{{pathKey: key.Value, key: keyResult, value: valueResult}}, errs
}

func postprocessMappingChildren(defs *definitions.Definitions, children []mappingChild) (*yaml.Node, error) {
	switch len(children) {
	case 0:
		// No children: pass definitions to parent and return nil.
		return nil, defs.PassToParent()
	case 1:
		if children[0].value == nil {
			// Only one child, which has no value: pass definitions to parent and return nil.
			return nil, defs.PassToParent()
		}

		if children[0].key == nil {
			// Only one child, which is a single value: replace the whole mapping with the single evaluated value.
			return children[0].value, nil
		}

		return &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: []*yaml.Node{children[0].key, children[0].value},
		}, nil
	}

	allAreSequences := !slices.ContainsFunc(children, func(child mappingChild) bool {
		return child.key != nil || (child.value != nil && child.value.Kind != yaml.SequenceNode)
	})
	if allAreSequences {
		childrenNodes := []*yaml.Node{}
		for _, child := range children {
			childrenNodes = append(childrenNodes, child.value.Content...)
		}

		if len(childrenNodes) == 0 {
			//nolint:nilnil // Nil as a value is useful here.
			return nil, nil
		}

		return &yaml.Node{Kind: yaml.SequenceNode, Content: childrenNodes}, nil
	}

	errs := errors.Multi()

	result := &yaml.Node{Kind: yaml.MappingNode, Content: make([]*yaml.Node, 0, 2*len(children))}
	for _, child := range children {
		// When there is no value, skip the child.
		if child.value == nil {
			continue
		}

		// If the child has no key, it is a problem...
		if child.key == nil {
			errs.Add(errors.Locate(child.pathKey, child.value.Line, child.value.Column, errors.ErrUnkeyedValue))

			continue
		}

		if err := convert.MergeIntoMapping(result, child.key, child.value); err != nil {
			errs.Add(errors.Locate(child.pathKey, child.key.Line, child.key.Column, err))
		}
	}

	return result, errs
}
