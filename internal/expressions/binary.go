package expressions

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) binary(defs node.Definitions, expr *ast.BinaryExpr) (node.Node, error) {
	left, lErr := s.resolve(defs, expr.X)
	right, rErr := s.resolve(defs, expr.Y)

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

	// case token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
	// 	return nil, fmt.Errorf("binary expressions %w", errors.ErrNotImplemented)

	case token.LAND, token.LOR:
		return s.logical(left, right, expr.Op), nil

	case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return node.NewBool(s.pos, s.compare(left, right, expr.Op)), nil

	default:
		return nil, fmt.Errorf("%q: %w in binary expression", expr.Op.String(), errors.ErrUnsupportedOperator)
	}
}
