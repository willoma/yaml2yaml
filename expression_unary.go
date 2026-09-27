package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/errors"
)

func (n *node) resolveUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	value, vErr := n.resolveExpr(expr.X)

	switch expr.Op {
	case token.NOT:
		bVal, ok := value.(bool)
		if !ok {
			return false, errors.Join(vErr, ErrNotBoolean)
		}

		return !bVal, nil
	default:
		return nil, errors.Join(vErr, fmt.Errorf("operator %q %w", expr.Op.String(), ErrUnsupported))
	}
}
