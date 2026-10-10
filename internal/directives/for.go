package directives

import (
	"fmt"

	"github.com/willoma/yaml2yaml/errors"
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

func (f *dirFor) Pos() [2]int {
	return f.sequence.Pos()
}

func (f *dirFor) String() string {
	return NameFor + " " + f.sequence.String()
}

func (f *dirFor) Process(defs node.Definitions) (node.Node, error) {
	seq, ok := f.sequence.(*node.SequenceSource)
	if !ok {
		seqNode, err := f.sequence.Process(defs)
		if err != nil {
			return nil, fmt.Errorf("sequence: %w", err)
		}

		defNode := defs.Get(seqNode.String())
		if defNode == nil {
			return nil, fmt.Errorf("sequence: %w", errors.ErrEmpty)
		}

		seq, ok = defNode.(*node.SequenceSource)
		if !ok {
			return nil, fmt.Errorf(`sequence: %w`, errors.ErrNotASequence)
		}
	}

	var result node.Node

	switch value := f.value.(type) {
	case *node.SequenceSource:
		result = make(node.Sequence, 0, len(seq.Value)*len(value.Value))
	case *node.MappingSource:
		result = make(node.Mapping, 0, len(seq.Value)*len(value.Value))
	default:
		return nil, fmt.Errorf(`content: %w`, errors.ErrNotAMappingOrSequence)
	}

	errs := errors.NewJoined()

	for i, item := range seq.Value {
		childDefs := defsWithArgs(defs, item)

		childDefs.Set("i", node.NewInt(f.sequence.Pos(), int64(i)))

		childResult, err := f.value.Process(childDefs)
		if err != nil {
			errs.Add(errors.LocateByIndex(i, f.value.Pos(), err))

			continue
		}

		result = append(result, childResult)
	}

	return result, errs.Compact()
}
