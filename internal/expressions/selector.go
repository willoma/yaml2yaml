package expressions

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) selector(defs node.Definitions, expr *ast.SelectorExpr) (node.Node, error) {
	value, err := s.resolve(defs, expr.X)
	if err != nil {
		return nil, err
	}

	if expr.Sel == nil {
		return nil, errors.ErrEmptySelector
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
