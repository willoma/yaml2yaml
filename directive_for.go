package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func doDirFor(defs *definitions.Definitions, key, value *yaml.Node) (*yaml.Node, error) {
	name, err := evaluateString(defs, key)

	defNode := defs.Get(name)
	if defNode == nil {
		return nil, err
	}

	errs := errors.Multi()

	errs.Add(err)

	if defNode.Kind != yaml.SequenceNode {
		errs.Add(fmt.Errorf(`expression: %w`, errors.ErrNotASequence))
	}

	if value == nil || value.Kind != yaml.SequenceNode && value.Kind != yaml.MappingNode {
		errs.Add(fmt.Errorf(`content: %w`, errors.ErrNotAMappingOrSequence))
	}

	if !errs.IsNil() {
		return nil, errs
	}

	result := &yaml.Node{
		Kind:    value.Kind,
		Tag:     value.Tag,
		Content: make([]*yaml.Node, 0, len(defNode.Content)*len(value.Content)),
		Line:    value.Line,
		Column:  value.Column,
	}

	for i, seqNode := range defNode.Content {
		childDefs := defs.Child()
		idx := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(i)}

		//nolint:errcheck // child has just been created, i cannot already be defined.
		childDefs.Set("i", idx)
		errs.Add(childDefs.SetFromNode(seqNode))

		childResult, err := process(childDefs, "["+strconv.Itoa(i)+"]", value)

		errs.Add(err)

		result.Content = append(result.Content, childResult.Content...)
	}

	return result, errs
}
