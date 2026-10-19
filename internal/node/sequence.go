package node

import (
	"iter"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

type (
	SequenceSource struct {
		pos   [2]int
		value []Source
	}

	Sequence []Node
)

func NewSequenceSource(pos [2]int, length int) *SequenceSource {
	return &SequenceSource{
		pos:   pos,
		value: make([]Source, 0, length),
	}
}

func (s *SequenceSource) Pos() [2]int {
	return s.pos
}

func (s *SequenceSource) String() string {
	result := make([]string, len(s.value))
	for i, item := range s.value {
		result[i] = item.String()
	}

	return "[" + strings.Join(result, ", ") + "]"
}

func (s *SequenceSource) Append(item Source) {
	s.value = append(s.value, item)
}

func (s *SequenceSource) Size() int {
	return len(s.value)
}

func (s *SequenceSource) Get(i int) (Source, error) {
	if len(s.value) <= i {
		return nil, errors.ErrOutOfBounds
	}

	return s.value[i], nil
}

func (s *SequenceSource) EqualsSource(other Source) bool {
	otherSeq, ok := other.(*SequenceSource)
	if !ok {
		return false
	}

	if len(s.value) != len(otherSeq.value) {
		return false
	}

	for i := range s.value {
		if !s.value[i].EqualsSource(otherSeq.value[i]) {
			return false
		}
	}

	return true
}

func (s *SequenceSource) Process(defs Definitions) (Node, error) {
	var errs errors.JoinedError

	result := make(Sequence, 0, len(s.value))

	for i, item := range s.value {
		processed, err := item.Process(defs)
		if err != nil {
			errs = errs.Append(errors.LocateByIndex(i, s.pos, err))

			continue
		}

		result = append(result, processed)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SequenceSource) Iter() iter.Seq2[int, Source] {
	return func(yield func(int, Source) bool) {
		for i, item := range s.value {
			if !yield(i, item) {
				return
			}
		}
	}
}

func (s Sequence) Bool() (value, ok bool) {
	return len(s) > 0, true
}

func (s Sequence) Int() (int64, bool) {
	return 0, false
}

func (s Sequence) Float() (float64, bool) {
	return 0, false
}

func (s Sequence) Binary() ([]byte, bool) {
	return nil, false
}

func (s Sequence) String() string {
	result := make([]string, len(s))
	for i, item := range s {
		result[i] = item.String()
	}

	return "[" + strings.Join(result, ", ") + "]"
}

func (s Sequence) EqualsNode(other Node) bool {
	otherSeq, ok := other.(Sequence)
	if !ok {
		return false
	}

	if len(s) != len(otherSeq) {
		return false
	}

	for i := range s {
		if !s[i].EqualsNode(otherSeq[i]) {
			return false
		}
	}

	return true
}

func (s Sequence) ToYAML() *yaml.Node {
	content := make([]*yaml.Node, len(s))
	for i, item := range s {
		content[i] = item.ToYAML()
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Content: content}
}
