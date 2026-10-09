package directives

import "github.com/willoma/yaml2yaml/internal/node"

type definitionsWithArgs struct {
	parent      node.Definitions
	args        map[string]node.Source
	definitions map[string]node.Source
}

func defsWithArgs(parent node.Definitions, args node.Source) *definitionsWithArgs {
	defs := &definitionsWithArgs{
		parent:      parent,
		args:        make(map[string]node.Source),
		definitions: make(map[string]node.Source),
	}

	if args == nil {
		return defs
	}

	mapping, ok := args.(*node.MappingSource)
	if !ok {
		defs.args["_"] = args

		return defs
	}

	for _, item := range mapping.Value {
		entry, ok := item.(*node.MappingEntrySource)
		if !ok {
			// Ignore non-raw entries (ie. directives)
			continue
		}

		key, err := entry.Key.Process(parent)
		if err != nil {
			// Ignore erroneous keys
			continue
		}

		defs.args[key.String()] = entry.Value
	}

	return defs
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
