package definitions

import (
	"fmt"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/node"
)

var (
	ErrAlreadyDefined = fmt.Errorf("already defined")
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

	var errs []error

	for name, value := range d.definitions {
		err := d.parent.Set(name, value)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
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
		return fmt.Errorf("%q %w", name, ErrAlreadyDefined)
	}

	d.definitions[name] = value

	return nil
}

// SetFromNode sets the value of the definition with the given name from a YAML node.
func (d *Definitions) SetFromNode(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return d.Set("_", value)
	}

	var errs []error

	for i := 0; i < len(value.Content); i += 2 {
		key := node.ToString(value.Content[i])
		if err := d.Set(key, value.Content[i+1]); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
