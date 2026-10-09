package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) index(defs node.Definitions, expr *ast.IndexExpr) (node.Node, error) {
	errs := errors.NewJoined()

	value, err := s.resolve(defs, expr.X)
	if err != nil {
		errs.Add(fmt.Errorf("expression: %w", err))
	}

	indexNode, err := s.resolve(defs, expr.Index)
	if err != nil {
		errs.Add(fmt.Errorf("index: %w", err))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xSeq, ok := value.(node.Sequence)
	if !ok {
		errs.Add(fmt.Errorf("expression: %w", errors.ErrNotASequence))
	}

	index, ok := indexNode.Int()
	if !ok {
		errs.Add(fmt.Errorf("index: %w", errors.ErrNotANumber))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if len(xSeq) <= int(index) {
		return nil, errors.ErrOutOfBounds
	}

	return xSeq[index], nil
}

func (s *Source) indexList(defs node.Definitions, expr *ast.IndexListExpr) (node.Node, error) {
	errs := errors.NewJoined()

	value, err := s.resolve(defs, expr.X)
	if err != nil {
		errs.Add(fmt.Errorf("expression: %w", err))
	}

	var indices []int

	for i, indexExpr := range expr.Indices {
		indexNode, err := s.resolve(defs, indexExpr)
		if err != nil {
			errs.Add(fmt.Errorf("index #%d: %w", i, err))

			continue
		}

		index, ok := indexNode.Int()
		if !ok {
			errs.Add(fmt.Errorf("index: %w", errors.ErrNotANumber))

			continue
		}

		indices = append(indices, int(index))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xSeq, ok := value.(node.Sequence)
	if !ok {
		return nil, fmt.Errorf("expression: %w", errors.ErrNotASequence)
	}

	result := make(node.Sequence, len(indices))

	for i, index := range indices {
		if len(xSeq) <= index {
			errs.Add(fmt.Errorf("%d: %w", index, errors.ErrOutOfBounds))
		}

		result[i] = xSeq[index]
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}
