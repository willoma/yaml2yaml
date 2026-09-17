package yaml2yaml

import (
	"go.yaml.in/yaml/v4"
)

// mappingChild is the result of evaluating a mapping child.
//
//   - kvs is not nil if child type is mapping entry or entries
//   - value is not nil if child type is single value
//   - both are nil if child is to be ignored
type mappingChild struct {
	path  string
	kvs   []*yaml.Node
	value *yaml.Node
}

func (n *node) evaluateMappingChild(key, value *node) (mappingChild, error) {
	if key.Kind != yaml.ScalarNode {
		keyResult, kErr := key.process()
		valueResult, vErr := value.process()

		return mappingChild{kvs: []*yaml.Node{keyResult, valueResult}}, joinErrors(kErr, vErr)
	}

	return n.mappingDirective(key, value)
}

func (n *node) mappingNoDirective(key, value *node) (mappingChild, error) {
	keyResult, kErr := key.process()
	valueResult, vErr := value.process()

	return mappingChild{kvs: []*yaml.Node{keyResult, valueResult}}, joinErrors(kErr, vErr)
}
