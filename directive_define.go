package yaml2yaml

import "github.com/willoma/yaml2yaml/errors"

func (n *node) defineMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	dErr := n.definitions.Set(name, value.Node)
	if dErr != nil {
		dErr = key.locatedError(dErr)
	}

	return nil, errors.Join(kErr, dErr)
}
