package node

import (
	oerrors "errors"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

type (
	MappingSource struct {
		pos   [2]int
		Value []Source
	}

	MappingEntrySource struct {
		Key   Source
		Value Source
	}

	Mapping []mappingEntry

	mappingEntry struct {
		pos   [2]int
		Key   Node
		Value Node
	}
)

func NewMappingSource(pos [2]int, length int) *MappingSource {
	return &MappingSource{
		pos:   pos,
		Value: make([]Source, 0, length),
	}
}

func (m *MappingSource) Pos() [2]int {
	return m.pos
}

func (m *MappingSource) String() string {
	result := make([]string, len(m.Value))
	for i, item := range m.Value {
		result[i] = item.String()
	}

	return "{" + strings.Join(result, ", ") + "}"
}

func (m *MappingSource) EqualsSource(other Source) bool {
	otherMap, ok := other.(*MappingSource)
	if !ok {
		return false
	}

	if len(m.Value) != len(otherMap.Value) {
		return false
	}

	for i := range m.Value {
		if !m.Value[i].EqualsSource(otherMap.Value[i]) {
			return false
		}
	}

	return true
}

func (m *MappingSource) Process(parentDefs Definitions) (Node, error) {
	defs := childDefs(parentDefs)
	processed := NewProcessedEntries(len(m.Value))

	for _, item := range m.Value {
		itemResult, err := item.Process(defs)

		processed.Append(item.String(), item.Pos(), itemResult, err)
	}

	result, err := processed.Result()
	if oerrors.Is(err, errors.ErrEmpty) {
		defs.passToParent()
	}

	return result, err
}

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

func (m *MappingEntrySource) EqualsSource(other Source) bool {
	if otherEntry, ok := other.(*MappingEntrySource); ok {
		return m.Key.EqualsSource(otherEntry.Key) && m.Value.EqualsSource(otherEntry.Value)
	}

	return false
}

func (m *MappingEntrySource) Process(defs Definitions) (Node, error) {
	var errs errors.JoinedError

	key, err := m.Key.Process(defs)
	if err != nil {
		errs = errs.Append(err)
	}

	value, err := m.Value.Process(defs)
	if err != nil {
		errs = errs.Append(err)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return mappingEntry{pos: m.Pos(), Key: key, Value: value}, nil
}

func (m Mapping) Bool() (value, ok bool) {
	return len(m) > 0, true
}

func (m Mapping) Int() (int64, bool) {
	return 0, false
}

func (m Mapping) Float() (float64, bool) {
	return 0, false
}

func (m Mapping) Binary() ([]byte, bool) {
	return nil, false
}

func (m Mapping) String() string {
	result := make([]string, len(m))
	for i, item := range m {
		result[i] = item.String()
	}

	return "{" + strings.Join(result, ", ") + "}"
}

func (m Mapping) EqualsNode(other Node) bool {
	otherMap, ok := other.(Mapping)
	if !ok {
		return false
	}

	if len(m) != len(otherMap) {
		return false
	}

	for i := range m {
		if !m[i].EqualsNode(otherMap[i]) {
			return false
		}
	}

	return true
}

func (m Mapping) ToYAML() *yaml.Node {
	content := make([]*yaml.Node, 0, len(m)*2)

	for _, item := range m {
		if item.Key == nil || item.Value == nil {
			continue
		}

		content = append(content, item.Key.ToYAML(), item.Value.ToYAML())
	}

	return &yaml.Node{Kind: yaml.MappingNode, Content: content}
}

func (m Mapping) Merge(other Mapping) (Mapping, error) {
	var errs errors.JoinedError

	for _, child := range other {
		newM, err := m.mergeEntry(child)
		if err != nil {
			errs = append(errs, errors.LocateByKey(child.String(), child.pos, err))

			continue
		}

		m = newM
	}

	return m, errs.Compact()
}

func (m Mapping) mergeEntry(value mappingEntry) (Mapping, error) {
	for i, entry := range m {
		if entry.Key.EqualsNode(value.Key) {
			newEntry, err := entry.merge(value)
			if err != nil {
				return nil, err
			}

			m[i] = newEntry

			return m, nil
		}
	}

	return append(m, value), nil
}

func (m mappingEntry) Bool() (value, ok bool) {
	return m.Value.Bool()
}

func (m mappingEntry) Int() (int64, bool) {
	return m.Value.Int()
}

func (m mappingEntry) Float() (float64, bool) {
	return m.Value.Float()
}

func (m mappingEntry) String() string {
	return m.Key.String()
}

func (m mappingEntry) EqualsNode(other Node) bool {
	otherEntry, ok := other.(mappingEntry)
	if !ok {
		return false
	}

	return m.Key.EqualsNode(otherEntry.Key) && m.Value.EqualsNode(otherEntry.Value)
}

func (m mappingEntry) ToYAML() *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{m.Key.ToYAML(), m.Value.ToYAML()},
	}
}

func (m mappingEntry) merge(other mappingEntry) (mappingEntry, error) {
	switch value := m.Value.(type) {
	case Mapping:
		oMapping, ok := other.Value.(Mapping)
		if !ok {
			return mappingEntry{}, errors.ErrNotAMapping
		}

		newValue, err := value.Merge(oMapping)
		if err != nil {
			return mappingEntry{}, err
		}

		m.Value = newValue

		return m, nil
	case Sequence:
		oSequence, ok := other.Value.(Sequence)
		if !ok {
			return mappingEntry{}, errors.ErrNotASequence
		}

		value = append(value, oSequence...)
		m.Value = value

		return m, nil
	default:
		return mappingEntry{}, errors.ErrMultipleScalarValues
	}
}
