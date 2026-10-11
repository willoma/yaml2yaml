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

	switch value := f.value.(type) {
	case *node.SequenceSource:
		return f.process(defs, seq, make(node.Sequence, 0, len(seq.Value)*len(value.Value)))
	case *node.MappingSource:
		return f.process(defs, seq, make(node.Mapping, 0, len(seq.Value)*len(value.Value)))
	default:
		return f.process(defs, seq, make(node.Sequence, 0, len(seq.Value)))
	}
}

func (f *dirFor) process(defs node.Definitions, seq *node.SequenceSource, initial node.Node) (node.Node, error) {
	result := initial

	errs := errors.NewJoined()

	for i, item := range seq.Value {
		childDefs := defsWithArgs(defs, item)

		childDefs.args["i"] = node.NewInt(f.sequence.Pos(), int64(i))

		childResult, err := f.value.Process(childDefs)
		if err != nil {
			errs.Add(errors.LocateByIndex(i, f.value.Pos(), err))

			continue
		}

		result, err = result.Merge(childResult)
		if err != nil {
			errs.Add(errors.LocateByIndex(i, f.value.Pos(), err))
		}
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	switch result := result.(type) {
	case node.Mapping:
		if len(result) == 0 {
			return nil, errors.ErrEmpty
		}
	case node.Sequence:
		if len(result) == 0 {
			return nil, errors.ErrEmpty
		}
	}

	return result, nil
}
