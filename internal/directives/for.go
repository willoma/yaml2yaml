package directives

import (
	oerrors "errors"
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

func (f *dirFor) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirFor); ok {
		return f.sequence.EqualsSource(other.sequence) && f.value.EqualsSource(other.value)
	}

	return false
}

func (f *dirFor) Process(defs node.Definitions) (node.Node, error) {
	seq, err := f.processSeq(defs)
	if err != nil {
		return nil, fmt.Errorf("sequence: %w", err)
	}

	processed := node.NewProcessedEntries(len(seq.Value))

	for i, item := range seq.Value {
		childDefs := defsWithArgs(defs, item)

		childDefs.args["i"] = node.NewInt(f.sequence.Pos(), int64(i))

		childResult, err := f.value.Process(childDefs)
		if oerrors.Is(err, errors.ErrEmpty) {
			childDefs.passToParent()
		}

		processed.AppendIndexed(i, f.value.Pos(), childResult, err)
	}

	return processed.Result()
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
