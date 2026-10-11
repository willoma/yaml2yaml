package node

import (
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

type (
	MappingSource struct {
		pos   [2]int
		Value []Source
	}

	Mapping []MappingEntry

	processedMappingEntry struct {
		src   Source
		value Node
		err   error
	}
)

func NewMappingSource(pos [2]int, length int) *MappingSource {
	return &MappingSource{
		pos:   pos,
		Value: make([]Source, 0, length),
	}
}

func (p processedMappingEntry) LocateOwnError() error {
	return errors.LocateByKey(p.src.String(), p.src.Pos(), p.err)
}

func (p processedMappingEntry) LocateError(err error) error {
	return errors.LocateByKey(p.src.String(), p.src.Pos(), err)
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

func (m *MappingSource) Process(parentDefs Definitions) (Node, error) {
	defs := childDefs(parentDefs)
	processed := make([]processedMappingEntry, 0, len(m.Value))

	var mappingsCount, othersCount, sequencesCount int

	for _, item := range m.Value {
		itemResult, err := item.Process(defs)
		if err == nil {
			switch itemResult.(type) {
			case Mapping, MappingEntry:
				mappingsCount++
			case Sequence:
				sequencesCount++
			default:
				othersCount++
			}
		}

		processed = append(processed, processedMappingEntry{item, itemResult, err})
	}

	errs := errors.NewJoined()

	var result Node

	switch {
	case mappingsCount > 0:
		result = m.postprocess(Mapping{}, processed, errs)
	case sequencesCount > 0:
		result = m.postprocess(Sequence{}, processed, errs)
	case othersCount > 0:
		result = m.postprocess(Bool{}, processed, errs)
	default:
		// No content: ignore the result, it would not have changed anyway...
		m.postprocess(Bool{}, processed, errs)
		defs.passToParent()
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.ErrEmpty
	}

	return result, nil
}

func (m *MappingSource) postprocess(node Node, processed []processedMappingEntry, errs *errors.Joined) Node {
	var err error

	for _, entry := range processed {
		if entry.err != nil {
			errs.Add(entry.LocateOwnError())

			continue
		}

		node, err = node.Merge(entry.value)
		if err != nil {
			errs.Add(entry.LocateError(err))
		}
	}

	return node
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

func (m Mapping) Equals(other Node) bool {
	if otherMap, ok := other.(Mapping); ok {
		return slices.Equal(m, otherMap)
	}

	return false
}

func (m Mapping) Merge(other Node) (Node, error) {
	switch other := other.(type) {
	case MappingEntry:
		return m.addEntry(other)
	case Mapping:
		errs := errors.NewJoined()

		for _, child := range other {
			newM, err := m.Merge(child)
			errs.Add(err)

			//nolint:forcetypeassert // The result of [Mapping.Merge] is always a [Mapping].
			m = newM.(Mapping)
		}

		return m, errs.Compact()
	default:
		return m, errors.ErrUnkeyedValue
	}
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

func (m Mapping) addEntry(processed MappingEntry) (Mapping, error) {
	for i, entry := range m {
		if !entry.Key.Equals(processed.Key) {
			continue
		}

		valueMap, ok := entry.Value.(Mapping)
		if !ok {
			m[i].Value = processed.Value

			return m, nil
		}

		var err error

		m[i].Value, err = valueMap.Merge(processed.Value)

		return m, err
	}

	m = append(m, processed)

	return m, nil
}
