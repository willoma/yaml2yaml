package yaml2yaml

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

func (n *node) includeMapping(key, value *node) ([]mappingChild, error) {
	keyValue, kErr := key.evaluateValue()
	if kErr != nil {
		kErr = n.locatedError(kErr)
	}

	path := anyToString(keyValue)

	includedValue, rErr := n.readInclude(path)
	if rErr != nil {
		rErr = n.locatedError(rErr)
	}

	result, cErr := nodeToMappingChild(key.path, includedValue)
	if cErr != nil {
		cErr = n.locatedError(cErr)
	}

	var vErr error
	if value != nil && value.Tag != "!!null" {
		vErr = n.locatedError(fmt.Errorf("%w on include", ErrUnexpectedValue))
	}

	return result, joinErrors(kErr, rErr, cErr, vErr)
}

func (n *node) includeScalar() (*yaml.Node, error) {
	value, eErr := n.evaluateValue()
	path := anyToString(value)
	includedValue, rErr := n.readInclude(path)

	return includedValue, joinErrors(eErr, rErr)
}

func (n *node) readInclude(path string) (*yaml.Node, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, n.locatedError(err)
	}

	if stat.IsDir() {
		return n.readDirIncludes(path)
	}

	return n.readFileInclude(path)
}

func (n *node) readDirIncludes(path string) (*yaml.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, n.locatedError(err)
	}

	var (
		result *yaml.Node
		errs   []error
	)

	for _, entry := range entries {
		var (
			value *yaml.Node
			err   error
		)

		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			value, err = n.readDirIncludes(entryPath)
		} else {
			value, err = n.readFileInclude(entryPath)
		}

		if err != nil {
			errs = append(errs, err)
		}

		if value != nil {
			result, err = mergeNodes(result, value)
			if err != nil {
				errs = append(errs, n.locatedError(err))
			}
		}
	}

	return result, joinErrors(errs...)
}

func (n *node) readFileInclude(path string) (*yaml.Node, error) {
	node, err := readFile(path)
	if err != nil {
		return nil, n.locatedError(err)
	}

	return n.fileChild(path, node).process()
}
