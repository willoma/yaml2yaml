package yaml2yaml

import (
	"fmt"
	"go/ast"
	"go/parser"
)

func (n *node) resolveExpression(src string) (any, error) {
	expr, err := parser.ParseExpr(src)
	if err != nil {
		return nil, n.locatedError(fmt.Errorf("parsing: %w", err))
	}

	result, err := n.resolveExpr(expr)
	if err != nil {
		err = n.locatedError(err)
	}

	return result, err
}

//nolint:funlen // Exhaustive switch on expression types.
func (n *node) resolveExpr(expr ast.Expr) (any, error) {
	// These are all the possible expression types, according to go/ast/ast.go.
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return nil, ErrSyntaxError

	case *ast.Ident:
		return n.resolveIdent(expr)

	case *ast.Ellipsis:
		return nil, fmt.Errorf("ellipsis are %w", ErrUnsupported)

	case *ast.BasicLit:
		return n.resolveBasicLit(expr)

	case *ast.FuncLit:
		return nil, fmt.Errorf("functions are %w", ErrUnsupported)

	case *ast.CompositeLit:
		return nil, fmt.Errorf("composites are %w", ErrUnsupported)

	case *ast.ParenExpr:
		return n.resolveExpr(expr.X)

	case *ast.SelectorExpr:
		return n.resolveSelectorExpr(expr)

	case *ast.IndexExpr:
		return nil, fmt.Errorf("index expressions %w", ErrNotImplemented)

	case *ast.IndexListExpr:
		return nil, fmt.Errorf("index list expressions %w", ErrNotImplemented)

	case *ast.SliceExpr:
		return nil, fmt.Errorf("slice expressions %w", ErrNotImplemented)

	case *ast.TypeAssertExpr:
		return nil, fmt.Errorf("type assertions %w", ErrNotImplemented)

	case *ast.CallExpr:
		return nil, fmt.Errorf("functions are %w", ErrUnsupported)

	case *ast.StarExpr:
		return nil, fmt.Errorf("pointers are %w", ErrUnsupported)

	case *ast.UnaryExpr:
		return n.resolveUnaryExpr(expr)

	case *ast.BinaryExpr:
		return n.resolveBinaryExpr(expr)

	case *ast.KeyValueExpr:
		return nil, fmt.Errorf("key-value pairs are %w", ErrUnsupported)

	case *ast.ArrayType:
		return nil, fmt.Errorf("arrays are %w", ErrUnsupported)

	case *ast.StructType:
		return nil, fmt.Errorf("structs are %w", ErrUnsupported)

	case *ast.FuncType:
		return nil, fmt.Errorf("functions are %w", ErrUnsupported)

	case *ast.InterfaceType:
		return nil, fmt.Errorf("interfaces are %w", ErrUnsupported)

	case *ast.MapType:
		return nil, fmt.Errorf("maps are %w", ErrUnsupported)

	case *ast.ChanType:
		return nil, fmt.Errorf("chans are %w", ErrUnsupported)

	default:
		return nil, fmt.Errorf("%w expression type %T", ErrInvalid, expr)
	}
}
