package convert

import (
	"slices"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

// MergeNodes merges two YAML nodes recursively, if they are compatible.
// It does not return errors on [yaml.StreamNode] or default kinds,
// because it runs on processed YAML nodes only
// (ie. those kinds have already been excluded).
func MergeNodes(left, right *yaml.Node) (*yaml.Node, error) {
	if right == nil {
		return left, nil
	}

	if left == nil {
		return right, nil
	}

	switch left.Kind {
	case yaml.DocumentNode:
		return MergeNodes(left.Content[0], right)
	case yaml.AliasNode:
		return MergeNodes(left.Alias, right)
	case yaml.SequenceNode:
		return mergeSequenceNodes(left, right)
	case yaml.MappingNode:
		return mergeMappingNodes(left, right)
	case yaml.ScalarNode:
		return mergeScalarNodes(left, right)
	default:
		return right, nil
	}
}

func mergeSequenceNodes(left, right *yaml.Node) (*yaml.Node, error) {
	switch right.Kind {
	case yaml.DocumentNode:
		return mergeSequenceNodes(left, right.Content[0])

	case yaml.AliasNode:
		return mergeSequenceNodes(left, right.Alias)

	case yaml.SequenceNode:
		content := slices.Clone(left.Content)
		content = append(content, right.Content...)

		return &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: content,
		}, nil

	case yaml.MappingNode:
		return left, errors.ErrCannotMergeMappingIntoSequence

	case yaml.ScalarNode:
		content := slices.Clone(left.Content)
		content = append(content, right)

		return &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: content,
		}, nil

	default:
		return left, nil
	}
}

func mergeMappingNodes(left, right *yaml.Node) (*yaml.Node, error) {
	switch right.Kind {
	case yaml.DocumentNode:
		return mergeMappingNodes(left, right.Content[0])

	case yaml.AliasNode:
		return mergeMappingNodes(left, right.Alias)

	case yaml.SequenceNode:
		return left, errors.ErrCannotMergeSequenceIntoMapping

	case yaml.MappingNode:
		return mergeMappingContent(left.Content, right.Content)

	case yaml.ScalarNode:
		return left, errors.ErrCannotMergeScalarIntoMapping

	default:
		return left, nil
	}
}

func mergeScalarNodes(left, right *yaml.Node) (*yaml.Node, error) {
	switch right.Kind {
	case yaml.DocumentNode:
		return mergeScalarNodes(left, right.Content[0])

	case yaml.AliasNode:
		return mergeScalarNodes(left, right.Alias)

	case yaml.SequenceNode:
		return left, errors.ErrCannotMergeSequenceIntoScalar

	case yaml.MappingNode:
		return left, errors.ErrCannotMergeMappingIntoSequence

	case yaml.ScalarNode:
		// Only replace the scalar with its new value...
		return right, nil

	default:
		return left, nil
	}
}

func mergeMappingContent(left, right []*yaml.Node) (*yaml.Node, error) {
	errs := errors.Multi()

	result := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: slices.Clone(left),
	}

	for i := 0; i < len(right); i += 2 {
		key := right[i]
		value := right[i+1]

		if err := MergeIntoMapping(result, key, value); err != nil {
			errs.Add(errors.Locate(NodeToString(key), key.Line, key.Column, err))
		}
	}

	return result, errs
}
