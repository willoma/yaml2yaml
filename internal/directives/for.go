package directives

import (
	oerrors "errors"
	"fmt"
	"maps"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameFor = "!for"

type dirFor struct {
	sequence *expressions.Source
	value    node.Source
}

func For(sequence, value node.Source) (*dirFor, error) {
	seqExpr, err := expressions.Parse(sequence.Pos(), sequence.String())
	if err != nil {
		return nil, err
	}

	return &dirFor{sequence: seqExpr, value: value}, nil
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
	seq, err := f.getSeq(defs)
	if err != nil {
		return nil, fmt.Errorf("sequence: %w", err)
	}

	processed := node.NewProcessedEntries(seq.Size())

	for i, item := range seq.Iter() {
		childDefs, err := defsWithArgs(defs, item)
		if err != nil {
			processed.AppendIndexed(i, item.Pos(), nil, err)

			continue
		}

		childDefs.args["i"] = node.NewInt(f.sequence.Pos(), int64(i))

		childResult, err := f.value.Process(childDefs)
		if oerrors.Is(err, errors.ErrEmpty) {
			childDefs.passToParent()
		}

		processed.AppendIndexed(i, item.Pos(), childResult, err)
	}

	return processed.Result()
}

func (f *dirFor) SourceMap(defs node.Definitions) (map[string]node.Source, error) {
	seq, err := f.getSeq(defs)
	if err != nil {
		return nil, fmt.Errorf("sequence: %w", err)
	}

	var errs errors.JoinedError

	result := make(map[string]node.Source)

	for i, item := range seq.Iter() {
		childDefs, err := defsWithArgs(defs, item)
		if err != nil {
			errs = append(errs, errors.LocateByIndex(i, item.Pos(), err))

			continue
		}

		childDefs.args["i"] = node.NewInt(f.sequence.Pos(), int64(i))

		sourceMapper, ok := f.value.(node.SourceMapper)
		if !ok {
			errs = append(errs, errors.LocateByIndex(i, item.Pos(), errors.ErrNotAMapping))
		}

		mapSrc, err := sourceMapper.SourceMap(defs)
		if err != nil {
			errs = append(errs, errors.LocateByIndex(i, item.Pos(), err))

			continue
		}

		maps.Copy(result, mapSrc)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}

func (f *dirFor) getSeq(defs node.Definitions) (*node.SequenceSource, error) {
	src, err := f.sequence.Source(defs)
	if err != nil {
		return nil, err
	}

	seq, ok := src.(*node.SequenceSource)
	if !ok {
		return nil, errors.ErrNotASequence
	}

	return seq, nil
}
