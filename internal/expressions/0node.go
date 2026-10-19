package expressions

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/willoma/yaml2yaml/internal/node"
)

type Source struct {
	pos  [2]int
	src  string
	expr ast.Expr
}

func Parse(pos [2]int, src string) (*Source, error) {
	expr, err := parser.ParseExpr(src)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return &Source{pos: pos, src: src, expr: expr}, nil
}

func (s *Source) Pos() [2]int {
	return s.pos
}

func (s *Source) String() string {
	return s.src
}

func (s *Source) EqualsSource(other node.Source) bool {
	if otherSource, ok := other.(*Source); ok {
		return s.src == otherSource.src
	}

	return false
}

func (s *Source) Process(defs node.Definitions) (node.Node, error) {
	return s.resolve(defs, s.expr)
}

func (s *Source) PopulateDefs(defs node.Definitions) error {
	return nil
}
