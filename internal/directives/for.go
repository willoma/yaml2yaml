package directives

import (
	"errors"
	"fmt"

	yerrors "github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameFor = "!for"

type dirFor struct {
	sequence node.Source
	value    node.Source
}

func For(sequence, value node.Source) *dirFor {
	return &dirFor{sequence: sequence, value: value}
}

func (f *dirFor) Process(defs node.Definitions) (node.Node, error) {
	seq, ok := f.sequence.(node.SequenceSource)
	if !ok {
		seqNode, err := f.sequence.Process(defs)
		if err != nil {
			return nil, fmt.Errorf("sequence: %w", err)
		}

		defNode := defs.Get(seqNode.String())
		if defNode == nil {
			return nil, fmt.Errorf("sequence: %w", yerrors.ErrEmpty)
		}

		seq, ok = defNode.(node.SequenceSource)
		if !ok {
			return nil, fmt.Errorf(`sequence: %w`, yerrors.ErrNotASequence)
		}
	}

	var result node.Node

	switch value := f.value.(type) {
	case node.SequenceSource:
		result = make(node.Sequence, 0, len(seq)*len(value))
	case node.MappingSource:
		result = make(node.Mapping, 0, len(seq)*len(value))
	default:
		return nil, fmt.Errorf(`content: %w`, yerrors.ErrNotAMappingOrSequence)
	}

	var errs []error

	for i, item := range seq {
		childDefs := defsWithArgs(defs, item)

		childDefs.Set("i", node.Int(i))

		childResult, err := f.value.Process(childDefs)
		if err != nil {
			if errors.Is(err, yerrors.ErrEmpty) {
				continue
			}
			// TODO Locate error
			errs = append(errs, err)

			continue
		}

		result, err = result.Merge(childResult)
		if err != nil {
			// TODO Locate error
			errs = append(errs, err)
		}
	}

	return result, errors.Join(errs...)
}
