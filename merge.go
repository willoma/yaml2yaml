package yaml2yaml

import (
	"fmt"
	"slices"

	"go.yaml.in/yaml/v4"
)

func mergeNodes(left, right *yaml.Node) (*yaml.Node, error) {
	if right == nil {
		return left, nil
	}

	if left == nil {
		return right, nil
	}

	switch left.Kind {
	case yaml.DocumentNode:
		return mergeNodes(left.Content[0], right)
	case yaml.AliasNode:
		return mergeNodes(left.Alias, right)
	case yaml.SequenceNode:
		return mergeSequenceNodes(left, right)
	case yaml.MappingNode:
		return mergeMappingNodes(left, right)
	case yaml.ScalarNode:
		return mergeScalarNodes(left, right)
	case yaml.StreamNode:
		return right, fmt.Errorf("merge left: %w kind stream", ErrUnsupported)
	default:
		return right, fmt.Errorf("merge left: %w unnamed kind", ErrUnsupported)
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
		return left, fmt.Errorf("%w mapping into sequence", ErrCannotMerge)

	case yaml.ScalarNode:
		return left, fmt.Errorf("%w scalar into sequence", ErrCannotMerge)

	case yaml.StreamNode:
		return left, fmt.Errorf("merge right: %w kind stream", ErrUnsupported)

	default:
		return left, fmt.Errorf("merge right: %w unnamed kind", ErrUnsupported)
	}
}

func mergeMappingNodes(left, right *yaml.Node) (*yaml.Node, error) {
	switch right.Kind {
	case yaml.DocumentNode:
		return mergeMappingNodes(left, right.Content[0])

	case yaml.AliasNode:
		return mergeMappingNodes(left, right.Alias)

	case yaml.SequenceNode:
		return left, fmt.Errorf("%w sequence into mapping", ErrCannotMerge)

	case yaml.MappingNode:
		content, err := mergeMappingContent(left.Content, right.Content)

		return &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: content,
		}, err

	case yaml.ScalarNode:
		return left, fmt.Errorf("%w scalar into mapping", ErrCannotMerge)

	case yaml.StreamNode:
		return left, fmt.Errorf("merge right: %w kind stream", ErrUnsupported)

	default:
		return left, fmt.Errorf("merge right: %w unnamed kind", ErrUnsupported)
	}
}

func mergeScalarNodes(left, right *yaml.Node) (*yaml.Node, error) {
	switch right.Kind {
	case yaml.DocumentNode:
		return mergeScalarNodes(left, right.Content[0])

	case yaml.AliasNode:
		return mergeScalarNodes(left, right.Alias)

	case yaml.SequenceNode:
		return left, fmt.Errorf("%w sequence into scalar", ErrCannotMerge)

	case yaml.MappingNode:
		return left, fmt.Errorf("%w mapping into sequence", ErrCannotMerge)

	case yaml.ScalarNode:
		// Only replace the scalar by its new value...
		return right, nil

	case yaml.StreamNode:
		return left, fmt.Errorf("merge right: %w kind stream", ErrUnsupported)

	default:
		return left, fmt.Errorf("merge right: %w unnamed kind", ErrUnsupported)
	}
}

func mergeMappingContent(left, right []*yaml.Node) ([]*yaml.Node, error) {
	var (
		err  error
		errs []error
	)

	result := make([]*yaml.Node, 0, len(left))
	remainingCandidates := make([]*yaml.Node, 0, len(right))

	for i := 0; i < len(left); i += 2 {
		value := left[i+1]

		for j := 0; j < len(right); j += 2 {
			if nodesEqual(left[i], right[j]) {
				value, err = mergeNodes(value, right[j+1])
				if err != nil {
					keyAny, kErr := nodeToAny(value.Content[i])
					if kErr != nil {
						errs = append(errs, kErr)
					}

					errs = append(errs, fmt.Errorf("%s: %w", anyToString(keyAny), err))
				}
			} else {
				remainingCandidates = append(remainingCandidates, right[j], right[j+1])
			}
		}

		result = append(result, left[i], value)
		right = remainingCandidates
		remainingCandidates = remainingCandidates[:0]
	}

	result = append(result, right...)

	return result, joinErrors(errs...)
}

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
