package yaml2yaml

func (n *node) ifMapping(key, value *node) ([]mappingChild, error) {
	expr, eErr := key.resolveExpression(key.Value)

	if !anyToBool(expr) {
		return nil, eErr
	}

	resultValue, vErr := value.process()
	result, cErr := nodeToMappingChild(key.path, resultValue)

	return result, joinErrors(eErr, vErr, cErr)
}
