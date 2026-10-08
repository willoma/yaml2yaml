package node

import (
	"errors"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"

	yerrors "github.com/willoma/yaml2yaml/internal/errors"
)

type (
	MappingSource []Source
	Mapping       []MappingEntry
)

func (m MappingSource) Process(parentDefs Definitions) (Node, error) {
	var errs []error

	processed := make([]Node, 0, len(m))
	defs := childDefs(parentDefs)

	for _, item := range m {
		itemResult, err := item.Process(defs)
		if err != nil {
			if errors.Is(err, yerrors.ErrEmpty) {
				continue
			}

			// TODO Locate error
			errs = append(errs, err)
		}

		processed = append(processed, itemResult)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	switch len(processed) {
	case 0:
		// No children, pass definitions to parent and ignore the mapping altogether.
		defs.passToParent()

		return nil, yerrors.ErrEmpty
	case 1:
		return processed[0], nil
	}

	var (
		result Node = Mapping{}
		err    error
	)

	for _, child := range processed {
		result, err = result.Merge(child)
		if err != nil {
			// TODO Locate error
			errs = append(errs, err)
		}
	}

	return result, errors.Join(errs...)
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
	result := make([]string, 0, len(m))
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
		var errs []error

		for _, child := range other {
			newM, err := m.Merge(child)
			if err != nil {
				// TODO Locate error
				errs = append(errs, err)
			}

			//nolint:forcetypeassert // The result of [Mapping.Merge] is always a [Mapping].
			m = newM.(Mapping)
		}

		return m, errors.Join(errs...)
	default:
		return m, yerrors.ErrUnkeyedValue
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
