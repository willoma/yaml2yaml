package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func resolveExpression(defs *definitions.Definitions, src string) (any, error) {
	expr, err := parser.ParseExpr(src)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	return resolveExpr(defs, expr)
}

//nolint:funlen // Exhaustive switch on expression types.
func resolveExpr(defs *definitions.Definitions, expr ast.Expr) (any, error) {
	// These are all the possible expression types, according to go/ast/ast.go.
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return nil, errors.ErrSyntaxError

	case *ast.Ident:
		return resolveIdent(defs, expr)

	case *ast.Ellipsis:
		return nil, fmt.Errorf("ellipsis: %w", errors.ErrUnsupportedExpression)

	case *ast.BasicLit:
		return resolveBasicLit(expr)

	case *ast.FuncLit:
		return nil, fmt.Errorf("function: %w", errors.ErrUnsupportedExpression)

	case *ast.CompositeLit:
		return nil, fmt.Errorf("composite: %w", errors.ErrUnsupportedExpression)

	case *ast.ParenExpr:
		return resolveExpr(defs, expr.X)

	case *ast.SelectorExpr:
		return resolveSelectorExpr(defs, expr)

	case *ast.IndexExpr:
		return nil, fmt.Errorf("index expressions %w", errors.ErrNotImplemented)

	case *ast.IndexListExpr:
		return nil, fmt.Errorf("index list expressions %w", errors.ErrNotImplemented)

	case *ast.SliceExpr:
		return nil, fmt.Errorf("slice expressions %w", errors.ErrNotImplemented)

	case *ast.TypeAssertExpr:
		return nil, fmt.Errorf("type assertions %w", errors.ErrNotImplemented)

	case *ast.CallExpr:
		return nil, fmt.Errorf("function: %w", errors.ErrUnsupportedExpression)

	case *ast.StarExpr:
		return nil, fmt.Errorf("pointer: %w", errors.ErrUnsupportedExpression)

	case *ast.UnaryExpr:
		return resolveUnaryExpr(defs, expr)

	case *ast.BinaryExpr:
		return resolveBinaryExpr(defs, expr)

	case *ast.KeyValueExpr:
		return nil, fmt.Errorf("key-value pair: %w", errors.ErrUnsupportedExpression)

	case *ast.ArrayType:
		return nil, fmt.Errorf("array: %w", errors.ErrUnsupportedExpression)

	case *ast.StructType:
		return nil, fmt.Errorf("struct: %w", errors.ErrUnsupportedExpression)

	case *ast.FuncType:
		return nil, fmt.Errorf("function: %w", errors.ErrUnsupportedExpression)

	case *ast.InterfaceType:
		return nil, fmt.Errorf("interface: %w", errors.ErrUnsupportedExpression)

	case *ast.MapType:
		return nil, fmt.Errorf("map: %w", errors.ErrUnsupportedExpression)

	case *ast.ChanType:
		return nil, fmt.Errorf("chan: %w", errors.ErrUnsupportedExpression)

	default:
		return nil, fmt.Errorf("%w type %T", errors.ErrInvalidExpression, expr)
	}
}
