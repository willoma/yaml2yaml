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

func (s *SequenceSource) EqualsSource(other Source) bool {
	otherSeq, ok := other.(*SequenceSource)
	if !ok {
		return false
	}

	if len(s.Value) != len(otherSeq.Value) {
		return false
	}

	for i := range s.Value {
		if !s.Value[i].EqualsSource(otherSeq.Value[i]) {
			return false
		}
	}

	return true
}

func (s *SequenceSource) Process(defs Definitions) (Node, error) {
	var errs errors.JoinedError

	result := make(Sequence, 0, len(s.Value))

	for i, item := range s.Value {
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
