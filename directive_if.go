package yaml2yaml

import (
	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func doDirIf(defs *definitions.Definitions, key, value *yaml.Node) (*yaml.Node, error) {
	expr, err := resolveExpression(defs, key.Value)

	if !convert.AnyToBool(expr) {
		return nil, err
	}

	errs := errors.Multi()

	errs.Add(err)

	result, err := processTransparent(defs, value)

	errs.Add(err)

	return result, errs
}
