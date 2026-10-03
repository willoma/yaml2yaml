package node

import (
	"encoding/base64"
	"math"
	"strconv"

	"go.yaml.in/yaml/v4"
)

// scalar nodes are sources and results at the same time,
// because there is no special processing to do on them.
type (
	Bool   bool
	Int    int64
	Float  float64
	String string
	Binary []byte
)

func (b Bool) Process(Definitions) (Node, error) {
	return b, nil
}

func (b Bool) Bool() (value, ok bool) {
	return bool(b), true
}

func (b Bool) Int() (int64, bool) {
	if b {
		return 1, false
	}

	return 0, false
}

func (b Bool) Float() (float64, bool) {
	if b {
		return 1, false
	}

	return 0, false
}

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}

func (b Bool) Equals(other Node) bool {
	if otherBool, ok := other.(Bool); ok {
		return b == otherBool
	}

	return false
}

func (b Bool) Merge(other Node) (Node, error) {
	return other, nil
}

func (b Bool) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(bool(b))}
}

func (i Int) Process(Definitions) (Node, error) {
	return i, nil
}

func (i Int) Bool() (value, ok bool) {
	return i > 0, true
}

func (i Int) Int() (int64, bool) {
	return int64(i), true
}

func (i Int) Float() (float64, bool) {
	return float64(i), true
}

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int) Equals(other Node) bool {
	if otherInt, ok := other.Int(); ok {
		return int64(i) == otherInt
	}

	return false
}

func (i Int) Merge(other Node) (Node, error) {
	return other, nil
}

func (i Int) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(int64(i), 10)}
}

func (f Float) Process(Definitions) (Node, error) {
	return f, nil
}

func (f Float) Bool() (value, ok bool) {
	return f > 0, true
}

func (f Float) Int() (int64, bool) {
	truncated := math.Trunc(float64(f))

	if truncated == float64(f) {
		return int64(truncated), true
	}

	return 0, false
}

func (f Float) Float() (float64, bool) {
	return float64(f), true
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (f Float) Equals(other Node) bool {
	if otherFloat, ok := other.Float(); ok {
		return float64(f) == otherFloat
	}

	return false
}

func (f Float) Merge(other Node) (Node, error) {
	return other, nil
}

func (f Float) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: strconv.FormatFloat(float64(f), 'f', -1, 64)}
}

func (s String) Process(Definitions) (Node, error) {
	return s, nil
}

func (s String) Bool() (value, ok bool) {
	if value, err := strconv.ParseBool(string(s)); err == nil {
		return value, true
	}

	return false, false
}

func (s String) Int() (int64, bool) {
	if value, err := strconv.ParseInt(string(s), 10, 64); err == nil {
		return value, true
	}

	return 0, false
}

func (s String) Float() (float64, bool) {
	if value, err := strconv.ParseFloat(string(s), 64); err == nil {
		return value, true
	}

	return 0, false
}

func (s String) String() string {
	return string(s)
}

func (s String) Equals(other Node) bool {
	if sInt, ok := s.Int(); ok {
		if otherInt, ok := other.Int(); ok {
			return sInt == otherInt
		}
	}

	if sFloat, ok := s.Float(); ok {
		if otherFloat, ok := other.Float(); ok {
			return sFloat == otherFloat
		}
	}

	if sBool, ok := s.Bool(); ok {
		if otherBool, ok := other.Bool(); ok {
			return sBool == otherBool
		}
	}

	return string(s) == other.String()
}

func (s String) Merge(other Node) (Node, error) {
	return other, nil
}

func (s String) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: string(s)}
}

func (b Binary) Process(Definitions) (Node, error) {
	return b, nil
}

func (b Binary) Bool() (value, ok bool) {
	return false, false
}

func (b Binary) Int() (int64, bool) {
	return 0, false
}

func (b Binary) Float() (float64, bool) {
	return 0, false
}

func (b Binary) String() string {
	return string(b)
}

func (b Binary) Equals(other Node) bool {
	return string(b) == other.String()
}

func (b Binary) Merge(other Node) (Node, error) {
	return other, nil
}

func (b Binary) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!binary", Value: base64.StdEncoding.EncodeToString(b)}
}
