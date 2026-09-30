package yaml2yaml

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func resolveSelectorExpr(defs *definitions.Definitions, expr *ast.SelectorExpr) (any, error) {
	errs := errors.Multi()
	if expr.Sel == nil {
		errs.Add(errors.ErrEmptySelector)
	}

	value, err := resolveExpr(defs, expr.X)
	if err != nil {
		errs.Add(err)

		return nil, errs
	}

	xMap, ok := value.(map[string]any)
	if !ok {
		errs.Add(fmt.Errorf("left: %w", errors.ErrNotAMapping))

		return nil, errs
	}

	return xMap[expr.Sel.Name], nil
}
