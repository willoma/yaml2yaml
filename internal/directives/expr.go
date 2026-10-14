package directives

import (
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameExpr = "!expr"

func Expr(expr node.Source) (*expressions.Source, error) {
	return expressions.Parse(expr.Pos(), expr.String())
}
