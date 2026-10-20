package expressions

import (
	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

type (
	funLen     struct{ baseFun }
	funIsValue struct{ baseFun }
	funValues  struct{ baseFun }
)

//nolint:gochecknoglobals // This never changes.
var funs = map[string]exprFunction{
	"len":      {funLen{}},
	"is_value": {funIsValue{}},
	"values":   {funValues{}},
}

func (f funLen) callNode(args []node.Node) (node.Node, error) {
	if len(args) != 1 {
		return nil, errors.ErrInvalidArgumentsCount
	}

	var length int64

	switch arg := args[0].(type) {
	case node.Mapping:
		length = int64(len(arg))
	case node.Sequence:
		length = int64(len(arg))
	case nil:
		length = 0
	default:
		length = 1
	}

	return node.NewInt([2]int{}, length), nil
}

func (f funLen) callSource(args []node.Source) (node.Source, error) {
	if len(args) != 1 {
		return nil, errors.ErrInvalidArgumentsCount
	}

	var length int64

	switch arg := args[0].(type) {
	case *node.MappingSource:
		length = int64(arg.Size())
	case *node.SequenceSource:
		length = int64(arg.Size())
	case nil:
		length = 0
	default:
		length = 1
	}

	return node.NewInt([2]int{}, length), nil
}

func (f funIsValue) callNode(args []node.Node) (node.Node, error) {
	if len(args) != 1 {
		return nil, errors.ErrInvalidArgumentsCount
	}

	var result bool

	switch args[0].(type) {
	case node.Bool, node.Int, node.Float, node.String, node.Binary:
		result = true
	}

	return node.NewBool([2]int{}, result), nil
}

func (f funIsValue) callSource(args []node.Source) (node.Source, error) {
	if len(args) != 1 {
		return nil, errors.ErrInvalidArgumentsCount
	}

	var result bool

	switch args[0].(type) {
	case node.Bool, node.Int, node.Float, node.String, node.Binary:
		result = true
	}

	return node.NewBool([2]int{}, result), nil
}

func (f funValues) callNode(args []node.Node) (node.Node, error) {
	if len(args) != 1 {
		return nil, errors.ErrInvalidArgumentsCount
	}

	argMap, ok := args[0].(node.Mapping)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	result := make(node.Sequence, 0, len(argMap))

	for _, item := range argMap {
		result = append(result, item.Value())
	}

	return result, nil
}

func (f funValues) callSource(args []node.Source) (node.Source, error) {
	if len(args) != 1 {
		return nil, errors.ErrInvalidArgumentsCount
	}

	argMap, ok := args[0].(*node.MappingSource)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	result := node.NewSequenceSource([2]int{}, argMap.Size())

	for item := range argMap.Iter() {
		entry, ok := item.(*node.MappingEntrySource)
		if !ok {
			continue
		}

		result.Append(entry.Value())
	}

	return result, nil
}
