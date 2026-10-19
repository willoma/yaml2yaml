package directives

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
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

func (i *dirInclude) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirInclude); ok {
		return i.basePath == other.basePath && i.path.EqualsSource(other.path)
	}

	return false
}

func (i *dirInclude) Process(defs node.Definitions) (node.Node, error) {
	path, isDir, err := i.getPath(defs)
	if err != nil {
		return nil, err
	}

	if isDir {
		return i.dirNode(defs, path)
	}

	return i.fileNode(defs, path)
}

func (i *dirInclude) SourceMap(defs node.Definitions) (map[string]node.Source, error) {
	path, isDir, err := i.getPath(defs)
	if err != nil {
		return nil, err
	}

	if isDir {
		return i.dirSourceMap(defs, path)
	}

	return i.fileSourceMap(defs, path)
}

func (i *dirInclude) getPath(defs node.Definitions) (path string, isDir bool, err error) {
	pathNode, err := i.path.Process(defs)
	if err != nil {
		return "", false, fmt.Errorf("path: %w", err)
	}

	path = pathNode.String()

	if !filepath.IsAbs(path) {
		path = filepath.Join(i.basePath, pathNode.String())
	}

	stat, err := os.Stat(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return "", false, err
	}

	return path, stat.IsDir(), nil
}

func (i *dirInclude) dirNode(defs node.Definitions, path string) (node.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return nil, err
	}

	processed := node.NewProcessedEntries(len(entries))

	for _, entry := range entries {
		var value node.Node

		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			value, err = i.dirNode(defs, entryPath)
		} else {
			value, err = i.fileNode(defs, entryPath)
		}

		processed.Append(entry.Name(), i.Pos(), value, err)
	}

	return processed.Result()
}

func (i *dirInclude) dirSourceMap(defs node.Definitions, path string) (map[string]node.Source, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		//nolint:wrapcheck // Original error is OK.
		return nil, err
	}

	var errs errors.JoinedError

	result := make(map[string]node.Source)

	for _, entry := range entries {
		var thisResult map[string]node.Source

		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			thisResult, err = i.dirSourceMap(defs, entryPath)
		} else {
			thisResult, err = i.fileSourceMap(defs, entryPath)
		}

		if err != nil {
			errs = append(errs, errors.LocateByKey(entry.Name(), i.Pos(), err))

			continue
		}

		maps.Copy(result, thisResult)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}

func (i *dirInclude) fileNode(defs node.Definitions, path string) (node.Node, error) {
	sourceNode, err := i.fileSource(path)
	if err != nil {
		return nil, err
	}

	return sourceNode.Process(defs)
}

func (i *dirInclude) fileSourceMap(defs node.Definitions, path string) (map[string]node.Source, error) {
	sourceNode, err := i.fileSource(path)
	if err != nil {
		return nil, err
	}

	sourceMapper, ok := sourceNode.(node.SourceMapper)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	return sourceMapper.SourceMap(defs)
}

func (i *dirInclude) fileSource(path string) (node.Source, error) {
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

	return sourceNode, nil
}
