package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) selectorNode(defs node.Definitions, expr *ast.SelectorExpr) (node.Node, error) {
	var errs errors.JoinedError

	value, err := s.resolveNode(defs, expr.X)
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

	return xMap.GetValue(expr.Sel.Name)
}

func (s *Source) selectorSource(defs node.Definitions, expr *ast.SelectorExpr) (node.Source, error) {
	var errs errors.JoinedError

	value, err := s.resolveSource(defs, expr.X)
	if err != nil {
		errs = errs.Append(fmt.Errorf("expression: %w", err))
	}

	if expr.Sel == nil {
		errs = errs.Append(fmt.Errorf("selector: %w", errors.ErrEmpty))
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	xMap, ok := value.(node.SourceMapper)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	entries, err := xMap.SourceMap(defs)
	if err != nil {
		return nil, err
	}

	result, ok := entries[expr.Sel.Name]
	if !ok {
		return nil, errors.ErrEmpty
	}

	return result, nil
}
