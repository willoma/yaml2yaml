package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) indexList(defs node.Definitions, expr *ast.IndexListExpr) (node.Node, error) {
	var errs errors.JoinedError

	value, err := s.resolveNode(defs, expr.X)
	if err != nil {
		errs = errs.Append(fmt.Errorf("expression: %w", err))
	}

	var indices []int

	for i, indexExpr := range expr.Indices {
		indexNode, err := s.resolveNode(defs, indexExpr)
		if err != nil {
			errs = errs.Append(fmt.Errorf("index #%d: %w", i, err))

			continue
		}

		index, ok := indexNode.Int()
		if !ok {
			errs = errs.Append(fmt.Errorf("index: %w", errors.ErrNotANumber))

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
			errs = errs.Append(fmt.Errorf("%d: %w", index, errors.ErrOutOfBounds))

			continue
		}

		result[i] = xSeq[index]
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}
