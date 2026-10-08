package expressions

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/willoma/yaml2yaml/internal/node"
)

type Source struct {
	expr ast.Expr
}

func Parse(src string) (*Source, error) {
	expr, err := parser.ParseExpr(src)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return &Source{expr: expr}, nil
}

func (s Source) Process(defs node.Definitions) (node.Node, error) {
	return resolve(defs, s.expr)
}
