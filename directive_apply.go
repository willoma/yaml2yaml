package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

func (n *node) applyMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	node := n.getDefinition(name)
	if node == nil {
		return nil, joinErrors(kErr, fmt.Errorf("%q is %w", name, ErrNotDefined))
	}

	resultValue := n.keyedChild(key.Tag, key.Value, node)
	aErr := resultValue.applyDefinitions(value.Node)

	processedResultValue, vErr := resultValue.process()
	result, cErr := nodeToMappingChild(key.path, processedResultValue)

	return result, joinErrors(kErr, aErr, vErr, cErr)
}

func (n *node) applyScalar() (*yaml.Node, error) {
	value, eErr := n.evaluateValue()
	name := anyToString(value)

	node := n.getDefinition(name)
	if node == nil {
		return anyToNode(n.Node, value), n.locatedError(fmt.Errorf("%q is %w", name, ErrNotDefined))
	}

	result, err := n.sameLevel(node).process()

	return result, joinErrors(eErr, err)
}
