package yaml2yaml

import (
	"go.yaml.in/yaml/v4"
)

func (n *node) sequenceLength() int {
	return len(n.Content)
}

func (n *node) processSequence() (*yaml.Node, error) {
	var errs []error

	children := make([]*yaml.Node, 0, n.sequenceLength())

	for i, child := range n.Content {
		childResult, err := n.indexedChild(i, child).process()
		if err != nil {
			errs = append(errs, err)
		}

		if childResult != nil {
			children = append(children, childResult)
		}
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Content: children}, joinErrors(errs...)
}
