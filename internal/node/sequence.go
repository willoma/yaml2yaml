package node

import (
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

type (
	SequenceSource struct {
		pos   [2]int
		Value []Source
	}

	Sequence []Node
)

func NewSequenceSource(pos [2]int, length int) *SequenceSource {
	return &SequenceSource{
		pos:   pos,
		Value: make([]Source, 0, length),
	}
}

func (s *SequenceSource) Pos() [2]int {
	return s.pos
}

func (s *SequenceSource) String() string {
	result := make([]string, len(s.Value))
	for i, item := range s.Value {
		result[i] = item.String()
	}

	return "[" + strings.Join(result, ", ") + "]"
}

func (s *SequenceSource) Process(defs Definitions) (Node, error) {
	errs := errors.NewJoined()

	result := make(Sequence, 0, len(s.Value))

	for i, item := range s.Value {
		processed, err := item.Process(defs)
		if err != nil {
			errs.Add(errors.LocateByIndex(i, s.pos, err))

			continue
		}

		result = append(result, processed)
	}

	return result, errs.Compact()
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

func (s Sequence) Equals(other Node) bool {
	otherSeq, ok := other.(Sequence)
	if !ok {
		return false
	}

	if len(s) != len(otherSeq) {
		return false
	}

	for i := range s {
		if !s[i].Equals(otherSeq[i]) {
			return false
		}
	}

	return true
}

func (s Sequence) Merge(other Node) (Node, error) {
	if otherSeq, ok := other.(Sequence); ok {
		s = append(s, otherSeq...)
	} else {
		s = append(s, other)
	}

	return s, nil
}

func (s Sequence) ToYAML() *yaml.Node {
	content := make([]*yaml.Node, len(s))
	for i, item := range s {
		content[i] = item.ToYAML()
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Content: content}
}
