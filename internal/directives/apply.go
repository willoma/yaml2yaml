package directives

import (
	"errors"
	"fmt"

	yerrors "github.com/willoma/yaml2yaml/internal/errors"
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

func (a *dirApply) Process(defs node.Definitions) (node.Node, error) {
	nameNode, err := a.name.Process(defs)
	if err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}

	def := defs.Get(nameNode.String())
	if def == nil {
		return nil, yerrors.ErrNotDefined
	}

	childDefs := defsWithArgs(defs, a.args)

	result, err := def.Process(childDefs)
	if err != nil {
		if errors.Is(err, yerrors.ErrEmpty) {
			childDefs.passToParent()
		}

		return nil, err
	}

	return result, nil
}
