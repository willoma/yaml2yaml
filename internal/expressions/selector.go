package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) selector(defs node.Definitions, expr *ast.SelectorExpr) (node.Node, error) {
	errs := errors.NewJoined()

	value, err := s.resolve(defs, expr.X)
	if err != nil {
		errs.Add(fmt.Errorf("expression: %w", err))
	}

	if expr.Sel == nil {
		return nil, fmt.Errorf("selector: %w", errors.ErrEmpty)
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
