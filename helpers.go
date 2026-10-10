package yaml2yaml

import (
	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
)

func locateMappingEntryError(i int, basePath string, src *yaml.Node, err error) error {
	key, kErr := parseNode(basePath, src)
	if kErr != nil {
		return errors.LocateByIndex(i/2, pos(src), err)
	}

	return errors.LocateByKey(key.String(), pos(src), err)
}

func pos(src *yaml.Node) [2]int {
	return [2]int{src.Line, src.Column}
}
