package expressions

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) resolve(defs node.Definitions, expr ast.Expr) (node.Node, error) {
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return nil, errors.ErrSyntaxError

	case *ast.Ident:
		return s.ident(defs, expr)

	case *ast.BasicLit:
		return s.basicLit(expr), nil

	case *ast.ParenExpr:
		return s.resolve(defs, expr.X)

	case *ast.SelectorExpr:
		return s.selector(defs, expr)

	case *ast.IndexExpr:
		return s.index(defs, expr)

	case *ast.IndexListExpr:
		return s.indexList(defs, expr)

	case *ast.SliceExpr:
		return s.slice(defs, expr)

	case *ast.CallExpr:
		return s.call(defs, expr)

	case *ast.UnaryExpr:
		return s.unary(defs, expr)

	case *ast.BinaryExpr:
		return s.binary(defs, expr)

	case nil:
		return nil, errors.ErrEmpty

	default:
		return nil, errors.ErrUnsupportedExpression
	}
}
