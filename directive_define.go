package yaml2yaml

import (
	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func doDirDefine(defs *definitions.Definitions, key, value *yaml.Node) error {
	errs := errors.Multi()

	name, err := evaluateString(defs, key)

	errs.Add(err)

	err = defs.Set(name, value)

	errs.Add(err)

	return errs
}
