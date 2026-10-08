package directives

import (
	"fmt"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameIf = "!if"

type dirIf struct {
	condition *expressions.Source
	value     node.Source
}

func If(condition, value node.Source) (node.Source, error) {
	cond, ok := condition.(node.String)
	if !ok {
		return nil, fmt.Errorf("condition: %w", errors.ErrInvalidExpression)
	}

	expr, err := expressions.Parse(cond.String())
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	return &dirIf{condition: expr, value: value}, nil
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
