package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/token"
)

func (n *node) resolveBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	left, lErr := n.resolveExpr(expr.X)
	if lErr != nil {
		lErr = fmt.Errorf("left: %w", lErr)
	}

	right, rErr := n.resolveExpr(expr.Y)
	if rErr != nil {
		rErr = fmt.Errorf("right: %w", rErr)
	}

	if err := joinErrors(lErr, rErr); err != nil {
		return nil, err
	}

	var (
		result any
		opErr  error
	)

	switch expr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		result, lErr, opErr, rErr = resolveMath(left, right, expr.Op)

	case token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
		opErr = fmt.Errorf("binary expressions %w", ErrNotImplemented)

	case token.LAND, token.LOR:
		result, lErr, opErr, rErr = resolveAndOr(left, right, expr.Op)

	case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		result = resolveComparison(left, right, expr.Op)

	default:
		opErr = fmt.Errorf("operator %q %w", expr.Op.String(), ErrUnsupported)

	}

	return result, joinErrors(lErr, opErr, rErr)
}

func resolveAndOr(left, right any, op token.Token) (result any, lErr, opErr, rErr error) {
	bLeft, ok := left.(bool)
	if !ok {
		lErr = fmt.Errorf("left: %w", ErrNotBoolean)
	}

	bRight, ok := right.(bool)
	if !ok {
		rErr = fmt.Errorf("right: %w", ErrNotBoolean)
	}

	if op == token.LAND {
		result = bLeft && bRight
	} else {
		result = bLeft || bRight
	}

	return result, lErr, opErr, rErr
}
