package yaml2yaml

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"
)

func (n *node) mappingDirective(key, value *node) ([]mappingChild, error) {
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
		keyResult, kErr := key.process()
		valueResult, vErr := value.process()

		return []mappingChild{{path: key.path, key: keyResult, value: valueResult}}, joinErrors(kErr, vErr)
	}
}

func (n *node) defineMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	dErr := n.setDefinition(name, value.Node)
	if dErr != nil {
		dErr = key.locatedError(dErr)
	}

	return nil, joinErrors(kErr, dErr)
}

func (n *node) applyMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	node, dErr := n.getDefinition(name)
	if dErr != nil {
		return nil, joinErrors(kErr, dErr)
	}

	resultValue := n.keyedChild(key.Tag, key.Value, node)
	aErr := resultValue.applyDefinitions(value.Node)

	processedResultValue, vErr := resultValue.process()
	result, cErr := nodeToMappingChild(key.path, processedResultValue)

	return result, joinErrors(kErr, aErr, vErr, cErr)
}

func (n *node) ifMapping(key, value *node) ([]mappingChild, error) {
	expr, eErr := key.resolveExpression(key.Value)

	if !anyToBool(expr) {
		return nil, eErr
	}

	resultValue, vErr := value.process()
	result, cErr := nodeToMappingChild(key.path, resultValue)

	return result, joinErrors(eErr, vErr, cErr)
}

func (n *node) forMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	name := anyToString(keyValue)

	seq, dErr := n.getDefinition(name)
	if dErr != nil {
		return nil, joinErrors(kErr, dErr)
	}

	if seq.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%w: !for expression must be a sequence", ErrInvalid)
	}

	if value.Kind != yaml.SequenceNode && value.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%w: !for content must be a sequence or a mapping", ErrInvalid)
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

		if err := child.applyDefinitions(seqNode); err != nil {
			errs = append(errs, err)
		}

		childResult, err := child.process()
		if err != nil {
			errs = append(errs, err)
		}

		result.Content = append(result.Content, childResult.Content...)
	}

	children, cErr := nodeToMappingChild(name, result)
	errs = append(errs, cErr)

	return children, joinErrors(errs...)
}
