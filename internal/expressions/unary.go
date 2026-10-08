package expressions

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func unary(defs node.Definitions, expr *ast.UnaryExpr) (node.Node, error) {
	value, err := resolve(defs, expr.X)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case token.NOT:
		ok, _ := value.Bool()

		return node.Bool(!ok), nil
	default:
		return nil, fmt.Errorf("%q: %w in unary expression", expr.Op.String(), errors.ErrUnsupportedOperator)
	}
}
