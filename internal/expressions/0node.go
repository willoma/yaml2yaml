package expressions

import (
	"fmt"
	"go/parser"

	"github.com/willoma/yaml2yaml/internal/node"
)

type Source string

func (s Source) Process(defs node.Definitions) (node.Node, error) {
	expr, err := parser.ParseExpr(string(s))
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return Resolve(defs, expr)
}
