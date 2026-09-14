package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/token"
)

func (n *node) resolveUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	value, vErr := n.resolveExpr(expr.X)

	switch expr.Op {
	case token.NOT:
		bVal, ok := value.(bool)
		if !ok {
			return false, joinErrors(vErr, ErrNotBoolean)
		}

		return !bVal, nil
	default:
		return nil, joinErrors(vErr, fmt.Errorf("operator %q %w", expr.Op.String(), ErrUnsupported))

	}
}
