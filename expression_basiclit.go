package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func (n *node) resolveBasicLit(expr *ast.BasicLit) (any, error) {
	switch expr.Kind {
	case token.INT:
		return strconv.Atoi(expr.Value)
	case token.FLOAT:
		return strconv.ParseFloat(expr.Value, 64)
	case token.IMAG:
		return strconv.ParseComplex(expr.Value, 128)
	case token.CHAR:
		return strconv.Unquote(expr.Value)
	case token.STRING:
		return strconv.Unquote(expr.Value)
	default:
		return nil, fmt.Errorf("basic literal kind %s %w", expr.Kind.String(), ErrUnsupported)
	}
}
