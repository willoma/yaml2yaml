package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) indexNode(defs node.Definitions, expr *ast.IndexExpr) (node.Node, error) {
	var errs errors.JoinedError

	value, err := s.resolveNode(defs, expr.X)
	if err != nil {
		errs = errs.Append(fmt.Errorf("expression: %w", err))
	}

	indexNode, err := s.resolveNode(defs, expr.Index)
	if err != nil {
		errs = errs.Append(fmt.Errorf("index: %w", err))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xSeq, ok := value.(node.Sequence)
	if !ok {
		errs = errs.Append(fmt.Errorf("expression: %w", errors.ErrNotASequence))
	}

	index, ok := indexNode.Int()
	if !ok {
		errs = errs.Append(fmt.Errorf("index: %w", errors.ErrNotANumber))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if len(xSeq) <= int(index) {
		return nil, errors.ErrOutOfBounds
	}

	return xSeq[index], nil
}

func (s *Source) indexSource(defs node.Definitions, expr *ast.IndexExpr) (node.Source, error) {
	var errs errors.JoinedError

	value, err := s.resolveSource(defs, expr.X)
	if err != nil {
		errs = errs.Append(fmt.Errorf("expression: %w", err))
	}

	indexNode, err := s.resolveNode(defs, expr.Index)
	if err != nil {
		errs = errs.Append(fmt.Errorf("index: %w", err))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xSeq, ok := value.(*node.SequenceSource)
	if !ok {
		errs = errs.Append(fmt.Errorf("expression: %w", errors.ErrNotASequence))
	}

	index, ok := indexNode.Int()
	if !ok {
		errs = errs.Append(fmt.Errorf("index: %w", errors.ErrNotANumber))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return xSeq.Get(int(index))
}
