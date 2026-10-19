package directives

import (
	"fmt"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameIf = "!if"

type dirIf struct {
	condition *expressions.Source
	value     node.Source
}

func If(condition, value node.Source) (*dirIf, error) {
	expr, err := expressions.Parse(condition.Pos(), condition.String())
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	return &dirIf{condition: expr, value: value}, nil
}

func (i *dirIf) Pos() [2]int {
	return i.condition.Pos()
}

func (i *dirIf) String() string {
	return NameIf + " " + i.condition.String()
}

func (i *dirIf) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirIf); ok {
		return i.condition.EqualsSource(other.condition) && i.value.EqualsSource(other.value)
	}

	return false
}

func (i *dirIf) Process(defs node.Definitions) (node.Node, error) {
	cond, err := i.condition.Process(defs)
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	if ok, _ := cond.Bool(); !ok {
		return nil, errors.ErrEmpty
	}

	return i.value.Process(defs)
}

func (i *dirIf) SourceMap(defs node.Definitions) (map[string]node.Source, error) {
	cond, err := i.condition.Process(defs)
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	if ok, _ := cond.Bool(); !ok {
		return nil, errors.ErrEmpty
	}

	sourceMapper, ok := i.value.(node.SourceMapper)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	return sourceMapper.SourceMap(defs)
}
