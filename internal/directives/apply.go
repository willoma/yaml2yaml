package directives

import (
	oerrors "errors"
	"fmt"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameApply = "!apply"

type dirApply struct {
	name node.Source
	args node.Source
}

func Apply(name, args node.Source) *dirApply {
	return &dirApply{name: name, args: args}
}

func (a *dirApply) Pos() [2]int {
	return a.name.Pos()
}

func (a *dirApply) String() string {
	return NameApply + " " + a.name.String()
}

func (a *dirApply) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirApply); ok {
		return a.name.EqualsSource(other.name) && a.args.EqualsSource(other.args)
	}

	return false
}

func (a *dirApply) Process(defs node.Definitions) (node.Node, error) {
	nameNode, err := a.name.Process(defs)
	if err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}

	name := nameNode.String()

	def := defs.Get(name)
	if def == nil {
		return nil, fmt.Errorf("%q %w", name, errors.ErrNotDefined)
	}

	childDefs := defsWithArgs(defs, a.args)

	result, err := def.Process(childDefs)
	if err != nil {
		if oerrors.Is(err, errors.ErrEmpty) {
			childDefs.passToParent()
		}

		return nil, err
	}

	return result, nil
}

func (a *dirApply) PopulateDefs(defs node.Definitions) error {
	nameNode, err := a.expr.Process(defs)
	if err != nil {
		return fmt.Errorf("name: %w", err)
	}

	name := nameNode.String()

	def := defs.Get(name)
	if def == nil {
		return fmt.Errorf("%q %w", name, errors.ErrNotDefined)
	}

	return def.PopulateDefs(defs)
}
