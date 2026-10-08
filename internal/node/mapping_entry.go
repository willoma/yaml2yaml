package node

import (
	"errors"

	"go.yaml.in/yaml/v4"
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

func (m MappingEntrySource) Process(defs Definitions) (Node, error) {
	key, kErr := m.Key.Process(defs)
	value, vErr := m.Value.Process(defs)

	return MappingEntry{Key: key, Value: value}, errors.Join(kErr, vErr)
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
	return m.Key.String() + ": " + m.Value.String()
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
