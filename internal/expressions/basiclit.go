package expressions

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/willoma/yaml2yaml/internal/node"
)

func basicLit(expr *ast.BasicLit) (node.Node, error) {
	switch expr.Kind {
	case token.INT:
		if value, err := strconv.ParseInt(expr.Value, 10, 64); err == nil {
			return node.Int(value), nil
		}
	case token.FLOAT:
		if value, err := strconv.ParseFloat(expr.Value, 64); err == nil {
			return node.Float(value), nil
		}
	case token.CHAR, token.STRING:
		if value, err := strconv.Unquote(expr.Value); err == nil {
			return node.String(value), nil
		}
	}

	return node.String(expr.Value), nil
}
