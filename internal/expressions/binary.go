package expressions

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func binary(defs node.Definitions, expr *ast.BinaryExpr) (node.Node, error) {
	var errs []error

	left, err := Resolve(defs, expr.X)
	if err != nil {
		errs = append(errs, fmt.Errorf("left: %w", err))
	}

	right, err := Resolve(defs, expr.Y)
	if err != nil {
		errs = append(errs, fmt.Errorf("right: %w", err))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	if left == nil && right == nil {
		//nolint:nilnil // Nil as a value is okay here.
		return nil, nil
	}

	if left == nil {
		return right, nil
	}

	if right == nil {
		return left, nil
	}

	switch expr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return binMath(left, right, expr.Op)

	// case token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
	// 	return nil, fmt.Errorf("binary expressions %w", errors.ErrNotImplemented)

	case token.LAND, token.LOR:
		return node.Bool(logical(left, right, expr.Op)), nil

	case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return node.Bool(compare(left, right, expr.Op)), nil

	default:
		return nil, fmt.Errorf("%q: %w in binary expression", expr.Op.String(), errors.ErrUnsupportedOperator)
	}
}
