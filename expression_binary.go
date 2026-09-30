package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func resolveBinaryExpr(defs *definitions.Definitions, expr *ast.BinaryExpr) (any, error) {
	errs := errors.Multi()

	left, err := resolveExpr(defs, expr.X)
	if err != nil {
		errs.Add(fmt.Errorf("left: %w", err))
	}

	right, err := resolveExpr(defs, expr.Y)
	if err != nil {
		errs.Add(fmt.Errorf("right: %w", err))
	}

	if !errs.IsNil() {
		return nil, errs
	}

	switch expr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return resolveMath(left, right, expr.Op)

	case token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
		return nil, fmt.Errorf("binary expressions %w", errors.ErrNotImplemented)

	case token.LAND, token.LOR:
		return resolveAndOr(left, right, expr.Op), nil

	case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return resolveComparison(left, right, expr.Op), nil

	default:
		return nil, fmt.Errorf("%q: %w in binary expression", expr.Op.String(), errors.ErrUnsupportedOperator)
	}
}

func resolveAndOr(left, right any, op token.Token) bool {
	bLeft := convert.AnyToBool(left)
	bRight := convert.AnyToBool(right)

	if op == token.LAND {
		return bLeft && bRight
	}

	return bLeft || bRight
}
