package expressions

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/willoma/yaml2yaml/internal/node"
)

type nodeAndSource interface {
	node.Node
	node.Source
}

//nolint:ireturn // This is only to easily use this for getting [node.Node] or [node.Source].
func (s *Source) basicLit(expr *ast.BasicLit) nodeAndSource {
	switch expr.Kind {
	case token.INT:
		if value, err := strconv.ParseInt(expr.Value, 10, 64); err == nil {
			return node.NewInt(s.pos, value)
		}
	case token.FLOAT:
		if value, err := strconv.ParseFloat(expr.Value, 64); err == nil {
			return node.NewFloat(s.pos, value)
		}
	case token.CHAR, token.STRING:
		if value, err := strconv.Unquote(expr.Value); err == nil {
			return node.NewString(s.pos, value)
		}
	}

	return node.NewString(s.pos, expr.Value)
}
