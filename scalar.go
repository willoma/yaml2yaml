package yaml2yaml

import (
	"go.yaml.in/yaml/v4"
)

func (n *node) processScalar() (*yaml.Node, error) {
	return n.scalarDirective()
}

func (n *node) scalarNoDirective() (*yaml.Node, error) {
	value, err := n.evaluateValue()

	return anyToNode(n.Node, value), err
}
