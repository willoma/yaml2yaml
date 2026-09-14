package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"
)

// mappingChild is the result of evaluating a mapping child.
//
//   - kvs is not nil if child type is mapping entry or entries
//   - value is not nil if child type is single value
//   - both are nil if child is to be ignored
type mappingChild struct {
	path  string
	kvs   []*yaml.Node
	value *yaml.Node
}

func (n *node) evaluateMappingChild(key, value *node) (mappingChild, error) {
	if key.Kind != yaml.ScalarNode {
		keyResult, kErr := key.process()
		valueResult, vErr := value.process()

		return mappingChild{kvs: []*yaml.Node{keyResult, valueResult}}, joinErrors(kErr, vErr)
	}

	switch key.Tag {
	case "!define":
		return n.mappingChildDefine(key, value)
	case "!apply":
		return n.mappingChildApply(key, value)
	case "!if":
		return n.mappingChildIf(key, value)
	case "!for":
		return n.mappingChildFor(key, value)
	default:
		keyResult, kErr := key.process()
		valueResult, vErr := value.process()

		return mappingChild{kvs: []*yaml.Node{keyResult, valueResult}}, joinErrors(kErr, vErr)
	}
}

func (n *node) mappingChildDefine(key, value *node) (mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	dErr := n.setDefinition(name, value.Node)
	if dErr != nil {
		dErr = key.locatedError(dErr)
	}

	return mappingChild{}, joinErrors(kErr, dErr)
}
func (n *node) mappingChildApply(key, value *node) (mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	node, dErr := n.getDefinition(name)
	if dErr != nil {
		return mappingChild{}, joinErrors(kErr, dErr)
	}

	resultValue := n.keyedChild(key.Tag, key.Value, node)

	errs := []error{kErr}

	if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content); i += 2 {
			keyAny, err := nodeToAny(value.Content[i])
			if err != nil {
				errs = append(errs, err)
			}

			key := anyToString(keyAny)

			resultValue.setDefinition(key, value.Content[i+1])
		}
	} else {
		resultValue.setDefinition("_", value.Node)
	}

	result, rErr := resultValue.toMappingChild()

	errs = append(errs, rErr)

	return result, joinErrors(errs...)
}

func (n *node) mappingChildIf(key, value *node) (mappingChild, error) {
	expr, eErr := key.resolveExpression(key.Value)

	eval, ok := expr.(bool)
	if !ok {
		err := key.locatedError(fmt.Errorf("%w: !if expression must evaluate to a boolean, got %q", ErrInvalid, expr))

		return mappingChild{}, joinErrors(eErr, err)
	}

	if !eval {
		return mappingChild{}, eErr
	}

	result, vErr := value.toMappingChild()

	return result, joinErrors(eErr, vErr)
}

func (n *node) mappingChildFor(key, value *node) (mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	seq, dErr := n.getDefinition(name)
	if dErr != nil {
		return mappingChild{}, joinErrors(kErr, dErr)
	}

	if seq.Kind != yaml.SequenceNode {
		return mappingChild{}, fmt.Errorf("%w: !for expression must be a sequence", ErrInvalid)
	}

	if value.Kind != yaml.SequenceNode && value.Kind != yaml.MappingNode {
		return mappingChild{}, fmt.Errorf("%w: !for content must be a sequence or a mapping", ErrInvalid)
	}

	errs := []error{kErr}

	result := &yaml.Node{
		Kind:    value.Kind,
		Tag:     value.Tag,
		Content: make([]*yaml.Node, 0, len(seq.Content)*len(value.Content)),
		Line:    value.Line,
		Column:  value.Column,
	}

	for i, seqNode := range seq.Content {
		child := n.indexedChild(i, value.Node)
		idx := scalarNode(seqNode, "!!int", strconv.Itoa(i))

		child.setDefinition("i", idx)
		child.setDefinition("_", seqNode)

		childResult, err := child.process()
		if err != nil {
			errs = append(errs, err)
		}

		result.Content = append(result.Content, childResult.Content...)
	}

	if value.Kind == yaml.MappingNode {
		return mappingChild{kvs: result.Content}, joinErrors(errs...)
	}

	return mappingChild{value: result}, joinErrors(errs...)
}
