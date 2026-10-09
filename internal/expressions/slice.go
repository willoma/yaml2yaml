package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

//nolint:funlen // It is sequential, easy to understand.
func (s *Source) slice(defs node.Definitions, expr *ast.SliceExpr) (node.Node, error) {
	errs := errors.NewJoined()

	value, err := s.resolve(defs, expr.X)
	if err != nil {
		errs.Add(fmt.Errorf("expression: %w", err))
	}

	startNode, err := s.resolve(defs, expr.Low)
	if err != nil {
		errs.Add(fmt.Errorf("start: %w", err))
	}

	endNode, err := s.resolve(defs, expr.High)
	if err != nil {
		errs.Add(fmt.Errorf("end: %w", err))
	}

	if expr.Slice3 {
		errs.Add(fmt.Errorf("3-index slice: %w", errors.ErrUnsupportedExpression))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xSeq, ok := value.(node.Sequence)
	if !ok {
		errs.Add(fmt.Errorf("expression: %w", errors.ErrNotASequence))
	}

	start, ok := startNode.Int()
	if !ok {
		errs.Add(fmt.Errorf("start: %w", errors.ErrNotANumber))
	}

	end, ok := endNode.Int()
	if !ok {
		errs.Add(fmt.Errorf("end: %w", errors.ErrNotANumber))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if len(xSeq) <= int(start) {
		errs.Add(fmt.Errorf("start: %w", errors.ErrOutOfBounds))
	}

	if len(xSeq) < int(end) {
		errs.Add(fmt.Errorf("end: %w", errors.ErrOutOfBounds))
	}

	if start >= end {
		errs.Add(fmt.Errorf("%w: start is not before end", errors.ErrInvalidExpression))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return xSeq[start:end], nil
}
