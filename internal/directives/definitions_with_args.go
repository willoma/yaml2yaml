package directives

import "github.com/willoma/yaml2yaml/internal/node"

type definitionsWithArgs struct {
	parent      node.Definitions
	args        map[string]node.Source
	definitions map[string]node.Source
}

func defsWithArgs(parent node.Definitions, args node.Source) (*definitionsWithArgs, error) {
	defs := &definitionsWithArgs{
		parent:      parent,
		definitions: make(map[string]node.Source),
	}

	if args == nil {
		return defs, nil
	}

	defs.definitions["_"] = args

	if err := args.PopulateDefs(defs); err != nil {
		return nil, err
	}

	defs.args = defs.definitions
	defs.definitions = make(map[string]node.Source)

	return defs, nil
}

func (d *definitionsWithArgs) Get(name string) node.Source {
	if value, ok := d.args[name]; ok {
		return value
	}

	if value, ok := d.definitions[name]; ok {
		return value
	}

	return d.parent.Get(name)
}

func (d *definitionsWithArgs) Set(name string, value node.Source) {
	d.definitions[name] = value
}

func (d *definitionsWithArgs) passToParent() {
	for name, value := range d.definitions {
		d.parent.Set(name, value)
	}
}
