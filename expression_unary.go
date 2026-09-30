package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func resolveUnaryExpr(defs *definitions.Definitions, expr *ast.UnaryExpr) (any, error) {
	errs := errors.Multi()

	value, err := resolveExpr(defs, expr.X)

	errs.Add(err)

	switch expr.Op {
	case token.NOT:
		return !convert.AnyToBool(value), errs
	default:
		errs.Add(fmt.Errorf("%q: %w in unary expression", expr.Op.String(), errors.ErrUnsupportedOperator))

		return nil, errs
	}
}
