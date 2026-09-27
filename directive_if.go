package yaml2yaml

import "github.com/willoma/yaml2yaml/errors"

func (n *node) ifMapping(key, value *node) ([]mappingChild, error) {
	expr, eErr := key.resolveExpression(key.Value)

	if !anyToBool(expr) {
		return nil, eErr
	}

	resultValue, vErr := value.process()
	result, cErr := nodeToMappingChild(key.path, resultValue)

	return result, errors.Join(eErr, vErr, cErr)
}
