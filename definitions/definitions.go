package definitions

import (
	"fmt"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/errors"
)

type Definitions struct {
	parent      *Definitions
	definitions map[string]*yaml.Node
}

// New returns a new, empty [Definitions].
func New() *Definitions {
	return &Definitions{
		definitions: make(map[string]*yaml.Node),
	}
}

// Child returns a new [Definitions] with the current [Definitions] as its parent.
func (d *Definitions) Child() *Definitions {
	return &Definitions{
		parent:      d,
		definitions: make(map[string]*yaml.Node),
	}
}

// Sibling returns a new [Definitions] with the current [Definitions]' parent as its parent.
func (d *Definitions) Sibling() *Definitions {
	return &Definitions{
		parent:      d.parent,
		definitions: make(map[string]*yaml.Node),
	}
}

// PassToParent passes all definitions to the parent, if any.
func (d *Definitions) PassToParent() error {
	if d.parent == nil {
		return nil
	}

	errs := errors.Multi()

	for name, value := range d.definitions {
		errs.Add(d.parent.Set(name, value))
	}

	return errs
}

// Get returns the value of the definition with the given name, if any.
// It first checks the current [Definitions], then its parent, if any.
func (d *Definitions) Get(name string) *yaml.Node {
	if value, ok := d.definitions[name]; ok {
		return value
	}

	if d.parent == nil {
		return nil
	}

	return d.parent.Get(name)
}

// Set sets the value of the definition with the given name.
func (d *Definitions) Set(name string, value *yaml.Node) error {
	if d.definitions[name] != nil {
		return fmt.Errorf("%q %w", name, errors.ErrAlreadyDefined)
	}

	d.definitions[name] = value

	return nil
}

// SetFromNode sets the value of the definition with the given name from a YAML node.
func (d *Definitions) SetFromNode(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return d.Set("_", value)
	}

	errs := errors.Multi()

	for i := 0; i < len(value.Content); i += 2 {
		errs.Add(d.Set(
			convert.NodeToString(value.Content[i]),
			value.Content[i+1],
		))
	}

	return errs
}
