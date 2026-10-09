package directives

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
	"github.com/willoma/yaml2yaml/internal/read"
)

const NameInclude = "!include"

type (
	ParseNode func(basePath string, src *yaml.Node) (node.Source, error)

	dirInclude struct {
		basePath  string
		path      node.Source
		parseNode ParseNode
	}
)

func Include(basePath string, path node.Source, parseNode ParseNode) *dirInclude {
	return &dirInclude{basePath: basePath, path: path, parseNode: parseNode}
}

func (i *dirInclude) Pos() [2]int {
	return i.path.Pos()
}

func (i *dirInclude) String() string {
	return NameInclude + " " + i.path.String()
}

func (i *dirInclude) Process(defs node.Definitions) (node.Node, error) {
	pathNode, err := i.path.Process(defs)
	if err != nil {
		return nil, fmt.Errorf("path: %w", err)
	}

	path := pathNode.String()

	if !filepath.IsAbs(path) {
		path = filepath.Join(i.basePath, pathNode.String())
	}

	stat, err := os.Stat(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return nil, err
	}

	if stat.IsDir() {
		return i.dir(defs, path)
	}

	return i.file(defs, path)
}

func (i *dirInclude) dir(defs node.Definitions, path string) (node.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return nil, err
	}

	var result node.Node

	errs := errors.NewJoined()

	for _, entry := range entries {
		var value node.Node

		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			value, err = i.dir(defs, entryPath)
		} else {
			value, err = i.file(defs, entryPath)
		}

		if err != nil {
			errs.Add(err)

			continue
		}

		result, err = result.Merge(value)
		errs.Add(err)
	}

	if result == nil {
		return nil, errors.ErrEmpty
	}

	return result, errs.Compact()
}

func (i *dirInclude) file(defs node.Definitions, path string) (node.Node, error) {
	src, err := read.File(path)
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(path)

	sourceNode, err := i.parseNode(dirPath, src)
	if err != nil {
		return nil, err
	}

	if sourceNode == nil {
		return nil, errors.ErrEmpty
	}

	return sourceNode.Process(defs)
}
