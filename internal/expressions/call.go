package expressions

import (
	"fmt"
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) call(defs node.Definitions, expr *ast.CallExpr) (node.Node, error) {
	var errs errors.JoinedError

	fun, err := s.resolveNode(&definitionsWithFuns{parent: defs}, expr.Fun)
	if err != nil {
		errs = errs.Append(fmt.Errorf("function: %w", err))
	} else if fun == nil {
		errs = errs.Append(fmt.Errorf("function: %w", errors.ErrEmpty))
	}

	args := make([]node.Node, len(expr.Args))
	for i, arg := range expr.Args {
		argNode, err := s.resolveNode(defs, arg)
		if err != nil {
			errs = errs.Append(fmt.Errorf("argument %d: %w", i+1, err))
		}

		args[i] = argNode
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	callFun, ok := fun.(callableFunction)
	if !ok {
		return nil, errors.ErrNotAFunction
	}

	return callFun.Call(args)
}
