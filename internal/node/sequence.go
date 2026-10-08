package node

import (
	"errors"
	"strings"

	"go.yaml.in/yaml/v4"
)

type (
	SequenceSource []Source
	Sequence       []Node
)

func (s SequenceSource) Process(defs Definitions) (Node, error) {
	var errs []error

	result := make(Sequence, 0, len(s))

	for _, item := range s {
		processed, err := item.Process(defs)
		if err != nil {
			// TODO Locate error
			errs = append(errs, err)
		}

		if processed == nil {
			continue
		}

		result = append(result, processed)
	}

	return result, errors.Join(errs...)
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
