package yaml2yaml

import (
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

func doDirInclude(defs *definitions.Definitions, key, value *yaml.Node) (*yaml.Node, error) {
	errs := errors.Multi()

	filePath, err := evaluateString(defs, key)

	errs.Add(err)

	stat, err := os.Stat(filePath)
	if err != nil {
		errs.Add(err)

		return nil, errs
	}

	var result *yaml.Node

	if stat.IsDir() {
		result, err = readDirIncludes(defs, filePath)
	} else {
		result, err = readFileInclude(defs, filePath)
	}

	errs.Add(err)

	if value != nil && value.Tag != "!!null" {
		errs.Add(errors.ErrUnexpectedValue)
	}

	return result, errs
}

func readDirIncludes(defs *definitions.Definitions, path string) (*yaml.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return nil, err
	}

	var result *yaml.Node

	errs := errors.Multi()

	for _, entry := range entries {
		var value *yaml.Node

		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			value, err = readDirIncludes(defs, entryPath)
		} else {
			value, err = readFileInclude(defs, entryPath)
		}

		errs.Add(err)

		result, err = convert.MergeNodes(result, value)

		errs.Add(err)
	}

	return result, errs
}

func readFileInclude(defs *definitions.Definitions, path string) (*yaml.Node, error) {
	node, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return process(defs.Child(), "("+path+")", node)
}
