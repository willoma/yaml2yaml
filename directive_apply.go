package yaml2yaml

import (
	"fmt"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

func (n *node) applyMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	node := n.definitions.Get(name)
	if node == nil {
		return nil, errors.Join(kErr, fmt.Errorf("%q is %w", name, ErrNotDefined))
	}

	resultValue := n.keyedChild(key.Tag, key.Value, node)
	aErr := resultValue.definitions.SetFromNode(value.Node)

	processedResultValue, vErr := resultValue.process()
	result, cErr := nodeToMappingChild(key.path, processedResultValue)

	return result, errors.Join(kErr, aErr, vErr, cErr)
}

func (n *node) applyScalar() (*yaml.Node, error) {
	value, eErr := n.evaluateValue()
	name := anyToString(value)

	node := n.definitions.Get(name)
	if node == nil {
		return anyToNode(n.Node, value), n.locatedError(fmt.Errorf("%q is %w", name, ErrNotDefined))
	}

	result, err := n.sameLevel(node).process()

	return result, errors.Join(eErr, err)
}
