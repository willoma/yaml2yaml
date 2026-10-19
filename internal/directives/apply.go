package directives

import (
	oerrors "errors"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameApply = "!apply"

type dirApply struct {
	expr *expressions.Source
	args node.Source
}

func Apply(expr, args node.Source) (*dirApply, error) {
	exprSrc, err := expressions.Parse(expr.Pos(), expr.String())
	if err != nil {
		return nil, err
	}

	return &dirApply{expr: exprSrc, args: args}, nil
}

func (a *dirApply) Pos() [2]int {
	return a.expr.Pos()
}

func (a *dirApply) String() string {
	return NameApply + " " + a.expr.String()
}

func (a *dirApply) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirApply); ok {
		return a.expr.EqualsSource(other.expr) && a.args.EqualsSource(other.args)
	}

	return false
}

func (a *dirApply) Process(defs node.Definitions) (node.Node, error) {
	childDefs, err := defsWithArgs(defs, a.args)
	if err != nil {
		return nil, err
	}

	result, err := a.expr.Process(childDefs)
	if err != nil {
		if oerrors.Is(err, errors.ErrEmpty) {
			childDefs.passToParent()
		}

		return nil, err
	}

	return result, nil
}

func (a *dirApply) SourceMap(defs node.Definitions) (map[string]node.Source, error) {
	childDefs, err := defsWithArgs(defs, a.args)
	if err != nil {
		return nil, err
	}

	result, err := a.expr.SourceMap(childDefs)
	if err != nil {
		return nil, err
	}

	return result, nil
}
