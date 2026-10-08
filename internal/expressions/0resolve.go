package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

//nolint:funlen // Exhaustive switch on expression types.
func resolve(defs node.Definitions, expr ast.Expr) (node.Node, error) {
	// These are all the possible expression types, according to go/ast/ast.go.
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return nil, errors.ErrSyntaxError

	case *ast.Ident:
		return ident(defs, expr)

	case *ast.Ellipsis:
		return nil, fmt.Errorf("ellipsis: %w", errors.ErrUnsupportedExpression)

	case *ast.BasicLit:
		return basicLit(expr)

	case *ast.FuncLit:
		return nil, fmt.Errorf("function: %w", errors.ErrUnsupportedExpression)

	case *ast.CompositeLit:
		return nil, fmt.Errorf("composite: %w", errors.ErrUnsupportedExpression)

	case *ast.ParenExpr:
		return resolve(defs, expr.X)

	case *ast.SelectorExpr:
		return selector(defs, expr)

		// 	case *ast.IndexExpr:
		// 		return nil, fmt.Errorf("index expressions %w", errors.ErrNotImplemented)

		// 	case *ast.IndexListExpr:
		// 		return nil, fmt.Errorf("index list expressions %w", errors.ErrNotImplemented)

		// 	case *ast.SliceExpr:
		// 		return nil, fmt.Errorf("slice expressions %w", errors.ErrNotImplemented)

		// 	case *ast.TypeAssertExpr:
		// 		return nil, fmt.Errorf("type assertions %w", errors.ErrNotImplemented)

	case *ast.CallExpr:
		return nil, fmt.Errorf("function: %w", errors.ErrUnsupportedExpression)

	case *ast.StarExpr:
		return nil, fmt.Errorf("pointer: %w", errors.ErrUnsupportedExpression)

	case *ast.UnaryExpr:
		return unary(defs, expr)

	case *ast.BinaryExpr:
		return binary(defs, expr)

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
