package yaml2yaml

import (
	"fmt"
	"go/ast"
)

func (n *node) resolveSelectorExpr(expr *ast.SelectorExpr) (any, error) {
	var sErr error
	if expr.Sel == nil {
		sErr = fmt.Errorf("%w: selector expression is empty", ErrInvalid)
	}

	x, xErr := n.resolveExpr(expr.X)
	if xErr != nil {
		return nil, joinErrors(sErr, xErr)
	}

	xMap, ok := x.(map[string]any)
	if !ok {
		return nil, joinErrors(sErr, fmt.Errorf("left: %w", ErrNotMapping))
	}

	return xMap[expr.Sel.Name], nil
}
