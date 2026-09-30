package yaml2yaml

import (
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func processSequence(defs *definitions.Definitions, node *yaml.Node) (*yaml.Node, error) {
	errs := errors.Multi()

	children := make([]*yaml.Node, 0, len(node.Content))

	for i, child := range node.Content {
		childResult, err := process(defs.Child(), "["+strconv.Itoa(i)+"]", child)

		errs.Add(err)

		if childResult != nil {
			children = append(children, childResult)
		}
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Content: children}, errs
}
