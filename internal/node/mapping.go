package node

import (
	oerrors "errors"
	"iter"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

type (
	MappingSource struct {
		pos   [2]int
		value []Source
	}

	MappingEntrySource struct {
		key   Source
		value Source
	}

	Mapping []MappingEntry

	MappingEntry struct {
		pos   [2]int
		key   Node
		value Node
	}
)

func NewMappingSource(pos [2]int, length int) *MappingSource {
	return &MappingSource{
		pos:   pos,
		value: make([]Source, 0, length),
	}
}

func (m *MappingSource) Pos() [2]int {
	return m.pos
}

func (m *MappingSource) String() string {
	result := make([]string, len(m.value))
	for i, item := range m.value {
		result[i] = item.String()
	}

	return "{" + strings.Join(result, ", ") + "}"
}

func (m *MappingSource) Size() int {
	return len(m.value)
}

func (m *MappingSource) EqualsSource(other Source) bool {
	otherMap, ok := other.(*MappingSource)
	if !ok {
		return false
	}

	if len(m.value) != len(otherMap.value) {
		return false
	}

	for i := range m.value {
		if !m.value[i].EqualsSource(otherMap.value[i]) {
			return false
		}
	}

	return true
}

func (m *MappingSource) Process(parentDefs Definitions) (Node, error) {
	defs := childDefs(parentDefs)
	processed := NewProcessedEntries(len(m.value))

	for _, item := range m.value {
		itemResult, err := item.Process(defs)

		processed.Append(item.String(), item.Pos(), itemResult, err)
	}

	result, err := processed.Result()
	if oerrors.Is(err, errors.ErrEmpty) {
		defs.passToParent()
	}

	return result, err
}

func (m *MappingSource) SourceMap(defs Definitions) (map[string]Source, error) {
	defsMap := make(map[string]Source, len(m.value))

	var errs errors.JoinedError

	for _, item := range m.value {
		switch dItem := item.(type) {
		case *MappingEntrySource:
			key, err := dItem.key.Process(defs)
			if err != nil {
				errs = append(errs, errors.LocateByKey(item.String(), item.Pos(), err))

				continue
			}

			defsMap[key.String()] = dItem.value
		case SourceMapper:
			itemDefs, err := dItem.SourceMap(defs)
			if err != nil {
				errs = append(errs, errors.LocateByKey(item.String(), item.Pos(), err))

				continue
			}

			for name, value := range itemDefs {
				// Only use the definition if it does not already exist (avoid overwriting a first-level definition).
				if _, exists := defsMap[name]; !exists {
					defsMap[name] = value
				}
			}
		}
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return defsMap, nil
}

func (m *MappingSource) Append(item Source, overwrite bool) {
	if newEntry, ok := item.(*MappingEntrySource); ok {
		for _, item := range m.value {
			entry, ok := item.(*MappingEntrySource)
			if !ok {
				continue
			}

			if entry.key.EqualsSource(newEntry.key) {
				if overwrite {
					entry.value = newEntry.value
				}

				return
			}
		}
	}

	m.value = append(m.value, item)
}

func (m *MappingSource) Clone() *MappingSource {
	return &MappingSource{
		pos:   m.pos,
		value: slices.Clone(m.value),
	}
}

func (m *MappingSource) Merge(other *MappingSource) {
	for _, item := range other.value {
		m.Append(item, true)
	}
}

func (m *MappingSource) AppendAlias(item Source) bool {
	switch item := item.(type) {
	case *Alias:
		aliasMap, ok := item.target.(*MappingSource)
		if !ok {
			return false
		}

		for _, item := range aliasMap.value {
			m.Append(item, false)
		}

		return true

	case *SequenceSource:
		var items []Source

		for _, entry := range item.value {
			alias, ok := entry.(*Alias)
			if !ok {
				return false
			}

			aliasMap, ok := alias.target.(*MappingSource)
			if !ok {
				return false
			}

			items = append(items, aliasMap.value...)
		}

		for _, item := range items {
			m.Append(item, false)
		}

		return true

	default:
		return false
	}
}

func (m *MappingSource) Iter() iter.Seq[Source] {
	return func(yield func(Source) bool) {
		for _, item := range m.value {
			if !yield(item) {
				return
			}
		}
	}
}

func NewMappingEntrySource(key, value Source) *MappingEntrySource {
	return &MappingEntrySource{
		key:   key,
		value: value,
	}
}

func (m *MappingEntrySource) Pos() [2]int {
	return m.key.Pos()
}

func (m *MappingEntrySource) String() string {
	return m.key.String()
}

func (m *MappingEntrySource) EqualsSource(other Source) bool {
	if otherEntry, ok := other.(*MappingEntrySource); ok {
		return m.key.EqualsSource(otherEntry.key) && m.value.EqualsSource(otherEntry.value)
	}

	return false
}

func (m *MappingEntrySource) Process(defs Definitions) (Node, error) {
	var errs errors.JoinedError

	key, err := m.key.Process(defs)
	if err != nil {
		errs = errs.Append(err)
	}

	value, err := m.value.Process(defs)
	if err != nil {
		errs = errs.Append(err)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return MappingEntry{pos: m.Pos(), key: key, value: value}, nil
}

func (m *MappingEntrySource) Value() Source {
	return m.value
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
		if item.key == nil || item.value == nil {
			continue
		}

		content = append(content, item.key.ToYAML(), item.value.ToYAML())
	}

	return &yaml.Node{Kind: yaml.MappingNode, Content: content}
}

func (m Mapping) GetValue(key string) (Node, error) {
	for _, entry := range m {
		if entry.key.String() == key {
			return entry.value, nil
		}
	}

	return nil, errors.ErrEmpty
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

func (m Mapping) mergeEntry(value MappingEntry) (Mapping, error) {
	for i, entry := range m {
		if entry.key.EqualsNode(value.key) {
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

func (m MappingEntry) Bool() (value, ok bool) {
	return m.value.Bool()
}

func (m MappingEntry) Int() (int64, bool) {
	return m.value.Int()
}

func (m MappingEntry) Float() (float64, bool) {
	return m.value.Float()
}

func (m MappingEntry) String() string {
	return m.key.String()
}

func (m MappingEntry) EqualsNode(other Node) bool {
	otherEntry, ok := other.(MappingEntry)
	if !ok {
		return false
	}

	return m.key.EqualsNode(otherEntry.key) && m.value.EqualsNode(otherEntry.value)
}

func (m MappingEntry) ToYAML() *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{m.key.ToYAML(), m.value.ToYAML()},
	}
}

func (m MappingEntry) Key() Node {
	return m.key
}

func (m MappingEntry) Value() Node {
	return m.value
}

func (m MappingEntry) merge(other MappingEntry) (MappingEntry, error) {
	switch value := m.value.(type) {
	case Mapping:
		oMapping, ok := other.value.(Mapping)
		if !ok {
			return MappingEntry{}, errors.ErrNotAMapping
		}

		newValue, err := value.Merge(oMapping)
		if err != nil {
			return MappingEntry{}, err
		}

		m.value = newValue

		return m, nil
	case Sequence:
		oSequence, ok := other.value.(Sequence)
		if !ok {
			return MappingEntry{}, errors.ErrNotASequence
		}

		value = append(value, oSequence...)
		m.value = value

		return m, nil
	default:
		return MappingEntry{}, errors.ErrMultipleScalarValues
	}
}
