package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

func (n *node) forMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	seq := n.definitions.Get(name)
	if seq == nil {
		return nil, nil
	}

	if seq.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%w: !for expression must be a sequence", ErrInvalid)
	}

	if value.Kind != yaml.SequenceNode && value.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%w: !for content must be a sequence or a mapping", ErrInvalid)
	}

	errs := []error{kErr}

	result := &yaml.Node{
		Kind:    value.Kind,
		Tag:     value.Tag,
		Content: make([]*yaml.Node, 0, len(seq.Content)*len(value.Content)),
		Line:    value.Line,
		Column:  value.Column,
	}

	for i, seqNode := range seq.Content {
		child := n.indexedChild(i, value.Node)
		idx := scalarNode(seqNode, "!!int", strconv.Itoa(i))

		//nolint:errcheck // child has just been created, i is not already defined
		child.definitions.Set("i", idx)

		if err := child.definitions.SetFromNode(seqNode); err != nil {
			errs = append(errs, err)
		}

		childResult, err := child.process()
		if err != nil {
			errs = append(errs, err)
		}

		result.Content = append(result.Content, childResult.Content...)
	}

	children, cErr := nodeToMappingChild(name, result)
	errs = append(errs, cErr)

	return children, errors.Join(errs...)
}
