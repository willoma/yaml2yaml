package yaml2yaml

func (n *node) defineMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	dErr := n.setDefinition(name, value.Node)
	if dErr != nil {
		dErr = key.locatedError(dErr)
	}

	return nil, joinErrors(kErr, dErr)
}
