package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) selector(defs node.Definitions, expr *ast.SelectorExpr) (node.Node, error) {
	var errs errors.JoinedError

	value, err := s.resolve(defs, expr.X)
	if err != nil {
		errs = errs.Append(fmt.Errorf("expression: %w", err))
	}

	if expr.Sel == nil {
		errs = errs.Append(fmt.Errorf("selector: %w", errors.ErrEmpty))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xMap, ok := value.(node.Mapping)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	for _, entry := range xMap {
		if entry.Key.String() == expr.Sel.Name {
			return entry.Value, nil
		}
	}

	return nil, errors.ErrEmpty
}
