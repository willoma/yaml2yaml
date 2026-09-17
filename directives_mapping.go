package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"
)

func (n *node) mappingDirective(key, value *node) (mappingChild, error) {
	switch key.Tag {
	case dirDefine:
		return n.defineMapping(key, value)
	case dirApply:
		return n.applyMapping(key, value)
	case dirIf:
		return n.ifMapping(key, value)
	case dirFor:
		return n.forMapping(key, value)
	default:
		return n.mappingNoDirective(key, value)
	}
}

func (n *node) defineMapping(key, value *node) (mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	dErr := n.setDefinition(name, value.Node)
	if dErr != nil {
		dErr = key.locatedError(dErr)
	}

	return mappingChild{}, joinErrors(kErr, dErr)
}

func (n *node) applyMapping(key, value *node) (mappingChild, error) {
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

			if err := resultValue.setDefinition(key, value.Content[i+1]); err != nil {
				errs = append(errs, err)
			}
		}
	} else {
		//nolint:errcheck // resultValue has just been created, _ is not already defined
		resultValue.setDefinition("_", value.Node)
	}

	result, rErr := resultValue.toMappingChild()

	errs = append(errs, rErr)

	return result, joinErrors(errs...)
}

func (n *node) ifMapping(key, value *node) (mappingChild, error) {
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

func (n *node) forMapping(key, value *node) (mappingChild, error) {
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

		//nolint:errcheck // child has just been created, i is not already defined
		child.setDefinition("i", idx)

		//nolint:errcheck // child has just been created, _ is not already defined
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
