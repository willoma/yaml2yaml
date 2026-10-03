package directives

import (
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
	"github.com/willoma/yaml2yaml/internal/read"
)

const NameInclude = "!include"

type Include struct {
	BasePath  string
	Path      node.Source
	ParseNode func(basePath string, src *yaml.Node) node.Source
}

func (i *Include) Process(defs node.Definitions) (node.Node, error) {
	pathNode, pErr := i.Path.Process(defs)

	path := filepath.Join(i.BasePath, pathNode.String())

	stat, sErr := os.Stat(path)
	if sErr != nil {
		return nil, errors.Join(pErr, sErr)
	}

	var (
		result node.Node
		rErr   error
	)

	if stat.IsDir() {
		result, rErr = i.dir(defs, path)
	} else {
		result, rErr = i.file(defs, path)
	}

	return result, errors.Join(pErr, rErr)
}

func (i *Include) dir(defs node.Definitions, path string) (node.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return nil, err
	}

	var result node.Node

	var errs []error

	for _, entry := range entries {
		var value node.Node

		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			value, err = i.dir(defs, entryPath)
		} else {
			value, err = i.file(defs, entryPath)
		}

		if err != nil {
			errs = append(errs, err)
		}

		result, err = result.Merge(value)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return result, errors.Join(errs...)
}

func (i *Include) file(defs node.Definitions, path string) (node.Node, error) {
	src, err := read.File(path)
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(path)

	sourceNode := i.ParseNode(dirPath, src)

	if sourceNode == nil {
		//nolint:nilnil // Nil as a value is okay here.
		return nil, nil
	}

	return sourceNode.Process(defs)
}
