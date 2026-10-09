package node

import (
	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/errors"
)

type (
	MappingEntrySource struct {
		Key   Source
		Value Source
	}

	// MappingEntry represents a single key-value pair in a YAML mapping.
	// If used by itself, it will be converted to a YAML mapping node with a single entry.
	MappingEntry struct {
		Key   Node
		Value Node
	}
)

func NewMappingEntrySource(key, value Source) *MappingEntrySource {
	return &MappingEntrySource{
		Key:   key,
		Value: value,
	}
}

func (m *MappingEntrySource) Pos() [2]int {
	return m.Key.Pos()
}

func (m *MappingEntrySource) String() string {
	return m.Key.String()
}

func (m *MappingEntrySource) Process(defs Definitions) (Node, error) {
	errs := errors.NewJoined()

	key, err := m.Key.Process(defs)
	errs.Add(err)

	value, err := m.Value.Process(defs)
	errs.Add(err)

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return MappingEntry{Key: key, Value: value}, nil
}

func (m MappingEntry) Bool() (value, ok bool) {
	return m.Value.Bool()
}

func (m MappingEntry) Int() (int64, bool) {
	return m.Value.Int()
}

func (m MappingEntry) Float() (float64, bool) {
	return m.Value.Float()
}

func (m MappingEntry) String() string {
	return m.Key.String()
}

func (m MappingEntry) Equals(other Node) bool {
	otherEntry, ok := other.(MappingEntry)
	if !ok {
		return false
	}

	return m.Key.Equals(otherEntry.Key) && m.Value.Equals(otherEntry.Value)
}

func (m MappingEntry) Merge(other Node) (Node, error) {
	return Mapping{m}.Merge(other)
}

func (m MappingEntry) ToYAML() *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{m.Key.ToYAML(), m.Value.ToYAML()},
	}
}
