package expressions

import (
	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

type (
	callableFunction interface {
		node.Node
		Call(args []node.Node) (node.Node, error)
	}
	exprFunction struct{ callableFunction }

	baseFun    struct{}
	funLen     struct{ baseFun }
	funIsValue struct{ baseFun }

	definitionsWithFuns struct {
		parent node.Definitions
	}
)

//nolint:gochecknoglobals // This never changes.
var funs = map[string]exprFunction{
	"len":      {funLen{}},
	"is_value": {funIsValue{}},
}

func (d *definitionsWithFuns) Get(name string) node.Source {
	if fun, ok := funs[name]; ok {
		return fun
	}

	return d.parent.Get(name)
}

func (d *definitionsWithFuns) Set(name string, value node.Source) {
	d.parent.Set(name, value)
}

func (f exprFunction) Pos() [2]int                                 { return [2]int{} }
func (f exprFunction) String() string                              { return "function call" }
func (f exprFunction) EqualsSource(node.Source) bool               { return false }
func (f exprFunction) Process(node.Definitions) (node.Node, error) { return f, nil }
func (f exprFunction) PopulateDefs(node.Definitions) error         { return nil }

func (f baseFun) Bool() (value, ok bool)    { return false, false }
func (f baseFun) Int() (int64, bool)        { return 0, false }
func (f baseFun) Float() (float64, bool)    { return 0, false }
func (f baseFun) String() string            { return "function call" }
func (f baseFun) EqualsNode(node.Node) bool { return false }
func (f baseFun) ToYAML() *yaml.Node        { return nil }

func (f funLen) Call(args []node.Node) (node.Node, error) {
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

func (f funIsValue) Call(args []node.Node) (node.Node, error) {
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
