package expressions

import (
	"fmt"
	"go/ast"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

type (
	callableFunction interface {
		node.Node
		callNode(args []node.Node) (node.Node, error)
		callSource(args []node.Source) (node.Source, error)
	}
	exprFunction struct{ callableFunction }

	baseFun struct{}

	definitionsWithFuns struct {
		parent node.Definitions
	}
)

//nolint:dupl // Yes, it is very similar to callSource.
func (s *Source) callNode(defs node.Definitions, expr *ast.CallExpr) (node.Node, error) {
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

	return callFun.callNode(args)
}

//nolint:dupl // Yes, it is very similar to callNode.
func (s *Source) callSource(defs node.Definitions, expr *ast.CallExpr) (node.Source, error) {
	var errs errors.JoinedError

	fun, err := s.resolveNode(&definitionsWithFuns{parent: defs}, expr.Fun)
	if err != nil {
		errs = errs.Append(fmt.Errorf("function: %w", err))
	} else if fun == nil {
		errs = errs.Append(fmt.Errorf("function: %w", errors.ErrEmpty))
	}

	args := make([]node.Source, len(expr.Args))
	for i, arg := range expr.Args {
		argNode, err := s.resolveSource(defs, arg)
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

	return callFun.callSource(args)
}

func (d *definitionsWithFuns) Get(name string) node.Source {
	if fun, ok := funs[name]; ok {
		return fun
	}

	return d.parent.Get(name)
}

func (d *definitionsWithFuns) Set(name string, value node.Source) {
	d.parent.Set(name, value)
}

func (f exprFunction) Pos() [2]int                                 { return [2]int{} }
func (f exprFunction) String() string                              { return "function call" }
func (f exprFunction) EqualsSource(node.Source) bool               { return false }
func (f exprFunction) Process(node.Definitions) (node.Node, error) { return f, nil }

func (f baseFun) Bool() (value, ok bool)    { return false, false }
func (f baseFun) Int() (int64, bool)        { return 0, false }
func (f baseFun) Float() (float64, bool)    { return 0, false }
func (f baseFun) String() string            { return "function call" }
func (f baseFun) EqualsNode(node.Node) bool { return false }
func (f baseFun) ToYAML() *yaml.Node        { return nil }
