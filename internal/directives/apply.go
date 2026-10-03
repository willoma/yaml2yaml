package directives

import (
	"fmt"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameApply = "!apply"

type Apply struct {
	Name node.Source
	Args node.Source
}

func (a *Apply) Process(defs node.Definitions) (node.Node, error) {
	var errs []error

	nameNode, err := a.Name.Process(defs)
	if err != nil {
		errs = append(errs, err)
	}

	name := nameNode.String()

	def := defs.Get(name)
	if def == nil {
		err := fmt.Errorf("%q %w", name, errors.ErrNotDefined)

		errs = append(errs, err)

		return nil, errors.Join(errs...)
	}

	childDefs := defsWithArgs(defs, a.Args)

	result, err := def.Process(childDefs)
	if err != nil {
		errs = append(errs, err)
	}

	if result == nil {
		childDefs.passToParent()
	}

	return result, errors.Join(errs...)
}
