package yaml2yaml

import (
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
)

func evaluateString(defs *definitions.Definitions, node *yaml.Node) (string, error) {
	if expr, ok := strings.CutPrefix(node.Value, "?"); ok {
		value, err := resolveExpression(defs, expr)

		return convert.AnyToString(value), err
	}

	return convert.NodeToString(node), nil
}
