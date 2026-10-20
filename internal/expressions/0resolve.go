package expressions

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) resolveNode(defs node.Definitions, expr ast.Expr) (node.Node, error) {
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return nil, errors.ErrSyntaxError

	case *ast.Ident:
		return s.identNode(defs, expr)

	case *ast.BasicLit:
		return s.basicLit(expr), nil

	case *ast.ParenExpr:
		return s.resolveNode(defs, expr.X)

	case *ast.SelectorExpr:
		return s.selectorNode(defs, expr)

	case *ast.IndexExpr:
		return s.indexNode(defs, expr)

	case *ast.IndexListExpr:
		return s.indexList(defs, expr)

	case *ast.SliceExpr:
		return s.slice(defs, expr)

	case *ast.CallExpr:
		return s.callNode(defs, expr)

	case *ast.UnaryExpr:
		return s.unary(defs, expr)

	case *ast.BinaryExpr:
		return s.binaryNode(defs, expr)

	case nil:
		return nil, errors.ErrEmpty

	default:
		return nil, errors.ErrUnsupportedExpression
	}
}

func (s *Source) resolveSource(defs node.Definitions, expr ast.Expr) (node.Source, error) {
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return nil, errors.ErrSyntaxError

	case *ast.Ident:
		return s.ident(defs, expr)

	case *ast.BasicLit:
		return s.basicLit(expr), nil

	case *ast.ParenExpr:
		return s.resolveSource(defs, expr.X)

	case *ast.SelectorExpr:
		return s.selectorSource(defs, expr)

	case *ast.IndexExpr:
		return s.indexSource(defs, expr)

	case *ast.IndexListExpr:
		return nil, errors.ErrNotAMapping

	case *ast.SliceExpr:
		return nil, errors.ErrNotAMapping

	case *ast.CallExpr:
		return s.callSource(defs, expr)

	case *ast.UnaryExpr:
		return nil, errors.ErrNotAMapping

	case *ast.BinaryExpr:
		return s.binarySource(defs, expr)

	case nil:
		return nil, errors.ErrEmpty

	default:
		return nil, errors.ErrUnsupportedExpression
	}
}
