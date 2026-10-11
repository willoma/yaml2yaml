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
	seq, err := f.processSeq(defs)
	if err != nil {
		return nil, fmt.Errorf("sequence: %w", err)
	}

	result := f.initialResult(seq)

	var errs errors.JoinedError

	for i, item := range seq.Value {
		childDefs := defsWithArgs(defs, item)

		childDefs.args["i"] = node.NewInt(f.sequence.Pos(), int64(i))

		childResult, err := f.value.Process(childDefs)
		if err != nil {
			errs = errs.Append(errors.LocateByIndex(i, f.value.Pos(), err))

			continue
		}

		result, err = result.Merge(childResult)
		if err != nil {
			errs = errs.Append(errors.LocateByIndex(i, f.value.Pos(), err))
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

func (f *dirFor) processSeq(defs node.Definitions) (*node.SequenceSource, error) {
	if seq, ok := f.sequence.(*node.SequenceSource); ok {
		return seq, nil
	}

	name, err := f.sequence.Process(defs)
	if err != nil {
		return nil, err
	}

	defNode := defs.Get(name.String())
	if defNode == nil {
		return nil, errors.ErrEmpty
	}

	seq, ok := defNode.(*node.SequenceSource)
	if !ok {
		return nil, errors.ErrNotASequence
	}

	return seq, nil
}

func (f *dirFor) initialResult(seq *node.SequenceSource) node.Node {
	switch value := f.value.(type) {
	case *node.SequenceSource:
		return make(node.Sequence, 0, len(seq.Value)*len(value.Value))
	case *node.MappingSource:
		return make(node.Mapping, 0, len(seq.Value)*len(value.Value))
	default:
		return make(node.Sequence, 0, len(seq.Value))
	}
}
