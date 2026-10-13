package node

import (
	"encoding/base64"
	"math"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

// scalar nodes are sources and results at the same time,
// because there is no special processing to do on them.
type (
	Bool struct {
		pos   [2]int
		value bool
	}

	Int struct {
		pos   [2]int
		value int64
	}

	Float struct {
		pos   [2]int
		value float64
	}

	String struct {
		pos   [2]int
		value string
	}

	Binary struct {
		pos   [2]int
		Value []byte
	}
)

func NewBool(pos [2]int, value bool) Bool {
	return Bool{pos: pos, value: value}
}

func (b Bool) Pos() [2]int {
	return b.pos
}

func (b Bool) Process(Definitions) (Node, error) {
	return b, nil
}

func (b Bool) Bool() (value, ok bool) {
	return b.value, true
}

func (b Bool) Int() (int64, bool) {
	if b.value {
		return 1, false
	}

	return 0, false
}

func (b Bool) Float() (float64, bool) {
	if b.value {
		return 1, false
	}

	return 0, false
}

func (b Bool) String() string {
	return strconv.FormatBool(b.value)
}

func (b Bool) Equals(other Node) bool {
	if otherBool, ok := other.(Bool); ok {
		return b == otherBool
	}

	return false
}

func (b Bool) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(b.value)}
}

func NewInt(pos [2]int, value int64) Int {
	return Int{pos: pos, value: value}
}

func (i Int) Pos() [2]int {
	return i.pos
}

func (i Int) Process(Definitions) (Node, error) {
	return i, nil
}

func (i Int) Bool() (value, ok bool) {
	return i.value > 0, true
}

func (i Int) Int() (int64, bool) {
	return i.value, true
}

func (i Int) Float() (float64, bool) {
	return float64(i.value), true
}

func (i Int) String() string {
	return strconv.FormatInt(i.value, 10)
}

func (i Int) Equals(other Node) bool {
	if otherInt, ok := other.Int(); ok {
		return i.value == otherInt
	}

	return false
}

func (i Int) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(i.value, 10)}
}

func NewFloat(pos [2]int, value float64) Float {
	return Float{pos: pos, value: value}
}

func (f Float) Pos() [2]int {
	return f.pos
}

func (f Float) Process(Definitions) (Node, error) {
	return f, nil
}

func (f Float) Bool() (value, ok bool) {
	return f.value > 0, true
}

func (f Float) Int() (int64, bool) {
	truncated := math.Trunc(f.value)

	if truncated == f.value {
		return int64(truncated), true
	}

	return 0, false
}

func (f Float) Float() (float64, bool) {
	return f.value, true
}

func (f Float) String() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

func (f Float) Equals(other Node) bool {
	if otherFloat, ok := other.Float(); ok {
		return f.value == otherFloat
	}

	return false
}

func (f Float) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: strconv.FormatFloat(f.value, 'f', -1, 64)}
}

func NewString(pos [2]int, value string) String {
	return String{pos: pos, value: value}
}

func (s String) Pos() [2]int {
	return s.pos
}

func (s String) Process(Definitions) (Node, error) {
	return s, nil
}

func (s String) Bool() (value, ok bool) {
	if value, err := strconv.ParseBool(s.value); err == nil {
		return value, true
	}

	return false, false
}

func (s String) Int() (int64, bool) {
	if value, err := strconv.ParseInt(s.value, 10, 64); err == nil {
		return value, true
	}

	return 0, false
}

func (s String) Float() (float64, bool) {
	if value, err := strconv.ParseFloat(s.value, 64); err == nil {
		return value, true
	}

	return 0, false
}

func (s String) String() string {
	return s.value
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

	return s.value == other.String()
}

func (s String) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Style: s.style(), Tag: "!!str", Value: s.value}
}

func (s String) style() yaml.Style {
	// Preserve YAML 1.1 compatibility for output format.
	if s.value == "on" || s.value == "off" {
		return yaml.SingleQuotedStyle
	}

	// Make multi-line strings easier to read
	if strings.Contains(s.value, "\n") {
		return yaml.LiteralStyle
	}

	return 0
}

func NewBinary(pos [2]int, value []byte) Binary {
	return Binary{pos: pos, Value: value}
}

func (b Binary) Pos() [2]int {
	return b.pos
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
	return string(b.Value)
}

func (b Binary) Equals(other Node) bool {
	return string(b.Value) == other.String()
}

func (b Binary) ToYAML() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!binary", Value: base64.StdEncoding.EncodeToString(b.Value)}
}
