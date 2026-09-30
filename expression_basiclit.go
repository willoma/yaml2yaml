package yaml2yaml

import (
	"go/ast"
	"go/token"
	"strconv"
)

//nolint:wrapcheck // Original errors are ok.
func resolveBasicLit(expr *ast.BasicLit) (any, error) {
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
		//nolint:nilnil // This will never happen, but we need to return something.
		return nil, nil
	}
}
