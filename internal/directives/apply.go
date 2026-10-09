package directives

import (
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
		if errors.IsEmpty(err) {
			childDefs.passToParent()
		}

		return nil, err
	}

	return result, nil
}
