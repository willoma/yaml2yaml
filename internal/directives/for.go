package directives

import (
	"fmt"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameFor = "!for"

type For struct {
	Sequence node.Source
	Value    node.Source
}

func (f *For) Process(defs node.Definitions) (node.Node, error) {
	var errs []error

	seq, ok := f.Sequence.(node.SequenceSource)
	if !ok {
		seqNode, err := f.Sequence.Process(defs)
		if err != nil {
			errs = append(errs, err)
		}

		defNode := defs.Get(seqNode.String())
		if defNode == nil {
			return nil, errors.Join(errs...)
		}

		seq, ok = defNode.(node.SequenceSource)
		if !ok {
			errs = append(errs, fmt.Errorf(`expression: %w`, errors.ErrNotASequence))

			return nil, errors.Join(errs...)
		}
	}

	result, err := f.process(defs, seq)
	if err != nil {
		errs = append(errs, err)
	}

	return result, errors.Join(errs...)
}

func (f *For) process(defs node.Definitions, seq node.SequenceSource) (node.Node, error) {
	var result node.Node

	switch value := f.Value.(type) {
	case node.SequenceSource:
		result = make(node.Sequence, 0, len(seq)*len(value))
	case node.MappingSource:
		result = make(node.Mapping, 0, len(seq)*len(value))
	default:
		return nil, fmt.Errorf(`content: %w`, errors.ErrNotAMappingOrSequence)
	}

	var errs []error

	for i, item := range seq {
		childDefs := defsWithArgs(defs, item)

		childDefs.Set("i", node.Int(i))

		childResult, err := f.Value.Process(childDefs)
		if err != nil {
			errs = append(errs, err)
		}

		result, err = result.Merge(childResult)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return result, errors.Join(errs...)
}
