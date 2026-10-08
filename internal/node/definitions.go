package node

type Definitions interface {
	Get(name string) Source
	Set(name string, value Source)
}

type childDefinitions struct {
	parent      Definitions
	definitions map[string]Source
}

func childDefs(parent Definitions) *childDefinitions {
	return &childDefinitions{
		parent:      parent,
		definitions: make(map[string]Source),
	}
}

// Get returns the value of the definition with the given name, if any.
// It first checks the current [Definitions], then its parent, if any.
func (d *childDefinitions) Get(name string) Source {
	if value, ok := d.definitions[name]; ok {
		return value
	}

	return d.parent.Get(name)
}

// Set sets the value of the definition with the given name.
func (d *childDefinitions) Set(name string, value Source) {
	d.definitions[name] = value
}

func (d *childDefinitions) passToParent() {
	for name, value := range d.definitions {
		d.parent.Set(name, value)
	}
}
