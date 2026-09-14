package yaml2yaml

import (
	"go.yaml.in/yaml/v4"
)

func (n *node) processScalar() (*yaml.Node, error) {
	value, eErr := n.evaluateValue()

	switch n.Tag {
	case "!apply":
		name := anyToString(value)

		node, err := n.getDefinition(name)
		if err != nil {
			return anyToNode(n.Node, value), n.locatedError(err)
		}

		result, err := n.child("", node).process()

		return result, joinErrors(eErr, err)
	default:
		return anyToNode(n.Node, value), eErr
	}
}
