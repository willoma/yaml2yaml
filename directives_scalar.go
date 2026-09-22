package yaml2yaml

import "go.yaml.in/yaml/v4"

func (n *node) scalarDirective() (*yaml.Node, error) {
	switch n.Tag {
	case dirApply:
		return n.applyScalar()
	default:
		return n.scalarNoDirective()
	}
}

func (n *node) applyScalar() (*yaml.Node, error) {
	value, eErr := n.evaluateValue()
	name := anyToString(value)

	node, err := n.getDefinition(name)
	if err != nil {
		return anyToNode(n.Node, value), n.locatedError(err)
	}

	result, err := n.sameLevel(node).process()

	return result, joinErrors(eErr, err)
}
