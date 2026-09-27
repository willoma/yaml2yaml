package yaml2yaml

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
)

func (n *node) resolveSelectorExpr(expr *ast.SelectorExpr) (any, error) {
	var sErr error
	if expr.Sel == nil {
		sErr = fmt.Errorf("%w: selector expression is empty", ErrInvalid)
	}

	x, xErr := n.resolveExpr(expr.X)
	if xErr != nil {
		return nil, errors.Join(sErr, xErr)
	}

	xMap, ok := x.(map[string]any)
	if !ok {
		return nil, errors.Join(sErr, fmt.Errorf("left: %w", ErrNotMapping))
	}

	return xMap[expr.Sel.Name], nil
}
