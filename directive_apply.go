package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func doDirApply(defs *definitions.Definitions, key, value *yaml.Node) (*yaml.Node, error) {
	errs := errors.Multi()

	name, err := evaluateString(defs, key)

	errs.Add(err)

	defNode := defs.Get(name)
	if defNode == nil {
		errs.Add(fmt.Errorf("%q: %w", name, errors.ErrNotDefined))

		return nil, errs
	}

	var childDefs *definitions.Definitions
	if value == nil {
		childDefs = defs.Sibling()
	} else {
		childDefs = defs.Child()
		errs.Add(childDefs.SetFromNode(value))
	}

	result, err := processTransparent(childDefs, defNode)

	errs.Add(err)

	return result, errs
}
