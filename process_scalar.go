package yaml2yaml

import (
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
)

func processScalar(defs *definitions.Definitions, node *yaml.Node) (*yaml.Node, error) {
	result, processed, err := directiveScalar(defs, node)
	if processed {
		return result, err
	}

	if expr, ok := strings.CutPrefix(node.Value, "?"); ok {
		resultAny, err := resolveExpression(defs, expr)

		return convert.AnyToNode(resultAny), err
	}

	return &yaml.Node{
		Kind:   yaml.ScalarNode,
		Tag:    node.Tag,
		Value:  node.Value,
		Line:   node.Line,
		Column: node.Column,
	}, nil
}
