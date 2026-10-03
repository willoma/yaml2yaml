package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func selector(defs node.Definitions, expr *ast.SelectorExpr) (node.Node, error) {
	var errs []error

	value, err := Resolve(defs, expr.X)
	if err != nil {
		errs = append(errs, fmt.Errorf("left: %w", err))
	}

	if expr.Sel == nil {
		errs = append(errs, errors.ErrEmptySelector)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	xMap, ok := value.(node.Mapping)
	if !ok {
		return nil, fmt.Errorf("left: %w", errors.ErrNotAMapping)
	}

	for _, entry := range xMap {
		if entry.Key.String() == expr.Sel.Name {
			return entry.Value, nil
		}
	}

	//nolint:nilnil // Nil as a value is okay here.
	return nil, nil
}
