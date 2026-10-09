package expressions

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) basicLit(expr *ast.BasicLit) (node.Node, error) {
	switch expr.Kind {
	case token.INT:
		if value, err := strconv.ParseInt(expr.Value, 10, 64); err == nil {
			return node.NewInt(s.pos, value), nil
		}
	case token.FLOAT:
		if value, err := strconv.ParseFloat(expr.Value, 64); err == nil {
			return node.NewFloat(s.pos, value), nil
		}
	case token.CHAR, token.STRING:
		if value, err := strconv.Unquote(expr.Value); err == nil {
			return node.NewString(s.pos, value), nil
		}
	}

	return node.NewString(s.pos, expr.Value), nil
}
