package node

import (
	oerrors "errors"

	"github.com/willoma/yaml2yaml/errors"
)

type (
	ProcessedEntries struct {
		entries []processedEntry

		errorsCount    int
		mappingsCount  int
		scalarsCount   int
		sequencesCount int
	}

	processedEntry struct {
		index int
		key   string
		pos   [2]int
		value Node
		err   error
	}
)

func NewProcessedEntries(size int) *ProcessedEntries {
	return &ProcessedEntries{
		entries: make([]processedEntry, 0, size),
	}
}

func (p *ProcessedEntries) Append(key string, pos [2]int, value Node, err error) {
	if oerrors.Is(err, errors.ErrEmpty) {
		return
	}

	p.entries = append(p.entries, processedEntry{key: key, pos: pos, value: value, err: err})

	if err != nil {
		p.errorsCount++

		return
	}

	switch value.(type) {
	case Mapping, MappingEntry:
		p.mappingsCount++
	case Sequence:
		p.sequencesCount++
	default:
		p.scalarsCount++
	}
}

func (p *ProcessedEntries) AppendIndexed(index int, pos [2]int, value Node, err error) {
	p.entries = append(p.entries, processedEntry{index: index, pos: pos, value: value, err: err})

	if err != nil {
		p.errorsCount++

		return
	}

	switch value.(type) {
	case Mapping, MappingEntry:
		p.mappingsCount++
	case Sequence:
		p.sequencesCount++
	default:
		p.scalarsCount++
	}
}

func (p *ProcessedEntries) Result() (Node, error) {
	switch {
	case p.errorsCount > 0:
		return nil, p.errors()
	case p.mappingsCount > 0:
		return p.mapping()
	case p.sequencesCount > 0:
		return p.sequence()
	case p.scalarsCount == 1:
		return p.scalar()
	case p.scalarsCount > 1:
		return nil, p.multipleScalars()
	default:
		return nil, errors.ErrEmpty
	}
}

func (p *ProcessedEntries) errors() error {
	var errs errors.JoinedError

	for _, entry := range p.entries {
		if entry.err != nil {
			errs = errs.Append(entry.locateError(entry.err))
		}
	}

	return errs.Compact()
}

func (p *ProcessedEntries) mapping() (Node, error) {
	var errs errors.JoinedError

	result := make(Mapping, 0, p.mappingsCount)

	for _, entry := range p.entries {
		switch value := entry.value.(type) {
		case MappingEntry:
			newResult, err := result.mergeEntry(value)
			if err != nil {
				errs = errs.Append(entry.locateError(err))

				continue
			}

			result = newResult
		case Mapping:
			newResult, err := result.Merge(value)
			if err != nil {
				errs = errs.Append(entry.locateError(err))

				continue
			}

			result = newResult
		default:
			errs = errs.Append(entry.locateError(errors.ErrUnkeyedValue))
		}
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.ErrEmpty
	}

	return result, nil
}

func (p *ProcessedEntries) sequence() (Node, error) {
	var errs errors.JoinedError

	result := make(Sequence, 0, len(p.entries))

	for _, entry := range p.entries {
		entrySeq, ok := entry.value.(Sequence)
		if !ok {
			errs = append(errs, entry.locateError(errors.ErrNotASequence))

			continue
		}

		result = append(result, entrySeq...)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.ErrEmpty
	}

	return result, nil
}

func (p *ProcessedEntries) scalar() (Node, error) {
	for _, entry := range p.entries {
		return entry.value, nil
	}

	return nil, errors.ErrEmpty
}

func (p *ProcessedEntries) multipleScalars() error {
	errs := make(errors.JoinedError, len(p.entries))

	for i, entry := range p.entries {
		errs[i] = entry.locateError(errors.ErrMultipleScalarValues)
	}

	return errs.Compact()
}

func (p processedEntry) locateError(err error) error {
	if p.key != "" {
		return errors.LocateByKey(p.key, p.pos, err)
	}

	return errors.LocateByIndex(p.index, p.pos, err)
}
