package expressions

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) binaryNode(defs node.Definitions, expr *ast.BinaryExpr) (node.Node, error) {
	left, lErr := s.resolveNode(defs, expr.X)
	right, rErr := s.resolveNode(defs, expr.Y)

	switch {
	case lErr != nil && rErr != nil:
		return nil, fmt.Errorf("left: %w, right: %w", lErr, rErr)
	case lErr != nil:
		return nil, fmt.Errorf("left: %w", lErr)
	case rErr != nil:
		return nil, fmt.Errorf("right: %w", rErr)
	}

	switch expr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return s.binMath(left, right, expr.Op)

	case token.LAND, token.LOR:
		return s.binLogical(left, right, expr.Op), nil

	case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return node.NewBool(s.pos, s.binCompare(left, right, expr.Op)), nil

	default:
		return nil, fmt.Errorf("%q: %w in binary expression", expr.Op.String(), errors.ErrUnsupportedOperator)
	}
}

func (s *Source) binarySource(defs node.Definitions, expr *ast.BinaryExpr) (node.Source, error) {
	left, lErr := s.resolveSource(defs, expr.X)
	right, rErr := s.resolveSource(defs, expr.Y)

	switch {
	case lErr != nil && rErr != nil:
		return nil, fmt.Errorf("left: %w, right: %w", lErr, rErr)
	case lErr != nil:
		return nil, fmt.Errorf("left: %w", lErr)
	case rErr != nil:
		return nil, fmt.Errorf("right: %w", rErr)
	}

	leftMap, lOK := left.(*node.MappingSource)
	rightMap, rOK := right.(*node.MappingSource)

	switch {
	case !lOK && !rOK:
		return nil, fmt.Errorf("left: %w, right: %w", errors.ErrNotAMapping, errors.ErrNotAMapping)
	case !lOK:
		return nil, fmt.Errorf("left: %w", errors.ErrNotAMapping)
	case !rOK:
		return nil, fmt.Errorf("right: %w", errors.ErrNotAMapping)
	}

	// The only operation supported is adding two maps...
	if expr.Op != token.ADD {
		return nil, errors.ErrUnsupportedOperator
	}

	leftMap.Merge(rightMap)

	return leftMap, nil
}
